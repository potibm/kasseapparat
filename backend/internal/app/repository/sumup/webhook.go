package sumup

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	StatusSuccessful TransactionStatus = "successful"
	StatusFailed     TransactionStatus = "failed"
)

type SumupTransactionWebhookPayload struct {
	ID        uuid.UUID                          `json:"id"`
	EventType string                             `json:"event_type"`
	Payload   SumupTransactionWebhookPayloadData `json:"payload"`
	Timestamp time.Time                          `json:"timestamp"`
}

type SumupTransactionWebhookPayloadData struct {
	ClientTransactionID uuid.UUID         `json:"client_transaction_id"`
	MerchantCode        string            `json:"merchant_code"`
	Status              TransactionStatus `json:"status"`
	TransactionID       *uuid.UUID        `json:"transaction_id"`
}
