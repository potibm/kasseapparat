package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // for dev only
	},
}

func (h *Handler) HandleTransactionWebSocket(c *gin.Context) {
	transactionID, conn, ok := h.upgradeAndRegister(c)
	if !ok {
		return
	}

	defer conn.Close()
	defer unregisterConnection(transactionID)

	if err := h.sendInitialStatus(conn, transactionID); err != nil {
		log.Println("Sending initial status failed:", err)
		return
	}

	h.listenAndHandleMessages(conn, transactionID)
}

func (h *Handler) upgradeAndRegister(c *gin.Context) (uuid.UUID, *websocket.Conn, bool) {
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return uuid.Nil, nil, false
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return uuid.Nil, nil, false
	}

	log.Printf("WebSocket connected for transaction: %s", transactionID)

	registerConnection(transactionID, conn)

	return transactionID, conn, true
}

func (h *Handler) sendInitialStatus(conn *websocket.Conn, transactionID uuid.UUID) error {
	purchase, err := h.sqliteRepository.GetPurchaseByID(transactionID)
	if err != nil {
		log.Println("Failed to get purchase by ID:", err)

		_ = conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "failed to retrieve purchase",
		})

		return err
	}

	return conn.WriteJSON(map[string]interface{}{
		"type":   "status_update",
		"status": purchase.Status,
	})
}

func (h *Handler) listenAndHandleMessages(conn *websocket.Conn, transactionID uuid.UUID) {
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WS read error:", err)
			break
		}

		msgType, _ := msg["type"].(string)
		log.Printf("WS message for %s: %v", transactionID, msg)

		switch msgType {
		case "cancel_payment":
			h.handleCancelPayment(conn, msg, transactionID)
		case "ping":
			_ = conn.WriteJSON(map[string]interface{}{
				"type": "pong",
			})

		default:
			_ = conn.WriteJSON(map[string]interface{}{
				"type":           "error",
				"message":        "unknown command",
				"transaction_id": transactionID,
			})
		}
	}
}

func (h *Handler) handleCancelPayment(conn *websocket.Conn, msg map[string]interface{}, transactionID uuid.UUID) {
	readerID, ok := msg["reader_id"].(string)
	if !ok || readerID == "" {
		_ = conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "reader_id missing or invalid",
		})

		return
	}

	err := h.sumupRepository.CreateReaderTerminateAction(readerID)
	if err != nil {
		_ = conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "failed to cancel payment",
		})
	} else {
		_ = conn.WriteJSON(map[string]interface{}{
			"type":           "cancel_ack",
			"transaction_id": transactionID,
		})
	}
}
