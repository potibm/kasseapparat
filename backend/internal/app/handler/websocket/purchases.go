package websocket

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (h *Handler) HandleTransactionWebSocket(c *gin.Context) {
	protocols := websocket.Subprotocols(c.Request)
	if len(protocols) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})

		return
	}

	tokenStr := protocols[0]

	token, err := h.jwtMiddleware.ParseTokenString(tokenStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})

		return
	}

	claims := jwt.ExtractClaimsFromToken(token)
	identity := claims[h.jwtMiddleware.IdentityKey]
	c.Set("identity", identity)

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

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)

		return uuid.Nil, nil, false
	}

	log.Printf("WebSocket connected for transaction: %s", transactionID)

	if !registerConnection(transactionID, conn) {
		msg := websocket.FormatCloseMessage(CloseTooManyConnections, "connection limit reached")
		conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
		conn.Close()

		return uuid.Nil, nil, false
	}

	return transactionID, conn, true
}

func (h *Handler) sendInitialStatus(conn *websocket.Conn, transactionID uuid.UUID) error {
	purchase, err := h.sqliteRepository.GetPurchaseByID(transactionID)
	if err != nil {
		log.Println("Failed to get purchase by ID:", err)

		sendWSMessage(conn, "error", gin.H{"message": "failed to retrieve purchase"}, &uuid.Nil)

		return err
	}

	sendWSMessage(conn, "status_update", gin.H{"status": purchase.Status}, &transactionID)

	return nil
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
			sendWSMessage(conn, "ping_ack", gin.H{}, &uuid.Nil)

		default:
			sendWSMessage(conn, "error", gin.H{"message": "unknown command"}, &uuid.Nil)
		}
	}
}

func (h *Handler) handleCancelPayment(conn *websocket.Conn, msg map[string]interface{}, transactionID uuid.UUID) {
	readerID, ok := msg["reader_id"].(string)
	if !ok || readerID == "" {
		sendWSMessage(conn, "error", gin.H{"message": "reader_id missing or invalid"}, &transactionID)

		return
	}

	err := h.sumupRepository.CreateReaderTerminateAction(readerID)
	if err != nil {
		sendWSMessage(conn, "error", gin.H{"message": "failed to cancel payment"}, &transactionID)
	} else {
		sendWSMessage(conn, "cancel_ack", gin.H{"transaction_id": transactionID}, &transactionID)
	}
}
