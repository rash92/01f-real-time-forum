package main

import (
	"crypto/tls"
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("static/*.html"))

}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--reset" {
		dbmanagement.CreateDatabaseWithTables()
	}

	mux := http.NewServeMux()
	cert, _ := tls.LoadX509KeyPair("https/localhost.crt", "https/localhost.key")
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	path := "./static"
	fs := http.FileServer(http.Dir(path))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// handlers
	mux.HandleFunc("/", protectPostGetRequests(IndexHandler))
	// mux.HandleFunc("/posts", protectGetRequests(IndexHandler))
	mux.HandleFunc("/categories/", CategoriesHandler)
	mux.HandleFunc("/posts/", PostsHandler)

	// authentication handlers
	mux.HandleFunc("/login", protectGetRequests(LoginHandler))
	mux.HandleFunc("/authenticate", protectPostRequests(AuthenticateHandler))
	mux.HandleFunc("/logout", protectGetRequests(LogoutHandler))
	mux.HandleFunc("/register", protectGetRequests(RegisterHandler))
	mux.HandleFunc("/register_account", protectPostRequests(RegisterAccountHandler))

	// oauth handlers
	mux.HandleFunc("/google/login", protectPostRequests(GoogleLoginHandler))
	mux.HandleFunc("/google/callback", protectGetRequests(GoogleCallbackHandler))
	mux.HandleFunc("/github/login", protectPostRequests(GithubLoginHandler))
	mux.HandleFunc("/github/callback", protectGetRequests(GithubCallbackHandler))
	mux.HandleFunc("/facebook/login", protectPostRequests(FacebookLoginHandler))
	mux.HandleFunc("/facebook/callback", protectGetRequests(FacebookCallbackHandler))

	// forum handlers
	mux.HandleFunc("/forum", protectPostGetRequests(ForumHandler))
	mux.HandleFunc("/submitpost", protectPostGetRequests(SubmitPostHandler))
	mux.HandleFunc("/admin", protectPostGetRequests(AdminHandler))
	mux.HandleFunc("/user", protectPostGetRequests(UserHandler))
	mux.HandleFunc("/privacy_policy", protectGetRequests(PrivacyPolicyHandler))
	mux.HandleFunc("/error", protectGetRequests(ErrorHandler))
	mux.HandleFunc("/oautherror", protectGetRequests(OauthErrorHandler))

	dbmanagement.DeleteAllSessions()
	dbmanagement.ResetAllUserLoggedInStatus()
	dbmanagement.ResetAllTokens()
	// dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
