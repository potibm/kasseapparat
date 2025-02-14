package response

import (
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

type PurchaseItemResponse struct {
	ID              uint            `json:"id"`
	PurchaseID      uint            `json:"purchaseID"` // Foreign key to Purchase
	ProductID       uint            `json:"productID"`  // Foreign key to Product
	Product         ProductResponse `json:"product"`
	Quantity        int             `json:"quantity"`
	NetPrice        decimal.Decimal `json:"netPrice"`
	GrossPrice      decimal.Decimal ` json:"grossPrice"`
	VATRate         decimal.Decimal `json:"vatRate"`
	VATAmount       decimal.Decimal `json:"vatAmount"`
	TotalNetPrice   decimal.Decimal ` json:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal ` json:"totalGrossPrice"`
	TotalVATAmount  decimal.Decimal `json:"totalVatAmount"`
}

func ToPurchaseItemResponse(purchaseItem models.PurchaseItem) PurchaseItemResponse {
	response := PurchaseItemResponse{
		ID:              purchaseItem.ID,
		PurchaseID:      purchaseItem.PurchaseID,
		ProductID:       purchaseItem.ProductID,
		Product:         ToProductResponse(purchaseItem.Product),
		Quantity:        purchaseItem.Quantity,
		NetPrice:        purchaseItem.NetPrice,
		GrossPrice:      purchaseItem.GrossPrice(),
		TotalNetPrice:   purchaseItem.TotalNetPrice(),
		TotalGrossPrice: purchaseItem.TotalGrossPrice(),
		TotalVATAmount:  purchaseItem.TotalVATAmount(),
		VATRate:         purchaseItem.VATRate,
		VATAmount:       purchaseItem.VATAmount(),
	}

	return response
}

func ToPurchaseItemsResponse(purchaseItems []models.PurchaseItem) []PurchaseItemResponse {
	var responses []PurchaseItemResponse
	for _, purchaseItem := range purchaseItems {
		responses = append(responses, ToPurchaseItemResponse(purchaseItem))
	}

	return responses
}
