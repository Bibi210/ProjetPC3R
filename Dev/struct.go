package main

import (
	"encoding/json"
	"io"
	"net/http"
	"unsafe"
)

func Clone(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := make([]byte, len(s))
	copy(b, s)
	return unsafe.String(&b[0], len(b))
}

func parseRequestToBytes(r *http.Request) ([]byte, error) {
	body := r.Body
	defer body.Close()
	var buffer []byte = make([]byte, r.ContentLength)
	read, err := body.Read(buffer)
	if err != nil && err != io.EOF {
		return []byte{}, ServerRuntimeError("Can't Parse Request Body", err)
	}
	return buffer[:read], nil
}

func parseRequestToStruct(r *http.Request, t any) error {
	bytes, err := parseRequestToBytes(r)
	if err != nil {
		return err
	}
	return ServerRuntimeError("Can't parse JSON", bytesToStruct(bytes, t))
}

func bytesToStruct(bytes []byte, buffer any) error {
	return json.Unmarshal(bytes, buffer)
}

func stringToStruct(s string, buffer any) error {
	return bytesToStruct([]byte(s), buffer)
}

func structToJSON(t any) (string, error) {
	obj, err := json.Marshal(t)
	if err != nil {
		return "", ServerRuntimeError("Could not convert struct to JSON", err)
	}
	return string(obj), nil
}

func sendStruct(w http.ResponseWriter, t any) error {
	str, err := structToJSON(t)
	if err != nil {
		return err
	}
	io.WriteString(w, str)
	return nil
}

type AuthJSON struct {
	Login string
	Mdp   string
}
type Username struct {
	Username string
}
