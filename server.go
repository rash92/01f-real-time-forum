package main

import (
	"crypto/tls"
	"fmt"
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("static/*.html"))

}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	//mux.HandleFunc("/", protectPostGetRequests(IndexHandler))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "one-page.html", nil)
	})
	// mux.HandleFunc("/posts", protectGetRequests(IndexHandler))
	mux.HandleFunc("/categories/", CategoriesHandler)
	mux.HandleFunc("/posts/", PostsHandler)
	mux.HandleFunc("/ws", wsEndpoint)

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

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Check the headers for WebSocket upgrade
	if r.Header.Get("Upgrade") != "websocket" || !websocket.IsWebSocketUpgrade(r) {
		http.Error(w, "Invalid WebSocket request", http.StatusBadRequest)
		return
	}

	log.Println("WebSocket upgrade request received:", r.Header)

	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	log.Println("here")
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
