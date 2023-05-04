package Database

import (
	"Backend/Helpers"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const CreateComments = `CREATE TABLE IF NOT EXISTS Comments (
	ComID INTEGER PRIMARY KEY AUTOINCREMENT,
	Post INTEGER NOT NULL,
	Message INTEGER NOT NULL,
	FOREIGN KEY (Post)
		REFERENCES ShitPost(ShitPostID),
	FOREIGN KEY (Message)
		REFERENCES Msg(MsgID)
);`

type CommentRow struct {
	Comid int
	post  int
	msg   MsgRow
}

func (c *CommentRow) String() string {
	return fmt.Sprintf("ComID : %d | Post : %d | Msg : %v ", c.Comid, c.post, c.msg)
}

func ReadFromRowComment(c *sql.DB, row *sql.Rows) CommentRow {
	r := CommentRow{}
	msgID := 0
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&r.Comid, &r.post, &msgID))
	r.msg = GetMsg(c, msgID)
	return r
}

func SendComment(c *sql.DB, sender string, shitpostID int, content string) int {
	msgId := SendMsg(c, sender, content)
	GetShitPost(c, shitpostID) // Check if shitpost exist
	row := query(c, "INSERT INTO Comments (Post,Message) VALUES (?,?) RETURNING ComID", shitpostID, msgId)
	defer row.Close()
	var comId int
	if !row.Next() {
		Helpers.OnlyServerError("Comment doesn't exist")
	} else {
		err := row.Scan(&comId)
		if err != nil {
			Helpers.OnlyServerError("Can't read comment row")
		}
	}
	return comId
}

func DeleteComment(c *sql.DB, comID int) {
	executeRequest(c, "DELETE FROM Comments WHERE ComID = ?", comID)
}

func GetComment(c *sql.DB, comID int) CommentRow {
	rows := query(c, "SELECT * FROM Comments WHERE ComID = ?", comID)
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("Comment don't exist")
	}
	return ReadFromRowComment(c, rows)
}

func GetCommentAsJSON(c *sql.DB, comID int) Helpers.ResponseCommentJSON {
	com := GetComment(c, comID)
	return Helpers.ResponseCommentJSON{Id: com.msg.msgID, Msg: GetMsgAsJSON(c, com.msg.msgID), Upvotes: GetCommentVotesTotal(c, comID)}
}

func GetCommentListAsJSON(c *sql.DB, comIDs []int) []Helpers.ResponseCommentJSON {
	var out []Helpers.ResponseCommentJSON
	for _, comID := range comIDs {
		out = append(out, GetCommentAsJSON(c, comID))
	}
	return out
}

func showCommentTable(c *sql.DB) {
	fmt.Println("Comments Table :")
	rows := query(c, "SELECT * FROM Comments")
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ReadFromRowComment(c, rows))
	}
}
