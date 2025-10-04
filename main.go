package main

import (
	"log"
	"net/http"

	"mtwebviz/server"
	"mtwebviz/touch"
)

func main() {
	// Set up touch event callback to broadcast to WebSocket clients
	touch.SetEventCallback(server.BroadcastEvent)

	// Start broadcaster goroutine
	go server.Broadcaster()

	// Start multitouch tracking
	go touch.Start()

	// Setup HTTP handlers
	http.HandleFunc("/ws", server.HandleWebSocket)
	http.HandleFunc("/", server.HandleFrontend)

	log.Println("Starting WebSocket server on :8080")
	log.Println("Open http://localhost:8080 in your browser")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
