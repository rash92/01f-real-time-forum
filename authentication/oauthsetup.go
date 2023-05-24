package auth

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

var GoogleAuthURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
var GithubAuthURL = "https://api.github.com/user"
var FBAuthURL = "https://graph.facebook.com/v13.0/me?fields=id,name,email,picture&access_token&access_token="
var stateRandomiser = randomString(12)
var Randomstate = stateRandomiser

var GoogleClientID = "518245388319-qj8tmfb4ue1hfodjophsp9bprfe0om66.apps.googleusercontent.com"
var GoogleClientSecret = "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"

var GithubClientID = "7a5de3b35a748e59ec9b"
var GithubClientSecret = "5256fcbe219c58b572029d7443dc3a996c7d378a"

var FacebookClientID = "1481843748971935"
var FacebookClientSecret = "43dbe0fc14822f0b89eab4262a6e6978"

func GoogleSetupConfig() *oauth2.Config {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf := &oauth2.Config{
		ClientID:     GoogleClientID,
		ClientSecret: GoogleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://localhost:8080/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	return conf
}

func GithubSetupConfig() *oauth2.Config {
	// Oauth configuration for github
	conf := &oauth2.Config{
		ClientID:     GithubClientID,
		ClientSecret: GithubClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  "https://localhost:8080/github/callback",
		Scopes: []string{
			"user",
			"user:email",
		},
	}
	return conf
}

func FacebookSetupConfig() *oauth2.Config {
	// Oauth configuration for facebook
	conf := &oauth2.Config{
		ClientID:     FacebookClientID,
		ClientSecret: FacebookClientSecret,
		Endpoint:     facebook.Endpoint,
		RedirectURL:  "https://localhost:8080/facebook/callback",
		Scopes: []string{
			"email",
			"public_profile",
		},
	}
	return conf
}
