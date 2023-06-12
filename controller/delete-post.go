package controller

import (
	"encoding/json"
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type DeletePostFormData struct {
	PostID string `json:"editPost"`
}

func DeletePost(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("session error is: ", err)
	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		utils.HandleError("cannot select user from session id ", err)
		data.Cookie = sessionId

		data.UserInfo = user
		message := fmt.Sprintf("session id is: %v user info is: %v cookie data is: %v", sessionId, data.UserInfo, data.Cookie)
		utils.WriteMessageToLogFile(message)

		if r.Method == "POST" {

			var formData DeletePostFormData
			jerr := json.NewDecoder(r.Body).Decode(&formData)
			if jerr != nil {
				// Handle error
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			idToDelete := formData.PostID

			if idToDelete != "" {
				post, err := dbmanagement.SelectPostFromUUID(idToDelete)
				if err != nil {
					PageErrors(w, r, tmpl, 500, "Internal Server Error")
					return
				}
				dbmanagement.DeletePostWithUUID(idToDelete)
				message := fmt.Sprintf("Deleting post with id: %v and contents: %v", idToDelete, post)
				utils.WriteMessageToLogFile(message)
			}

		}

	}
}
