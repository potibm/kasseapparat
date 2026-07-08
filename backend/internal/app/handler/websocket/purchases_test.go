package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- MOCKS ---

// Mock for the SQLite Repository.
type mockSqliteRepo struct {
	mock.Mock
}

func (m *mockSqliteRepo) GetPurchaseByID(id uuid.UUID) (*models.Purchase, error) {
	args := m.Called(id)

	purchase, _ := args.Get(0).(*models.Purchase)

	return purchase, args.Error(1)
}

// Mock for the SumUp Repository.
type mockSumupRepo struct {
	sumup.RepositoryInterface
	mock.Mock
}

func (m *mockSumupRepo) CreateReaderTerminateAction(readerID string) error {
	args := m.Called(readerID)

	return args.Error(0)
}

// --- TEST SETUP ---

func setupTestServer(t *testing.T) (*Handler, *httptest.Server, *mockSqliteRepo, *mockSumupRepo) {
	gin.SetMode(gin.TestMode)

	mockSqlite := new(mockSqliteRepo)
	mockSumup := new(mockSumupRepo)

	handler := &Handler{
		sqliteRepository: mockSqlite,
		sumupRepository:  mockSumup,
		upgrader: websocket.Upgrader{
			// Allow testing from any origin.
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	router := gin.New()
	router.GET("/ws/:id", handler.HandleTransactionWebSocket)

	server := httptest.NewServer(router)

	return handler, server, mockSqlite, mockSumup
}

func TestHandleTransactionWebSocketInvalidUUID(t *testing.T) {
	//nolint:dogsled // intended use for the setup function
	_, server, _, _ := setupTestServer(t)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/invalid-uuid-format"

	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)

	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleTransactionWebSocketHappyPath(t *testing.T) {
	_, server, mockSqlite, mockSumup := setupTestServer(t)
	defer server.Close()

	transactionID := uuid.New()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + transactionID.String()

	// 1. Define Mock Expectations
	// Initial status call
	mockSqlite.On("GetPurchaseByID", transactionID).Return(&models.Purchase{Status: "pending"}, nil)
	// Cancel Payment Call
	mockSumup.On("CreateReaderTerminateAction", "reader-123").Return(nil)

	// 2. Setup connection
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

	defer conn.Close()

	// 3. Read initial status
	var initialMsg map[string]interface{}

	err = conn.ReadJSON(&initialMsg)
	require.NoError(t, err)
	assert.Equal(t, "pending", initialMsg["status"])

	// 4. Test ping
	err = conn.WriteJSON(map[string]interface{}{
		"type": "ping",
	})
	require.NoError(t, err)

	var pingAck map[string]interface{}

	err = conn.ReadJSON(&pingAck)
	require.NoError(t, err)
	assert.Equal(t, "ping_ack", pingAck["type"])

	// 5. Test cancel payment
	err = conn.WriteJSON(map[string]interface{}{
		"type":      "cancel_payment",
		"reader_id": "reader-123",
	})
	require.NoError(t, err)

	var cancelAck map[string]interface{}

	err = conn.ReadJSON(&cancelAck)
	require.NoError(t, err)
	assert.Equal(t, "cancel_ack", cancelAck["type"])

	mockSqlite.AssertExpectations(t)
	mockSumup.AssertExpectations(t)
}
