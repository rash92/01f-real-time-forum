package controller

import (
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type SubmitData struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
	TitleName  string
	IsCorrect  bool
	IsEdit     bool
	EditPost   dbmanagement.Post
	Tags       string
	TagsList   []dbmanagement.Tag
}

func SubmitPost(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	auth.LoggedInStatus(w, r, tmpl, 1)

	data := SubmitData{}
	user := dbmanagement.User{}
	tags := []dbmanagement.Tag{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("Unable to get session from browser in SubmitPost function", err)
	user, err = dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("Unable to select user with sessionID in SubmitPost function", err)
	if r.Method == "POST" {
		err := dbmanagement.UpdateUserToken(user.UUID, 1)
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}

		idToEdit := r.FormValue("editpost")
		if idToEdit != "" {
			data.IsEdit = true
			data.EditPost, err = dbmanagement.SelectPostFromUUID(idToEdit)
			if err != nil {
				PageErrors(w, r, tmpl, 500, "Internal Server Error")
				return
			}
			tags = dbmanagement.SelectAllTagsFromPost(data.EditPost.UUID)
		}
	}
	data.TitleName = "Submit to Forum"
	data.Cookie = sessionId
	data.UserInfo = user
	tagsAsString := ""
	for _, v := range tags {
		tagsAsString += v.TagName
		tagsAsString += " "
	}
	data.Tags = tagsAsString
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "submitpost.html", data)
}
