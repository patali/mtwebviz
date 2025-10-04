package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"mtwebviz/touch"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.RWMutex
	broadcast = make(chan touch.FrameEvent, 100)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
)

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Printf("Client connected. Total clients: %d\n", len(clients))

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		log.Printf("Client disconnected. Total clients: %d\n", len(clients))
	}()

	// Keep connection alive and handle client messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// Broadcaster sends events to all connected clients
func Broadcaster() {
	for event := range broadcast {
		clientsMu.RLock()
		for client := range clients {
			err := client.WriteJSON(event)
			if err != nil {
				log.Printf("Write error: %v", err)
				client.Close()
				clientsMu.RUnlock()
				clientsMu.Lock()
				delete(clients, client)
				clientsMu.Unlock()
				clientsMu.RLock()
			}
		}
		clientsMu.RUnlock()
	}
}

// BroadcastEvent sends a touch event to all connected clients
func BroadcastEvent(event touch.FrameEvent) {
	select {
	case broadcast <- event:
	default:
		// Channel full, drop frame
	}
}
