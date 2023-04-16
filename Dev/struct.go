package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func parseRequestToBytes(r *http.Request) []byte {
	body := r.Body
	defer body.Close()
	var buffer []byte = make([]byte, r.ContentLength)
	read, err := body.Read(buffer)
	if err != nil && err != io.EOF {
		ServerRuntimeError("Can't Parse Request Body", err)
		return []byte{}
	}
	return buffer[:read]
}

func parseRequestToString(r *http.Request) string {
	return string(parseRequestToBytes(r))
}

func parseRequestToStruct(r *http.Request, t any) {
	bytes := parseRequestToBytes(r)
	bytesToStruct(bytes, t)
}

func bytesToStruct(bytes []byte, buffer any) {
	if json.Unmarshal(bytes, buffer) != nil {
		ServerRuntimeError("Can't parse JSON", nil)
	}
}

func stringToStruct(s string, buffer any) {
	bytesToStruct([]byte(s), buffer)
}

func structToJSON(t any) string {
	obj, err := json.Marshal(t)
	if err != nil {
		ServerRuntimeError("Could not convert struct to JSON", err)
	}
	return string(obj)
}

type Input struct {
	Content  any
	Username string
}

type Output struct {
	Success bool
	Message string
	Result  any
}

type JSONInput struct {
	Content any
}

type AuthJSON struct {
	Login string
	Mdp   string
}
type Username struct {
	Username string
}
