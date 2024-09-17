package order

import (
	"github.com/potibm/kasseapparat/internal/app/entities/product"
	"github.com/shopspring/decimal"
)

type OrderLineItem struct {
	ID         uint
	Order      Order
	Product    product.Product
	Quantity   uint
	Price      decimal.Decimal
	TotalPrice decimal.Decimal
}
