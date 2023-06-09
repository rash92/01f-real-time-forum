package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	pongWait = 10 * time.Second

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

	User         dbmanagement.User
	typing       chan bool
	typingStatus bool
	recipient    *Client
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
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() { // Same as POST
	defer func() {
		log.Println("closing at readpump")
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.HandleError("Unexpected Websocket Close", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var msg ReadMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			utils.HandleError("Error decoding JSON:", err)
			continue
		}

		switch msg.Type {
		case "recipientSelect":
			name, ok := msg.Info["name"].(string)
			if ok {
				userConnection, _ := dbmanagement.SelectUserFromName(name)
				recipientClient := c.hub.GetClientByUsername(name)
				c.recipient = recipientClient // This brings back hashed password, probably not necessary
				message := fmt.Sprintf("User Data: %v", userConnection)
				utils.WriteMessageToLogFile(message)
			}
		case "private":
			recipient, ok1 := msg.Info["recipient"].(string)
			receiver, _ := dbmanagement.SelectUserFromName(recipient)
			text, ok2 := msg.Info["text"].(string)
			user := c.User.UUID
			if ok1 && ok2 {
				message := fmt.Sprintf("Private Message: %s %s %s", user, receiver.UUID, text)
				utils.WriteMessageToLogFile(message)
			}
		case "typing":
			isTyping, ok2 := msg.Info["isTyping"].(bool)
			user := c.User.UUID
			if ok2 {
				message := fmt.Sprintf("typing: %s %v", user, isTyping)
				utils.WriteMessageToLogFile(message)
				c.typing <- isTyping
				c.hub.typingBroadcast <- c
			}
		default:
			c.hub.broadcast <- message
			message := fmt.Sprintf("Recieved %s", message)
			utils.WriteMessageToLogFile(message)
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

	// Create a timer to track typing state
	typingTimer := time.NewTimer(0)
	typingTimer.Stop()

	defer func() {
		log.Println("closing at writepump")
		ticker.Stop()
		onlineUsersTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
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
		case isTyping := <-c.typing:
			c.typingStatus = isTyping

			// Reset the typing timer whenever isTyping changes
			if isTyping {
				typingTimer.Reset(2 * time.Second)
			} else {
				typingTimer.Stop()
			}

			// Broadcast typing status to other clients in the hub
			message := WriteMessage{
				Type: "typing",
				Data: map[string]interface{}{
					"username": c.User.Name,
					"isTyping": isTyping,
				},
			}
			jsonMessage, _ := json.Marshal(message)
			// Check if recipient is available and has a valid connection
			if c.recipient != nil && c.recipient.send != nil {
				c.recipient.send <- jsonMessage
			}
			// c.send <- jsonMessage
		case <-typingTimer.C: // Handle the typing timer expiration
			c.typingStatus = false
			// Broadcast typing status to other clients in the hub
			message := WriteMessage{
				Type: "typing",
				Data: map[string]interface{}{
					"username": c.User.Name,
					"isTyping": false,
				},
			}
			jsonMessage, _ := json.Marshal(message)
			// Check if recipient is available and has a valid connection
			if c.recipient != nil && c.recipient.send != nil {
				c.recipient.send <- jsonMessage
			}
		case <-onlineUsersTicker.C:
			onlineUsersData := OnlineUsersHandler()
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))

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
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
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
		case <-ticker.C:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
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
	message := fmt.Sprintf("client User from DB %v", clientUser)
	utils.WriteMessageToLogFile(message)
	if err != nil || clientUser.Name == "" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), User: clientUser, typingStatus: false}
	client.typing = make(chan bool)
	client.hub.register <- client

	// Store the client object in the clientsByUsername map
	hub.clientsByUsername[clientUser.Name] = client

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
	go hub.BroadcastTypingStatus()
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
