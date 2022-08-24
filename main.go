package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

/*
	Configurations for upgrader
	Converts the connection from regular
	http request to one that can be used
	with websocket
*/
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Setting up for accepting incomming connections..
func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Convertign our connection to websocker one using upgrader..
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, "WS Connection Error: ", err)
		return
	}

	// Closing the connection at the end..
	defer ws.Close()

	/*
		Enabling the server to recieve as many
		messages from the client
		Hence infinite for loop..
	*/
	for {
		/*
			First argument is the type of message
			Ignoring the type (Text, Binary, Ping-Pong)
			Ping-Pong message is the one that ensures
			that server and client are both responding to
			each other..
		*/
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			// handleDisconnect(ws)
			break
		}
		// Converting the bytes message to string format..
		msg := string(bytes)
		// handleIncommingMessage(ws, msg)
	}
}

/*
	Same implementation as the socketHandler
	But echos the clients request back.
*/
func echoServer(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, "WS Connection Error: ", err)
		return
	}
	defer ws.Close()
	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			break
		}
		msg := string(bytes)
	}
}

func main() {
	http.HandleFunc("/verify", handler)
	http.ListenAndServe("", nil)
}
