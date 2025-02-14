package response

import "github.com/shopspring/decimal"

type ProductStats struct {
	ID              uint            `json:"id"`
	Name            string          `json:"name"`
	SoldItems       int             `json:"soldItems"`
	TotalNetPrice   decimal.Decimal `json:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal `json:"totalGrossPrice"`
}
