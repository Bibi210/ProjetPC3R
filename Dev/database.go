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
/* Table of Msg -> MsgID : PK(int) | Sender :  FK(UserID) | Content : string | Date : string | Upvotes : int */
/* Table of DM -> Receiver : FK(UserID) | Message : FK(MsgID)
/* Table of Comments -> Post : FK(ShitPost) | Message : FK(MsgID)  */

/* Get User Profile -> Select *(!Password) from Users */
/* Get User ID ->  Select UserID from Users where UserID = $1
/* Get User ShitPosts ->  Select * from ShitPost where Poster = Get User ID */
/* Get User Comments -> Select * from Comments where Poster = Get User ID*/

func addUserToDatabase(db *sql.DB, auth AuthJSON) {
	user := user_row{username: auth.Login, password: auth.Mdp}
	user.Create(db)
}

func getUserProfile(db *sql.DB, username username) UserProfileJSON {
	user := getUser(db, username)
	return UserProfileJSON{Username: user.username, UserID: user.userID}
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
	session  time.Time
	lastSeen time.Time
}

func (u *user_row) String() string {
	return fmt.Sprintf("UserID : %d | Username : %s | Password : %s | Session : %s | LastSeen : %s", u.userID, u.username, u.password, formatTime(u.session), formatTime(u.lastSeen))
}

func ReadFromRow(row *sql.Rows) user_row {
	u := user_row{}
	var lastSeen string
	var session string
	ServerRuntimeError("Error While Reading Row", row.Scan(&u.userID, &u.username, &u.password, &session, &lastSeen))
	u.lastSeen = parseTime(lastSeen)
	u.session = parseTime(session)
	return u
}

func (u *user_row) Create(c *sql.DB) {
	executeRequest(c, "INSERT INTO Users (Username, Password, Session, LastSeen) VALUES (?, ?, ?, ?)", u.username, u.password, formatTime(u.session), formatTime(u.lastSeen))
}

func (u *user_row) Update(c *sql.DB) {
	executeRequest(c, "UPDATE Users SET Username = ?, Password = ?, Session = ?, LastSeen = ? WHERE UserID = ?", u.username, u.password, formatTime(u.session), formatTime(u.lastSeen), u.userID)
}

func (u *user_row) UpdateSession(c *sql.DB) {
	executeRequest(c, "UPDATE Users SET Session = ?, LastSeen = ? WHERE UserID = ?", formatTime(u.session), formatTime(u.lastSeen), u.userID)
}

func (u user_row) Delete(c *sql.DB) {
	executeRequest(c, "DELETE FROM Users WHERE UserID = ?", u.userID)
}

func getUser(c *sql.DB, username username) user_row {
	rows := query(c, "SELECT * FROM Users WHERE Username = ?", username)
	defer rows.Close()
	if !rows.Next() {
		OnlyServerError("User don't exist")
	}
	return ReadFromRow(rows)
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

const createMsg = `CREATE TABLE IF NOT EXISTS Msg (
	MsgID INTEGER PRIMARY KEY AUTOINCREMENT,
	Sender INTEGER NOT NULL,
	Content TEXT NOT NULL,
	Date TEXT NOT NULL,
	Upvotes INTEGER NOT NULL,
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
	FOREIGN KEY (Post)
		REFERENCES ShitPost(ShitPostID),
	FOREIGN KEY (Message)
		REFERENCES Msg(MsgID)
);`

func openDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./information.db")
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
