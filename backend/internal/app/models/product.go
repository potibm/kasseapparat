package models

import (
	"github.com/shopspring/decimal"
)

// Product represents a product model.
type Product struct {
	GormOwnedModel

	Name                string          `json:"name"`
	NetPrice            decimal.Decimal `gorm:"type:TEXT"               json:"netPrice"`
	VATRate             decimal.Decimal `gorm:"type:TEXT;default:'0.0'" json:"vatRate"`
	WrapAfter           bool            `gorm:"default:false"           json:"wrapAfter"`
	Hidden              bool            `gorm:"default:false"           json:"hidden"`
	SoldOut             bool            `gorm:"default:false"           json:"soldOut"`
	ApiExport           bool            `gorm:"default:false"           json:"apiExport"`
	Pos                 int             `json:"pos"`
	TotalStock          int             `gorm:"default:0"               json:"totalStock"`
	UnitsSold           int             `gorm:"default:0"               json:"unitsSold"`
	SoldOutRequestCount int             `gorm:"default:0"               json:"soldOutRequestCount"`
	Guestlists          []Guestlist     `json:"guestlists"`
}

func (p Product) GrossPrice(decimalPlaces int32) decimal.Decimal {
	return p.NetPrice.Add(p.VATAmount(decimalPlaces)).Round(decimalPlaces)
}

func (p Product) VATAmount(decimalPlaces int32) decimal.Decimal {
	const hundred = 100

	return p.NetPrice.Mul(p.VATRate.Div(decimal.NewFromInt(hundred))).Round(decimalPlaces)
}
