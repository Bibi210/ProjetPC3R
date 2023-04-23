package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/* KeyUserName */

/* Table of Users  -> UserID : PK(int) | Username : string | Password : string | Session : Blob | LastSeen : string */
/* Table of ShitPost -> ShitPostID : PK(int) | Poster : FK(UserID) | Caption : string | URL : string | Date : string | Upvotes : int*/
/* Table of Msg -> MsgID : PK(int) | Sender :  FK(UserID) | Content : string | Date : string  */
/* Table of DM -> Receiver : FK(UserID) | Message : FK(MsgID)
/* Table of Comments -> Post : FK(ShitPost) | Message : FK(MsgID) | Upvotes : int  */

/* Get User Profile -> Select *(!Password) from Users */
/* Get User ID ->  Select UserID from Users where UserID = $1
/* Get User ShitPosts ->  Select * from ShitPost where Poster = Get User ID */
/* Get User Comments -> Select * from Comments where Poster = Get User ID*/

func addUserToDatabase(db *sql.DB, auth RequestAuthJSON) {
	user := user_row{username: auth.Login, password: auth.Mdp, Session: time.Now(), LastSeen: time.Now()}
	user.CreateUser(db)
}

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
	password string
	Session  time.Time
	LastSeen time.Time
}

func (u *user_row) String() string {
	return fmt.Sprintf("UserID : %d | Username : %s | Password : %s | Session : %s | LastSeen : %s", u.userID, u.username, u.password, formatTime(u.Session), formatTime(u.LastSeen))
}

func ReadFromRow(row *sql.Rows) user_row {
	u := user_row{}
	var lastSeen string
	var session string
	ServerRuntimeError("Error While Reading Row", row.Scan(&u.userID, &u.username, &u.password, &session, &lastSeen))

	u.LastSeen = parseTime(lastSeen)
	u.Session = parseTime(session)

	return u
}

func (u *user_row) CreateUser(c *sql.DB) {
	executeRequest(c, "INSERT INTO Users (Username, Password, Session, LastSeen) VALUES (?, ?, ?, ?)", u.username, u.password, formatTime(u.Session), formatTime(u.LastSeen))
}

func (u *user_row) UpdateUser(c *sql.DB) {
	executeRequest(c, "UPDATE Users SET Username = ?, Password = ?, Session = ?, LastSeen = ? WHERE UserID = ?", u.username, u.password, formatTime(u.Session), formatTime(u.LastSeen), u.userID)
}

func (u *user_row) UpdateUserSession(c *sql.DB) {
	executeRequest(c, "UPDATE Users SET Session = ?, LastSeen = ? WHERE UserID = ?", formatTime(u.Session), formatTime(u.LastSeen), u.userID)
}

func (u user_row) DeleteUser(c *sql.DB) {
	executeRequest(c, "DELETE FROM Users WHERE UserID = ?", u.userID)
}

func (u user_row) Public(db *sql.DB) ResponseUserProfileJSON {
	return ResponseUserProfileJSON{Username: u.username, UserID: u.userID, LastSeen: formatTime(u.LastSeen), Posts: GetUserShitPosts(db, username(u.username))}
}

func (u user_row) Private(db *sql.DB) ResponseUserProfileJSON {
	return ResponseUserProfileJSON{Username: u.username, UserID: u.userID, LastSeen: formatTime(u.LastSeen), Posts: GetUserShitPosts(db, username(u.username))}
}

func getUser(c *sql.DB, username username) user_row {
	rows := query(c, "SELECT * FROM Users WHERE Username = ?", username)
	defer rows.Close()
	if !rows.Next() {
		OnlyServerError("User don't exist")
	}
	return ReadFromRow(rows)
}

func getUserByID(c *sql.DB, userID int) user_row {
	rows := query(c, "SELECT * FROM Users WHERE UserID = ?", userID)
	defer rows.Close()
	if !rows.Next() {
		OnlyServerError("User don't exist")
	}
	return ReadFromRow(rows)
}

func isUserExist(c *sql.DB, username username) bool {
	rows := query(c, "SELECT * FROM Users WHERE Username = ?", username)
	defer rows.Close()
	return rows.Next()
}

func showUserTable() {
	db := openDatabase()
	defer closeDatabase(db)
	rows := query(db, "SELECT * FROM Users")
	defer rows.Close()
	for rows.Next() {
		log.Println(ReadFromRow(rows))
	}
}

const createShitPost = `CREATE TABLE IF NOT EXISTS ShitPost (
	ShitPostID INTEGER PRIMARY KEY AUTOINCREMENT,
	Poster INTEGER NOT NULL,
	Caption TEXT NOT NULL,
	URL TEXT NOT NULL,
	Date TEXT NOT NULL,
	Upvotes INTEGER NOT NULL,
	FOREIGN KEY (Poster) 
		REFERENCES Users(UserID)
);`

type saved_shitpost_row struct {
	shitpostID int
	poster     int
	caption    string
	url        string
	Date       time.Time
	upvotes    int
}

func (s *saved_shitpost_row) String() string {
	return fmt.Sprintf("ShitPostID : %d | Poster : %d | Caption : %s | URL : %s | Date : %s | Upvotes : %d", s.shitpostID, s.poster, s.caption, s.url, formatTime(s.Date), s.upvotes)
}

func ReadFromRowShitPost(row *sql.Rows) saved_shitpost_row {
	r := saved_shitpost_row{}
	var date string
	ServerRuntimeError("Error While Reading Row", row.Scan(&r.shitpostID, &r.poster, &r.caption, &r.url, &date, &r.upvotes))
	r.Date = parseTime(date)
	return r
}

func SaveShitPost(c *sql.DB, poster username, info SaveShitPostJSON) {
	userID := getUser(c, poster).userID
	executeRequest(c, "INSERT INTO ShitPost (Poster, Caption, URL, Date, Upvotes) VALUES (?, ?, ?, ?, ?)", userID, info.Caption, info.Url, formatTime(time.Now()), 0)
}

func DeleteShitPost(c *sql.DB, shitpostID int) {
	executeRequest(c, "DELETE FROM ShitPost WHERE ShitPostID = ?", shitpostID)
}

func GetShitPost(c *sql.DB, shitpostID int) ResponseSavedShitPostJSON {
	rows := query(c, "SELECT * FROM ShitPost WHERE ShitPostID = ?", shitpostID)
	defer rows.Close()
	if !rows.Next() {
		OnlyServerError("ShitPost don't exist")
	}
	v := ReadFromRowShitPost(rows)
	return ResponseSavedShitPostJSON{Caption: v.caption, Date: formatTime(v.Date), Upvotes: v.upvotes}
}

func GetUserShitPosts(c *sql.DB, username username) []ResponseSavedShitPostJSON {
	user := getUser(c, username)
	rows := query(c, "SELECT * FROM ShitPost WHERE Poster = ?", user.userID)
	defer rows.Close()
	var result []ResponseSavedShitPostJSON
	for rows.Next() {
		v := ReadFromRowShitPost(rows)
		result = append(result, ResponseSavedShitPostJSON{Caption: v.caption, Date: formatTime(v.Date), Upvotes: v.upvotes, Creator: user.username, Url: v.url})
	}
	return result
}

func GetAllShitPosts(c *sql.DB) []ResponseSavedShitPostJSON {
	rows := query(c, "SELECT * FROM ShitPost")
	defer rows.Close()
	var result []ResponseSavedShitPostJSON
	for rows.Next() {
		v := ReadFromRowShitPost(rows)
		user := getUserByID(c, v.poster)
		result = append(result, ResponseSavedShitPostJSON{Caption: v.caption, Date: formatTime(v.Date), Upvotes: v.upvotes, Creator: user.username, Url: v.url})
	}
	return result
}

func showShitPostTable() {
	db := openDatabase()
	defer closeDatabase(db)
	rows := query(db, "SELECT * FROM ShitPost")
	defer rows.Close()
	for rows.Next() {
		log.Println(ReadFromRowShitPost(rows))
	}
}

const createMsg = `CREATE TABLE IF NOT EXISTS Msg (
	MsgID INTEGER PRIMARY KEY AUTOINCREMENT,
	Sender INTEGER NOT NULL,
	Content TEXT NOT NULL,
	Date TEXT NOT NULL,
	FOREIGN KEY (Sender)
		REFERENCES Users(UserID)
);`

const createDM = `CREATE TABLE IF NOT EXISTS DM (
	Receiver INTEGER NOT NULL,
	Message INTEGER NOT NULL,
	FOREIGN KEY (Receiver)
		REFERENCES Users(UserID),
	FOREIGN KEY (Message)
		REFERENCES Msg(MsgID)
);`

const createComments = `CREATE TABLE IF NOT EXISTS Comments (
	Post INTEGER NOT NULL,
	Message INTEGER NOT NULL,
	Upvotes INTEGER NOT NULL,
	FOREIGN KEY (Post)
		REFERENCES ShitPost(ShitPostID),
	FOREIGN KEY (Message)
		REFERENCES Msg(MsgID)
);`

func openDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	ServerRuntimeError("Error While Opening Database", err)
	return db
}

func closeDatabase(db *sql.DB) {
	ServerRuntimeError("Error While Closing Database", db.Close())
}

func executeRequest(c *sql.DB, query string, args ...interface{}) sql.Result {
	res, err := c.Exec(query, args...)
	ServerRuntimeError("Error While Executing Query", err)
	return res
}

func query(c *sql.DB, query string, args ...interface{}) *sql.Rows {
	rows, err := c.Query(query, args...)
	ServerRuntimeError("Error While Querying Database", err)
	return rows
}

func createDatabase() {
	c := openDatabase()
	defer closeDatabase(c)
	executeRequest(c, createUsers)
	executeRequest(c, createShitPost)
	executeRequest(c, createMsg)
	executeRequest(c, createDM)
	executeRequest(c, createComments)

}

func deleteDatabase() {
	c := openDatabase()
	defer closeDatabase(c)
	executeRequest(c, "DROP TABLE Users")
	executeRequest(c, "DROP TABLE ShitPost")
	executeRequest(c, "DROP TABLE Msg")
	executeRequest(c, "DROP TABLE DM")
	executeRequest(c, "DROP TABLE Comments")
}

func shutdownDatabase(cleanDatabase bool) {
	if cleanDatabase {
		deleteDatabase()
	}
	log.Println("Database Shutdown")
}
