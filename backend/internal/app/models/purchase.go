package models

import "github.com/shopspring/decimal"

type Purchase struct {
	GormOwnedModel
	TotalPrice    decimal.Decimal `gorm:"type:TEXT" json:"totalPrice"`
	PurchaseItems []PurchaseItem  `gorm:"foreignKey:PurchaseID" json:"purchaseItems"`
}
