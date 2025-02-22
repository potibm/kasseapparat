package models

import "github.com/shopspring/decimal"

type PurchaseItem struct {
	GormModel
	PurchaseID uint            `json:"purchaseID"` // Foreign key to Purchase
	ProductID  uint            `json:"productID"`  // Foreign key to Product
	Product    Product         `json:"product"`
	Quantity   int             `json:"quantity"`
	NetPrice   decimal.Decimal `gorm:"type:TEXT"  json:"netPrice"`
	VATRate    decimal.Decimal `gorm:"type:TEXT"  json:"vatRate"`
}

func (pi PurchaseItem) GrossPrice() decimal.Decimal {
	return pi.NetPrice.Add(pi.VATAmount()).Round(getFractionDigitsMax())
}

func (pi PurchaseItem) VATAmount() decimal.Decimal {
	return pi.NetPrice.Mul(pi.vatRateAsPercentage()).Round(getFractionDigitsMax())
}

func (pi PurchaseItem) getQuantityAsDecimal() decimal.Decimal {
	return decimal.NewFromInt(int64(pi.Quantity))
}

func (pi PurchaseItem) TotalNetPrice() decimal.Decimal {
	return pi.NetPrice.Mul(pi.getQuantityAsDecimal()).Round(getFractionDigitsMax())
}

func (pi PurchaseItem) TotalGrossPrice() decimal.Decimal {
	return pi.GrossPrice().Mul(pi.getQuantityAsDecimal()).Round(getFractionDigitsMax())
}

func (pi PurchaseItem) TotalVATAmount() decimal.Decimal {
	return pi.VATAmount().Mul(pi.getQuantityAsDecimal()).Round(getFractionDigitsMax())
}

func (pi PurchaseItem) vatRateAsPercentage() decimal.Decimal {
	return pi.VATRate.Div(decimal.NewFromInt(100))
}
