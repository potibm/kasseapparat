package websocket

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/potibm/kasseapparat/internal/app/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

	slog.Info("Current WebSocket connection number", "count", len(connections.clients))

	if len(connections.clients) >= maxConnections {
		slog.Error("WebSocket connection limit reached", "max_connections", maxConnections)

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
		slog.Warn("No WebSocket client for transaction", "transaction_id", transactionID.String())

		return
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	client.lastSeen = time.Now()
	sendWSMessage(client.Conn, "status_update", gin.H{"status": string(status)}, transactionID)
}

func sendWSMessage(conn *websocket.Conn, msgType string, data gin.H, transactionID uuid.UUID) {
	payload := data
	payload["type"] = msgType

	if err := conn.WriteJSON(payload); err != nil {
		slog.Error(
			"WebSocket send error",
			"message_type",
			msgType,
			"error",
			err,
			"transaction_id",
			transactionID.String(),
		)

		closeMsg := websocket.FormatCloseMessage(
			websocket.CloseAbnormalClosure,
			"failed to send message",
		)

		_ = conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
		conn.Close()

		if transactionID != uuid.Nil {
			unregisterConnection(transactionID)
		}

		msgSentCounter.Add(context.Background(), 1,
			metric.WithAttributes(
				attribute.String("msg_type", msgType),
				attribute.String("transaction_id", transactionID.String()),
			),
		)
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

	slog.Debug("Checking for stale WebSocket connections", "total_connections", len(connections.clients))
	connections.Lock()

	for id, conn := range connections.clients {
		conn.mu.Lock()
		inactive := now.Sub(conn.lastSeen) > timeout
		conn.mu.Unlock()
		slog.Debug("Checking connection activity", "transaction_id", id, "last_seen", conn.lastSeen, "inactive", inactive)

		if inactive {
			staleIDs = append(staleIDs, id)
			staleConns = append(staleConns, conn.Conn)
		}
	}

	numRemoved := len(staleIDs)
	// Remove stale connections from the registry while holding the lock
	for _, id := range staleIDs {
		delete(connections.clients, id)
	}
	connections.Unlock()

	if numRemoved > 0 {
		activeConns.Add(context.Background(), int64(-numRemoved))
	}

	// Close WebSocket connections outside the lock to avoid blocking other operations
	for i, ws := range staleConns {
		slog.Info("Cleaning up stale WebSocket connection", "transaction_id", staleIDs[i])

		msg := websocket.FormatCloseMessage(CloseStaleConnection, "connection stale")
		_ = ws.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))

		ws.Close()
	}
}
