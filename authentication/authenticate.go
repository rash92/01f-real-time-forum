package auth

import (
	"encoding/json"
	"fmt"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"regexp"
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

type RegisterAccountFormData struct {
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Age       int    `json:"age,string"`
}

type AuthenticateFormData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type OauthAccount struct {
	Name, Email string
}

// Displays the log in page.
func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	sessionId, err := GetSessionFromBrowser(w, r)
	fmt.Println("SessionID:", sessionId)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		utils.HandleError("Unable to get session id from browser in login:", err)
	}
	_, err = dbmanagement.SelectUserFromSession(sessionId)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		LoggedInStatus(w, r, tmpl, 0)
		data := Data{}
		data.TitleName = "Login"
		data.IsCorrect = true
		data.TagsList = dbmanagement.SelectAllTags()
		// tmpl.ExecuteTemplate(w, "login.html", data)

		// Convert the struct to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			utils.HandleError("cannot marshal loggedinstatus data", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

/*
Authenticate user with credentials - If the username and password match an entry in the database then the user is redirected to the forum page,
otherwise the user stays on the log in page. Session Cookie is also set here.
*/
func Authenticate(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// Parse the JSON body into FormData struct
	var formData AuthenticateFormData
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		utils.HandleError("cannot marshal authenticate data", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userName := formData.UserName
	password := formData.Password
	var user dbmanagement.User
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, userName)
	if err != nil {
		utils.HandleError("Regex failed for username or email", err)
	}
	if matched {
		user, err = dbmanagement.SelectUserFromEmail(userName)
		utils.HandleError("Unable to get user from email error:", err)
	} else {
		user, err = dbmanagement.SelectUserFromName(userName)
		utils.HandleError("Unable to get user from name error, trying by email:", err)
	}
	if CompareHash(user.Password, password) && user.IsLoggedIn == 0 {
		err := CreateUserSession(w, r, user)
		utils.HandleError("Failed to create session in authenticate", err)
		// user.LimitTokens
		LimitRequests(w, r, user)
		err = dbmanagement.UpdateUserToken(user.UUID, 1)
		utils.HandleError("Unable to update users token", err)
		dbmanagement.UpdateUserLoggedInStatus(user.UUID, 1)
		utils.WriteMessageToLogFile(user.IsLoggedIn)
		data := Data{}
		data.IsCorrect = true
		data.TagsList = dbmanagement.SelectAllTags()

		// Convert the struct to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			utils.HandleError("cannot marshal compare hash data", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		if user.IsLoggedIn != 0 {
			utils.WriteMessageToLogFile("Already Logged In!")
			data := Data{}
			data.TitleName = "Login"
			data.IsCorrect = true
			data.IsLoggedIn = true
			data.TagsList = dbmanagement.SelectAllTags()
			// tmpl.ExecuteTemplate(w, "login.html", data)

			// Convert the struct to JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				utils.HandleError("cannot marshal login data", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		} else {
			utils.WriteMessageToLogFile("Incorrect Password!")
			data := Data{}
			data.TitleName = "Login"
			data.IsCorrect = false
			data.TagsList = dbmanagement.SelectAllTags()
			// tmpl.ExecuteTemplate(w, "login.html", data)

			// Convert the struct to JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				utils.HandleError("cannot marshal tag data", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
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
	// tmpl.ExecuteTemplate(w, "register.html", data)

	// Convert the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		utils.HandleError("cannot marshal registration data", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// Registers a user with given details then redirects to log in page.  Password is hashed here.
func RegisterAcount(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// Parse the JSON body into FormData struct
	var formData RegisterAccountFormData

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		// Handle error
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	data := Data{}
	if r.Method == "POST" {
		password := HashPassword(formData.Password)
		_, err := dbmanagement.InsertUser(formData.UserName, formData.Email, password, "user", 0, formData.FirstName, formData.LastName, formData.Gender, formData.Age)
		if err != nil {
			data.RegisterError = strings.Split(err.Error(), ".")[1]
			data.TagsList = dbmanagement.SelectAllTags()
			// tmpl.ExecuteTemplate(w, "register.html", data)
		}
	}
	// Convert the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		utils.HandleError("cannot marshal tag data", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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
