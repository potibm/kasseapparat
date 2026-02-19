package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/potibm/kasseapparat/internal/app/models"
)

const maxConnections = 100
const CloseTooManyConnections = 4001
const CloseStaleConnection = 4005

const CleanupStaleConnectionsInterval = 5 * time.Minute

type wsConnection struct {
	Conn     *websocket.Conn
	mu       sync.Mutex
	lastSeen time.Time
}

var connections = struct {
	sync.RWMutex

	clients map[string]*wsConnection // transactionID → conn
}{
	clients: make(map[string]*wsConnection),
}

func registerConnection(transactionID uuid.UUID, conn *websocket.Conn) bool {
	connections.Lock()
	defer connections.Unlock()

	log.Println("Current WebSocket connection number:", len(connections.clients))

	if len(connections.clients) >= maxConnections {
		log.Printf("WebSocket connection limit reached (%d)", maxConnections)

		return false
	}

	connections.clients[transactionID.String()] = &wsConnection{
		Conn:     conn,
		lastSeen: time.Now(),
	}

	return true
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

	client.lastSeen = time.Now()
	sendWSMessage(client.Conn, "status_update", gin.H{"status": string(status)}, &transactionID)
}

func sendWSMessage(conn *websocket.Conn, msgType string, data gin.H, transactionID *uuid.UUID) {
	payload := data
	payload["type"] = msgType

	if transactionID != nil {
		payload["transaction_id"] = transactionID.String()
	}

	if err := conn.WriteJSON(payload); err != nil {
		log.Printf("WebSocket send error [%s]: %v", msgType, err)

		closeMsg := websocket.FormatCloseMessage(
			websocket.CloseAbnormalClosure,
			"failed to send message",
		)

		_ = conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
		conn.Close()

		if transactionID != nil {
			unregisterConnection(*transactionID)
		}

		return
	}
}

func StartCleanupRoutine(timeout time.Duration) {
	go func() {
		ticker := time.NewTicker(CleanupStaleConnectionsInterval)
		defer ticker.Stop()

		for range ticker.C {
			cleanupStaleConnections(timeout)
		}
	}()
}

func cleanupStaleConnections(timeout time.Duration) {
	now := time.Now()

	var (
		staleIDs   []string
		staleConns []*websocket.Conn
	)

	connections.Lock()

	for id, conn := range connections.clients {
		conn.mu.Lock()
		inactive := now.Sub(conn.lastSeen) > timeout
		conn.mu.Unlock()

		if inactive {
			staleIDs = append(staleIDs, id)
			staleConns = append(staleConns, conn.Conn)
		}
	}
	// Remove stale connections from the registry while holding the lock
	for _, id := range staleIDs {
		delete(connections.clients, id)
	}

	connections.Unlock()

	// Close WebSocket connections outside the lock to avoid blocking other operations
	for i, ws := range staleConns {
		log.Printf("Cleaning up stale WebSocket connection: %s", staleIDs[i])

		msg := websocket.FormatCloseMessage(CloseStaleConnection, "connection stale")
		_ = ws.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))

		ws.Close()
	}
}
