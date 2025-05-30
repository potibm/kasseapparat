package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/shopspring/decimal"
)

type SumupTransactionReponse struct {
	ID              string          `json:"id"`
	TransactionCode string          `json:"transactionCode"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	CreatedAt       time.Time       `json:"createdAt"`
	Status          string          `json:"status"`
}

func (handler *Handler) GetSumupTransactions(c *gin.Context) {
	transactions, _ := handler.sumupRepository.GetTransactions()

	c.Header("X-Total-Count", strconv.Itoa(len(transactions)))
	c.JSON(http.StatusOK, toSumupTransactionResponses(transactions))
}

func toSumupTransactionResponse(c sumup.Transaction) SumupTransactionReponse {
	return SumupTransactionReponse{
		ID:              c.ID,
		TransactionCode: c.TransactionCode,
		Amount:          c.Amount,
		Currency:        c.Currency,
		CreatedAt:       c.CreatedAt,
		Status:          c.Status,
	}
}

func toSumupTransactionResponses(transactions []sumup.Transaction) []SumupTransactionReponse {
	responses := make([]SumupTransactionReponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = toSumupTransactionResponse(transaction)
	}

	return responses
}
