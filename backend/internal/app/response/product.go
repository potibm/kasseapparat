package response

import (
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

type ProductResponse struct {
	ID                  uint               `json:"id"`
	Name                string             `json:"name"`
	NetPrice            decimal.Decimal    `json:"netPrice"`
	GrossPrice          decimal.Decimal    `json:"grossPrice"`
	VATRate             decimal.Decimal    `json:"vatRate"`
	VATAmount           decimal.Decimal    `json:"vatAmount"`
	WrapAfter           bool               `json:"wrapAfter"`
	Hidden              bool               `json:"hidden"`
	SoldOut             bool               `json:"soldOut"`
	ApiExport           bool               `json:"apiExport"`
	Pos                 int                `json:"pos"`
	TotalStock          int                `json:"totalStock"`
	UnitsSold           int                `json:"unitsSold"`
	SoldOutRequestCount int                `json:"soldOutRequestCount"`
	Guestlists          []models.Guestlist `json:"guestlists"`
}

type ExtendedProductResponse struct {
	ProductResponse
	UnitsSold           int `json:"unitsSold"`
	SoldOutRequestCount int `json:"soldOutRequestCount"`
}

func ToProductResponse(product models.Product, decimalPlaces int32) ProductResponse {
	response := ProductResponse{
		ID:                  product.ID,
		Name:                product.Name,
		NetPrice:            product.NetPrice,
		GrossPrice:          product.GrossPrice(decimalPlaces),
		VATRate:             product.VATRate,
		VATAmount:           product.VATAmount(decimalPlaces),
		WrapAfter:           product.WrapAfter,
		Hidden:              product.Hidden,
		SoldOut:             product.SoldOut,
		ApiExport:           product.ApiExport,
		Pos:                 product.Pos,
		TotalStock:          product.TotalStock,
		UnitsSold:           product.UnitsSold,
		SoldOutRequestCount: product.SoldOutRequestCount,
		Guestlists:          product.Guestlists,
	}

	return response
}

func ToExtendedProductResponse(product models.Product, unitsSold int, soldOutRequestCount int, decimalPlaces int32) ExtendedProductResponse {
	response := ExtendedProductResponse{
		ProductResponse:     ToProductResponse(product, decimalPlaces),
		UnitsSold:           unitsSold,
		SoldOutRequestCount: soldOutRequestCount,
	}

	return response
}

func ToProductResponses(products []models.Product, decimalPlaces int32) []ProductResponse {
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = ToProductResponse(product, decimalPlaces)
	}

	return productResponses
}
