package controller

import (
	"encoding/json"
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"time"
)

type CommentFormData struct {
	NewCommentContent  string `json:"commentContent"`
	DeleteComment      string `json:"deleteComment"`
	EditCommentContent string `json:"editComment"`
	CommentUUID        string `json:"commentUUID"`
}

func Comment(w http.ResponseWriter, r *http.Request, tmpl *template.Template, postid string) {
	data := Data{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("session error is: ", err)
	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		utils.HandleError("cannot select user from session ", err)
		data.Cookie = sessionId

		data.UserInfo = user
		message := fmt.Sprintf("session id is: %v user info is: %v cookie data is: %v", sessionId, data.UserInfo, data.Cookie)
		utils.WriteMessageToLogFile(message)
		if r.Method == "POST" {
			err := dbmanagement.UpdateUserToken(user.UUID, 1)
			if err != nil {
				http.Redirect(w, r, "/error", http.StatusSeeOther)
				return
			}

			var formData CommentFormData
			jerr := json.NewDecoder(r.Body).Decode(&formData)
			if jerr != nil {
				utils.HandleError("Invalid request payload", err)
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			comment := formData.NewCommentContent
			deleteComment := formData.DeleteComment
			editComment := formData.EditCommentContent
			commentuuid := formData.CommentUUID

			if CheckInputs(comment) {
				userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
				utils.HandleError("Unable to get user with uuid in all posts", err)
				thisComment := dbmanagement.InsertComment(comment, postid, userFromUUID.UUID, 0, 0, time.Now())
				post, err := dbmanagement.SelectPostFromUUID(postid)
				if err != nil {
					PageErrors(w, r, tmpl, 500, "Internal Server Error")
					return
				}
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, postid, thisComment.UUID, user.UUID, 0, "")
			}
			if deleteComment != "" {
				dbmanagement.DeleteFromTableWithUUID("Comments", deleteComment)
			}
			if editComment != "" {
				dbmanagement.UpdateComment(commentuuid, editComment, postid, user.UUID, dbmanagement.SelectCommentFromUUID(commentuuid).Likes, dbmanagement.SelectCommentFromUUID(commentuuid).Dislikes, time.Now())
			}
		}
	}
}
