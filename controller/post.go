package controller

import (
	"encoding/json"
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"strings"
)

type PostData struct {
	Post             dbmanagement.Post
	Comments         []dbmanagement.Comment
	NumOfComments    int
	Cookie           string
	UserInfo         dbmanagement.User
	TitleName        string
	HasNotifications bool
	Notifications    []dbmanagement.Notification
	TagsList         []dbmanagement.Tag
}

func CheckInputs(str string) bool {
	spl := strings.Fields(str)
	return len(spl) > 0
}

func Post(w http.ResponseWriter, r *http.Request, tmpl *template.Template, postid string) {
	data := Data{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("session error is: ", err)
	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
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

			idToReport := r.FormValue("reportpost")

			if idToReport != "" {
				dbmanagement.CreateAdminRequest(user.UUID, user.Name, idToReport, "", "", "this post has been reported by a moderator")
				message := fmt.Sprintf("a post has been reported with id: %v", idToReport)
				utils.WriteMessageToLogFile(message)

			}

		}

		utils.HandleError("Unable to get user", err)
		post, err := dbmanagement.SelectPostFromUUID(postid)
		if err != nil {
			PageErrors(w, r, tmpl, 500, "Internal Server Error")
			return
		}
		comments := dbmanagement.SelectAllCommentsFromPost(postid)
		for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
			comments[i], comments[j] = comments[j], comments[i]
		}

		data := PostData{}
		data.Cookie = sessionId
		user.Notifications = dbmanagement.SelectAllNotificationsFromUser(user.UUID)
		data.UserInfo = user
		data.Post = post
		data.Comments = append(data.Comments, comments...)
		data.NumOfComments = len(comments)
		data.TagsList = dbmanagement.SelectAllTags()
		message = fmt.Sprintf("%v, %v", data.UserInfo.Name, data.Post.OwnerId)
		utils.WriteMessageToLogFile(message)
		//tmpl.ExecuteTemplate(w, "post.html", data)

		// Convert the struct to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			utils.HandleError("cannot marshal post data ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
