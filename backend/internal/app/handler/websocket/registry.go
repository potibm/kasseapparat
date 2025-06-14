package websocket

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

func PushUpdate(transactionID uuid.UUID, status string) {
	connections.RLock()
	client, ok := connections.clients[transactionID.String()]
	connections.RUnlock()

	if !ok {
		return
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	_ = client.Conn.WriteJSON(map[string]interface{}{
		"type":   "status_update",
		"status": status,
	})
}
