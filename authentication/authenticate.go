package auth

import (
	"fmt"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"strings"
)

type Data struct {
	ListOfData    []dbmanagement.Post
	Cookie        string
	UserInfo      dbmanagement.User
	TitleName     string
	IsCorrect     bool
	IsLoggedIn    bool
	RegisterError string
	TagsList      []dbmanagement.Tag
}

type OauthAccount struct {
	Name, Email string
}

// Displays the log in page.
func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	sessionId, err := GetSessionFromBrowser(w, r)
	utils.HandleError("Unable to get session id from browser in login:", err)
	_, err = dbmanagement.SelectUserFromSession(sessionId)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		LoggedInStatus(w, r, tmpl, 0)
		data := Data{}
		data.TitleName = "Login"
		data.IsCorrect = true
		data.TagsList = dbmanagement.SelectAllTags()
		tmpl.ExecuteTemplate(w, "login.html", data)
	}
}

/*
Authenticate user with credentials - If the username and password match an entry in the database then the user is redirected to the forum page,
otherwise the user stays on the log in page. Session Cookie is also set here.
*/
func Authenticate(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	userName := r.FormValue("user_name")
	password := r.FormValue("password")

	user, err := dbmanagement.SelectUserFromName(userName)
	utils.HandleError("Unable to get user error:", err)

	if CompareHash(user.Password, password) && user.IsLoggedIn == 0 {
		err := CreateUserSession(w, r, user)
		utils.HandleError("Failed to create session in authenticate", err)
		// user.LimitTokens
		LimitRequests(w, r, user)
		err = dbmanagement.UpdateUserToken(user.UUID, 1)
		utils.HandleError("Unable to update users token", err)
		dbmanagement.UpdateUserLoggedInStatus(user.UUID, 1)
		utils.WriteMessageToLogFile(user.IsLoggedIn)
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	} else {
		if user.IsLoggedIn != 0 {
			utils.WriteMessageToLogFile("Already Logged In!")
			data := Data{}
			data.TitleName = "Login"
			data.IsCorrect = true
			data.IsLoggedIn = true
			data.TagsList = dbmanagement.SelectAllTags()
			tmpl.ExecuteTemplate(w, "login.html", data)
		} else {
			utils.WriteMessageToLogFile("Incorrect Password!")
			data := Data{}
			data.TitleName = "Login"
			data.IsCorrect = false
			data.TagsList = dbmanagement.SelectAllTags()
			tmpl.ExecuteTemplate(w, "login.html", data)
		}
	}
}

// Logs user out
func Logout(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	utils.WriteMessageToLogFile("Logging out")
	cookie, err := r.Cookie("session")
	utils.HandleError("Failed to get cookie", err)
	session := cookie.Value
	user, _ := dbmanagement.SelectUserFromSession(session)
	dbmanagement.UpdateUserLoggedInStatus(user.UUID, 0)
	utils.WriteMessageToLogFile(user.IsLoggedIn)
	if err != http.ErrNoCookie {
		err := dbmanagement.DeleteSessionByUUID(session)
		utils.HandleError("Failed to get cookie", err)
	}
	clearcookie := http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &clearcookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Displays the register page
func Register(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	LoggedInStatus(w, r, tmpl, 0)
	data := Data{}
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "register.html", data)
}

// Registers a user with given details then redirects to log in page.  Password is hashed here.
func RegisterAcount(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		email := r.FormValue("email")
		password := HashPassword(r.FormValue("password"))
		_, err := dbmanagement.InsertUser(userName, email, password, "user", 0)
		data := Data{}
		if err != nil {
			data.RegisterError = strings.Split(err.Error(), ".")[1]
			data.TagsList = dbmanagement.SelectAllTags()
			tmpl.ExecuteTemplate(w, "register.html", data)
		}
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Checks whether the user is logged in or not for displaying certain pages
func LoggedInStatus(w http.ResponseWriter, r *http.Request, tmpl *template.Template, desiredStatus int) {
	cookie, err := r.Cookie("session")
	message := fmt.Sprint("Current Cookie: ", cookie)
	utils.WriteMessageToLogFile(message)
	utils.HandleError("Failed to get cookie", err)
	session := cookie.Value
	user, _ := dbmanagement.SelectUserFromSession(session)
	if user.IsLoggedIn != desiredStatus {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
