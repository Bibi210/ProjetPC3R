package Database

import (
	"Backend/Helpers"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const createUsers = `CREATE TABLE IF NOT EXISTS Users (
	UserID INTEGER PRIMARY KEY AUTOINCREMENT,
	Username TEXT NOT NULL UNIQUE,
	Password TEXT NOT NULL,
	Session BLOB,
	LastSeen TEXT
);`

type user_row struct {
	userID   int
	username string
	Password string
	Session  time.Time
	LastSeen time.Time
}

func (u *user_row) String() string {
	return fmt.Sprintf("UserID : %d | Username : %s | Password : %s | Session : %s | LastSeen : %s", u.userID, u.username, u.Password, Helpers.FormatTime(u.Session), Helpers.FormatTime(u.LastSeen))
}

func ReadFromRowUser(row *sql.Rows) user_row {
	u := user_row{}
	var lastSeen string
	var session string
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&u.userID, &u.username, &u.Password, &session, &lastSeen))

	u.LastSeen = Helpers.ParseTime(lastSeen)
	u.Session = Helpers.ParseTime(session)

	return u
}

func AddUserToDatabase(db *sql.DB, Login string, Mdp string) {
	user := user_row{username: Login, Password: Mdp, Session: time.Now(), LastSeen: time.Now()}
	executeRequest(db, "INSERT INTO Users (Username, Password, Session, LastSeen) VALUES (?, ?, ?, ?)", user.username, user.Password, Helpers.FormatTime(user.Session), Helpers.FormatTime(user.LastSeen))
}

func (u *user_row) UpdateUserSession(c *sql.DB) {
	executeRequest(c, "UPDATE Users SET Session = ?, LastSeen = ? WHERE UserID = ?", Helpers.FormatTime(u.Session), Helpers.FormatTime(u.LastSeen), u.userID)
}

func (u user_row) DeleteUser(c *sql.DB) {
	executeRequest(c, "DELETE FROM Users WHERE UserID = ?", u.userID)
}

func (u user_row) Public(db *sql.DB) Helpers.ResponseUserPublicProfileJSON {
	return Helpers.ResponseUserPublicProfileJSON{Username: u.username, LastSeen: Helpers.FormatTime(u.LastSeen), Posts: GetUserShitPostsIDS(db, u.username), Comments: GetUserCommentsIDS(db, u.userID), VotedPosts: GetUserVotedShitPostsIDS(db, u.userID), VotedComments: GetUserVotedCommentsIDS(db, u.userID)}
}

func (u user_row) Private(db *sql.DB) Helpers.ResponseUserPublicProfileJSON {
	return Helpers.ResponseUserPublicProfileJSON{Username: u.username, LastSeen: Helpers.FormatTime(u.LastSeen), Posts: GetUserShitPostsIDS(db, u.username), Comments: GetUserCommentsIDS(db, u.userID)}
}

func GetUser(c *sql.DB, username string) user_row {
	rows := query(c, "SELECT * FROM Users WHERE Username = ?", username)
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("User don't exist")
	}
	return ReadFromRowUser(rows)
}

func IsUserExist(c *sql.DB, username string) bool {
	rows := query(c, "SELECT * FROM Users WHERE Username = ?", username)
	defer rows.Close()
	return rows.Next()
}

func GetUserShitPostsIDS(c *sql.DB, username string) []int {
	rows := query(c, "SELECT ShitPostID FROM ShitPost WHERE Poster = ?", GetUser(c, username).userID)
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&id))
		ids = append(ids, id)
	}
	return ids
}

func GetUserMessages(c *sql.DB, id int) []msg_row {
	rows := query(c, "SELECT * FROM Msg WHERE Sender = ?", id)
	defer rows.Close()
	var msgs []msg_row
	for rows.Next() {
		msgs = append(msgs, ReadFromRowMsg(rows))
	}
	return msgs
}

func GetUserComments(c *sql.DB, id int) []comment_row {
	msgs := GetUserMessages(c, id)
	var comments []comment_row
	for _, msg := range msgs {
		rows := query(c, "SELECT * FROM Comments WHERE Message = ?", msg.msgID)
		defer rows.Close()
		for rows.Next() {
			comments = append(comments, ReadFromRowComment(c, rows))
		}
	}
	return comments
}

func GetUserCommentsIDS(c *sql.DB, id int) []int {
	comments := GetUserComments(c, id)
	var ids []int
	for _, comment := range comments {
		ids = append(ids, comment.Comid)
	}
	return ids
}

func GetAllUsersIDS(c *sql.DB) []int {
	rows := query(c, "SELECT UserID FROM Users")
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&id))
		ids = append(ids, id)
	}
	return ids
}

func GetAllUsersAsJSON(c *sql.DB) []Helpers.ResponseUserPublicProfileJSON {
	rows := query(c, "SELECT * FROM Users")
	defer rows.Close()
	var users []Helpers.ResponseUserPublicProfileJSON
	for rows.Next() {
		users = append(users, ReadFromRowUser(rows).Public(c))
	}
	return users
}

func SearchUser(c *sql.DB, username string) []string {
	rows := query(c, "SELECT * FROM Users WHERE Username LIKE ?", "%"+username+"%")
	defer rows.Close()
	var users []string
	for rows.Next() {
		users = append(users, ReadFromRowUser(rows).username)
	}
	return users
}

func showUserTable(db *sql.DB) {
	fmt.Println("Users Table :")
	rows := query(db, "SELECT * FROM Users")
	defer rows.Close()
	for rows.Next() {
		log.Println(ReadFromRowUser(rows))
	}
}