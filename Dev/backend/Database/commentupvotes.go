package Database

import (
	"Backend/Helpers"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const createCommentUpvotes = `CREATE TABLE IF NOT EXISTS CommentUpvotes (
	CommentUpvotesID INTEGER PRIMARY KEY AUTOINCREMENT,
	Comment INTEGER NOT NULL,
	Upvoter INTEGER NOT NULL,
	Date TEXT NOT NULL,
	Vote INTEGER NOT NULL,
	FOREIGN KEY (Comment)
		REFERENCES Comments(CommentID),
	FOREIGN KEY (Upvoter)
		REFERENCES Users(UserID)
);`

type saved_commentupvotes_row struct {
	commentupvotesID int
	comment          int
	upvoter          int
	Date             time.Time
	vote             int
}

func (s *saved_commentupvotes_row) String() string {
	return fmt.Sprintf("CommentUpvotesID : %d | Comment : %d | Upvoter : %d | Date : %s | Vote : %d", s.commentupvotesID, s.comment, s.upvoter, Helpers.FormatTime(s.Date), s.vote)
}

func ReadFromRowCommentUpvotes(row *sql.Rows) saved_commentupvotes_row {
	r := saved_commentupvotes_row{}
	var date string
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&r.commentupvotesID, &r.comment, &r.upvoter, &date, &r.vote))
	r.Date = Helpers.ParseTime(date)
	return r
}

func SaveCommentUpvotes(c *sql.DB, upvoter string, commentID int, vote int) Helpers.ResponseUpvoteJSON {
	upvoterID := GetUser(c, upvoter).userID
	if vote != 1 && vote != -1 {
		Helpers.OnlyServerError("Invalid Vote Value")
	}
	rows, err := c.Query("SELECT Vote FROM CommentUpvotes WHERE Upvoter = ? AND Comment = ?", upvoterID, commentID)
	Helpers.ServerRuntimeError("Error While Querying CommentUpvotes", err)
	defer rows.Close()
	if rows.Next() {
		var value int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&value))
		if value == vote {
			Helpers.OnlyServerError("Already Voted")
		}
	}
	_, err = c.Exec("INSERT INTO CommentUpvotes (Comment, Upvoter, Date, Vote) VALUES (?, ?, ?, ?)", commentID, upvoterID, Helpers.FormatTime(time.Now()), vote)
	Helpers.ServerRuntimeError("Error While Saving CommentUpvotes", err)
	return Helpers.ResponseUpvoteJSON{Acceptedvalue: vote, PostVotes: GetCommentVotesTotal(c, commentID)}
}

func GetCommentUpvotes(c *sql.DB, commentID int) []saved_commentupvotes_row {
	rows, err := c.Query("SELECT * FROM CommentUpvotes WHERE Comment = ?", commentID)
	Helpers.ServerRuntimeError("Error While Querying CommentUpvotes", err)
	defer rows.Close()
	var result []saved_commentupvotes_row
	for rows.Next() {
		result = append(result, ReadFromRowCommentUpvotes(rows))
	}
	return result
}

func GetCommentVotesTotal(c *sql.DB, commentID int) int {
	rows, err := c.Query("SELECT Vote FROM CommentUpvotes WHERE Comment = ?", commentID)
	Helpers.ServerRuntimeError("Error While Querying CommentUpvotes", err)
	defer rows.Close()
	var result int = 0
	if rows.Next() {
		var value int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&value))
		result += value
	}
	return result
}

func GetCommentVotesByUser(c *sql.DB, upvoter int) []saved_commentupvotes_row {
	rows, err := c.Query("SELECT * FROM CommentUpvotes WHERE Upvoter = ?", upvoter)
	Helpers.ServerRuntimeError("Error While Querying CommentUpvotes", err)
	defer rows.Close()
	var result []saved_commentupvotes_row
	for rows.Next() {
		result = append(result, ReadFromRowCommentUpvotes(rows))
	}
	return result
}

func GetUserVotedCommentsIDS(c *sql.DB, upvoter int) []int {
	rows := GetCommentVotesByUser(c, upvoter)
	var result []int
	for _, row := range rows {
		result = append(result, row.comment)
	}
	return result
}

func showCommentUpvotesTable(c *sql.DB) {
	log.Println("CommentUpvotes")
	rows, err := c.Query("SELECT * FROM CommentUpvotes")
	Helpers.ServerRuntimeError("Error While Querying CommentUpvotes", err)
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ReadFromRowCommentUpvotes(rows))
	}
}
