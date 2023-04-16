package main

import (
	"io"
	"log"
	"net/http"
)

func authUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	cookie, err := r.Cookie("Session")
	if err != nil {
		http.Error(w, OnlyServerError("Need authentification").Error(), http.StatusBadRequest)
		return "", false
	}
	username, err := verifySession(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return "", false
	}
	return username, true
}

func LoginWithRemember(w http.ResponseWriter, r *http.Request) {
	err := checkMethod(r.Method, acceptableMethods{Put: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	var auth AuthJSON
	err = parseRequestToStruct(r, &auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusPartialContent)
		return
	}
	tokenstr, err := loginAccount(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	cookie_lifetime := int(sessionDuration.Seconds())
	cookie := http.Cookie{Name: "Session", Value: tokenstr, MaxAge: cookie_lifetime}
	http.SetCookie(w, &cookie)

}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	err := checkMethod(r.Method, acceptableMethods{Post: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	var auth AuthJSON
	err = parseRequestToStruct(r, &auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusPartialContent)
		return
	}
	err = addToDatabase(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	err := checkMethod(r.Method, acceptableMethods{Get: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	username, ok := authUser(w, r)
	if !ok {
		return
	}
	profile, err := getUserData(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	sendStruct(w, profile)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := checkMethod(r.Method, acceptableMethods{Put: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	username, ok := authUser(w, r)
	if !ok {
		return
	}
	err = logoutAccount(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	cookie := http.Cookie{Name: "Session", Value: "", MaxAge: -1}
	http.SetCookie(w, &cookie)
	io.WriteString(w, "logged out")
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	err := checkMethod(r.Method, acceptableMethods{Delete: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	username, ok := authUser(w, r)
	if !ok {
		return
	}
	err = logoutAccount(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = deleteFromDatabase(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	cookie := http.Cookie{Name: "Session", Value: "", MaxAge: -1}
	http.SetCookie(w, &cookie)
	io.WriteString(w, "Account Deleted")
}

func wrapHandler(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		f(w, r)
	}
}

func HandlersMap() map[string]func(http.ResponseWriter, *http.Request) {
	handlers := make(map[string]func(http.ResponseWriter, *http.Request))
	handlers["/login"] = LoginWithRemember
	handlers["/logout"] = Logout
	handlers["/get_profile"] = GetProfile
	handlers["/create_account"] = CreateAccount
	handlers["/delete_account"] = DeleteAccount

	for k, v := range handlers {
		handlers[k] = wrapHandler(v)
	}
	return handlers
}
