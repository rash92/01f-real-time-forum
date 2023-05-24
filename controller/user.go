package controller

import (
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type UserData struct {
	UserPosts         []dbmanagement.Post
	LikedUserPosts    []dbmanagement.Post
	DislikedUserPosts []dbmanagement.Post
	UserComments      []dbmanagement.Comment
	UserInfo          dbmanagement.User
	TitleName         string
	Cookie            string
	TagsList          []dbmanagement.Tag
}

func User(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	auth.LoggedInStatus(w, r, tmpl, 1)
	data := UserData{}
	SessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("Unable to find user session id", err)

	data.UserInfo, err = dbmanagement.SelectUserFromSession(SessionId)
	utils.HandleError("Unable to find user session id", err)

	if SessionId == "" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	data.UserInfo, err = dbmanagement.SelectUserFromSession(SessionId)
	data.UserInfo.Notifications = dbmanagement.SelectAllNotificationsFromUser(data.UserInfo.UUID)
	utils.HandleError("Could not get user session in user", err)
	data.UserPosts, err = dbmanagement.SelectAllPostsFromUser(data.UserInfo.Name)
	if err != nil {
		PageErrors(w, r, tmpl, 500, "Internal Server Error")
		return
	}
	data.UserComments = dbmanagement.SelectAllCommentsFromUser(data.UserInfo.Name)
	utils.WriteMessageToLogFile(data.UserComments)
	data.TitleName = "Welcome"

	if r.Method == "POST" {
		err := dbmanagement.UpdateUserToken(data.UserInfo.UUID, 1)
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
		postIdToDelete := r.FormValue("deletepost")
		if postIdToDelete != "" {
			dbmanagement.DeletePostWithUUID(postIdToDelete)
		}
		commentIdToDelete := r.FormValue("deletecomment")
		if commentIdToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("Comments", commentIdToDelete)
		}
		userIdToRequestModerator := r.FormValue("request to become moderator")
		if userIdToRequestModerator != "" {
			newrequest := dbmanagement.CreateAdminRequest(userIdToRequestModerator, data.UserInfo.Name, "", "", "", "this user is asking to become a moderator")
			utils.WriteMessageToLogFile("new request description is: " + newrequest.Description)
		}
		http.Redirect(w, r, "/user", http.StatusFound)

	}

	data.UserInfo.Notifications = dbmanagement.SelectAllNotificationsFromUser(data.UserInfo.UUID)
	utils.HandleError("Could not get user session in user", err)
	data.UserPosts, err = dbmanagement.SelectAllPostsFromUser(data.UserInfo.Name)
	if err != nil {
		PageErrors(w, r, tmpl, 500, "Internal Server Error")
		return
	}
	data.LikedUserPosts, err = dbmanagement.SelectAllLikedPostsFromUser(data.UserInfo)
	if err != nil {
		PageErrors(w, r, tmpl, 500, "Internal Server Error")
		return
	}
	data.DislikedUserPosts, err = dbmanagement.SelectAllDislikedPostsFromUser(data.UserInfo)
	if err != nil {
		PageErrors(w, r, tmpl, 500, "Internal Server Error")
		return
	}
	data.UserComments = dbmanagement.SelectAllCommentsFromUser(data.UserInfo.UUID)
	data.TitleName = data.UserInfo.Name
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "user.html", data)
}
