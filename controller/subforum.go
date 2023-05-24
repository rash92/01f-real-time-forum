package controller

import (
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"time"
)

type SubData struct {
	SubName    string
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
	TitleName  string
	TagsList   []dbmanagement.Tag
}

func SubForum(w http.ResponseWriter, r *http.Request, tmpl *template.Template, tag string) {
	data := SubData{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	if sessionId == "" {
		err := auth.CreateUserSession(w, r, dbmanagement.User{})
		if err != nil {
			utils.HandleError("Unable to create visitor session", err)
		} else {
			sessionId, _ = auth.GetSessionFromBrowser(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		data.Cookie = sessionId
		filterOrder := false
		data.UserInfo = user
		if r.Method == "POST" {
			err := dbmanagement.UpdateUserToken(user.UUID, 1)
			if err != nil {
				http.Redirect(w, r, "/error", http.StatusSeeOther)
				return
			}
			content := r.FormValue("post")
			like := r.FormValue("like")
			dislike := r.FormValue("dislike")
			filter := r.FormValue("filter")
			if filter == "oldest" {
				filterOrder = true
			}
			if content != "" {
				userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
				utils.HandleError("Unable get user with UUID in all Subforum function", err)
				dbmanagement.InsertPost("", content, userFromUUID.Name, 0, 0, time.Now(), "")
				if !ExistingTag(tag) {
					dbmanagement.InsertTag(tag)
				}
			}
			if like != "" {
				dbmanagement.AddReactionToPost(user.UUID, like, 1)
				post, err := dbmanagement.SelectPostFromUUID(like)
				if err != nil {
					PageErrors(w, r, tmpl, 500, "Internal Server Error")
					return
				}
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, like, "", user.UUID, 1, "")
			}
			if dislike != "" {
				dbmanagement.AddReactionToPost(user.UUID, dislike, -1)
				post, err := dbmanagement.SelectPostFromUUID(dislike)
				if err != nil {
					PageErrors(w, r, tmpl, 500, "Internal Server Error")
					return
				}
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, dislike, "", user.UUID, -1, "")
			}

			idToDelete := r.FormValue("deletepost")
			if idToDelete != "" {
				dbmanagement.DeleteFromTableWithUUID("Posts", idToDelete)
			}
			http.Redirect(w, r, "/categories/"+tag, http.StatusFound)
		}

		utils.HandleError("Unable to select user using sessionid", err)
		posts, err := dbmanagement.SelectAllPostsFromTag(tag)
		if err != nil {
			PageErrors(w, r, tmpl, 500, "Internal Server Error")
			return
		}
		if !filterOrder {
			for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}

		data.SubName = tag
		data.Cookie = sessionId
		user.Notifications = dbmanagement.SelectAllNotificationsFromUser(user.UUID)
		data.UserInfo = user
		data.ListOfData = posts
		data.TagsList = dbmanagement.SelectAllTags()
		tmpl.ExecuteTemplate(w, "subforum.html", data)
	}
}
