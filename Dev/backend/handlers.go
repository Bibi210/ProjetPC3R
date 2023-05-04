package main

import (
	"Backend/Database"
	"Backend/Helpers"
	"database/sql"
	"io"
	"log"
	"net/http"
)

func getToken(r *http.Request) tokenString {
	cookie, err := r.Cookie("Session")
	if err != nil {
		return ""
	}
	return tokenString(cookie.Value)
}
func parseRequest(r *http.Request) serviceInput {
	defer r.Body.Close()
	read, err := io.ReadAll(r.Body)
	Helpers.ServerRuntimeError("Can't Parse Request Body", err)
	return serviceInput{read, getToken(r)}
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
func getClientRequest(s serviceInput, buffer any) {
	Helpers.BytesToStruct(s.msg, buffer)
}

type tokenString string
type username string

type serviceInput struct {
	msg     []byte
	session tokenString
}
type ServiceOutput struct {
	msg            Helpers.OutputJSON
	newTokenString tokenString
}

type HttpValidHandler func(http.ResponseWriter, *http.Request)
type basicServiceFunc func(serviceInput) ServiceOutput
type dataServiceFunc func(*sql.DB, serviceInput) ServiceOutput
type authServiceFunc func(username, *sql.DB, serviceInput) ServiceOutput

type Service interface {
	ToHandler() HttpValidHandler
	acceptableMethods() Helpers.AcceptableMethods
}

type BasicService struct {
	service basicServiceFunc
	methods Helpers.AcceptableMethods
}

func (h BasicService) ToHandler() HttpValidHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer Helpers.ErrorCatcher(w)
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		Helpers.CheckMethod(r.Method, h.methods)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Content-Type", "application/json")
		output := h.service(parseRequest(r))
		outputJson := Helpers.StructToJSON(output.msg)
		if output.newTokenString != "" {
			// Write new token to cookie
			cookieLifetime := int(sessionDuration.Seconds())
			cookie := http.Cookie{Name: "Session", Value: string(output.newTokenString), MaxAge: cookieLifetime}
			http.SetCookie(w, &cookie)
		}
		_, err := io.WriteString(w, outputJson)
		if err != nil {
			return
		}
	}
}

func (h BasicService) acceptableMethods() Helpers.AcceptableMethods {
	return h.methods
}

type DataBasedService struct {
	service dataServiceFunc
	methods Helpers.AcceptableMethods
}

func (h DataBasedService) ToHandler() HttpValidHandler {
	f := func(input serviceInput) ServiceOutput {
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

func (h AuthServiceHandle) ToHandler() HttpValidHandler {
	f := func(db *sql.DB, input serviceInput) ServiceOutput {
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

func LoginWithRemember(db *sql.DB, input serviceInput) ServiceOutput {
	var auth Helpers.RequestAuthJSON
	getClientRequest(input, &auth)
	token := loginAccount(db, auth)
	msg := Helpers.OutputJSON{Success: true, Message: "Logged in", Result: token}
	return ServiceOutput{msg, token}
}

func CreateAccount(db *sql.DB, input serviceInput) ServiceOutput {
	var auth Helpers.RequestAuthJSON
	getClientRequest(input, &auth)
	Database.AddUserToDatabase(db, auth.Login, auth.Mdp)
	result := Helpers.OutputJSON{Success: true, Message: "Account Created"}
	return ServiceOutput{msg: result}
}

func GetPrivateProfile(name username, db *sql.DB, _ serviceInput) ServiceOutput {
	result := Helpers.OutputJSON{Success: true, Message: "Profile", Result: Database.GetUser(db, string(name)).Private(db)}
	return ServiceOutput{msg: result}
}

func GetPublicProfile(db *sql.DB, input serviceInput) ServiceOutput {
	var profile Helpers.RequestPublicUserProfileJSON
	getClientRequest(input, &profile)
	result := Helpers.OutputJSON{Success: true, Message: "Profile", Result: Database.GetUser(db, profile.Username).Public(db)}
	return ServiceOutput{msg: result}
}

func Logout(name username, db *sql.DB, _ serviceInput) ServiceOutput {
	logoutAccount(db, name)
	return ServiceOutput{msg: Helpers.OutputJSON{Success: true, Message: "Logged out"}}
}

func DeleteAccount(name username, db *sql.DB, _ serviceInput) ServiceOutput {
	logoutAccount(db, name)
	Database.GetUser(db, string(name)).DeleteUser(db)
	return ServiceOutput{msg: Helpers.OutputJSON{Success: true, Message: "Deleted Account"}}
}

func RandomShitPost(serviceInput) ServiceOutput {
	response, err := http.Get("https://api.thedailyshitpost.net/random")
	Helpers.ServerRuntimeError("Error while getting shitPost", err)

	var shitPostJSON Helpers.APIRandomShitPostJSON
	parseResponseToStruct(response, &shitPostJSON)
	if shitPostJSON.Error != "False" {
		Helpers.OnlyServerError("Remote Error while getting RandomShitPost")
	}
	return ServiceOutput{msg: Helpers.OutputJSON{Success: true, Message: "Random ShitPost", Result: shitPostJSON.Url}}
}

func SavePost(name username, db *sql.DB, input serviceInput) ServiceOutput {
	var post Helpers.RequestSaveShitPostJSON
	getClientRequest(input, &post)
	return ServiceOutput{msg: Helpers.OutputJSON{Success: true, Message: "Saved ShitPost", Result: Database.SaveShitPost(db, string(name), post.Url, post.Caption)}}
}

func GetSavedPost(db *sql.DB, input serviceInput) ServiceOutput {
	var post Helpers.RequestOnShitPostJSON
	getClientRequest(input, &post)
	result := Helpers.OutputJSON{Success: true, Message: "ShitPost Retrieved", Result: Database.GetShitPostAsJSON(db, post.ShitPostId)}
	return ServiceOutput{msg: result}
}

func GetSavedPosts(db *sql.DB, input serviceInput) ServiceOutput {
	var shitPostListJSON Helpers.RequestOnShitPostListJSON
	getClientRequest(input, &shitPostListJSON)
	result := Helpers.OutputJSON{Success: true, Message: "ShitPost List", Result: Database.GetShitPostListAsJSON(db, shitPostListJSON.ShitPostIds)}
	return ServiceOutput{msg: result}
}

func PostComment(name username, db *sql.DB, input serviceInput) ServiceOutput {
	var comment Helpers.RequestSendCommentJSON
	getClientRequest(input, &comment)
	id := Database.SendComment(db, string(name), comment.ShitPostId, comment.Content)
	return ServiceOutput{msg: Helpers.OutputJSON{Success: true, Message: "Posted Comment", Result: Helpers.ResponseSaveJSON{Id: id}}}
}

func GetSingleComment(db *sql.DB, input serviceInput) ServiceOutput {
	var comment Helpers.RequestOnCommentJSON
	getClientRequest(input, &comment)
	result := Helpers.OutputJSON{Success: true, Message: "Comment", Result: Database.GetCommentAsJSON(db, comment.CommentId)}
	return ServiceOutput{msg: result}
}

func GetComments(db *sql.DB, input serviceInput) ServiceOutput {
	var CommentLs Helpers.RequestOnCommentListJSON
	getClientRequest(input, &CommentLs)
	result := Helpers.OutputJSON{Success: true, Message: "Retrieved Comments", Result: Database.GetCommentListAsJSON(db, CommentLs.CommentIds)}
	return ServiceOutput{msg: result}
}

func PostCommentVote(name username, db *sql.DB, input serviceInput) ServiceOutput {
	var vote Helpers.RequestCommentVoteJSON
	getClientRequest(input, &vote)
	Database.SaveCommentUpvotes(db, string(name), vote.CommentId, vote.Value)
	return ServiceOutput{msg: Helpers.OutputJSON{Success: true, Message: "Voted"}}
}

func PostShitPostVote(name username, db *sql.DB, input serviceInput) ServiceOutput {
	var vote Helpers.RequestShitPostVoteJSON
	getClientRequest(input, &vote)
	Database.SavePostUpvotes(db, string(name), vote.ShitPostId, vote.Value)
	return ServiceOutput{msg: Helpers.OutputJSON{Success: true, Message: "Voted"}}
}

func Search(db *sql.DB, input serviceInput) ServiceOutput {
	var search Helpers.RequestSearchJSON
	getClientRequest(input, &search)
	output := Helpers.ResponseSearchJSON{ShitPosts: Database.SearchShitPost(db, search.Query), Users: Database.SearchUser(db, search.Query)}
	result := Helpers.OutputJSON{Success: true, Message: "Search", Result: output}
	return ServiceOutput{msg: result}
}

func GetTopUsers(db *sql.DB, input serviceInput) ServiceOutput {
	var top Helpers.RequestTopJSON
	getClientRequest(input, &top)
	result := Helpers.OutputJSON{Success: true, Message: "Top Users", Result: Database.GetTopUsersIDS(db, top.Count)}
	return ServiceOutput{msg: result}
}

func GetTopShitPosts(db *sql.DB, input serviceInput) ServiceOutput {
	var top Helpers.RequestTopJSON
	getClientRequest(input, &top)
	result := Helpers.OutputJSON{Success: true, Message: "Top ShitPosts", Result: Database.GetTopPostsIDs(db, top.Count)}
	return ServiceOutput{msg: result}
}

type FrontHandler struct {
	methods Helpers.AcceptableMethods
}

func (h FrontHandler) acceptableMethods() Helpers.AcceptableMethods {
	return h.methods
}

func (h FrontHandler) ToHandler() HttpValidHandler {
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
	handlers["/api/get_top_shitposts"] = DataBasedService{GetTopShitPosts, Helpers.AcceptableMethods{Put: true}}

	frontend := FrontHandler{Helpers.AcceptableMethods{Get: true}}

	handlers["/"] = frontend
	return handlers
}
