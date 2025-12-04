package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gin-project/config"
	"gin-project/models"

	"github.com/redis/go-redis/v9"
)

type ConnectionManager struct {
	redis    *redis.Client
	serverID string
}

var connManager *ConnectionManager

func init() {
	// Generate unique server ID (you could use hostname, UUID, etc.)
	hostname, _ := os.Hostname()
	serverID := fmt.Sprintf("%s-%d", hostname, time.Now().Unix())

	connManager = &ConnectionManager{
		serverID: serverID,
	}
}

func initConnectionManager() {
	connManager.redis = config.GetRedisClient()

	// Start listening for notifications
	go connManager.listenForNotifications()
}

// listenForNotifications subscribes to Redis pub/sub for incoming notifications
func (cm *ConnectionManager) listenForNotifications() {
	if cm.redis == nil {
		return
	}

	ctx := context.Background()
	pattern := "ws:notify:user:*"

	pubsub := cm.redis.PSubscribe(ctx, pattern)
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		// Extract user ID from channel name "ws:notify:user:123"
		parts := strings.Split(msg.Channel, ":")
		if len(parts) != 4 {
			continue
		}

		userID, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			continue
		}

		// Check if user is connected to this server
		if client := getLocalClient(uint(userID)); client != nil {
			// Parse the message
			var notification map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
				log.Printf("Error unmarshaling notification: %v", err)
				continue
			}

			// Don't send notifications from the same server (avoid loops)
			if fromServer, ok := notification["from_server"].(string); ok && fromServer == cm.serverID {
				continue
			}

			// Send to local client
			select {
			case client.Send <- []byte(msg.Payload):
				log.Printf("Redis notification delivered to user %d", userID)
			default:
				// Client channel is blocked, disconnect
				log.Printf("Client %d channel blocked, disconnecting", userID)
				close(client.Send)
				// Remove from local map and Redis
				removeLocalClient(uint(userID))
				cm.RemoveConnection(uint(userID))
			}
		}
	}
}

// StoreConnection stores a user's connection info in Redis
func (cm *ConnectionManager) StoreConnection(userID uint) error {
	if cm.redis == nil {
		// Fall back to in-memory if Redis is not available
		return nil
	}

	ctx := context.Background()
	key := fmt.Sprintf("ws:user:%d", userID)

	connectionInfo := map[string]interface{}{
		"server_id":    cm.serverID,
		"connected_at": time.Now().Unix(),
		"status":       "online",
	}

	data, err := json.Marshal(connectionInfo)
	if err != nil {
		return err
	}

	// Store with expiration (heartbeat mechanism)
	return cm.redis.Set(ctx, key, data, 5*time.Minute).Err()
}

// RemoveConnection removes a user's connection from Redis
func (cm *ConnectionManager) RemoveConnection(userID uint) error {
	if cm.redis == nil {
		return nil
	}

	ctx := context.Background()
	key := fmt.Sprintf("ws:user:%d", userID)
	return cm.redis.Del(ctx, key).Err()
}

// IsUserOnlineRedis checks if a user is online using Redis
func (cm *ConnectionManager) IsUserOnlineRedis(userID uint) bool {
	if cm.redis == nil {
		// Fall back to in-memory check
		return IsUserOnline(userID)
	}

	ctx := context.Background()
	key := fmt.Sprintf("ws:user:%d", userID)

	exists, err := cm.redis.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("Error checking user online status: %v", err)
		return false
	}

	return exists > 0
}

// GetUserServer returns which server a user is connected to
func (cm *ConnectionManager) GetUserServer(userID uint) (string, error) {
	if cm.redis == nil {
		// If user is connected to current server
		if IsUserOnline(userID) {
			return cm.serverID, nil
		}
		return "", fmt.Errorf("user not connected")
	}

	ctx := context.Background()
	key := fmt.Sprintf("ws:user:%d", userID)

	data, err := cm.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("user not connected")
	} else if err != nil {
		return "", err
	}

	var connectionInfo map[string]interface{}
	if err := json.Unmarshal([]byte(data), &connectionInfo); err != nil {
		return "", err
	}

	serverID, ok := connectionInfo["server_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid server_id in connection info")
	}

	return serverID, nil
}

// SendNotificationViaRedis sends notification using Redis pub/sub if user is on different server
func SendNotificationViaRedis(userID uint, notification models.Notification) error {
	if connManager.redis == nil {
		// Fall back to local sending
		SendNotificationToUser(userID, notification)
		return nil
	}

	// Check if user is connected locally first
	if IsUserOnline(userID) {
		SendNotificationToUser(userID, notification)
		return nil
	}

	// Check if user is connected to another server
	serverID, err := connManager.GetUserServer(userID)
	if err != nil {
		return fmt.Errorf("user not connected: %v", err)
	}

	if serverID == connManager.serverID {
		// User is on this server, send locally
		SendNotificationToUser(userID, notification)
		return nil
	}

	// User is on different server, use Redis pub/sub
	return publishNotification(userID, notification)
}

// publishNotification publishes notification to Redis for other servers
func publishNotification(userID uint, notification models.Notification) error {
	ctx := context.Background()
	channel := fmt.Sprintf("ws:notify:user:%d", userID)

	data, err := json.Marshal(map[string]interface{}{
		"type":         "notification",
		"user_id":      userID,
		"notification": notification,
		"from_server":  connManager.serverID,
	})
	if err != nil {
		return err
	}

	return connManager.redis.Publish(ctx, channel, data).Err()
}

// StartNotificationSubscriber starts listening for Redis notifications
func StartNotificationSubscriber() {
	if connManager.redis == nil {
		log.Println("Redis not available, skipping notification subscriber")
		return
	}

	initConnectionManager()

	ctx := context.Background()
	pattern := "ws:notify:user:*"

	pubsub := connManager.redis.PSubscribe(ctx, pattern)
	defer pubsub.Close()

	log.Printf("Server %s started listening for Redis notifications", connManager.serverID)

	for msg := range pubsub.Channel() {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
			log.Printf("Error unmarshaling notification: %v", err)
			continue
		}

		// Skip if message is from this server
		fromServer, ok := data["from_server"].(string)
		if ok && fromServer == connManager.serverID {
			continue
		}

		// Extract user ID and notification
		userIDFloat, ok := data["user_id"].(float64)
		if !ok {
			log.Printf("Invalid user_id in notification")
			continue
		}
		userID := uint(userIDFloat)

		// Check if user is connected to this server
		if IsUserOnline(userID) {
			notificationData, ok := data["notification"]
			if !ok {
				log.Printf("Invalid notification data")
				continue
			}

			// Send notification locally
			SendNotificationToUser(userID, notificationData)
		}
	}
}

// Heartbeat function to keep connection alive in Redis
func StartHeartbeat() {
	if connManager.redis == nil {
		return
	}

	ticker := time.NewTicker(2 * time.Minute)
	go func() {
		for range ticker.C {
			// Update all connected users' heartbeat
			hub.mutex.RLock()
			for userID := range hub.clients {
				connManager.StoreConnection(userID)
			}
			hub.mutex.RUnlock()
		}
	}()
}

// GetOnlineUsersFromRedis returns all online users across all servers
func GetOnlineUsersFromRedis() ([]uint, error) {
	if connManager.redis == nil {
		return GetConnectedUsers(), nil
	}

	ctx := context.Background()
	pattern := "ws:user:*"

	keys, err := connManager.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var userIDs []uint
	for _, key := range keys {
		// Extract user ID from key "ws:user:123"
		parts := strings.Split(key, ":")
		if len(parts) == 3 {
			if userID, err := strconv.ParseUint(parts[2], 10, 32); err == nil {
				userIDs = append(userIDs, uint(userID))
			}
		}
	}

	return userIDs, nil
}

// SendNotification sends a notification to a user via Redis pub/sub
func (cm *ConnectionManager) SendNotification(userID uint, notification interface{}) error {
	if cm.redis == nil {
		return fmt.Errorf("Redis not available")
	}

	ctx := context.Background()
	channel := fmt.Sprintf("ws:notify:user:%d", userID)

	message, err := json.Marshal(map[string]interface{}{
		"type":        "notification",
		"data":        notification,
		"from_server": cm.serverID,
		"timestamp":   time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	return cm.redis.Publish(ctx, channel, message).Err()
}
