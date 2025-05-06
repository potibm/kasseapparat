package models

import "github.com/shopspring/decimal"

type Purchase struct {
	GormOwnedModel
	TotalNetPrice   decimal.Decimal `gorm:"type:TEXT"             json:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal `gorm:"type:TEXT"             json:"totalGrossPrice"`
	PurchaseItems   []PurchaseItem  `gorm:"foreignKey:PurchaseID" json:"purchaseItems"`
	PaymentMethod   string          `gorm:"type:TEXT"             json:"paymentMethod"`
}
