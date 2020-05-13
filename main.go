package main

import (
	"fmt"
	"log"
	"net/http"
)

// PORT => Server port
const PORT = ":5000"

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	// Handlers
	http.HandleFunc("/", serveHomePage)
	// Start the server
	fmt.Printf("Server running!\nhttp://localhost%s/\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
