package main

import (
	"io"
	"log"
	"net/http"
)

type HttpHandler func(http.ResponseWriter, *http.Request)
type AuthServiceFunc func(string, http.ResponseWriter, *http.Request) Output
type ServiceFunc func(http.ResponseWriter, *http.Request) Output

type ServerHandle interface { // interface for handlers
	Handle() HttpHandler
	AcceptableMethods() acceptableMethods
}

func errorCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		err := r.(error)
		outmsg := structToJSON(Output{Success: false, Message: err.Error()})
		log.Printf("Error: %s", err.Error())
		io.WriteString(w, outmsg)
	}
}

func wrapHandler(f ServiceFunc, accepted acceptableMethods) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer errorCatcher(w)
		checkMethod(r.Method, accepted)
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		outStruct := f(w, r)
		outmsg := structToJSON(outStruct)
		if r.URL.Path == "/login" {
			cookie_lifetime := int(sessionDuration.Seconds())
			cookie := http.Cookie{Name: "Session", Value: outStruct.Result.(string), MaxAge: cookie_lifetime}
			http.SetCookie(w, &cookie)
		}
		io.WriteString(w, outmsg)
	}
}

type BasicService struct {
	handler ServiceFunc
	methods acceptableMethods
}

func (h BasicService) Handle() HttpHandler {
	return wrapHandler(h.handler, h.methods)
}

func (h BasicService) AcceptableMethods() acceptableMethods {
	return h.methods
}

type AuthService struct {
	handler AuthServiceFunc
	methods acceptableMethods
}

func authWrapper(toWrap AuthService) HttpHandler {
	return wrapHandler(func(w http.ResponseWriter, r *http.Request) Output {
		username := authUser(w, r)
		outmsg := toWrap.handler(username, w, r)
		if isUserConnected(username) {
			token := extendSession(username)
			cookie_lifetime := int(sessionDuration.Seconds())
			cookie := http.Cookie{Name: "Session", Value: token, MaxAge: cookie_lifetime}
			http.SetCookie(w, &cookie)
		} else {
			cookie := http.Cookie{Name: "Session", Value: "", MaxAge: -1}
			http.SetCookie(w, &cookie)
		}
		return outmsg
	}, toWrap.methods)
}

func (h AuthService) AcceptableMethods() acceptableMethods {
	return h.methods
}

func (h AuthService) Handle() HttpHandler {
	return authWrapper(h)
}

func authUser(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("Session")
	if err != nil {
		OnlyServerError("Need authentification")
	}
	username := verifySession(cookie.Value)
	return username
}

func LoginWithRemember(w http.ResponseWriter, r *http.Request) Output {
	var auth AuthJSON
	parseRequestToStruct(r, &auth)
	token := loginAccount(auth)
	return Output{Success: true, Message: "Logged in", Result: token}
}

func CreateAccount(w http.ResponseWriter, r *http.Request) Output {
	var auth AuthJSON
	parseRequestToStruct(r, &auth)
	addToDatabase(auth)
	return Output{Success: true, Message: "Account Created"}
}

func GetProfile(username string, w http.ResponseWriter, r *http.Request) Output {
	return Output{Success: true, Message: "Profile", Result: getUserData(username)}
}

func Logout(username string, w http.ResponseWriter, r *http.Request) Output {
	logoutAccount(username)
	return Output{Success: true, Message: "Logged out"}
}

func DeleteAccount(username string, w http.ResponseWriter, r *http.Request) Output {
	logoutAccount(username)
	deleteFromDatabase(username)
	return Output{Success: true, Message: "Account deleted"}
}

func RandomShitPost(w http.ResponseWriter, r *http.Request) Output {
	response, err := http.Get("https://api.thedailyshitpost.net/random")
	if err != nil {
		ServerRuntimeError("Error while getting shitpost", err)
	}
	var shitpost RandomShitPostJSON
	parseResponseToStruct(response, &shitpost)
	if shitpost.Error {
		OnlyServerError("Error while getting RandomShitpost")
	}
	return Output{Success: true, Message: "Random Shitpost", Result: shitpost.Url}
}

func HandlersMap() map[string]ServerHandle {
	handlers := make(map[string]ServerHandle)
	handlers["/login"] = BasicService{LoginWithRemember, acceptableMethods{Put: true}}
	handlers["/create_account"] = BasicService{CreateAccount, acceptableMethods{Post: true}}
	handlers["/get_profile"] = AuthService{GetProfile, acceptableMethods{Get: true}}
	handlers["/logout"] = AuthService{Logout, acceptableMethods{Put: true}}
	handlers["/delete_account"] = AuthService{DeleteAccount, acceptableMethods{Delete: true}}
	handlers["/random_shitpost"] = BasicService{RandomShitPost, acceptableMethods{Get: true}}
	return handlers
}
