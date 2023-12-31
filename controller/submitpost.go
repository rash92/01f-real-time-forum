package controller

import (
	"encoding/json"
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

type SubmitPostFormData struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	EditPost string   `json:"editpost"`
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

		// Parse the JSON body into FormData struct
		var formData SubmitPostFormData
		jerr := json.NewDecoder(r.Body).Decode(&formData)
		if jerr != nil {
			// Handle error
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		SubmissionHandler(w, r, user, formData, tmpl)
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
		utils.HandleError("cannot handle submitpost data ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
