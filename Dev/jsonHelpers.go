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
