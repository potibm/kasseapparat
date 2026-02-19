package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PurchaseItem struct {
	GormModel

	PurchaseID uuid.UUID       `gorm:"type:text"             json:"purchaseID"` // Foreign key to Purchase
	Purchase   Purchase        `gorm:"foreignKey:PurchaseID" json:"-"`
	ProductID  uint            `json:"productID"` // Foreign key to Product
	Product    Product         `gorm:"foreignKey:ProductID"  json:"product"`
	Quantity   int             `json:"quantity"`
	NetPrice   decimal.Decimal `gorm:"type:TEXT"             json:"netPrice"`
	VATRate    decimal.Decimal `gorm:"type:TEXT"             json:"vatRate"`
}

func (pi PurchaseItem) GrossPrice(decimalPlaces int32) decimal.Decimal {
	return pi.NetPrice.Add(pi.VATAmount(decimalPlaces)).Round(decimalPlaces)
}

func (pi PurchaseItem) VATAmount(decimalPlaces int32) decimal.Decimal {
	return pi.NetPrice.Mul(pi.vatRateAsPercentage()).Round(decimalPlaces)
}

func (pi PurchaseItem) getQuantityAsDecimal() decimal.Decimal {
	return decimal.NewFromInt(int64(pi.Quantity))
}

func (pi PurchaseItem) TotalNetPrice(decimalPlaces int32) decimal.Decimal {
	return pi.NetPrice.Mul(pi.getQuantityAsDecimal()).Round(decimalPlaces)
}

func (pi PurchaseItem) TotalGrossPrice(decimalPlaces int32) decimal.Decimal {
	return pi.GrossPrice(decimalPlaces).Mul(pi.getQuantityAsDecimal()).Round(decimalPlaces)
}

func (pi PurchaseItem) TotalVATAmount(decimalPlaces int32) decimal.Decimal {
	return pi.VATAmount(decimalPlaces).Mul(pi.getQuantityAsDecimal()).Round(decimalPlaces)
}

func (pi PurchaseItem) vatRateAsPercentage() decimal.Decimal {
	const hundred = 100

	return pi.VATRate.Div(decimal.NewFromInt(hundred))
}
