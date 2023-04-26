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

type msg_row struct {
	msgID   int
	sender  int
	content string
	Date    time.Time
}

func (m *msg_row) String() string {
	return fmt.Sprintf("MsgID : %d | Sender : %d | Content : %s | Date : %s", m.msgID, m.sender, m.content, Helpers.FormatTime(m.Date))
}

func ReadFromRowMsg(row *sql.Rows) msg_row {
	r := msg_row{}
	var date string
	Helpers.ServerRuntimeError("Error While Reading Row", row.Scan(&r.msgID, &r.sender, &r.content, &date))
	r.Date = Helpers.ParseTime(date)
	return r
}

func SendMsg(c *sql.DB, sender string, content string) msg_row {
	userID := GetUser(c, sender).userID
	executeRequest(c, "INSERT INTO Msg (Sender, Content, Date) VALUES (?, ?, ?)", userID, content, Helpers.FormatTime(time.Now()))
	rows := query(c, "SELECT * FROM Msg WHERE MsgID = (SELECT MAX(MsgID) FROM Msg)")
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("Msg don't exist")
	}
	return ReadFromRowMsg(rows)
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
	v := ReadFromRowMsg(rows)
	return Helpers.ResponseMsgJSON{Content: v.content, Date: Helpers.FormatTime(v.Date)}
}

func GetMsg(c *sql.DB, msgID int) msg_row {
	rows := query(c, "SELECT * FROM Msg WHERE MsgID = ?", msgID)
	defer rows.Close()
	if !rows.Next() {
		Helpers.OnlyServerError("Msg don't exist")
	}
	return ReadFromRowMsg(rows)
}

func showMsgTable(c *sql.DB) {
	fmt.Println("Msg Table :")
	rows := query(c, "SELECT * FROM Msg")
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ReadFromRowMsg(rows))
	}
}
