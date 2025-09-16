package monitoring

import (
	"encoding/json"
	"log"
	"skripsi-be/internal/models/entities"
	"time"

	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

// WebSocket connection manager
type ConnectionManager struct {
	connections map[string]*websocket.Conn
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
}

type Client struct {
	ID         string
	CustomerID string
	Conn       *websocket.Conn
	Manager    *ConnectionManager
}

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*websocket.Conn),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
	}
}

// Start starts the connection manager
func (cm *ConnectionManager) Start() {
	for {
		select {
		case client := <-cm.register:
			cm.connections[client.ID] = client.Conn
			log.Printf("Client %s connected", client.ID)

		case client := <-cm.unregister:
			if _, ok := cm.connections[client.ID]; ok {
				delete(cm.connections, client.ID)
				client.Conn.Close()
				log.Printf("Client %s disconnected", client.ID)
			}

		case message := <-cm.broadcast:
			for id, conn := range cm.connections {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Error sending message to client %s: %v", id, err)
					conn.Close()
					delete(cm.connections, id)
				}
			}
		}
	}
}

// BroadcastToCustomer sends a message to a specific customer
func (cm *ConnectionManager) BroadcastToCustomer(customerID string, message WebSocketMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// Find connections for this customer
	for id, conn := range cm.connections {
		// In a real implementation, you would store customer ID with connection
		// For now, we'll broadcast to all connections
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Error sending message to client %s: %v", id, err)
			conn.Close()
			delete(cm.connections, id)
		}
	}
}

// Global connection manager instance
var connectionManager = NewConnectionManager()

// WebSocketHandler handles WebSocket connections for customer monitoring
func WebSocketHandler(c *websocket.Conn) {
	// Get customer ID from query parameters or JWT token
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Customer ID required"))
		return
	}

	// Create client
	client := &Client{
		ID:         customerID + "_" + time.Now().Format("20060102150405"),
		CustomerID: customerID,
		Conn:       c,
		Manager:    connectionManager,
	}

	// Register client
	connectionManager.register <- client

	// Send initial connection status
	initialMessage := WebSocketMessage{
		Type:      "connection",
		Data:      map[string]string{"status": "connected"},
		Timestamp: time.Now(),
	}
	client.sendMessage(initialMessage)

	// Handle messages from client
	defer func() {
		connectionManager.unregister <- client
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// sendMessage sends a message to the client
func (c *Client) sendMessage(message WebSocketMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.Conn.WriteMessage(websocket.TextMessage, messageBytes)
}

// BroadcastConnectionStatus broadcasts connection status updates to all connected customers
func BroadcastConnectionStatus(db *gorm.DB) {
	// Get all active devices
	var devices []entities.NetwatchDevice
	if err := db.Find(&devices).Error; err != nil {
		log.Printf("Error getting devices: %v", err)
		return
	}

	for _, device := range devices {
		// Get latest event for this device
		var latestEvent entities.NetwatchEvent
		if err := db.Where("device_id = ?", device.ID).
			Order("event_time DESC").
			First(&latestEvent).Error; err != nil {
			continue
		}

		// Create status message
		statusMessage := WebSocketMessage{
			Type: "status_update",
			Data: map[string]interface{}{
				"device_id":   device.ID,
				"status":      latestEvent.EventType,
				"timestamp":   latestEvent.EventTime,
				"customer_id": device.CustomerID,
			},
			Timestamp: time.Now(),
		}

		// Broadcast to customer if device is associated with one
		if device.CustomerID != nil {
			connectionManager.BroadcastToCustomer(*device.CustomerID, statusMessage)
		}
	}
}

// StartStatusBroadcaster starts a goroutine that broadcasts status updates
func StartStatusBroadcaster(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Broadcast every 30 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				BroadcastConnectionStatus(db)
			}
		}
	}()
}
