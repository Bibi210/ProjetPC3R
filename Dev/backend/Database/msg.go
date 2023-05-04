package Database

import (
	"Backend/Helpers"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const createMsg = `CREATE TABLE IF NOT EXISTS Msg (
	MsgID INTEGER PRIMARY KEY AUTOINCREMENT,
	Sender INTEGER NOT NULL,
	Content TEXT NOT NULL,
	Date TEXT NOT NULL,
	FOREIGN KEY (Sender)
		REFERENCES Users(UserID)
);`

type MsgRow struct {
	msgID   int
	sender  string
	content string
	Date    time.Time
}

func (m *MsgRow) String() string {
	return fmt.Sprintf("MsgID : %d | Sender : %s | Content : %s | Date : %s", m.msgID, m.sender, m.content, Helpers.FormatTime(m.Date))
}

func ReadFromRowMsg(c *sql.DB, row *sql.Rows) MsgRow {
	r := MsgRow{}
	var date string
	var userId int
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&r.msgID, &userId, &r.content, &date))
	r.Date = Helpers.ParseTime(date)
	r.sender = GetUserByID(c, userId).username
	return r
}

func SendMsg(c *sql.DB, sender string, content string) int {
	userID := GetUser(c, sender).userID
	row := query(c, "INSERT INTO Msg (Sender, Content, Date) VALUES (?, ?, ?) RETURNING MsgID", userID, content, Helpers.FormatTime(time.Now()))
	defer row.Close()
	var msgId int
	if !row.Next() {
		Helpers.OnlyServerError("Msg doesn't exist")
	} else {
		err := row.Scan(&msgId)
		if err != nil {
			Helpers.OnlyServerError("Can't read msg row")
		}
	}
	return msgId
}

func DeleteMsg(c *sql.DB, msgID int) {
	executeRequest(c, "DELETE FROM Msg WHERE MsgID = ?", msgID)
}

func GetMsgAsJSON(c *sql.DB, msgID int) Helpers.ResponseMsgJSON {
	rows := query(c, "SELECT * FROM Msg WHERE MsgID = ?", msgID)
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("Msg don't exist")
	}
	v := ReadFromRowMsg(c, rows)
	return Helpers.ResponseMsgJSON{Sender: v.sender, Content: v.content, Date: Helpers.FormatTime(v.Date)}
}

func GetMsg(c *sql.DB, msgID int) MsgRow {
	rows := query(c, "SELECT * FROM Msg WHERE MsgID = ?", msgID)
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("Msg doesn't exist")
	}
	return ReadFromRowMsg(c, rows)
}

func showMsgTable(c *sql.DB) {
	fmt.Println("Msg Table :")
	rows := query(c, "SELECT * FROM Msg")
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ReadFromRowMsg(c, rows))
	}
}
