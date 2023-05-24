package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"strings"
	"time"
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
			comment := r.FormValue("comment")
			like := r.FormValue("like")
			dislike := r.FormValue("dislike")
			commentlike := r.FormValue("commentlike")
			commentdislike := r.FormValue("commentdislike")
			idToDelete := r.FormValue("deletepost")
			idToReport := r.FormValue("reportpost")
			deleteComment := r.FormValue("deletecomment")
			editComment := r.FormValue("editcomment")
			commentuuid := r.FormValue("commentuuid")
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
			if commentlike != "" {
				dbmanagement.AddReactionToComment(user.UUID, commentlike, 1)
				comment := dbmanagement.SelectCommentFromUUID(commentlike)
				receiverId, _ := dbmanagement.SelectUserFromName(comment.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, "", commentlike, user.UUID, 1, "")
			}
			if commentdislike != "" {
				dbmanagement.AddReactionToComment(user.UUID, commentdislike, -1)
				comment := dbmanagement.SelectCommentFromUUID(commentdislike)
				receiverId, _ := dbmanagement.SelectUserFromName(comment.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, "", commentdislike, user.UUID, -1, "")
			}

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
			if idToReport != "" {
				dbmanagement.CreateAdminRequest(user.UUID, user.Name, idToReport, "", "", "this post has been reported by a moderator")
				fmt.Println("a post has been reported with id: ", idToReport)
			}
			if deleteComment != "" {
				dbmanagement.DeleteFromTableWithUUID("Comments", deleteComment)
			}
			if editComment != "" {
				dbmanagement.UpdateComment(commentuuid, editComment, postid, user.UUID, dbmanagement.SelectCommentFromUUID(editComment).Likes, dbmanagement.SelectCommentFromUUID(editComment).Dislikes, time.Now())
			}
			http.Redirect(w, r, "/posts/"+postid, http.StatusFound)
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
		tmpl.ExecuteTemplate(w, "post.html", data)
	}
}
