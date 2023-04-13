package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"GET"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	io.WriteString(w, "Hello, HTTP!\n")
}

func Root(w http.ResponseWriter, r *http.Request) {
	log.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"GET"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	io.WriteString(w, "This is my website!\n")
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"POST"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	var auth Auth
	err = parseRequestToStruct(r, &auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusPartialContent)
		return
	}
	tokenstr, err := loginAccount(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else {
		io.WriteString(w, fmt.Sprintf("This is your token : %s \n", tokenstr))
	}
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	log.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"POST"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	var auth Auth
	err = parseRequestToStruct(r, &auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusPartialContent)
		return
	}

	if addToDatabase(auth) != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"GET"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	io.WriteString(w, "logged out")
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	log.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"GET"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	tokenStr := r.Header["Token"][0]
	if r.Header["Token"] == nil {
		http.Error(w, OnlyServerError("Can not find token in header").Error(), http.StatusBadRequest)
		return
	}
	profile, err := getUserData(tokenStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	str, err := structToJSON(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusPartialContent)
		return
	}
	io.WriteString(w, str)
}

func HandlersMap() map[string]func(http.ResponseWriter, *http.Request) {
	handlers := make(map[string]func(http.ResponseWriter, *http.Request))
	handlers["/hello"] = Hello
	handlers["/"] = Root
	handlers["/login"] = Login
	handlers["/get_profile"] = GetProfile
	handlers["/create_account"] = CreateAccount
	handlers["/logout"] = Logout
	return handlers
}
