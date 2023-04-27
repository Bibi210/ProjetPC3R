package Database

import (
	"Backend/Helpers"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const createPostUpvotes = `CREATE TABLE IF NOT EXISTS PostUpvotes (
	PostUpvotesID INTEGER PRIMARY KEY AUTOINCREMENT,
	Post INTEGER NOT NULL,
	Upvoter INTEGER NOT NULL,
	Date TEXT NOT NULL,
	Vote INTEGER NOT NULL,
	FOREIGN KEY (Post)
		REFERENCES ShitPost(ShitPostID),
	FOREIGN KEY (Upvoter)
		REFERENCES Users(UserID)
);`

type saved_postupvotes_row struct {
	postupvotesID int
	post          int
	upvoter       int
	Date          time.Time
	vote          int
}

func (s *saved_postupvotes_row) String() string {
	return fmt.Sprintf("PostUpvotesID : %d | Post : %d | Upvoter : %d | Date : %s | Vote : %d", s.postupvotesID, s.post, s.upvoter, Helpers.FormatTime(s.Date), s.vote)
}

func ReadFromRowPostUpvotes(row *sql.Rows) saved_postupvotes_row {
	r := saved_postupvotes_row{}
	var date string
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&r.postupvotesID, &r.post, &r.upvoter, &date, &r.vote))
	r.Date = Helpers.ParseTime(date)
	return r
}

func SavePostUpvotes(c *sql.DB, upvoter string, postID int, vote int) Helpers.ResponseUpvoteJSON {
	upvoterID := GetUser(c, upvoter).userID
	rows, err := c.Query("SELECT Vote FROM PostUpvotes WHERE Upvoter = ? AND Post = ?", upvoterID, postID)
	Helpers.ServerRuntimeError("Error While Querying PostUpvotes", err)
	defer rows.Close()
	if rows.Next() {
		var value int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&value))
		if value == vote {
			Helpers.OnlyServerError("Already Voted")
		}
	}
	_, err = c.Exec("INSERT INTO PostUpvotes (Post, Upvoter, Date, Vote) VALUES (?, ?, ?, ?)", postID, upvoterID, Helpers.FormatTime(time.Now()), vote)
	Helpers.ServerRuntimeError("Error While Saving PostUpvotes", err)
	return Helpers.ResponseUpvoteJSON{Acceptedvalue: vote, PostVotes: GetPostVotesTotal(c, postID)}
}

func GetPostUpvotes(c *sql.DB, postID int) []saved_postupvotes_row {
	rows, err := c.Query("SELECT * FROM PostUpvotes WHERE Post = ?", postID)
	Helpers.ServerRuntimeError("Error While Querying PostUpvotes", err)
	defer rows.Close()
	var result []saved_postupvotes_row
	for rows.Next() {
		result = append(result, ReadFromRowPostUpvotes(rows))
	}
	return result
}

func GetPostVotesTotal(c *sql.DB, postID int) int {
	rows, err := c.Query("SELECT Vote FROM PostUpvotes WHERE Post = ?", postID)
	Helpers.ServerRuntimeError("Error While Querying PostUpvotes", err)
	defer rows.Close()
	var result int = 0
	for rows.Next() {
		var value int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&value))
		result += value
	}
	return result
}

func GetPostUpvotesByUser(c *sql.DB, upvoter int) []saved_postupvotes_row {
	rows, err := c.Query("SELECT * FROM PostUpvotes WHERE Upvoter = ?", upvoter)
	Helpers.ServerRuntimeError("Error While Querying PostUpvotes", err)
	defer rows.Close()
	var result []saved_postupvotes_row
	for rows.Next() {
		result = append(result, ReadFromRowPostUpvotes(rows))
	}
	return result
}

func GetUserVotedShitPostsIDS(c *sql.DB, upvoter int) []int {
	rows := GetPostUpvotesByUser(c, upvoter)
	var result []int
	for _, row := range rows {
		result = append(result, row.post)
	}
	return result
}

func showPostUpvotesTable(c *sql.DB) {
	log.Println("Showing PostUpvotes")
	rows, err := c.Query("SELECT * FROM PostUpvotes")
	Helpers.ServerRuntimeError("Error While Querying PostUpvotes", err)
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ReadFromRowPostUpvotes(rows))
	}
}
