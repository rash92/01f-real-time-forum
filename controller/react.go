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

type ReactPostFormData struct {
	IsComment bool   `json:"isComment"`
	ID        string `json:"id"`
	Like      bool   `json:"like"`
	Dislike   bool   `json:"dislike"`
}

func ReactPost(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	fmt.Println("reacting")
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	if sessionId == "" {
		err := auth.CreateUserSession(w, r, dbmanagement.User{})
		if err != nil {
			utils.HandleError("Unable to create visitor session", err)
		} else {
			sessionId, _ = auth.GetSessionFromBrowser(w, r)

		}
	}

	if err == nil && r.Method == "POST" {
		user := dbmanagement.User{}
		user, err = dbmanagement.SelectUserFromSession(sessionId)

		// Parse the JSON body into FormData struct
		var formData ReactPostFormData
		err := json.NewDecoder(r.Body).Decode(&formData)
		if err != nil {
			// Handle error
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		isComment := formData.IsComment
		id := formData.ID
		like := formData.Like
		dislike := formData.Dislike

		fmt.Println(formData)

		if like {
			if isComment {
				dbmanagement.AddReactionToComment(user.UUID, id, 1)
				comment := dbmanagement.SelectCommentFromUUID(id)
				receiverId, _ := dbmanagement.SelectUserFromName(comment.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, "", id, user.UUID, 1, "")
			} else {
				dbmanagement.AddReactionToPost(user.UUID, id, 1)
				post, err := dbmanagement.SelectPostFromUUID(id)
				if err != nil {
					PageErrors(w, r, tmpl, 500, "Internal Server Error")
					return
				}
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, id, "", user.UUID, 1, "")
			}

		}
		if dislike {
			if isComment {
				dbmanagement.AddReactionToComment(user.UUID, id, -1)
				comment := dbmanagement.SelectCommentFromUUID(id)
				receiverId, _ := dbmanagement.SelectUserFromName(comment.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, "", id, user.UUID, -1, "")
			} else {
				dbmanagement.AddReactionToPost(user.UUID, id, -1)
				post, err := dbmanagement.SelectPostFromUUID(id)
				if err != nil {
					PageErrors(w, r, tmpl, 500, "Internal Server Error")
				}
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, id, "", user.UUID, -1, "")
			}
		}
	}

}
