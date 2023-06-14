package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"log"
	"sort"
	"time"
)

// create insert query
func InsertTextInChat(Text ChatText) error {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	//utils.WriteMessageToLogFile("Inserting text msg in ChatDB...")
	log.Println("Inserting text msg in ChatDB...")

	//check whether a conversation exists already

	UUID, exists := SelectChatId(Text.SenderId, Text.ReceiverId)

	if !exists {
		//generate UUID

		UUID = GenerateUUIDString()
		log.Println("Generating a UUID of value: ", UUID)
	}

	query := "INSERT INTO Chat(uuid, sender, receiver, text, time) VALUES(?, ?, ?, ?, ?);"
	statement, err := db.Prepare(query)
	if err != nil {
		// utils.HandleError("Chat INSERT Prepare failed: ", err)
		log.Println("Failed to prepare INSERT statement in Chat", err)
		return err
	}

	_, err = statement.Exec(UUID, Text.SenderId, Text.ReceiverId, Text.Content, Text.Time)
	if err != nil {
		// utils.HandleError("Chat INSERT statement execution failed: ", err)
		log.Println("Failed to execute INSERT statement in Chat", err)
		return err
	}

	return nil
}

// create select query
func SelectAllChat(ChatId string) ChatBox {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	query := fmt.Sprintf("SELECT sender, receiver, text, time FROM Chat WHERE uuid = '%s';", ChatId)

	row, err := db.Query(query)
	// utils.HandleError("ChatBox SELECT query failed: ", err)
	if err != nil {
		log.Println("ChatBox SELECT query failed: ", err)
	}
	defer row.Close()

	var chat ChatBox
	var texts []ChatText //collection of all individual text messages which will be sorted chronologically and appended to the above variable

	for row.Next() {
		var text ChatText
		err := row.Scan(&text.SenderId, &text.ReceiverId, &text.Content, &text.Time)
		if err != nil {
			log.Println("Chat SELECT query failed to scan: ", err)
			return ChatBox{}
		}
		texts = append(texts, text)
	}

	sort.Slice(texts, func(i, j int) bool {

		t1, err := time.Parse(time.RFC3339, texts[i].Time)
		if err != nil {
			log.Fatal("SelectAllChat sorting slice failed to parse t1: ", texts[i].Time)
			return false
		}
		t2, err := time.Parse(time.RFC3339, texts[j].Time)
		if err != nil {
			log.Fatal("SelectAllChat sorting slice failed to parse t2: ", texts[j].Time)
			return false
		}
		return t1.Before(t2)
	})
	chat.UUID = ChatId
	chat.Content = texts

	return chat
}

func SelectChatId(senderId, receiverId string) (string, bool) {

	var UUID string

	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	query := "SELECT uuid FROM Chat WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?) LIMIT 1"

	err := db.QueryRow(query, senderId, receiverId, receiverId, senderId).Scan(&UUID)

	if err != nil {
		log.Println("Chat UUID SELECT query failed: ", err)
		return "", false
	}

	return UUID, true
}

// to be implemented later
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
