package Database

import (
	"Backend/Helpers"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const createShitPost = `CREATE TABLE IF NOT EXISTS ShitPost (
	ShitPostID INTEGER PRIMARY KEY AUTOINCREMENT,
	Poster INTEGER NOT NULL,
	Caption TEXT NOT NULL,
	URL TEXT NOT NULL,
	Date TEXT NOT NULL,
	FOREIGN KEY (Poster) 
		REFERENCES Users(UserID)
);`

type saved_shitpost_row struct {
	shitpostID int
	poster     int
	caption    string
	url        string
	Date       time.Time
}

func (s *saved_shitpost_row) String() string {
	return fmt.Sprintf("ShitPostID : %d | Poster : %d | Caption : %s | URL : %s | Date : %s", s.shitpostID, s.poster, s.caption, s.url, Helpers.FormatTime(s.Date))
}

func ReadFromRowShitPost(row *sql.Rows) saved_shitpost_row {
	r := saved_shitpost_row{}
	var date string
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&r.shitpostID, &r.poster, &r.caption, &r.url, &date))
	r.Date = Helpers.ParseTime(date)
	return r
}

func SaveShitPost(c *sql.DB, poster string, url string, caption string) Helpers.ResponseSavedShitPostJSON {
	userID := GetUser(c, poster).userID
	executeRequest(c, "INSERT INTO ShitPost (Poster, Caption, URL, Date) VALUES (?, ?, ?, ?)", userID, caption, url, Helpers.FormatTime(time.Now()))
	rows := query(c, "SELECT last_insert_rowid()")
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("Can't get last inserted row")
	}
	var shitpostID int
	Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&shitpostID))
	return GetShitPostAsJSON(c, shitpostID)
}

func DeleteShitPost(c *sql.DB, shitpostID int) {
	executeRequest(c, "DELETE FROM ShitPost WHERE ShitPostID = ?", shitpostID)
}

func GetShitPost(c *sql.DB, shitpostID int) saved_shitpost_row {
	rows := query(c, "SELECT * FROM ShitPost WHERE ShitPostID = ?", shitpostID)
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("ShitPost don't exist")
	}
	return ReadFromRowShitPost(rows)
}

func GetShitPostCommentsIds(c *sql.DB, shitpostID int) []int {
	rows := query(c, "SELECT ComID FROM Comments WHERE Post = ?", shitpostID)
	defer rows.Close()
	var result []int
	for rows.Next() {
		var commentID int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&commentID))
		result = append(result, commentID)
	}
	return result
}

func GetShitPostAsJSON(c *sql.DB, shitpostID int) Helpers.ResponseSavedShitPostJSON {
	shitpost := GetShitPost(c, shitpostID)
	return Helpers.ResponseSavedShitPostJSON{
		Url:        shitpost.url,
		Caption:    shitpost.caption,
		Creator:    shitpost.poster,
		Date:       Helpers.FormatTime(shitpost.Date),
		Upvotes:    GetPostVotesTotal(c, shitpostID),
		CommentIds: GetShitPostCommentsIds(c, shitpostID),
	}
}

func GetShitPostListAsJSON(c *sql.DB, shitpostID []int) []Helpers.ResponseSavedShitPostJSON {
	var result []Helpers.ResponseSavedShitPostJSON
	for _, id := range shitpostID {
		result = append(result, GetShitPostAsJSON(c, id))
	}
	return result
}

func GetShitPostComments(c *sql.DB, shitpostID int) []Helpers.ResponseCommentJSON {
	ids := GetShitPostCommentsIds(c, shitpostID)
	var result []Helpers.ResponseCommentJSON
	for _, id := range ids {
		result = append(result, GetCommentAsJSON(c, id))
	}
	return result
}

func GetAllShitPostsID(c *sql.DB) []int {
	rows := query(c, "SELECT ShitPostID FROM ShitPost")
	defer rows.Close()
	var result []int
	for rows.Next() {
		var shitpostID int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&shitpostID))
		result = append(result, shitpostID)
	}
	return result
}

func SearchShitPost(c *sql.DB, input string) []int {
	input = "%" + input + "%"
	rows := query(c, "SELECT ShitPostID FROM ShitPost WHERE Caption LIKE ? OR URL LIKE ?", input, input)
	defer rows.Close()
	var result []int
	for rows.Next() {
		var shitpostID int
		Helpers.ServerRuntimeError("Error While Reading Row", rows.Scan(&shitpostID))
		result = append(result, shitpostID)
	}
	return result
}

func GetTopPostsIDs(c *sql.DB, limit int) []int {
	ids := GetAllShitPostsID(c)
	sort.Slice(ids, func(i, j int) bool {
		return GetPostVotesTotal(c, ids[i]) > GetPostVotesTotal(c, ids[j])
	},
	)
	if len(ids) < limit {
		return ids
	}
	return ids[:limit]
}

func showShitPostTable(c *sql.DB) {
	fmt.Println("ShitPost Table : ")
	rows := query(c, "SELECT * FROM ShitPost")
	defer rows.Close()
	for rows.Next() {
		log.Println(ReadFromRowShitPost(rows))
	}
}
