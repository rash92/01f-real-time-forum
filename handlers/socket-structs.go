package handlers

import "github.com/gorilla/websocket"

type SocketEventStruct struct {
	EventName    string
	EventPayload interface{}
}

type JoinDisconnectPayload struct {
	UserID string
	Users  []UserStruct
}

type UserStruct struct {
	Username string
	UserID   string
}

type Client struct {
	hub                 *Hub
	webSocketConnection *websocket.Conn
	send                chan SocketEventStruct
	username            string
	userID              string
}
