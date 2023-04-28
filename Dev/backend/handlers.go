package main

import (
	"Backend/Database"
	"Backend/Helpers"
	"database/sql"
	"io"
	"log"
	"net/http"
)

func getToken(r *http.Request) token_string {
	cookie, err := r.Cookie("Session")
	if err != nil {
		return ""
	}
	return token_string(cookie.Value)
}
func parseRequest(r *http.Request) service_input {
	defer r.Body.Close()
	read, err := io.ReadAll(r.Body)
	Helpers.ServerRuntimeError("Can't Parse Request Body", err)
	return service_input{read, getToken(r)}
}
func parseResponseToStruct(r *http.Response, t any) {
	if r == nil || r.StatusCode != 200 {
		Helpers.OnlyServerError("Request Failed with status code: " + string(r.Status))
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	Helpers.ServerRuntimeError("Can't Parse Response Body", err)
	Helpers.BytesToStruct(body, t)
}
func getClientRequest(s service_input, buffer any) {
	Helpers.BytesToStruct(s.msg, buffer)
}

type token_string string
type username string

type service_input struct {
	msg     []byte
	session token_string
}
type service_output struct {
	msg            Helpers.OutputJSON
	newTokenString token_string
}

type httpValidHandler func(http.ResponseWriter, *http.Request)
type basicServiceFunc func(service_input) service_output
type dataServiceFunc func(*sql.DB, service_input) service_output
type authServiceFunc func(username, *sql.DB, service_input) service_output

type Service interface {
	ToHandler() httpValidHandler
	acceptableMethods() Helpers.AcceptableMethods
}

type BasicService struct {
	service basicServiceFunc
	methods Helpers.AcceptableMethods
}

func (h BasicService) ToHandler() httpValidHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer Helpers.ErrorCatcher(w)
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		Helpers.CheckMethod(r.Method, h.methods)
		w.Header().Set("Content-Type", "application/json")
		output := h.service(parseRequest(r))
		outmsg := Helpers.StructToJSON(output.msg)
		if output.newTokenString != "" {
			// Write new token to cookie
			cookie_lifetime := int(sessionDuration.Seconds())
			cookie := http.Cookie{Name: "Session", Value: string(output.newTokenString), MaxAge: cookie_lifetime}
			http.SetCookie(w, &cookie)
		}
		io.WriteString(w, outmsg)
	}
}

func (h BasicService) acceptableMethods() Helpers.AcceptableMethods {
	return h.methods
}

type DataBasedService struct {
	service dataServiceFunc
	methods Helpers.AcceptableMethods
}

func (h DataBasedService) ToHandler() httpValidHandler {
	f := func(input service_input) service_output {
		db := Database.OpenDatabase()
		defer Database.CleanCloser(db)
		Database.ShowDatabase(db)
		return h.service(db, input)

	}
	return BasicService{f, h.acceptableMethods()}.ToHandler()
}

func (h DataBasedService) acceptableMethods() Helpers.AcceptableMethods {
	return h.methods
}

type AuthServiceHandle struct {
	handler authServiceFunc
	methods Helpers.AcceptableMethods
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

func (h AuthServiceHandle) acceptableMethods() Helpers.AcceptableMethods {
	return h.methods
}

func LoginWithRemember(db *sql.DB, input service_input) service_output {
	var auth Helpers.RequestAuthJSON
	getClientRequest(input, &auth)
	token := loginAccount(db, auth)
	msg := Helpers.OutputJSON{Success: true, Message: "Logged in", Result: token}
	return service_output{msg, token}
}

func CreateAccount(db *sql.DB, input service_input) service_output {
	var auth Helpers.RequestAuthJSON
	getClientRequest(input, &auth)
	Database.AddUserToDatabase(db, auth.Login, auth.Mdp)
	result := Helpers.OutputJSON{Success: true, Message: "Account Created"}
	return service_output{msg: result}
}

func GetPrivateProfile(name username, db *sql.DB, _ service_input) service_output {
	result := Helpers.OutputJSON{Success: true, Message: "Profile", Result: Database.GetUser(db, string(name)).Private(db)}
	return service_output{msg: result}
}

func GetPublicProfile(db *sql.DB, input service_input) service_output {
	var profile Helpers.RequestPublicUserProfileJSON
	getClientRequest(input, &profile)
	result := Helpers.OutputJSON{Success: true, Message: "Profile", Result: Database.GetUser(db, profile.Username).Public(db)}
	return service_output{msg: result}
}

func Logout(name username, db *sql.DB, input service_input) service_output {
	logoutAccount(db, name)
	return service_output{msg: Helpers.OutputJSON{Success: true, Message: "Logged out"}}
}

func DeleteAccount(name username, db *sql.DB, input service_input) service_output {
	logoutAccount(db, name)
	Database.GetUser(db, string(name)).DeleteUser(db)
	return service_output{msg: Helpers.OutputJSON{Success: true, Message: "Deleted Account"}}
}

func RandomShitPost(service_input) service_output {
	response, err := http.Get("https://api.thedailyshitpost.net/random")
	Helpers.ServerRuntimeError("Error while getting shitpost", err)

	var shitpost Helpers.APIRandomShitPostJSON
	parseResponseToStruct(response, &shitpost)
	if shitpost.Error != "False" {
		Helpers.OnlyServerError("Remote Error while getting RandomShitpost")
	}
	return service_output{msg: Helpers.OutputJSON{Success: true, Message: "Random Shitpost", Result: shitpost.Url}}
}

func SavePost(name username, db *sql.DB, input service_input) service_output {
	var post Helpers.RequestSaveShitPostJSON
	getClientRequest(input, &post)
	return service_output{msg: Helpers.OutputJSON{Success: true, Message: "Saved Shitpost", Result: Database.SaveShitPost(db, string(name), post.Url, post.Caption)}}
}

func GetSavedPost(db *sql.DB, input service_input) service_output {
	var post Helpers.RequestOnShitPostJSON
	getClientRequest(input, &post)
	result := Helpers.OutputJSON{Success: true, Message: "Shitpost Retrived", Result: Database.GetShitPostAsJSON(db, post.ShitPostId)}
	return service_output{msg: result}
}

func GetSavedPosts(db *sql.DB, input service_input) service_output {
	var shitpostLs Helpers.RequestOnShitPostListJSON
	getClientRequest(input, &shitpostLs)
	result := Helpers.OutputJSON{Success: true, Message: "Shitpost List", Result: Database.GetShitPostListAsJSON(db, shitpostLs.ShitPostIds)}
	return service_output{msg: result}
}

func PostComment(name username, db *sql.DB, input service_input) service_output {
	var comment Helpers.RequestSendCommentJSON
	getClientRequest(input, &comment)
	Database.SendComment(db, string(name), comment.ShitPostId, comment.Content)
	return service_output{msg: Helpers.OutputJSON{Success: true, Message: "Posted Comment"}}
}

func GetSingleComment(db *sql.DB, input service_input) service_output {
	var comment Helpers.RequestOnCommentJSON
	getClientRequest(input, &comment)
	result := Helpers.OutputJSON{Success: true, Message: "Comment", Result: Database.GetCommentAsJSON(db, comment.CommentId)}
	return service_output{msg: result}
}

func GetComments(db *sql.DB, input service_input) service_output {
	var CommentLs Helpers.RequestOnCommentListJSON
	getClientRequest(input, &CommentLs)
	result := Helpers.OutputJSON{Success: true, Message: "Retrived Comments", Result: Database.GetCommentListAsJSON(db, CommentLs.CommentIds)}
	return service_output{msg: result}
}

func PostCommentVote(name username, db *sql.DB, input service_input) service_output {
	var vote Helpers.RequestCommentVoteJSON
	getClientRequest(input, &vote)
	Database.SaveCommentUpvotes(db, string(name), vote.CommentId, vote.Value)
	return service_output{msg: Helpers.OutputJSON{Success: true, Message: "Voted"}}
}

func PostShitPostVote(name username, db *sql.DB, input service_input) service_output {
	var vote Helpers.RequestShitPostVoteJSON
	getClientRequest(input, &vote)
	Database.SavePostUpvotes(db, string(name), vote.ShitPostId, vote.Value)
	return service_output{msg: Helpers.OutputJSON{Success: true, Message: "Voted"}}
}

func Search(db *sql.DB, input service_input) service_output {
	var search Helpers.RequestSearchJSON
	getClientRequest(input, &search)
	output := Helpers.ResponseSearchJSON{ShitPosts: Database.SearchShitPost(db, search.Query), Users: Database.SearchUser(db, search.Query)}
	result := Helpers.OutputJSON{Success: true, Message: "Search", Result: output}
	return service_output{msg: result}
}

func GetTopUsers(db *sql.DB, input service_input) service_output {
	var top Helpers.RequestTopUsersJSON
	getClientRequest(input, &top)
	result := Helpers.OutputJSON{Success: true, Message: "Top Users", Result: Database.GetTopUsersIDS(db, top.Count)}
	return service_output{msg: result}
}

type FrontHandler struct {
	methods Helpers.AcceptableMethods
}

func (h FrontHandler) acceptableMethods() Helpers.AcceptableMethods {
	return h.methods
}

func (h FrontHandler) ToHandler() httpValidHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/dist/index.html")
	}
}

func HandlersMap() map[string]Service {
	handlers := make(map[string]Service)
	handlers["/api/login"] = DataBasedService{LoginWithRemember, Helpers.AcceptableMethods{Put: true}}
	handlers["/api/create_account"] = DataBasedService{CreateAccount, Helpers.AcceptableMethods{Post: true}}
	handlers["/api/get_private_profile"] = AuthServiceHandle{GetPrivateProfile, Helpers.AcceptableMethods{Get: true}}
	handlers["/api/get_public_profile"] = DataBasedService{GetPublicProfile, Helpers.AcceptableMethods{Put: true}}
	handlers["/api/logout"] = AuthServiceHandle{Logout, Helpers.AcceptableMethods{Get: true}}
	handlers["/api/delete_account"] = AuthServiceHandle{DeleteAccount, Helpers.AcceptableMethods{Delete: true}}
	handlers["/api/random_shitpost"] = BasicService{RandomShitPost, Helpers.AcceptableMethods{Get: true}}
	handlers["/api/save_shitpost"] = AuthServiceHandle{SavePost, Helpers.AcceptableMethods{Post: true}}
	handlers["/api/get_saved_shitpost"] = DataBasedService{GetSavedPost, Helpers.AcceptableMethods{Put: true}}
	handlers["/api/post_comment"] = AuthServiceHandle{PostComment, Helpers.AcceptableMethods{Post: true}}
	handlers["/api/get_comment"] = DataBasedService{GetSingleComment, Helpers.AcceptableMethods{Put: true}}
	handlers["/api/post_comment_vote"] = AuthServiceHandle{PostCommentVote, Helpers.AcceptableMethods{Post: true}}
	handlers["/api/post_shitpost_vote"] = AuthServiceHandle{PostShitPostVote, Helpers.AcceptableMethods{Post: true}}
	handlers["/api/search"] = DataBasedService{Search, Helpers.AcceptableMethods{Put: true}}
	handlers["/api/get_comment_list"] = DataBasedService{GetComments, Helpers.AcceptableMethods{Put: true}}
	handlers["/api/get_saved_shitpost_list"] = DataBasedService{GetSavedPosts, Helpers.AcceptableMethods{Put: true}}
	handlers["/api/get_top_users"] = DataBasedService{GetTopUsers, Helpers.AcceptableMethods{Put: true}}

	frontend := FrontHandler{Helpers.AcceptableMethods{Get: true}}
	handlers["/login"] = frontend
	handlers["/create_account"] = frontend
	handlers["/profile"] = frontend
	handlers["/logout"] = frontend
	handlers["/delete_account"] = frontend
	handlers["/"] = frontend
	return handlers
}
