package Helpers

import (
	"encoding/json"
	"time"
)

func BytesToStruct(bytes []byte, buffer any) {
	ServerRuntimeError("JSON with wrong format", json.Unmarshal(bytes, buffer))
}

func StructToJSON(t any) string {
	obj, err := json.Marshal(t)
	ServerRuntimeError("Could not convert struct to JSON", err)
	return string(obj)
}

func ParseTime(timeStr string) time.Time {
	time, err := time.Parse(time.ANSIC, timeStr)
	ServerRuntimeError("Error while parsing time", err)
	return time
}

func FormatTime(t time.Time) string {
	return t.Format(time.ANSIC)
}

type PostIds []int
type CommentIds []int

/* Server Responses */

type OutputJSON struct {
	Success bool
	Message string
	Result  any
}

type APIRandomShitPostJSON struct {
	Error string
	Url   string
}

type ResponseUserPublicProfileJSON struct {
	Username      string
	Posts         PostIds
	Comments      CommentIds
	LastSeen      string
	VotedPosts    PostIds
	VotedComments CommentIds
}

type ResponseUserPrivateProfileJSON struct {
	PublicProfile ResponseUserPublicProfileJSON
}

type ResponseSavedShitPostJSON struct {
	Url        string
	Caption    string
	Creator    int
	Date       string
	Upvotes    int
	CommentIds CommentIds
}

type ResponseMsgJSON struct {
	Sender  int
	Content string
	Date    string
}

type ResponseCommentJSON struct {
	Msg     ResponseMsgJSON
	Upvotes int
}

type ResponseUpvoteJSON struct {
	Acceptedvalue int
	PostVotes     int
}

type ResponseSearchJSON struct {
	ShitPosts PostIds
	Users     []string
}

/* ClientMessages */
type RequestAuthJSON struct {
	Login string
	Mdp   string
}

type RequestPublicUserProfileJSON struct {
	Username string
}

type RequestSaveShitPostJSON struct {
	Url     string
	Caption string
}

type RequestSendCommentJSON struct {
	ShitPostId int
	Content    string
}

type RequestOnShitPostJSON struct {
	ShitPostId int
}

type RequestOnCommentJSON struct {
	CommentId int
}

type RequestOnCommentListJSON struct {
	CommentIds CommentIds
}

type RequestOnShitPostListJSON struct {
	ShitPostIds PostIds
}

type RequestSendDmJSON struct {
	To      int
	Content string
}

type RequestShitPostVoteJSON struct {
	ShitPostId int
	Value      int
}

type RequestCommentVoteJSON struct {
	CommentId int
	Value     int
}

type RequestSearchJSON struct {
	Query string
}
