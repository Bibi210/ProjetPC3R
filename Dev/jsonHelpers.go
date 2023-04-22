package main

import (
	"encoding/json"

	"net/http"
)

/* func parseRequestToStruct(r *http.Request, t any) {
	bytes := parseRequestToBytes(r)
	bytesToStruct(bytes, t)
}
 */
func parseResponseToStruct(r *http.Response, t any) {
	bytes := parseResponseToBytes(r)
	bytesToStruct(bytes, t)
}

func bytesToStruct(bytes []byte, buffer any) {
	if json.Unmarshal(bytes, buffer) != nil {
		ServerRuntimeError("JSON with wrong format", nil)
	}
}

func stringToStruct(s handlerInput, buffer any) {
	bytesToStruct([]byte(s.msg), buffer)
}

func structToJSON(t any) string {
	obj, err := json.Marshal(t)
	if err != nil {
		ServerRuntimeError("Could not convert struct to JSON", err)
	}
	return string(obj)
}

/* Server Responses */

type OutputJSON struct {
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

type UserProfileJSON struct {
	UserID   int
	Username string
}
