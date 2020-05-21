package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"

	"github.com/gorilla/websocket"
)

// PORT => Server port
const PORT = ":5000"

// BROADCAST_ROOM => Channel for broadcasting to all clients
const BROADCAST_ROOM = "broadcast"

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
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal("socketio.NewServer", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("conencted:", s.ID())
		s.Join(BROADCAST_ROOM)
		server.BroadcastToRoom("/", BROADCAST_ROOM, "new_user", "connected: "+s.ID())
		return nil
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
