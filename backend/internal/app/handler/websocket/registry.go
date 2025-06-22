package websocket

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/potibm/kasseapparat/internal/app/models"
)

type wsConnection struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

var connections = struct {
	sync.RWMutex
	clients map[string]*wsConnection // transactionID → conn
}{
	clients: make(map[string]*wsConnection),
}

func registerConnection(transactionID uuid.UUID, conn *websocket.Conn) {
	connections.Lock()
	defer connections.Unlock()

	connections.clients[transactionID.String()] = &wsConnection{Conn: conn}
}

func unregisterConnection(transactionID uuid.UUID) {
	connections.Lock()
	defer connections.Unlock()

	delete(connections.clients, transactionID.String())
}

func PushUpdate(transactionID uuid.UUID, status models.PurchaseStatus) {
	connections.RLock()
	client, ok := connections.clients[transactionID.String()]
	connections.RUnlock()

	if !ok {
		log.Printf("No WebSocket client for %s", transactionID)
		return
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	sendWSMessage(client.Conn, "status_update", gin.H{"status": string(status)}, &transactionID)
}

func sendWSMessage(conn *websocket.Conn, msgType string, data gin.H, transactionID *uuid.UUID) {
	payload := gin.H{
		"type": msgType,
	}
	for k, v := range data {
		payload[k] = v
	}

	if transactionID != nil {
		payload["transaction_id"] = transactionID.String()
	}

	if err := conn.WriteJSON(payload); err != nil {
		log.Printf("WebSocket send error [%s]: %v", msgType, err)
	}
}
