package websocket

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

	ctx := c.Request.Context()
	startTime := time.Now()

	defer func() {
		conn.Close()

		defer unregisterConnection(transactionID)

		activeConns.Add(ctx, -1)
		connDuration.Record(ctx, time.Since(startTime).Seconds())
		slog.Info("WebSocket disconnected", "transaction_id", transactionID.String())
	}()

	if err := h.sendInitialStatus(conn, transactionID); err != nil {
		slog.Warn("Sending initial status failed", "error", err)

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

	protocols := websocket.Subprotocols(c.Request)
	var responseHeader http.Header
	if len(protocols) > 0 {
		responseHeader = http.Header{"Sec-WebSocket-Protocol": {protocols[0]}}
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, responseHeader)
	if err != nil {
		slog.Warn("WebSocket upgrade failed", "transaction_id", transactionID.String(), "error", err)

		return uuid.Nil, nil, false
	}

	activeConns.Add(c.Request.Context(), 1)
	slog.Info("WebSocket connected", "transaction_id", transactionID.String())

	if !registerConnection(transactionID, conn) {
		msg := websocket.FormatCloseMessage(CloseTooManyConnections, "connection limit reached")
		_ = conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
		conn.Close()
		activeConns.Add(c.Request.Context(), -1)
		slog.Error("Connection limit reached for transaction", "transaction_id", transactionID.String())

		return uuid.Nil, nil, false
	}

	return transactionID, conn, true
}

func (h *Handler) sendInitialStatus(conn *websocket.Conn, transactionID uuid.UUID) error {
	purchase, err := h.sqliteRepository.GetPurchaseByID(transactionID)
	if err != nil {
		slog.Warn("Failed to get purchase by ID", "error", err)

		sendWSMessage(conn, "error", gin.H{"message": "failed to retrieve purchase"}, transactionID)

		return err
	}

	sendWSMessage(conn, "status_update", gin.H{"status": purchase.Status}, transactionID)

	return nil
}

func (h *Handler) listenAndHandleMessages(conn *websocket.Conn, transactionID uuid.UUID) {
	for {
		ctx := context.Background()

		var msg map[string]any
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
			) {
				slog.WarnContext(ctx, "Websocket read error", "transaction_id", transactionID.String(), "error", err)
			} else {
				slog.InfoContext(ctx, "WebSocket connection closed", "transaction_id", transactionID.String())
			}

			break
		}

		msgType, _ := msg["type"].(string)

		msgRecvCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("msg_type", msgType),
				attribute.String("transaction_id", transactionID.String()),
			),
		)

		slog.DebugContext(
			ctx,
			"Websocket message for transaction",
			"transaction_id",
			transactionID.String(),
			"message_type",
			msgType,
		)
		msgRecvCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("msg_type", msgType),
			),
		)

		switch msgType {
		case "cancel_payment":
			h.handleCancelPayment(conn, msg, transactionID)
		case "ping":
			sendWSMessage(conn, "ping_ack", gin.H{}, transactionID)

		default:
			sendWSMessage(conn, "error", gin.H{"message": "unknown command"}, transactionID)
		}
	}
}

func (h *Handler) handleCancelPayment(conn *websocket.Conn, msg map[string]any, transactionID uuid.UUID) {
	readerID, ok := msg["reader_id"].(string)
	if !ok || readerID == "" {
		sendWSMessage(conn, "error", gin.H{"message": "reader_id missing or invalid"}, transactionID)

		return
	}

	err := h.sumupRepository.CreateReaderTerminateAction(readerID)
	if err != nil {
		sendWSMessage(conn, "error", gin.H{"message": "failed to cancel payment"}, transactionID)
	} else {
		sendWSMessage(conn, "cancel_ack", gin.H{"transaction_id": transactionID}, transactionID)
	}
}
