package main

import (
	"fmt"
	"log"
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
func OnlyServerError(text string) error {
	out := &errorStruct{nil, text}
	log.Println(out.Error())
	return out
}
func ServerRuntimeError(text string, err error) error {
	if err == nil {
		return nil
	}
	out := &errorStruct{err, text}
	log.Println(out.Error())
	return out
}

type acceptableMethods struct {
	Get    bool
	Post   bool
	Put    bool
	Delete bool
}

func checkMethod(method string, acceptable acceptableMethods) error {
	switch method {
	case "GET":
		if acceptable.Get {
			return nil
		}
	case "POST":
		if acceptable.Post {
			return nil
		}
	case "PUT":
		if acceptable.Put {
			return nil
		}
	case "DELETE":
		if acceptable.Delete {
			return nil
		}
	}
	return OnlyServerError(fmt.Sprintf("Method %s is not supported. Accepted %v\n", method, acceptable))
}
