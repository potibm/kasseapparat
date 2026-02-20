package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PurchaseItem struct {
	GormModel

	PurchaseID uuid.UUID       `json:"purchaseID" gorm:"type:text"` // Foreign key to Purchase
	Purchase   Purchase        `json:"-"          gorm:"foreignKey:PurchaseID"`
	ProductID  uint            `json:"productID"` // Foreign key to Product
	Product    Product         `json:"product"    gorm:"foreignKey:ProductID"`
	Quantity   int             `json:"quantity"`
	NetPrice   decimal.Decimal `json:"netPrice"   gorm:"type:TEXT"`
	VATRate    decimal.Decimal `json:"vatRate"    gorm:"type:TEXT"`
}

func (pi PurchaseItem) GrossPrice(decimalPlaces int32) decimal.Decimal {
	return pi.NetPrice.Add(pi.VATAmount(decimalPlaces)).Round(decimalPlaces)
}

func (pi PurchaseItem) VATAmount(decimalPlaces int32) decimal.Decimal {
	return pi.NetPrice.Mul(pi.vatRateAsPercentage()).Round(decimalPlaces)
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

func (pi PurchaseItem) getQuantityAsDecimal() decimal.Decimal {
	return decimal.NewFromInt(int64(pi.Quantity))
}

func (pi PurchaseItem) vatRateAsPercentage() decimal.Decimal {
	const hundred = 100

	return pi.VATRate.Div(decimal.NewFromInt(hundred))
}
