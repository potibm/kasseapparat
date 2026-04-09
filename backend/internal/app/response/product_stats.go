package response

import "github.com/shopspring/decimal"

type ProductStats struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	SoldItems       uint            `json:"soldItems"`
	TotalNetPrice   decimal.Decimal `json:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal `json:"totalGrossPrice"`
}
