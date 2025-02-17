package response

import (
	"time"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

type PurchaseResponse struct {
	ID              uint                   `json:"id"`
	CreatedAt       time.Time              `json:"createdAt"`
	CreatedByID     *uint                  `json:"createdById"`
	CreatedBy       *models.User           `json:"createdBy"`
	TotalNetPrice   decimal.Decimal        `json:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal        `json:"totalGrossPrice"`
	TotalVatAmount  decimal.Decimal        `json:"totalVatAmount"`
	PurchaseItems   []PurchaseItemResponse `json:"purchaseItems"`
}

func ToPurchaseResponse(purchase models.Purchase) PurchaseResponse {
	response := PurchaseResponse{
		ID:              purchase.ID,
		CreatedAt:       purchase.CreatedAt,
		CreatedByID:     purchase.CreatedByID,
		CreatedBy:       purchase.CreatedBy,
		TotalNetPrice:   purchase.TotalNetPrice,
		TotalGrossPrice: purchase.TotalGrossPrice,
		TotalVatAmount:  purchase.TotalGrossPrice.Sub(purchase.TotalNetPrice),
		PurchaseItems:   ToPurchaseItemsResponse(purchase.PurchaseItems),
	}

	return response
}

func ToPurchasesResponse(purchases []models.Purchase) []PurchaseResponse {
	var responses []PurchaseResponse
	for _, purchase := range purchases {
		responses = append(responses, ToPurchaseResponse(purchase))
	}

	return responses
}
