package main

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Map containing information regarding the socket and username..
var connNameMap = make(map[*websocket.Conn]string)

// Map to check that people dont register with same user name..
var nameConnMap = make(map[string]*websocket.Conn)

func handleIncommingMessage(sender *websocket.Conn, msg string) {
	/*
		If the user is found in connection to name
		map then he/she can send message..
	*/
	if _, ok := connNameMap[sender]; ok {
		sendChatMessage(sender, msg)
		return
	}

	// Registering a new User..
	username := strings.TrimSpace(msg)
	if username == "" || username == "server" {
		sender.WriteJSON(newError("Invalid username."))
	}

	// Checking for the user in the map if it already exists
	if _, ok := nameConnMap[username]; ok {
		// Write a message that user already exixts..
		sender.WriteJSON(newError("User already exists."))
	}

	// sendUserList(sender)
	connNameMap[sender] = username
	nameConnMap[username] = sender

	m := newMessage(msgJoin, "Server", username)
	dispatch(m)
}

func disconnectionHandler(sender *websocket.Conn) {
	username, ok := connNameMap[sender]
	if !ok {
		return
	}

	m := newMessage(msgLeave, "Server", username)
	dispatch(m)
	/*
		Deleting the values from both maps
		So they can further be used by
		other users..
	*/
	delete(nameConnMap, username)
	delete(connNameMap, sender)
}

/*
	There are 5 type of message types
	1. Chat Message (Normal chat message)
	2. Join Message (When a new user joins)
	3. Leave Message (When a user leaves)
	4. Error Message (When there's an error performing some task)
	5. Users Message (When we are trying to send list of available users)
*/
type messageType string

const (
	msgChat  messageType = "Message"
	msgJoin  messageType = "Join"
	msgLeave messageType = "Leave"
	msgError messageType = "Error"
	msgUsers messageType = "Users"
)

/*
	Encoding our message with JSON
	Makes it easy to work on both the
	client side and server side..
*/
type message struct {
	Type    messageType `json:"type"`
	Sender  string      `json:"sender"`
	Context interface{} `json:"context"`
	Date    time.Time   `json:"date"`
	Success bool        `json:"success"`
}

/*
	Returning error json message in-case of error
	Setting JSON paramteres
*/
func newError(content string) message {
	return message{
		Type:    msgError,
		Sender:  "",
		Context: content,
		Date:    time.Now().UTC(),
		Success: false,
	}
}

/*
	Returning json message in-case of no error
	Setting JSON paramteres
*/
func newMessage(msgType messageType, sender string, content string) message {
	return message{
		Type:    msgType,
		Sender:  sender,
		Context: content,
		Date:    time.Now().UTC(),
		Success: true,
	}
}

// Sending messages to all the users..
func dispatch(m message) {
	for client := range connNameMap {
		client.WriteJSON(m)
	}
}

/*
	1. Composing a message
	2. Sending message
*/
func sendChatMessage(sender *websocket.Conn, msg string) {
	// Generating the message..
	m := newMessage(msgChat, connNameMap[sender], msg)
	// Dispatching message..
	dispatch(m)
}
