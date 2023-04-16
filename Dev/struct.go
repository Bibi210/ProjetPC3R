package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func parseRequestToBytes(r *http.Request) []byte {
	defer r.Body.Close()
	read, err := io.ReadAll(r.Body)
	if err != nil {
		ServerRuntimeError("Can't Parse Request Body", err)
	}
	return read
}

func parseResponseToBytes(r *http.Response) []byte {

	if r == nil || r.StatusCode != 200 {
		OnlyServerError("Request Failed with status code: " + string(r.Status))
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ServerRuntimeError("Can't Parse Response Body", err)
	}
	return body
}

func parseResponseToString(r *http.Response) string {
	return string(parseResponseToBytes(r))
}

func parseRequestToString(r *http.Request) string {
	return string(parseRequestToBytes(r))
}

func parseRequestToStruct(r *http.Request, t any) {
	bytes := parseRequestToBytes(r)
	bytesToStruct(bytes, t)
}

func parseResponseToStruct(r *http.Response, t any) {
	bytes := parseResponseToBytes(r)
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

type Output struct {
	Success bool
	Message string
	Result  any
}

type AuthJSON struct {
	Login string
	Mdp   string
}

type RandomShitPostJSON struct {
	Url   string
	Error bool
}
