package order

import "github.com/shopspring/decimal"

type Order struct {
	ID         uint
	TotalPrice decimal.Decimal
	LineItens  []OrderLineItem
}
