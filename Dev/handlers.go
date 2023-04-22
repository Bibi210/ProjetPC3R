package main

import (
	"database/sql"
	"io"
	"log"
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

func parseRequestToString(r *http.Request) string {
	return string(parseRequestToBytes(r))
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

func getToken(r *http.Request) tokenString {
	cookie, err := r.Cookie("Session")
	if err != nil {
		return ""
	}
	return tokenString(cookie.Value)
}
func parseRequest(r *http.Request) handlerInput {
	return handlerInput{parseRequestToString(r), getToken(r)}
}

type tokenString string
type handlerInput struct {
	msg     string
	session tokenString
}
type username string
type handlerOutput struct {
	msg            OutputJSON
	newTokenString tokenString
}

type HttpHandler func(http.ResponseWriter, *http.Request)
type ServiceFunc func(handlerInput) handlerOutput
type DataServiceFunc func(*sql.DB, handlerInput) handlerOutput
type AuthServiceFunc func(username, *sql.DB, handlerInput) handlerOutput

type ServerHandle interface { // interface for handlers
	Handle() HttpHandler
	AcceptableMethods() acceptableMethods
}

type ServiceHandle struct {
	handler ServiceFunc
	methods acceptableMethods
}

func (h ServiceHandle) Handle() HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer errorCatcher(w)
		checkMethod(r.Method, h.methods)
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		output := h.handler(parseRequest(r))
		outmsg := structToJSON(output.msg)
		if output.newTokenString != "" {
			cookie_lifetime := int(sessionDuration.Seconds())
			cookie := http.Cookie{Name: "Session", Value: string(output.newTokenString), MaxAge: cookie_lifetime}
			http.SetCookie(w, &cookie)
		}
		showUserTable()
		io.WriteString(w, outmsg)
	}
}

func (h ServiceHandle) AcceptableMethods() acceptableMethods {
	return h.methods
}

type DataServiceHandle struct {
	handler DataServiceFunc
	methods acceptableMethods
}

func (h DataServiceHandle) Handle() HttpHandler {
	f := func(input handlerInput) handlerOutput {
		db := openDatabase()
		defer cleanCloser(db)
		return h.handler(db, input)

	}
	return ServiceHandle{f, h.AcceptableMethods()}.Handle()
}

func (h DataServiceHandle) AcceptableMethods() acceptableMethods {
	return h.methods
}

type AuthServiceHandle struct {
	handler AuthServiceFunc
	methods acceptableMethods
}

func (h AuthServiceHandle) Handle() HttpHandler {
	f := func(db *sql.DB, input handlerInput) handlerOutput {
		username := verifySession(input.session)
		output := h.handler(username, db, input)
		if !isLogged(db, username) {
			output.newTokenString =
				extendSession(db, string(username))
		}

		return output
	}
	return DataServiceHandle{f, h.AcceptableMethods()}.Handle()
}

func (h AuthServiceHandle) AcceptableMethods() acceptableMethods {
	return h.methods
}

func LoginWithRemember(db *sql.DB, input handlerInput) handlerOutput {
	var auth AuthJSON
	stringToStruct(input, &auth)
	token := loginAccount(db, auth)
	msg := OutputJSON{Success: true, Message: "Logged in", Result: token}
	return handlerOutput{msg, token}
}

func CreateAccount(db *sql.DB, input handlerInput) handlerOutput {
	var auth AuthJSON
	stringToStruct(input, &auth)
	addToDatabase(db, auth)
	result := OutputJSON{Success: true, Message: "Account Created"}
	return handlerOutput{msg: result}
}

func GetPrivateProfile(name username, db *sql.DB, _ handlerInput) handlerOutput {
	result := OutputJSON{Success: true, Message: "Profile", Result: getUserData(db, name)}
	return handlerOutput{msg: result}
}

func Logout(name username, db *sql.DB, input handlerInput) handlerOutput {
	logoutAccount(db, name)
	return handlerOutput{msg: OutputJSON{Success: true, Message: "Logged out"}}
}

func DeleteAccount(name username, db *sql.DB, input handlerInput) handlerOutput {
	logoutAccount(db, name)
	deleteFromDatabase(db, name)
	return handlerOutput{msg: OutputJSON{Success: true, Message: "Deleted Account"}}
}

func RandomShitPost(handlerInput) handlerOutput {
	response, err := http.Get("https://api.thedailyshitpost.net/random")
	if err != nil {
		ServerRuntimeError("Error while getting shitpost", err)
	}
	var shitpost RandomShitPostJSON
	parseResponseToStruct(response, &shitpost)
	if shitpost.Error {
		OnlyServerError("Remote Error while getting RandomShitpost")
	}
	return handlerOutput{msg: OutputJSON{Success: true, Message: "Random Shitpost", Result: shitpost.Url}}
}

func HandlersMap() map[string]ServerHandle {
	handlers := make(map[string]ServerHandle)
	handlers["/login"] = DataServiceHandle{LoginWithRemember, acceptableMethods{Put: true}}
	handlers["/create_account"] = DataServiceHandle{CreateAccount, acceptableMethods{Post: true}}
	handlers["/get_private_profile"] = AuthServiceHandle{GetPrivateProfile, acceptableMethods{Get: true}}
	handlers["/logout"] = AuthServiceHandle{Logout, acceptableMethods{Put: true}}
	handlers["/delete_account"] = AuthServiceHandle{DeleteAccount, acceptableMethods{Delete: true}}
	handlers["/random_shitpost"] = ServiceHandle{RandomShitPost, acceptableMethods{Get: true}}
	return handlers
}
