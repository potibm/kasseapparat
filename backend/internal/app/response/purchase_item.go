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
	GrossPrice      decimal.Decimal `json:"grossPrice"`
	VATRate         decimal.Decimal `json:"vatRate"`
	VATAmount       decimal.Decimal `json:"vatAmount"`
	TotalNetPrice   decimal.Decimal `json:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal `json:"totalGrossPrice"`
	TotalVATAmount  decimal.Decimal `json:"totalVatAmount"`
}

func ToPurchaseItemResponse(purchaseItem models.PurchaseItem, decimalPlaces int32) PurchaseItemResponse {
	response := PurchaseItemResponse{
		ID:              purchaseItem.ID,
		PurchaseID:      purchaseItem.PurchaseID,
		ProductID:       purchaseItem.ProductID,
		Product:         ToProductResponse(purchaseItem.Product, decimalPlaces),
		Quantity:        purchaseItem.Quantity,
		NetPrice:        purchaseItem.NetPrice,
		GrossPrice:      purchaseItem.GrossPrice(decimalPlaces),
		TotalNetPrice:   purchaseItem.TotalNetPrice(decimalPlaces),
		TotalGrossPrice: purchaseItem.TotalGrossPrice(decimalPlaces),
		TotalVATAmount:  purchaseItem.TotalVATAmount(decimalPlaces),
		VATRate:         purchaseItem.VATRate,
		VATAmount:       purchaseItem.VATAmount(decimalPlaces),
	}

	return response
}

func ToPurchaseItemsResponse(purchaseItems []models.PurchaseItem, decimalPlaces int32) []PurchaseItemResponse {
	var responses []PurchaseItemResponse
	for _, purchaseItem := range purchaseItems {
		responses = append(responses, ToPurchaseItemResponse(purchaseItem, decimalPlaces))
	}

	return responses
}
