package ws

import (
	"bytes"
	"encoding/json"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	User dbmanagement.User
}

type ReadMessage struct {
	Type string                 `json:"type"`
	Info map[string]interface{} `json:"info"`
}

type WriteMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() { // Same as POST
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var msg ReadMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue
		}
		switch msg.Type {
		case "recipientSelect":
			name, ok := msg.Info["name"].(string)
			if ok {
				log.Printf("Name: %s", name)
				userConnection, _ := dbmanagement.SelectUserFromName(name) // This brings back hashed password, probably not necessary
				log.Printf("User Data: %v", userConnection)
			}
		case "private":
			recipient, ok1 := msg.Info["recipient"].(string)
			text, ok2 := msg.Info["text"].(string)
			if ok1 && ok2 {
				log.Printf("Private Message: %s %s", recipient, text)
			}
		default:
			c.hub.broadcast <- message
			log.Printf("Recieved %s", message)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	onlineCheckerTicker := time.NewTicker(1 * time.Second)
	onlineUsersTicker := time.NewTicker(1 * time.Second) // Update online users every 5 seconds
	defer func() {
		ticker.Stop()
		onlineUsersTicker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-onlineUsersTicker.C:
			onlineUsersData := OnlineUsersHandler()
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			message := WriteMessage{
				Type: "onlineUsers",
				Data: onlineUsersData,
			}
			jsonMessage, _ := json.Marshal(message)
			w.Write(jsonMessage)

			if err := w.Close(); err != nil {
				return
			}
		case <-onlineCheckerTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			message := WriteMessage{
				Type: "userInfo",
				Data: c.User,
			}
			jsonMessage, _ := json.Marshal(message)
			w.Write(jsonMessage)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// Get SessionId from browser and tie it to client
	SessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("Unable to find user session id", err)

	clientUser, err := dbmanagement.SelectUserFromSession(SessionId)
	utils.HandleError("Unable to find user session id", err)
	log.Printf("client User from DB %v", clientUser)
	if err != nil || clientUser.Name == "" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), User: clientUser}
	client.hub.register <- client

	// Initial Send of Client User Info
	userMessage := WriteMessage{
		Type: "userInfo",
		Data: client.User,
	}

	jClientUser, _ := json.Marshal(userMessage)
	client.send <- jClientUser

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

type BasicUserInfo struct {
	Name           string
	LoggedInStatus int
}

func OnlineUsersHandler() []BasicUserInfo {
	onlineUsers := dbmanagement.SelectAllUsers()
	userArr := []BasicUserInfo{}
	for _, user := range onlineUsers {
		userArr = append(userArr, BasicUserInfo{user.Name, user.IsLoggedIn})
	}
	// TO DO: SORTED ALPHABETICALLY WHEN NEW USER ELSE BY LAST CHATTED TO
	sort.Slice(userArr, func(i, j int) bool {
		return userArr[i].Name < userArr[j].Name
	})
	return userArr

}
