package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

// PORT => Server port
const PORT = ":5000"

// BroadcastRoom => Channel for broadcasting to all clients
const BroadcastRoom = "broadcast"

type user struct {
	Name  string `json:"Name"`
	Color string `json:"Color"`
}

// Message => Used to form messages
type Message struct {
	User user   `json:"User"`
	Body string `json:"Body"`
}

var messages []Message

var users = make(map[string]user)

func sendActiveUsers(server *socketio.Server) {
	var payload []user

	for _, v := range users {
		payload = append(payload, v)
	}

	stringified, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("sendActiveUsers:", err)
	}

	server.BroadcastToRoom("/", BroadcastRoom, "users", string(stringified))
}

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
		payload, _ := json.Marshal(messages)
		s.SetContext("")
		fmt.Println("conencted:", s.ID())
		s.Emit("messages", string(payload))
		return nil
	})

	server.OnEvent("/", "new_user", func(s socketio.Conn, msg string) {
		var newUser user
		json.Unmarshal([]byte(msg), &newUser)
		users[s.ID()] = newUser
		s.Join(BroadcastRoom)
		server.BroadcastToRoom("/", BroadcastRoom, "admin", newUser.Name+" joined the chat room.")
		sendActiveUsers(server)
	})

	server.OnEvent("/", "message_sent", func(s socketio.Conn, msg string) {
		unmarshaled := Message{
			User: users[s.ID()],
			Body: msg,
		}

		message, _ := json.Marshal(unmarshaled)
		storeNewMessage(message)

		server.BroadcastToRoom("/", BroadcastRoom, "new_message", string(message))
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		str := e.Error()

		if strings.Contains(str, "going away") {
			server.BroadcastToRoom("/", BroadcastRoom, "admin", users[s.ID()].Name+" has left chat room.")
			delete(users, s.ID())
			sendActiveUsers(server)
		}
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
