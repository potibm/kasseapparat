package models

import "github.com/shopspring/decimal"

type PurchaseItem struct {
	GormModel
	PurchaseID uint            `json:"purchaseID"` // Foreign key to Purchase
	ProductID  uint            `json:"productID"`  // Foreign key to Product
	Product    Product         `json:"product"`
	Quantity   int             `json:"quantity"`
	Price      decimal.Decimal `gorm:"type:TEXT" json:"price"`
	TotalPrice decimal.Decimal `gorm:"type:TEXT" json:"totalPrice"`
}
