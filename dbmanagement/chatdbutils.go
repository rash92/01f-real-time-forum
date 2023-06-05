package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"sort"
)

// create insert query
func InsertTextInChat(Text ChatText) error {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting text msg in ChatDB...")

	//check whether a conversation exists already

	UUID, exists := SelectChatId(Text.SenderId, Text.ReceiverId)

	if !exists {
		//generate UUID
		UUID = GenerateUUIDString()
	}

	query := "INSERT INTO Chat(uuid, sender, receiver, text, time) VALUES(?, ?, ?, ?, ?);"
	statement, err := db.Prepare(query)
	utils.HandleError("Chat INSERT Prepare failed: ", err)

	_, err = statement.Exec(UUID, Text.SenderId, Text.ReceiverId, Text.Content, Text.Time.String())
	utils.HandleError("Chat INSERT statement execution failed: ", err)

	return nil
}

// create select query
func SelectAllChat(ChatId string) ChatBox {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	query := fmt.Sprintf("SELECT (sender, receiver, text, time) FROM Chat WHERE uuid = '%s';", ChatId)

	row, err := db.Query(query)
	utils.HandleError("ChatBox SELECT query failed: ", err)
	defer row.Close()

	var chat ChatBox
	var texts []ChatText //collection of all individual text messages which will be sorted chronologically and appended to the above variable

	for row.Next() {
		var text ChatText
		err := row.Scan(&text.SenderId, &text.ReceiverId, &text.Content, &text.Time)
		if err != nil {
			utils.HandleError("Chat SELECT query failed to scan: ", err)
			return ChatBox{}
		}
		texts = append(texts, text)
	}

	sort.Slice(texts, func(i, j int) bool { return texts[i].Time.Before(texts[j].Time) })
	chat.UUID = ChatId
	chat.Content = texts

	return chat
}

func SelectChatId(senderId, receiverId string) (string, bool) {

	var UUID string

	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	query := fmt.Sprintf(" SELECT uuid FROM Chat WHERE (sender = '%s' AND receiver = '%s') OR (sender = '%s' AND receiver = '%s') LIMIT 1;", senderId, receiverId, receiverId, senderId)

	row, err := db.Query(query)
	utils.HandleError("Chat UUID SELECT query failed: ", err)
	defer row.Close()

	err = row.Scan(&UUID)
	if err != nil {
		utils.HandleError("Chat UUID value scanning failed: ", err)
		return "", false
	}
	return UUID, true
}

func DeleteText(Text ChatText) error {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Deleting text msg from ChatDB...")

	//check whether a conversation exists already

	UUID, exists := SelectChatId(Text.SenderId, Text.ReceiverId)

	if !exists {
		//return error message
		utils.HandleError("Chat UUID SELECT query failed: ", fmt.Errorf("chat does not exist"))
		return fmt.Errorf("chat does not exist")
	}

	query := fmt.Sprintf("DELETE FROM Chat WHERE uuid='%s' AND sender='%s'AND time='%s';", UUID, Text.SenderId, Text.Time)
	statement, err := db.Prepare(query)
	utils.HandleError("Chat DELETE Prepare failed: ", err)

	_, err = statement.Exec()
	utils.HandleError("Chat DELETE statement execution failed: ", err)

	return err
}
