package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// PORT => Server port
const PORT = ":5000"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Message => Used to form messages
type Message struct {
	User string `json:"User"`
	Body string `json:"Body"`
}

var messages []Message

func storeNewMessage(message []byte) {
	var msg Message
	json.Unmarshal(message, &msg)
	messages = append(messages, msg)
}

func serveWSEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("serveWSEndpoint", err)
		return
	}
	conn.WriteJSON(messages)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WriteMessage", err)
			return
		}
		storeNewMessage(p)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", serveWSEndpoint)
	// Start the server
	fmt.Printf("Server running!\nhttp://localhost%s/\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
