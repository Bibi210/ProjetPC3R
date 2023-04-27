package Database

import (
	"Backend/Helpers"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./Database/database.db")
	Helpers.ServerRuntimeError("Error While Opening Database", err)
	return db
}

func closeDatabase(db *sql.DB) {
	Helpers.ServerRuntimeError("Error While Closing Database", db.Close())
}

func executeRequest(c *sql.DB, query string, args ...interface{}) sql.Result {
	res, err := c.Exec(query, args...)
	Helpers.ServerRuntimeError("Error While Executing Query", err)
	return res
}

func query(c *sql.DB, query string, args ...interface{}) *sql.Rows {
	rows, err := c.Query(query, args...)
	Helpers.ServerRuntimeError("Error While Querying Database", err)
	return rows
}


func CreateDatabase() {
	c := OpenDatabase()
	defer closeDatabase(c)
	executeRequest(c, createUsers)
	executeRequest(c, createShitPost)
	executeRequest(c, createMsg)
	executeRequest(c, CreateComments)
	executeRequest(c, createCommentUpvotes)
	executeRequest(c, createPostUpvotes)
}

func DeleteDatabase() {
	c := OpenDatabase()
	defer closeDatabase(c)
	executeRequest(c, "DROP TABLE Users")
	executeRequest(c, "DROP TABLE ShitPost")
	executeRequest(c, "DROP TABLE Msg")
	executeRequest(c, "DROP TABLE Comments")
	executeRequest(c, "DROP TABLE CommentUpvotes")
	executeRequest(c, "DROP TABLE PostUpvotes")
}

func ShowDatabase(db *sql.DB) {
	showUserTable(db)
	showShitPostTable(db)
	showCommentTable(db)
	showMsgTable(db)
	showCommentUpvotesTable(db)
	showPostUpvotesTable(db)
}

func ShutdownDatabase(cleanDatabase bool) {
	if cleanDatabase {
		DeleteDatabase()
	}
	log.Println("Database Shutdown")
}

func CleanCloser(db *sql.DB) {
	if r := recover(); r != nil {
		err := r.(error)
		closeDatabase(db)
		panic(err)
	}
	closeDatabase(db)
}
