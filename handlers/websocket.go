package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"gin-project/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for development
		// In production, you should check the origin properly
		return true
	},
}

// InitializeWebSocketManager initializes the Redis WebSocket connection manager
func InitializeWebSocketManager() {
	initConnectionManager()
}

// getLocalClient returns a local client by user ID (for Redis manager)
func getLocalClient(userID uint) *Client {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()
	return hub.clients[userID]
}

// removeLocalClient removes a local client by user ID (for Redis manager)
func removeLocalClient(userID uint) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	if client, exists := hub.clients[userID]; exists {
		delete(hub.clients, userID)
		close(client.Send)
	}
}

// Client represents a WebSocket client
type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients mapped by user ID
	clients map[uint]*Client

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast channel
	broadcast chan []byte

	// Mutex for thread safety
	mutex sync.RWMutex
}

// Global hub instance
var hub = &Hub{
	clients:    make(map[uint]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan []byte),
}

// Initialize and start the hub
func init() {
	go hub.run()
}

// Hub run method
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.UserID] = client
			h.mutex.Unlock()

			// Store connection in Redis
			if connManager != nil {
				connManager.StoreConnection(client.UserID)
			}

			log.Printf("Client connected: User ID %d", client.UserID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)

				// Remove connection from Redis
				if connManager != nil {
					connManager.RemoveConnection(client.UserID)
				}

				log.Printf("Client disconnected: User ID %d", client.UserID)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for userID, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					delete(h.clients, userID)
					close(client.Send)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func SendNotificationToUser(userID uint, notification_data interface{}) {
	hub.mutex.RLock()
	client, exists := hub.clients[userID]
	hub.mutex.RUnlock()

	if exists {
		message, err := json.Marshal(map[string]interface{}{
			"type": "notification",
			"data": notification_data,
		})
		if err != nil {
			log.Printf("Error marshaling notification: %v", err)
			return
		}

		select {
		case client.Send <- message:
			log.Printf("Notification sent to user %d", userID)
		default:
			// Client's send channel is blocked, remove the client
			hub.mutex.Lock()
			delete(hub.clients, userID)
			close(client.Send)
			hub.mutex.Unlock()
		}
	}
}

// WebSocket endpoint handler
func HandleWebSocket(c *gin.Context) {
	// Get user ID from JWT token (set by middleware)
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create new client
	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	// Register client
	hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	// Set read deadline and pong handler for keepalive
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		// Read message from client
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming message (you can add custom logic here)
		log.Printf("Received message from user %d: %s", c.UserID, message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// GetConnectedUsers returns the list of currently connected user IDs
func GetConnectedUsers() []uint {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	users := make([]uint, 0, len(hub.clients))
	for userID := range hub.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline checks if a user is currently connected
func IsUserOnline(userID uint) bool {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	_, exists := hub.clients[userID]
	return exists
}

// SendNotificationToRoom sends a notification to all users in a specific room
func SendNotificationToRoom(roomID string, notification models.Notification) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	// Extract user IDs from room ID (format: "userID1_userID2")
	parts := strings.Split(roomID, "_")
	if len(parts) != 2 {
		log.Printf("Invalid room ID format: %s", roomID)
		return
	}

	// Convert parts to user IDs
	var userIDs []uint
	for _, part := range parts {
		if userIDStr := strings.TrimSpace(part); userIDStr != "" {
			var userID uint
			if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err == nil {
				userIDs = append(userIDs, userID)
			}
		}
	}

	// Send notification to each user in the room
	for _, userID := range userIDs {
		if client, exists := hub.clients[userID]; exists {
			message, err := json.Marshal(map[string]interface{}{
				"type": "notification",
				"data": notification,
			})
			if err != nil {
				log.Printf("Error marshaling notification: %v", err)
				continue
			}

			select {
			case client.Send <- message:
				log.Printf("Room notification sent to user %d via local connection", userID)
			default:
				// Client's send channel is blocked, remove the client
				delete(hub.clients, userID)
				close(client.Send)
				// Remove from Redis as well
				if connManager != nil {
					connManager.RemoveConnection(client.UserID)
				}
			}
		} else {
			// Try to send via Redis (cross-server communication)
			if connManager != nil {
				if err := connManager.SendNotification(userID, notification); err != nil {
					log.Printf("Failed to send notification via Redis to user %d: %v", userID, err)
				} else {
					log.Printf("Room notification sent to user %d via Redis", userID)
				}
			}
		}
	}
}
