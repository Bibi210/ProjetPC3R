package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const createDM = `CREATE TABLE IF NOT EXISTS DM (
	DmID INTEGER PRIMARY KEY AUTOINCREMENT,
	Receiver INTEGER NOT NULL,
	Message INTEGER NOT NULL,
	FOREIGN KEY (Receiver)
		REFERENCES Users(UserID),
	FOREIGN KEY (Message)
		REFERENCES Msg(MsgID)
);`

type dm_row struct {
	DMid     int
	receiver int
	msg      msg_row
}

func (d *dm_row) String() string {
	return fmt.Sprintf("Receiver : %d | Message : %v", d.receiver, d.msg)
}

func ReadFromRowDM(c *sql.DB, row *sql.Rows) dm_row {
	r := dm_row{}
	ServerRuntimeError("Error While Reading Row", row.Scan(&r.DMid, &r.receiver, &r.msg.msgID))
	r.msg = GetMsg(c, r.msg.msgID)
	return r
}

func SendDM(c *sql.DB, sender username, receiver username, content string) {
	msg := SendMsg(c, sender, content)
	executeRequest(c, "INSERT INTO DM (Receiver,Message) VALUES (?,?)", getUser(c, receiver).userID, msg.msgID)
}

func DeleteDM(c *sql.DB, dmID int) {
	executeRequest(c, "DELETE FROM DM WHERE DmID = ?", dmID)
}

func GetDM(c *sql.DB, dmID int) dm_row {
	rows := query(c, "SELECT * FROM DM WHERE DmID = ?", dmID)
	defer rows.Close()
	if !rows.Next() {
		OnlyServerError("DM don't exist")
	}
	return ReadFromRowDM(c, rows)
}

func showDMTable(c *sql.DB) {
	fmt.Println("DM Table :")
	rows := query(c, "SELECT * FROM DM")
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ReadFromRowDM(c, rows))
	}
}
