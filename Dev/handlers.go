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

func getToken(r *http.Request) token_string {
	cookie, err := r.Cookie("Session")
	if err != nil {
		return ""
	}
	return token_string(cookie.Value)
}
func parseRequest(r *http.Request) service_input {
	return service_input{parseRequestToString(r), getToken(r)}
}
func parseResponseToStruct(r *http.Response, t any) {
	bytes := parseResponseToBytes(r)
	bytesToStruct(bytes, t)
}

type token_string string
type username string

type service_input struct {
	msg     string
	session token_string
}
type service_output struct {
	msg            OutputJSON
	newTokenString token_string
}

type httpValidHandler func(http.ResponseWriter, *http.Request)
type basicServiceFunc func(service_input) service_output
type dataServiceFunc func(*sql.DB, service_input) service_output
type authServiceFunc func(username, *sql.DB, service_input) service_output

type Service interface {
	ToHandler() httpValidHandler
	acceptableMethods() AcceptableMethods
}

type basic_service struct {
	service basicServiceFunc
	methods AcceptableMethods
}

func (h basic_service) ToHandler() httpValidHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ErrorCatcher(w)
		CheckMethod(r.Method, h.methods)
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		output := h.service(parseRequest(r))
		outmsg := structToJSON(output.msg)
		if output.newTokenString != "" {
			// Write new token to cookie
			cookie_lifetime := int(sessionDuration.Seconds())
			cookie := http.Cookie{Name: "Session", Value: string(output.newTokenString), MaxAge: cookie_lifetime}
			http.SetCookie(w, &cookie)
		}
		showUserTable()
		io.WriteString(w, outmsg)
	}
}

func (h basic_service) acceptableMethods() AcceptableMethods {
	return h.methods
}

type DataBasedService struct {
	service dataServiceFunc
	methods AcceptableMethods
}

func (h DataBasedService) ToHandler() httpValidHandler {
	f := func(input service_input) service_output {
		db := openDatabase()
		defer CleanCloser(db)
		return h.service(db, input)

	}
	return basic_service{f, h.acceptableMethods()}.ToHandler()
}

func (h DataBasedService) acceptableMethods() AcceptableMethods {
	return h.methods
}

type AuthServiceHandle struct {
	handler authServiceFunc
	methods AcceptableMethods
}

func (h AuthServiceHandle) ToHandler() httpValidHandler {
	f := func(db *sql.DB, input service_input) service_output {
		username := verifySession(input.session)
		output := h.handler(username, db, input)
		if isLogged(db, username) {
			output.newTokenString =
				extendSession(db, string(username))
		}

		return output
	}
	return DataBasedService{f, h.acceptableMethods()}.ToHandler()
}

func (h AuthServiceHandle) acceptableMethods() AcceptableMethods {
	return h.methods
}

func LoginWithRemember(db *sql.DB, input service_input) service_output {
	var auth RequestAuthJSON
	getClientMessage(input, &auth)
	token := loginAccount(db, auth)
	msg := OutputJSON{Success: true, Message: "Logged in", Result: token}
	return service_output{msg, token}
}

func CreateAccount(db *sql.DB, input service_input) service_output {
	var auth RequestAuthJSON
	getClientMessage(input, &auth)
	addUserToDatabase(db, auth)
	result := OutputJSON{Success: true, Message: "Account Created"}
	return service_output{msg: result}
}

func GetPrivateProfile(name username, db *sql.DB, _ service_input) service_output {
	result := OutputJSON{Success: true, Message: "Profile", Result: getUser(db, name).Private()}
	return service_output{msg: result}
}

func GetPublicProfile(db *sql.DB, input service_input) service_output {
	var profile RequestPublicUserProfileJSON
	getClientMessage(input, &profile)
	result := OutputJSON{Success: true, Message: "Profile", Result: getUser(db, username(profile.Username)).Public()}
	return service_output{msg: result}
}

func Logout(name username, db *sql.DB, input service_input) service_output {
	logoutAccount(db, name)
	return service_output{msg: OutputJSON{Success: true, Message: "Logged out"}}
}

func DeleteAccount(name username, db *sql.DB, input service_input) service_output {
	logoutAccount(db, name)
	getUser(db, name).Delete(db)
	return service_output{msg: OutputJSON{Success: true, Message: "Deleted Account"}}
}

func RandomShitPost(service_input) service_output {
	response, err := http.Get("https://api.thedailyshitpost.net/random")
	ServerRuntimeError("Error while getting shitpost", err)

	var shitpost ResponseRandomShitPostJSON
	parseResponseToStruct(response, &shitpost)
	if shitpost.Error {
		OnlyServerError("Remote Error while getting RandomShitpost")
	}
	return service_output{msg: OutputJSON{Success: true, Message: "Random Shitpost", Result: shitpost.Url}}
}

func HandlersMap() map[string]Service {
	handlers := make(map[string]Service)
	handlers["/login"] = DataBasedService{LoginWithRemember, AcceptableMethods{Put: true}}
	handlers["/create_account"] = DataBasedService{CreateAccount, AcceptableMethods{Post: true}}
	handlers["/get_private_profile"] = AuthServiceHandle{GetPrivateProfile, AcceptableMethods{Get: true}}
	handlers["/get_public_profile"] = DataBasedService{GetPublicProfile, AcceptableMethods{Get: true}}
	handlers["/logout"] = AuthServiceHandle{Logout, AcceptableMethods{Put: true}}
	handlers["/delete_account"] = AuthServiceHandle{DeleteAccount, AcceptableMethods{Delete: true}}
	handlers["/random_shitpost"] = basic_service{RandomShitPost, AcceptableMethods{Get: true}}
	return handlers
}
