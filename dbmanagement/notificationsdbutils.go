package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
)

func AddNotification(receivingUserId, postId, commentId, sendingUserId string, reaction int, notificationStatement string) {
	receiverName, _ := SelectUserFromUUID(receivingUserId)
	senderName, _ := SelectUserFromUUID(sendingUserId)

	if receiverName.UUID == senderName.UUID {
		return
	}

	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting notification record...")

	UUID := GenerateUUIDString()

	if notificationStatement == "" {
		if postId != "" && reaction != 0 {
			if reaction == 1 {
				notificationStatement = "liked your post"
			} else {
				notificationStatement = "disliked your post"
			}
		} else if postId != "" && commentId != "" {
			notificationStatement = "commented on your post"
		} else if commentId != "" && reaction != 0 {
			if reaction == 1 {
				notificationStatement = "liked your comment"
			} else {
				notificationStatement = "disliked your comment"
			}
		}
	}

	insertNotificationData := "INSERT INTO Notifications(UUID, receivingUserId, postId, commentId, sendingUserId, reaction, notificationStatement) VALUES (?, ?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertNotificationData)
	utils.HandleError("Notification Prepare failed: ", err)

	_, err = statement.Exec(UUID, receiverName.UUID, postId, commentId, senderName.Name, reaction, notificationStatement)
	utils.HandleError("Notification Statement Exec failed: ", err)
}

func SelectAllNotificationsFromUser(receiver string) []Notification {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Notifications WHERE receivingUserId = ?", receiver)
	utils.HandleError("Notification from User query failed: ", err)
	defer row.Close()

	var allNotifications []Notification

	for row.Next() {
		var currentNotification Notification
		row.Scan(&currentNotification.UUID, &currentNotification.Receiver, &currentNotification.PostId, &currentNotification.CommentId, &currentNotification.Sender, &currentNotification.Reaction, &currentNotification.Statement)
		message := fmt.Sprint("current notification has post id: ", currentNotification.PostId)
		utils.WriteMessageToLogFile(message)
		allNotifications = append(allNotifications, currentNotification)
	}
	return allNotifications
}

func SelectAllNotificationsFromUUID(UUID string) []Notification {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	row, err := db.Query("SELECT * FROM Notifications WHERE postId = ?", UUID)
	utils.HandleError("Notification from UUID query failed: ", err)
	defer row.Close()
	var allNotifications []Notification
	for row.Next() {
		var currentNotification Notification
		row.Scan(&currentNotification.UUID, &currentNotification.Receiver, &currentNotification.PostId, &currentNotification.CommentId, &currentNotification.Sender, &currentNotification.Reaction, &currentNotification.Statement)
		allNotifications = append(allNotifications, currentNotification)
	}
	return allNotifications
}
