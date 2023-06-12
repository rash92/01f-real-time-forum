package ws

import (
	"encoding/json"
	"log"
)

type Hub struct {
	clients           map[*Client]bool
	broadcast         chan []byte
	register          chan *Client
	unregister        chan *Client
	typingBroadcast   chan *Client
	clientsByUsername map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:         make(chan []byte),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[*Client]bool),
		typingBroadcast:   make(chan *Client),
		clientsByUsername: map[string]*Client{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Println("closing at hun run function - case unregiser")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					log.Println("closing at hun run function - case h.broadcast")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// getClientByUsername retrieves the client object based on the username
func (h *Hub) GetClientByUsername(username string) *Client {
	for client := range h.clients {
		if client.User.Name == username {
			return client
		}
	}
	return nil
}

func (h *Hub) BroadcastTypingStatus() {
	// Continuously listen for clients whose typing status needs to be broadcasted
	for {
		client := <-h.typingBroadcast

		// Iterate over all connected clients in the hub
		for c := range h.clients {
			// Skip broadcasting to the client who triggered the typing status
			if c != client {
				// Create a message containing the typing status information
				message := WriteMessage{
					Type: "typing",
					Data: map[string]interface{}{
						"username": client.User.Name,
						"isTyping": client.typingStatus,
					},
				}
				jsonMessage, _ := json.Marshal(message)

				// Send the typing status message to the client
				c.send <- jsonMessage
			}
		}
	}
}
