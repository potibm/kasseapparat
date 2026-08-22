package response

import "github.com/shopspring/decimal"

type PaymentMethodStats struct {
	ID              string          `json:"id"`
	PaymentMethod   string          `json:"paymentMethod"`
	Name            string          `json:"name"`
	PurchaseCount   int             `json:"purchaseCount"`
	TotalNetPrice   decimal.Decimal `json:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal `json:"totalGrossPrice"`
}
