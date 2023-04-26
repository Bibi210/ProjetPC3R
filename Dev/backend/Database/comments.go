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
	Upvotes INTEGER NOT NULL,
	FOREIGN KEY (Post)
		REFERENCES ShitPost(ShitPostID),
	FOREIGN KEY (Message)
		REFERENCES Msg(MsgID)
);`

type comment_row struct {
	Comid   int
	post    int
	msg     msg_row
	upvotes int
}

func (c *comment_row) String() string {
	return fmt.Sprintf("ComID : %d | Post : %d | Msg : %v | Upvotes : %d", c.Comid, c.post, c.msg, c.upvotes)
}

func ReadFromRowComment(c *sql.DB, row *sql.Rows) comment_row {
	r := comment_row{}
	msgID := 0
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&r.Comid, &r.post, &msgID, &r.upvotes))
	r.msg = GetMsg(c, msgID)
	return r
}

func SendComment(c *sql.DB, sender string, shitpostID int, content string) {
	msg := SendMsg(c, sender, content)
	GetShitPost(c, shitpostID) // Check if shitpost exist
	executeRequest(c, "INSERT INTO Comments (Post,Message,Upvotes) VALUES (?,?,?)", shitpostID, msg.msgID, 0)
}

func DeleteComment(c *sql.DB, comID int) {
	executeRequest(c, "DELETE FROM Comments WHERE ComID = ?", comID)
}

func GetComment(c *sql.DB, comID int) comment_row {
	rows := query(c, "SELECT * FROM Comments WHERE ComID = ?", comID)
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("Comment don't exist")
	}
	return ReadFromRowComment(c, rows)
}

func GetCommentAsJSON(c *sql.DB, comID int) Helpers.ResponseCommentJSON {
	com := GetComment(c, comID)
	return Helpers.ResponseCommentJSON{Msg: GetMsgAsJSON(c, com.msg.msgID), Upvotes: com.upvotes}
}

func UpvoteComment(c *sql.DB, comID int) {
	executeRequest(c, "UPDATE Comments SET Upvotes = Upvotes + 1 WHERE ComID = ?", comID)
}

func showCommentTable(c *sql.DB) {
	fmt.Println("Comments Table :")
	rows := query(c, "SELECT * FROM Comments")
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ReadFromRowComment(c, rows))
	}
}
