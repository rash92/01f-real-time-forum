package dbmanagement

import (
	"database/sql"
	"forum/utils"
)

func CreateAdminRequest(requestFromId string, requestFromName string, reportedPostId string, reportedCommentId string, reportedUserId string, description string) AdminRequest {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting admin request record...")

	UUID := GenerateUUIDString()
	insertAdminRequestData := "INSERT INTO AdminRequests(UUID, requestfromid, requestfromname, reportedpostid, reportedcommentid, reporteduserid, description) VALUES (?, ?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertAdminRequestData)
	utils.HandleError("User Prepare failed: ", err)

	utils.WriteMessageToLogFile("admint request content is: " + description)

	_, err = statement.Exec(UUID, requestFromId, requestFromName, reportedPostId, reportedCommentId, reportedUserId, description)
	utils.HandleError("Statement Exec failed: ", err)

	return AdminRequest{UUID, requestFromId, requestFromName, reportedPostId, reportedCommentId, reportedUserId, description}
}

func SelectAllAdminRequests() []AdminRequest {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM AdminRequests")
	utils.HandleError("Admin Request query failed: ", err)
	defer row.Close()

	var allAdminRequests []AdminRequest
	for row.Next() {
		var currentAdminRequest AdminRequest
		row.Scan(&currentAdminRequest.UUID, &currentAdminRequest.RequestFromId, &currentAdminRequest.RequestFromName, &currentAdminRequest.ReportedPostId, &currentAdminRequest.ReportedCommentId, &currentAdminRequest.ReportedUserId, &currentAdminRequest.Description)
		allAdminRequests = append(allAdminRequests, currentAdminRequest)
	}
	return allAdminRequests
}

func SelectAdminRequestFromUUID(UUID string) AdminRequest {
	var adminRequest AdminRequest
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM AdminRequests WHERE uuid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(UUID).Scan(&adminRequest.UUID, &adminRequest.RequestFromId, &adminRequest.RequestFromName, &adminRequest.ReportedPostId, &adminRequest.ReportedCommentId, &adminRequest.ReportedUserId, &adminRequest.Description)
	utils.HandleError("Query Row failed: ", err)

	return adminRequest
}
