package controller

import (
	"encoding/json"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type EditPostIDData struct {
	PostID string `json:"editPostID"`
}

func EditPost(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
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

		// Parse the JSON body into FormData struct
		var editData EditPostIDData
		jerr := json.NewDecoder(r.Body).Decode(&editData)
		if jerr != nil {
			// Handle error
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		idToEdit := editData.PostID

		if idToEdit != "" {
			data.IsEdit = true
			data.EditPost, err = dbmanagement.SelectPostFromUUID(idToEdit)
			if err != nil {
				utils.HandleError("id to edit failed", err)
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
	//tmpl.ExecuteTemplate(w, "submitpost.html", data)

	// Convert the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Handle error
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
