package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
)

func InsertTag(tag string) Tag {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting tag record...")

	UUID := GenerateUUIDString()
	insertPostData := "INSERT INTO Tags(uuid, tagname) VALUES (?, ?)"
	statement, err := db.Prepare(insertPostData)
	utils.HandleError("User Prepare failed", err)

	_, err = statement.Exec(UUID, tag)
	utils.HandleError("Statement Exec failed: ", err)

	return Tag{UUID, tag}
}

func InsertTaggedPost(tagId string, postId string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	message := fmt.Sprintln("Inserting tagged post record...")
	utils.WriteMessageToLogFile(message)

	insertTaggedPost := "INSERT INTO TaggedPosts(tagId, postId) VALUES (?, ?)"
	statement, err := db.Prepare(insertTaggedPost)
	utils.HandleError("User Prepare failed: ", err)

	_, err = statement.Exec(tagId, postId)
	utils.HandleError("Statement Exec failed: ", err)
}

func UpdateTaggedPost(postId string) {
	message := fmt.Sprintln("Updating tagged post record...")
	utils.WriteMessageToLogFile(message)
	DeleteFromTableWithPostId("TaggedPosts", postId)
}

func SelectAllTags() []Tag {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Tags")
	utils.HandleError("All tags query failed", err)
	defer row.Close()

	var allTags []Tag
	for row.Next() {
		var currentTag Tag
		row.Scan(&currentTag.UUID, &currentTag.TagName)
		allTags = append(allTags, currentTag)
	}
	return allTags
}

func SelectTagFromUUID(UUID string) (Tag, error) {
	var tag Tag
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Tags WHERE uuid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(UUID).Scan(&tag.UUID, &tag.TagName)
	utils.HandleError("Query Row failed: ", err)

	return tag, err
}

func SelectTagFromName(tagName string) (Tag, error) {
	var tag Tag
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Tags WHERE TagName = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(tagName).Scan(&tag.UUID, &tag.TagName)
	utils.HandleError("Query Row failed: ", err)

	return tag, err
}

func SelectAllTagsFromPost(postId string) []Tag {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT tagId FROM TaggedPosts WHERE postId = ?", postId)
	utils.HandleError("Tag query failed: ", err)
	defer row.Close()

	var allTags []Tag

	for row.Next() {
		var currentTagId string
		var currentTag Tag
		row.Scan(&currentTagId)
		currentTag, err = SelectTagFromUUID(currentTagId)
		utils.HandleError("couldn't select tag", err)
		allTags = append(allTags, currentTag)
	}
	return allTags
}
