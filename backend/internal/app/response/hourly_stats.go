package response

import "github.com/shopspring/decimal"

type HourlyRevenueStats struct {
	ID              string          `json:"id"`
	TimeBucket      string          `json:"timeBucket"`
	PaymentMethod   string          `json:"paymentMethod"`
	Name            string          `json:"name"`
	TotalGrossPrice decimal.Decimal `json:"totalGrossPrice"`
}

type HourlyQuantityStats struct {
	ID          string `json:"id"`
	TimeBucket  string `json:"timeBucket"`
	ProductID   int    `json:"productId"`
	ProductName string `json:"productName"`
	Quantity    uint   `json:"quantity"`
}
