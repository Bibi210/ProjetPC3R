package main

import (
	"encoding/json"
)

func bytesToStruct(bytes []byte, buffer any) {
	if json.Unmarshal(bytes, buffer) != nil {
		ServerRuntimeError("JSON with wrong format", nil)
	}
}

func getClientMessage(s service_input, buffer any) {
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

type ResponseRandomShitPostJSON struct {
	Url   string
	Error bool
}

type ResponseUserProfileJSON struct {
	UserID   int
	Username string
}

/* ClientMessages */
type RequestAuthJSON struct {
	Login string
	Mdp   string
}

type RequestPublicUserProfileJSON struct {
	Username string
}
