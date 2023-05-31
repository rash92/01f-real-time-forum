package dbmanagement

import (
	"database/sql"
	"forum/utils"
)

// create table
func InsertTextInChat(ChatId int, Text string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting admin request record...")

	//check whether a conversation exists already
}

//create insert query
