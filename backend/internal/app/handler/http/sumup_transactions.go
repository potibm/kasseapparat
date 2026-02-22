package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	model "github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SumupTransactionResponse struct {
	ID              uuid.UUID                       `json:"id"`
	TransactionID   uuid.UUID                       `json:"transactionId"`
	PublicID        string                          `json:"publicId"`
	TransactionCode string                          `json:"transactionCode"`
	Amount          decimal.Decimal                 `json:"amount"`
	CardType        string                          `json:"cardType,omitempty"`
	Currency        string                          `json:"currency"`
	Events          []SumupTransactionEventResponse `json:"events,omitempty"`
	CreatedAt       time.Time                       `json:"createdAt"`
	Status          string                          `json:"status"`
}

type SumupTransactionEventResponse struct {
	ID        int             `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Type      string          `json:"type"`
	Amount    decimal.Decimal `json:"amount"`
	Status    string          `json:"status,omitempty"`
}

func (handler *Handler) GetSumupTransactions(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	defaultOldestTime := time.Now().Add(-10 * time.Minute)
	oldestTime := queryTime(c, "oldest_time", &defaultOldestTime)

	transactions, err := handler.sumupRepository.GetTransactions(oldestTime)
	if err != nil {
		_ = c.Error(InternalServerError.WithMsg("failed to retrieve transactions").WithCause(err))

		return
	}

	transactionsLen := len(transactions)

	if end > transactionsLen {
		end = transactionsLen
	}

	// limit the results based on start and end parameters
	if start < 0 || start >= end {
		c.Header("X-Total-Count", "0")
		c.JSON(http.StatusOK, []SumupTransactionResponse{})

		return
	}

	transactions = transactions[start:end]

	c.Header("X-Total-Count", strconv.Itoa(transactionsLen))
	c.JSON(http.StatusOK, toSumupTransactionResponses(transactions))
}

func (handler *Handler) GetSumupTransactionByID(c *gin.Context) {
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(InvalidRequest.WithMsg("Invalid ID").WithCause(err))

		return
	}

	transaction, err := handler.sumupRepository.GetTransactionById(transactionID)
	if err != nil {
		_ = c.Error(InternalServerError.WithMsg("failed to retrieve transaction").WithCause(err))

		return
	}

	if transaction == nil {
		_ = c.Error(NotFound.WithMsg("Transaction not found"))

		return
	}

	c.JSON(http.StatusOK, toSumupTransactionResponse(*transaction))
}

func (handler *Handler) GetSumupTransactionWebhook(c *gin.Context) {
	var payload sumup.SumupTransactionWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(InvalidRequest.WithCause(err))

		return
	}

	if payload.EventType != "solo.transaction.updated" {
		slog.Warn("Unsupported event type", "event_type", payload.EventType)

		_ = c.Error(InvalidRequest.WithMsg("unsupported event type"))

		return
	}

	purchase, err := handler.repo.GetPurchaseBySumupClientTransactionID(payload.Payload.ClientTransactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(NotFound.WithMsg("purchase not found"))

			return
		}

		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	if purchase.Status != model.PurchaseStatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "purchase is not in pending status"})

		return
	}

	ctx := c.Request.Context()

	switch payload.Payload.Status {
	case sumup.StatusSuccessful:
		// Update purchase status to confirmed
		slog.Info("Updating purchase status to confirmed", "transaction_id", payload.Payload.ClientTransactionID)
		handler.statusPublisher.PushUpdate(purchase.ID, model.PurchaseStatusConfirmed)

		_, err = handler.purchaseService.FinalizePurchase(ctx, purchase.ID)
	case sumup.StatusFailed:
		// Update purchase status to failed
		slog.Info("Updating purchase status to failed", "transaction_id", payload.Payload.ClientTransactionID)
		handler.statusPublisher.PushUpdate(purchase.ID, model.PurchaseStatusFailed)

		_, err = handler.purchaseService.FailPurchase(ctx, purchase.ID)
	default:
		slog.Warn("Unsupported status", "status", payload.Payload.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported status"})
	}

	if err != nil {
		slog.Error("Failed to update purchase status", "error", err)
		_ = c.Error(InternalServerError.WithMsg("failed to update purchase status").WithCause(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func toSumupTransactionResponse(c sumup.Transaction) SumupTransactionResponse {
	return SumupTransactionResponse{
		ID:              c.TransactionID,
		TransactionID:   c.TransactionID,
		PublicID:        c.ID,
		TransactionCode: c.TransactionCode,
		Amount:          c.Amount,
		Currency:        c.Currency,
		CreatedAt:       c.CreatedAt,
		Status:          c.Status,
		CardType:        c.CardType,
		Events:          toSumupTransactionEventResponses(c.Events),
	}
}

func toSumupTransactionEventResponses(events []sumup.TransactionEvent) []SumupTransactionEventResponse {
	responses := make([]SumupTransactionEventResponse, len(events))
	for i, event := range events {
		responses[i] = toSumupTransactionEventResponse(event)
	}

	return responses
}

func toSumupTransactionEventResponse(e sumup.TransactionEvent) SumupTransactionEventResponse {
	return SumupTransactionEventResponse{
		ID:        e.ID,
		Timestamp: e.Timestamp,
		Type:      e.Type,
		Amount:    e.Amount,
		Status:    e.Status,
	}
}

func toSumupTransactionResponses(transactions []sumup.Transaction) []SumupTransactionResponse {
	responses := make([]SumupTransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = toSumupTransactionResponse(transaction)
	}

	return responses
}
