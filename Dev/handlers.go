package main

import (
	"fmt"
	"io"
	/* 	"log" */
	"net/http"
	"time"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"GET"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	io.WriteString(w, "Hello, HTTP!\n")
}

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"GET"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	io.WriteString(w, "This is my website!\n")
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.URL)
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
	cookie := http.Cookie{Name: session, Value: "1", MaxAge: sessionTime, Expires: time.Now().Add(time.Hour)}
	http.SetCookie(w, &cookie)
	io.WriteString(w, fmt.Sprintf("This is %+v page!\n", auth))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.URL)
	err := checkMethod(r.Method, []string{"GET"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	cookie := http.Cookie{Name: session, MaxAge: -1, Value: "1", Expires: time.Unix(0, 0)}
	http.SetCookie(w, &cookie)
	io.WriteString(w, "logged out")
}

func HandlersMap() map[string]func(http.ResponseWriter, *http.Request) {
	handlers := make(map[string]func(http.ResponseWriter, *http.Request))
	handlers["/hello"] = Hello
	handlers["/"] = Root
	handlers["/login"] = Login
	handlers["/logout"] = Logout
	return handlers
}
