package main

import (
	"encoding/json"
	"io"
	"net/http"
)

var session = "session"
var sessionTime = 60 * 60

func parseRequestToStruct(r *http.Request, t any) error {
	body := r.Body
	defer body.Close()
	var buffer []byte = make([]byte, r.ContentLength)
	read, err := body.Read(buffer)
	if err != nil && err != io.EOF {
		return err
	}
	return ServerRuntimeError("Can't parse JSON", json.Unmarshal(buffer[:read], t))
}

func structToJSON(t any) (string, error) {
	obj, err := json.Marshal(t)
	if err != nil {
		return "", ServerRuntimeError("Could not convert struct to JSON", err)
	}
	return string(obj), nil
}

type Auth struct {
	Login string
	Mdp   string
}
