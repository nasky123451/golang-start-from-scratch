package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	conn *websocket.Conn // WebSocket connection
	send chan []byte     // Channel for sending messages
	id   int             // Unique client ID
}

// Hub manages all connected clients and broadcasts messages
type Hub struct {
	clients    map[int]*Client // Map to store all connected clients
	register   chan *Client    // Channel for registering new clients
	unregister chan *Client    // Channel for unregistering clients
	broadcast  chan []byte     // Channel for broadcasting messages
	mu         sync.Mutex      // Mutex to protect shared resources
}

// Create a new hub
func newHub() *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256), // Use buffered channel
	}
}

// Main hub loop for handling client registration, unregistration, and broadcasting
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.id] = client
			h.mu.Unlock()
			log.Printf("Client %d connected", client.id)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.id]; ok {
				close(client.send)
				delete(h.clients, client.id)
				log.Printf("Client %d disconnected", client.id)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for _, client := range h.clients {
				select {
				case client.send <- message: // Send to client
				default:
					close(client.send)           // Close channel if full
					delete(h.clients, client.id) // Remove client from active list
				}
			}
			h.mu.Unlock()
		}
	}
}

// Handle reading and writing for each client
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c // Unregister the client
		c.conn.Close()      // Close the WebSocket connection
	}()

	for {
		_, message, err := c.conn.ReadMessage() // Read message from the WebSocket connection
		if err != nil {
			log.Println("read error:", err) // Log read error
			break
		}
		hub.broadcast <- message // Send the message to the broadcast channel
	}
}

// Handle writing to the client
func (c *Client) writePump() {
	ticker := time.NewTicker(time.Second * 10) // Create a ticker to send ping messages
	defer func() {
		ticker.Stop()  // Stop the ticker on exit
		c.conn.Close() // Ensure the connection is closed on exit
	}()

	for {
		select {
		case message, ok := <-c.send: // Wait for a message to send
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) // If the send channel is closed, send a close message
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message) // Send the message
			if err != nil {
				log.Printf("Client %d write error: %v", c.id, err) // Log write error
				return
			}
		case <-ticker.C: // Every tick, send a ping message to keep the connection alive
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("Client %d ping error: %v", c.id, err) // Log ping error
				return
			}
		}
	}
}

// WebSocket upgrader with no origin checking
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handle new WebSocket connections
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request, id int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256), id: id}
	hub.register <- client

	go client.writePump()
	go client.readPump(hub)
}

// Monitor system resource usage periodically
func monitorResources() {
	ticker := time.NewTicker(time.Second * 10) // Report every 10 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Get memory statistics
		var m runtime.MemStats
		runtime.ReadMemStats(&m) // Read the memory statistics

		log.Printf("Memory Usage: Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v",
			m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC) // Log memory usage
		log.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
	}
}

// Convert bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func main() {
	// Define a flag to enable/disable resource monitoring
	enableMonitoring := flag.Bool("monitor", false, "Enable resource monitoring")

	// Parse command-line flags
	flag.Parse()

	hub := newHub()
	go hub.run()

	// If the monitoring flag is set, start monitoring resources
	if *enableMonitoring {
		log.Println("Resource monitoring enabled")
		go monitorResources()
	} else {
		log.Println("Resource monitoring disabled")
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientID := rand.Int() // Generate a random client ID
		serveWs(hub, w, r, clientID)
	})

	log.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
