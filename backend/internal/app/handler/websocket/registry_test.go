package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// reset the global registry before each test to ensure test isolation
func resetRegistry() {
	connections.Lock()
	defer connections.Unlock()

	connections.clients = make(map[string]*wsConnection)
}

func getServerSideConn(t *testing.T) (conn *websocket.Conn, cleanupFunc func()) {
	connCh := make(chan *websocket.Conn, 1)

	// dummy http server, to upgrade incoming connections to websockets 
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)

		connCh <- conn

		// keep connection open for testing
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}))

	// client connects to the test server to trigger the upgrade and get the server-side connection
	dialer := websocket.DefaultDialer
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	clientConn, _, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err)

	// fetch the server-side connection from the channel
	serverConn := <-connCh

	// cleanup function to close connections and server after test
	cleanup := func() {
		clientConn.Close()
		serverConn.Close()
		srv.Close()
	}

	return serverConn, cleanup
}

func TestRegisterAndUnregisterConnection(t *testing.T) {
	resetRegistry()

	transactionID := uuid.New()

	// 1. register successfully
	success := registerConnection(transactionID, nil)
	assert.True(t, success)

	connections.RLock()
	_, exists := connections.clients[transactionID.String()]
	connections.RUnlock()
	assert.True(t, exists, "Connection should be registered")

	// 2. Wieder abmelden
	unregisterConnection(transactionID)

	connections.RLock()
	_, exists = connections.clients[transactionID.String()]
	connections.RUnlock()
	assert.False(t, exists, "Connection should be unregistered")
}

func TestRegisterConnection_MaxLimit(t *testing.T) {
	resetRegistry()

	// fill the registry to its max capacity
	for i := 0; i < maxConnections; i++ {
		success := registerConnection(uuid.New(), nil)
		assert.True(t, success)
	}

	// now the registry is full, the next registration should fail
	success := registerConnection(uuid.New(), nil)
	assert.False(t, success, "Connection should be rejected due to maxConnections limit")
}

func TestPushUpdate_Success(t *testing.T) {
	resetRegistry()
	gin.SetMode(gin.TestMode)

	// get a real server-side connection to test the full flow of PushUpdate
	conn, cleanup := getServerSideConn(t)
	defer cleanup()

	transactionID := uuid.New()
	registerConnection(transactionID, conn)

	// manipulate lastSeen to test if it gets updated on PushUpdate
	connections.Lock()
	oldTime := time.Now().Add(-1 * time.Hour)
	connections.clients[transactionID.String()].lastSeen = oldTime
	connections.Unlock()

	// Push update
	assert.NotPanics(t, func() {
		PushUpdate(transactionID, models.PurchaseStatus("completed"))
	})

	// Check if lastSeen was updated
	connections.RLock()
	client := connections.clients[transactionID.String()]
	connections.RUnlock()

	assert.True(t, client.lastSeen.After(oldTime), "lastSeen should be updated to Now()")
}

func TestPushUpdate_UnknownTransaction(t *testing.T) {
	resetRegistry()
	
	assert.NotPanics(t, func() {
		PushUpdate(uuid.New(), models.PurchaseStatus("completed"))
	})
}

func TestPushUpdate_SendError(t *testing.T) {
	resetRegistry()

	// get a real server-side connection to test the full flow of PushUpdate
	conn, cleanup := getServerSideConn(t)
	conn.Close() // Provozieren des Fehlers

	defer cleanup()

	transactionID := uuid.New()
	registerConnection(transactionID, conn)

	assert.NotPanics(t, func() {
		PushUpdate(transactionID, models.PurchaseStatus("failed"))
	})

	// sendWSMessage should remove the faulty connection from the registry
	connections.RLock()
	_, exists := connections.clients[transactionID.String()]
	connections.RUnlock()

	assert.False(t, exists, "Faulty connection should be removed from the registry by sendWSMessage")
}

func TestCleanupStaleConnections(t *testing.T) {
	resetRegistry()

	// Get two server-side connections to test the cleanup of stale connections
	conn1, cleanup1 := getServerSideConn(t)
	defer cleanup1()

	conn2, cleanup2 := getServerSideConn(t)
	defer cleanup2()

	activeID := uuid.New()
	staleID := uuid.New()

	registerConnection(activeID, conn1)
	registerConnection(staleID, conn2)

	// Let the stale connection appear old by setting lastSeen to a time in the past
	connections.Lock()
	connections.clients[staleID.String()].lastSeen = time.Now().Add(-10 * time.Minute)
	connections.Unlock()

	// Initiate cleanup of stale connections with a threshold of 5 minutes
	cleanupStaleConnections(5 * time.Minute)

	// Check
	connections.RLock()
	_, activeExists := connections.clients[activeID.String()]
	_, staleExists := connections.clients[staleID.String()]
	connections.RUnlock()

	assert.True(t, activeExists, "Active connection should still exist")
	assert.False(t, staleExists, "Stale connection should have been deleted")
}
