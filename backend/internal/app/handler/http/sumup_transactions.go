package http

import (
	"errors"
	"log"
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

type SumupTransactionReponse struct {
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

	transactions, _ := handler.sumupRepository.GetTransactions(oldestTime)

	transactionsLen := len(transactions)

	if end > transactionsLen {
		end = transactionsLen
	}

	// limit the results based on start and end parameters
	if start < 0 || start >= end {
		c.Header("X-Total-Count", "0")
		c.JSON(http.StatusOK, []SumupTransactionReponse{})

		return
	}

	transactions = transactions[start:end]

	c.Header("X-Total-Count", strconv.Itoa(transactionsLen))
	c.JSON(http.StatusOK, toSumupTransactionResponses(transactions))
}

func (handler *Handler) GetSumupTransactionByID(c *gin.Context) {
	transactionID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Invalid ID"))
		return
	}

	transaction, err := handler.sumupRepository.GetTransactionById(transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transaction"})
		return
	}

	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, toSumupTransactionResponse(*transaction))
}

func (handler *Handler) GetSumupTransactionWebhook(c *gin.Context) {
	var payload sumup.SumupTransactionWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if payload.EventType != "solo.transaction.updated" {
		log.Println("unsupported event type:", payload.EventType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported event type"})

		return
	}

	purchase, err := handler.repo.GetPurchaseBySumupClientTransactionID(payload.Payload.ClientTransactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
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
		log.Println("updating purchase status to confirmed for transaction ID:", payload.Payload.ClientTransactionID)
		handler.statusPublisher.PushUpdate(purchase.ID, model.PurchaseStatusConfirmed)

		_, err = handler.purchaseService.FinalizePurchase(ctx, purchase.ID)
	case sumup.StatusFailed:
		// Update purchase status to failed
		log.Println("updating purchase status to failed for transaction ID:", payload.Payload.ClientTransactionID)
		handler.statusPublisher.PushUpdate(purchase.ID, model.PurchaseStatusFailed)

		_, err = handler.purchaseService.FailPurchase(ctx, purchase.ID)
	default:
		log.Println("unsupported status:", payload.Payload.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported status"})
	}

	if err != nil {
		log.Println("failed to update purchase status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update purchase status"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func toSumupTransactionResponse(c sumup.Transaction) SumupTransactionReponse {
	return SumupTransactionReponse{
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

func toSumupTransactionResponses(transactions []sumup.Transaction) []SumupTransactionReponse {
	responses := make([]SumupTransactionReponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = toSumupTransactionResponse(transaction)
	}

	return responses
}
