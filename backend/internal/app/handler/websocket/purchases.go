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
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connected for transaction: %s", transactionID)

	registerConnection(transactionID, conn)
	defer unregisterConnection(transactionID)

	purchase, err := h.sqliteRepository.GetPurchaseByID(transactionID)
	if err != nil {
		log.Println("Failed to get purchase by ID:", err)

		_ = conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "failed to retrieve purchase",
		})

		return
	}

	_ = conn.WriteJSON(map[string]interface{}{
		"type":   "status_update",
		"status": purchase.Status,
	})

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
			readerID, ok := msg["reader_id"].(string)
			if !ok || readerID == "" {
				_ = conn.WriteJSON(map[string]interface{}{
					"type":    "error",
					"message": "reader_id missing or invalid",
				})

				break
			}

			err = h.sumupRepository.CreateReaderTerminateAction(readerID)
			if err != nil {
				_ = conn.WriteJSON(map[string]interface{}{
					"type":    "error",
					"messahe": "failed to cancel payment",
				})
			} else {
				_ = conn.WriteJSON(map[string]interface{}{
					"type":           "cancel_ack",
					"transaction_id": transactionID,
				})
			}
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
