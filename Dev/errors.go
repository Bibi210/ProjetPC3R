package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
)

type error_struct struct {
	err error
	s   string
}

func (e *error_struct) Error() string {
	if e.err == nil {
		return fmt.Sprintf("ServerSIDE error: %s", e.s)
	}
	return fmt.Sprintf("ServerSIDE error: %s due to -> Runtime error %s", e.s, e.err.Error())
}
func OnlyServerError(text string) {
	out := &error_struct{nil, text}
	panic(out)
}
func ServerRuntimeError(text string, err error) {
	if err == nil {
		return
	}
	out := &error_struct{err, text}
	panic(out)
}

func ErrorCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		err := r.(error)
		outmsg := structToJSON(OutputJSON{Success: false, Message: err.Error()})
		log.Printf("Error: %s", err.Error())
		io.WriteString(w, outmsg)
	}
}

func CleanCloser(db *sql.DB) {
	if r := recover(); r != nil {
		err := r.(error)
		closeDatabase(db)
		panic(err)
	}
	closeDatabase(db)
}

type AcceptableMethods struct {
	Get    bool
	Post   bool
	Put    bool
	Delete bool
}

func CheckMethod(method string, acceptable AcceptableMethods) {
	switch method {
	case "GET":
		if acceptable.Get {
			return
		}
	case "POST":
		if acceptable.Post {
			return
		}
	case "PUT":
		if acceptable.Put {
			return
		}
	case "DELETE":
		if acceptable.Delete {
			return
		}
	}
	OnlyServerError(fmt.Sprintf("Method %s is not supported. Accepted %v\n", method, acceptable))
}
