package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type errorStruct struct {
	err error
	s   string
}

func (e *errorStruct) Error() string {
	if e.err == nil {
		return fmt.Sprintf("ServerSIDE error: %s", e.s)
	}
	return fmt.Sprintf("ServerSIDE error: %s due to -> Runtime error %s", e.s, e.err.Error())
}
func OnlyServerError(text string) {
	out := &errorStruct{nil, text}
	panic(out)
}
func ServerRuntimeError(text string, err error) {
	if err == nil {
		return
	}
	out := &errorStruct{err, text}
	panic(out)
}

func errorCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		err := r.(error)
		outmsg := structToJSON(Output{Success: false, Message: err.Error()})
		log.Printf("Error: %s", err.Error())
		io.WriteString(w, outmsg)
	}
}

type acceptableMethods struct {
	Get    bool
	Post   bool
	Put    bool
	Delete bool
}

func checkMethod(method string, acceptable acceptableMethods) {
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
