package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

// PORT => Server port
const PORT = ":5000"

// BroadcastRoom => Channel for broadcasting to all clients
const BroadcastRoom = "broadcast"

// Message => Used to form messages
type Message struct {
	User string `json:"User"`
	Body string `json:"Body"`
}

type user struct {
	Name  string
	Color string
}

var messages []Message

var users = make(map[string]user)

func storeNewMessage(message []byte) {
	var msg Message
	json.Unmarshal(message, &msg)
	messages = append(messages, msg)
}

func main() {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal("socketio.NewServer", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("conencted:", s.ID())
		return nil
	})

	server.OnEvent("/", "new_user", func(s socketio.Conn, msg string) {
		var newUser user
		json.Unmarshal([]byte(msg), &newUser)
		users[s.ID()] = newUser
		s.Join(BroadcastRoom)
		server.BroadcastToRoom("/", BroadcastRoom, "user_joined", newUser.Name+" joined the chat room.")
	})

	server.OnEvent("/", "message_sent", func(s socketio.Conn, msg string) {
		server.BroadcastToRoom("/", BroadcastRoom, "new_message", users[s.ID()].Name+" joined the chat room.")
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	go func() {
		log.Fatal(server.Serve())
	}()

	defer server.Close()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.Handle("/socket.io/", server)

	// Start the server
	fmt.Printf("Server running!\nhttp://localhost%s/\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
