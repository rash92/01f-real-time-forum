package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Data struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
	TitleName  string
	IsCorrect  bool
	TagsList   []dbmanagement.Tag
}

/*
Executes the forum.html template that includes all posts in the database.  SessionID is used the determine which user is currently using the website.
Also handles inserting a new post that updates in realtime.
*/
func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	if sessionId == "" {
		err := auth.CreateUserSession(w, r, dbmanagement.User{})
		if err != nil {
			utils.HandleError("Unable to create visitor session", err)
		} else {
			sessionId, _ = auth.GetSessionFromBrowser(w, r)

		}
	}

	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		utils.HandleError("Unable to get user", err)
		data.Cookie = sessionId
		filterOrder := false
		data.UserInfo = user
		message := fmt.Sprint("session id is: ", sessionId, "user info is: ", data.UserInfo, "cookie data is: ", data.Cookie)
		utils.WriteMessageToLogFile(message)

		if r.Method == "POST" {
			err := dbmanagement.UpdateUserToken(user.UUID, 1)
			if err != nil {
				http.Redirect(w, r, "/error", http.StatusSeeOther)
			}

			filter := r.FormValue("filter")
			if filter == "oldest" {
				filterOrder = true
			} else {
				SubmissionHandler(w, r, user, tmpl)
				http.Redirect(w, r, "/", http.StatusFound)
			}
		}

		posts, err := dbmanagement.SelectAllPosts()
		if err != nil {
			PageErrors(w, r, tmpl, 500, "Internal Server Error")
			return
		}

		if !filterOrder {
			for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
		data.Cookie = sessionId
		user.Notifications = dbmanagement.SelectAllNotificationsFromUser(user.UUID)
		data.UserInfo = user
		data.TitleName = "Forum"
		data.TagsList = dbmanagement.SelectAllTags()
		data.ListOfData = append(data.ListOfData, posts...)
		tmpl.ExecuteTemplate(w, "forum.html", data)
	}
}

func ExistingTag(tag string) bool {
	allTags := dbmanagement.SelectAllTags()
	for _, v := range allTags {
		if tag == v.TagName {
			return true
		}
	}
	return false
}

func UploadHandler(w http.ResponseWriter, r *http.Request, file multipart.File, fileHeader *multipart.FileHeader) string {
	err := os.MkdirAll("./static/uploads", os.ModePerm)
	if err != nil {
		utils.HandleError("error creating file directory for uploads", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	destinationFile, err := os.Create(fmt.Sprintf("./static/uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		utils.HandleError("error creating file for image", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, file)

	if err != nil {
		utils.HandleError("error copying file to destination", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}

	utils.WriteMessageToLogFile("file uploaded successfully")
	fileName := destinationFile.Name()[1:]
	return fileName
}

// followed this: https://freshman.tech/file-upload-golang/
func SubmissionHandler(w http.ResponseWriter, r *http.Request, user dbmanagement.User, tmpl *template.Template) {
	// 20 megabytes
	idToDelete := r.FormValue("deletepost")
	if idToDelete != "" {
		dbmanagement.DeletePostWithUUID(idToDelete)
	}
	notificationToDelete := r.FormValue("delete notification")
	if notificationToDelete != "" {
		dbmanagement.DeleteFromTableWithUUID("Notifications", notificationToDelete)
	}

	like := r.FormValue("like")
	dislike := r.FormValue("dislike")

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
		}
		receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
		dbmanagement.AddNotification(receiverId.UUID, dislike, "", user.UUID, -1, "")
	}

	maxSize := 20 * 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize))
	err := r.ParseMultipartForm(int64(maxSize))
	if err != nil {
		// only actual post submissions have multipart enabled, deleting, likes, dislikes aren't mulipart but that's already handled above so can end function
		if err.Error() == "request Content-Type isn't multipart/form-data" {
			return
		}
		utils.HandleError("error parsing form for image, likely too big", err)
		return
	}

	file, fileHeader, err := r.FormFile("submission-image")
	fileName := ""
	if err != nil {
		// if you were trying to make a post without an image it will log this 'error' but still submit the text and tags
		utils.HandleError("error retrieving file from form", err)
	} else {
		utils.WriteMessageToLogFile("trying to retrieve file...")
		defer file.Close()
		fileName = UploadHandler(w, r, file, fileHeader)
	}

	title := r.FormValue("submission-title")
	content := r.FormValue("post")
	tags := r.Form["tags"]
	edit := r.FormValue("editpost")

	if edit != "" {
		if CheckInputs(content) && CheckInputs(title) {
			userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
			utils.HandleError("cant get user with uuid in all posts", err)
			getLikes, err := dbmanagement.SelectPostFromUUID(edit)
			if err != nil {
				PageErrors(w, r, tmpl, 500, "Internal Server Error")
			}
			getDislikes, err := dbmanagement.SelectPostFromUUID(edit)
			if err != nil {
				PageErrors(w, r, tmpl, 500, "Internal Server Error")
			}
			editedPost, err := dbmanagement.UpdatePost(edit, title, content, userFromUUID.Name, getLikes.Likes, getDislikes.Dislikes, time.Now(), fileName)
			if err != nil {
				PageErrors(w, r, tmpl, 500, "Internal Server Error")
			}
			dbmanagement.UpdateTaggedPost(edit)
			for _, tag := range tags {
				InputTags(tag, editedPost)
				InputTags(tag, editedPost)
			}

		}
	} else {
		if CheckInputs(content) && CheckInputs(title) {
			userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
			utils.HandleError("cant get user with uuid in all posts", err)
			post, err := dbmanagement.InsertPost(title, content, userFromUUID.Name, 0, 0, time.Now(), fileName)
			if err != nil {
				PageErrors(w, r, tmpl, 500, "Internal Server Error")
			}
			for _, tag := range tags {
				InputTags(tag, post)
				InputTags(tag, post)
			}
		}
	}
}

func InputTags(tags string, post dbmanagement.Post) {
	if CheckInputs(tags) {
		tagslice := strings.Fields(tags)
		for _, tagname := range tagslice {
			tagname = strings.ToLower(tagname)
			if !ExistingTag(tagname) {
				dbmanagement.InsertTag(tagname)
			}
			tag, err := dbmanagement.SelectTagFromName(tagname)
			utils.HandleError("unable to retrieve tag id", err)
			dbmanagement.InsertTaggedPost(tag.UUID, post.UUID)
		}
	}
}
