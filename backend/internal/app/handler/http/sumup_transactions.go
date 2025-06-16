package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/shopspring/decimal"
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
