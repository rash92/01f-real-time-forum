package auth

import (
	"context"
	"errors"
	"forum/utils"
	"html/template"
	"net/http"
)

// Google Oauth
func GoogleLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	LoggedInStatus(w, r, tmpl, 0)
	googleConfig := GoogleSetupConfig()
	url := googleConfig.AuthCodeURL(Randomstate)
	// redirect to google login page
	http.Redirect(w, r, url, http.StatusFound)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// state
	state := r.FormValue("state")
	if state != Randomstate {
		utils.WriteMessageToLogFile("Google auth state error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// code
	code := r.FormValue("code")
	// configuration
	googleConfig := GoogleSetupConfig()

	// exchange code for token
	token, err := googleConfig.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	// use google api to get user info
	resp, err := http.Get(GoogleAuthURL + token.AccessToken)
	utils.HandleError("Failed to fetch user data from google:", err)

	// parse response
	value := ParseOauthResponse(resp)

	if value["email"] == nil || value["given_name"] == nil {
		http.Redirect(w, r, "/oautherror", http.StatusFound)
		utils.HandleError("unable to get ", errors.New("unable to get email address or name"))
		return
	}

	account := OauthAccount{
		Name:  utils.AssertString(value["given_name"]),
		Email: utils.AssertString(value["email"]),
	}
	// login and create session for user
	LoginUserWithOauth(w, r, tmpl, account)
}
