package controller

import (
	"encoding/json"
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"log"

	"forum/utils"
	"html/template"
	"net/http"
)

type AdminData struct {
	AllUsers         []dbmanagement.User
	AdminRequests    []dbmanagement.AdminRequest
	ReportedPosts    []dbmanagement.Post
	ReportedComments []dbmanagement.Comment
	ReportedUsers    []dbmanagement.User
	TitleName        string
	UserInfo         dbmanagement.User
	TagsList         []dbmanagement.Tag
}

type AdminFormData struct {
	SetToUserID         string `json:"set_to_user"`
	SetToModeratorID    string `json:"set_to_mod"`
	SetToAdminID        string `json:"set_to_admin"`
	DeleteUserID        string `json:"delete_user"`
	TagsToCreate        string `json:"tag_create"`
	DeleteTag           string `json:"delete_tag"`
	AllTagDelete        string `json:"delete_all_tag"`
	DeleteRequestID     string `json:"delete_request"`
	AcknowledgeReportID string `json:"acknowledge_report"`
	ResponseMessage     string `json:"response_message"`
}

// username: admin password: admin for existing user with admin permissions, can create and change other users to be admin while logged in as anyone who is admin
func Admin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	adminData := AdminData{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("Unable to get session from browser in admin handler", err)
	user, err := dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("Unable to user session from admin handler", err)
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		utils.WriteMessageToLogFile("Please log in as a user with admin permissions")
		return
	}

	loggedInAs, err := dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("Unable get logged in user in admin", err)
	if loggedInAs.Permission != "admin" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		utils.WriteMessageToLogFile("Please log in as a user with admin permissions")
		return
	}

	if r.Method == "POST" {
		// Parse the JSON body into FormData struct
		var formData AdminFormData
		jerr := json.NewDecoder(r.Body).Decode(&formData)
		if jerr != nil {
			// Handle error
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := dbmanagement.UpdateUserToken(user.UUID, 1)
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
		dbmanagement.UpdateUserToken(user.UUID, 1)
		userToChange := formData.SetToUserID
		if userToChange != "" {
			dbmanagement.UpdateUserPermissionFromUUID(userToChange, "user")
		}
		userToChange = formData.SetToModeratorID
		if userToChange != "" {
			dbmanagement.UpdateUserPermissionFromUUID(userToChange, "moderator")
		}
		userToChange = formData.SetToAdminID
		if userToChange != "" {
			dbmanagement.UpdateUserPermissionFromUUID(userToChange, "admin")
		}
		userToChange = formData.DeleteUserID
		if userToChange != "" {
			dbmanagement.DeleteFromTableWithUUID("Users", userToChange)
		}
		tagsToCreate := formData.TagsToCreate
		if tagsToCreate != "" {
			dbmanagement.InsertTag(tagsToCreate)
		}
		tagToChange := formData.DeleteTag
		if tagToChange != "" {
			dbmanagement.DeleteFromTableWithUUID("Tags", tagToChange)
		}
		tagToDeletePostsLinkedTo := formData.AllTagDelete
		if tagToDeletePostsLinkedTo != "" {
			dbmanagement.DeleteAllPostsWithTag(tagToDeletePostsLinkedTo)
		}
		adminRequestToDelete := formData.DeleteRequestID
		if adminRequestToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("AdminRequests", adminRequestToDelete)
		}
		adminRequestToAcknowledge := formData.AcknowledgeReportID
		if adminRequestToAcknowledge != "" {
			request := dbmanagement.SelectAdminRequestFromUUID(adminRequestToAcknowledge)
			responseMessage := formData.ResponseMessage
			if responseMessage == "" {
				dbmanagement.AddNotification(request.RequestFromId, request.ReportedPostId, "", loggedInAs.UUID, 0, "The admin has recieved your report")
			} else {
				dbmanagement.AddNotification(request.RequestFromId, request.ReportedPostId, "", loggedInAs.UUID, 0, responseMessage)
			}
		}
	}
	utils.HandleError("cant get user", err)
	adminData.AllUsers = dbmanagement.SelectAllUsers()
	adminData.AdminRequests = dbmanagement.SelectAllAdminRequests()
	for _, adminRequest := range adminData.AdminRequests {
		if adminRequest.ReportedPostId != "" {
			postdata, err := dbmanagement.SelectPostFromUUID(adminRequest.ReportedPostId)
			if err != nil {
				PageErrors(w, r, tmpl, 500, "Internal Server Error")
				return
			}
			adminData.ReportedPosts = append(adminData.ReportedPosts, postdata)

			message := fmt.Sprint("reported posts retrieved are: ", adminRequest)
			utils.WriteMessageToLogFile(message)
		}
		if adminRequest.ReportedCommentId != "" {
			adminData.ReportedComments = append(adminData.ReportedComments, dbmanagement.SelectCommentFromUUID(adminRequest.ReportedCommentId))
		}
		if adminRequest.ReportedUserId != "" {
			currentUser, err := dbmanagement.SelectUserFromUUID(adminRequest.ReportedUserId)
			utils.HandleError("couldn't select user when looking for reported users", err)
			adminData.ReportedUsers = append(adminData.ReportedUsers, currentUser)
		}
	}
	adminData.TitleName = "Admin"
	adminData.TagsList = dbmanagement.SelectAllTags()
	adminData.UserInfo = loggedInAs
	//tmpl.ExecuteTemplate(w, "admin.html", adminData)

	// Convert the struct to JSON
	jsonData, err := json.Marshal(adminData)
	if err != nil {
		// Handle error
		log.Println(err)

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
