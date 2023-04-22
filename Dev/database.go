package main

import (
	"fmt"
	"log"

	"time"

	"database/sql"

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

func addToDatabase(db *sql.DB, auth AuthJSON) {
	createUser(db, auth.Login, auth.Mdp)
}

func deleteFromDatabase(db *sql.DB, username username) {
	deleteUser(db, string(username))
}

func getUserData(db *sql.DB, username username) UserProfileJSON {
	user := getUser(db, string(username))
	return UserProfileJSON{Username: user.Username, UserID: user.UserID}
}

const createUsers = `CREATE TABLE IF NOT EXISTS Users (
	UserID INTEGER PRIMARY KEY AUTOINCREMENT,
	Username TEXT NOT NULL UNIQUE,
	Password TEXT NOT NULL,
	Session BLOB,
	LastSeen TEXT
);`

type User struct {
	UserID   int
	Username string
	Password string
	Session  time.Time
	LastSeen time.Time
}

func (u *User) String() string {
	return fmt.Sprintf("UserID : %d | Username : %s | Password : %s | Session : %s | LastSeen : %s", u.UserID, u.Username, u.Password, formatTime(u.Session), formatTime(u.LastSeen))
}

func ReadFromRow(row *sql.Rows) User {
	u := User{}
	var lastSeen string
	var session string
	ServerRuntimeError("Error While Reading Row", row.Scan(&u.UserID, &u.Username, &u.Password, &session, &lastSeen))
	u.LastSeen = parseTime(lastSeen)
	u.Session = parseTime(session)
	return u
}

func (u *User) InsertRow(c *sql.DB) {
	Execute(c, "INSERT INTO Users (Username, Password, Session, LastSeen) VALUES (?, ?, ?, ?)", u.Username, u.Password, formatTime(u.Session), formatTime(u.LastSeen))
}

func (u *User) UpdateRow(c *sql.DB) {
	Execute(c, "UPDATE Users SET Username = ?, Password = ?, Session = ?, LastSeen = ? WHERE UserID = ?", u.Username, u.Password, formatTime(u.Session), formatTime(u.LastSeen), u.UserID)
}

func (u *User) UpdateSession(c *sql.DB) {
	Execute(c, "UPDATE Users SET Session = ?, LastSeen = ? WHERE UserID = ?", formatTime(u.Session), formatTime(u.LastSeen), u.UserID)
}

func getUser(c *sql.DB, username string) User {
	rows := query(c, "SELECT * FROM Users WHERE Username = ?", username)
	defer rows.Close()
	if !rows.Next() {
		OnlyServerError("User don't exist")
	}
	return ReadFromRow(rows)
}

func deleteUser(c *sql.DB, username string) {
	Execute(c, "DELETE FROM Users WHERE Username = ?", username)
}

func createUser(co *sql.DB, username, password string) {
	user := User{Username: username, Password: password}
	user.InsertRow(co)
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
	err := db.Close()
	if err != nil {
		log.Println(err)
		ServerRuntimeError("Error While Closing Database", db.Close())
	}

}

func Execute(c *sql.DB, query string, args ...interface{}) sql.Result {
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
	Execute(c, createUsers)
	Execute(c, createShitPost)
	Execute(c, createMsg)
	Execute(c, createDM)
	Execute(c, createComments)

}

func deleteDatabase() {
	c := openDatabase()
	defer closeDatabase(c)
	Execute(c, "DROP TABLE Users")
	Execute(c, "DROP TABLE ShitPost")
	Execute(c, "DROP TABLE Msg")
	Execute(c, "DROP TABLE DM")
	Execute(c, "DROP TABLE Comments")
}

func shutdownDatabase(cleanDatabase bool) {
	if cleanDatabase {
		deleteDatabase()
	}
	log.Println("Database Shutdown")
}
