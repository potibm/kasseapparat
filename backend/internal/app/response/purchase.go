package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

type PurchaseResponse struct {
	ID                       uuid.UUID              `json:"id"`
	CreatedAt                time.Time              `json:"createdAt"`
	CreatedByID              *uint                  `json:"createdById"`
	CreatedBy                *models.User           `json:"createdBy"`
	PaymentMethod            models.PaymentMethod   `json:"paymentMethod"`
	TotalNetPrice            decimal.Decimal        `json:"totalNetPrice"`
	SumupTransactionID       uuid.UUID              `json:"sumupTransactionId,omitempty"`
	SumupClientTransactionID uuid.UUID              `json:"sumupClientTransactionId,omitempty"`
	TotalGrossPrice          decimal.Decimal        `json:"totalGrossPrice"`
	TotalVatAmount           decimal.Decimal        `json:"totalVatAmount"`
	PurchaseItems            []PurchaseItemResponse `json:"purchaseItems"`
	Status                   string                 `json:"status"`
}

func ToPurchaseResponse(purchase models.Purchase, decimalPlaces int32) PurchaseResponse {
	response := PurchaseResponse{
		ID:                       purchase.ID,
		CreatedAt:                purchase.CreatedAt,
		CreatedByID:              purchase.CreatedByID,
		CreatedBy:                purchase.CreatedBy,
		PaymentMethod:            purchase.PaymentMethod,
		TotalNetPrice:            purchase.TotalNetPrice,
		TotalGrossPrice:          purchase.TotalGrossPrice,
		TotalVatAmount:           purchase.TotalGrossPrice.Sub(purchase.TotalNetPrice),
		PurchaseItems:            ToPurchaseItemsResponse(purchase.PurchaseItems, decimalPlaces),
		Status:                   string(purchase.Status),
		SumupTransactionID:       uuid.Nil,
		SumupClientTransactionID: uuid.Nil,
	}

	if purchase.SumupTransactionID != nil {
		response.SumupTransactionID = *purchase.SumupTransactionID
	}

	if purchase.SumupClientTransactionID != nil {
		response.SumupClientTransactionID = *purchase.SumupClientTransactionID
	}

	return response
}

func ToPurchasesResponse(purchases []models.Purchase, decimalPlaces int32) []PurchaseResponse {
	responses := make([]PurchaseResponse, 0, len(purchases))

	for _, purchase := range purchases {
		responses = append(responses, ToPurchaseResponse(purchase, decimalPlaces))
	}

	return responses
}
