package models

import (
	"github.com/shopspring/decimal"
)

// Product represents a product model.
type Product struct {
	GormOwnedModel

	Name                string          `json:"name"`
	NetPrice            decimal.Decimal `json:"netPrice"            gorm:"type:TEXT"`
	VATRate             decimal.Decimal `json:"vatRate"             gorm:"type:TEXT;default:'0.0'"`
	WrapAfter           bool            `json:"wrapAfter"           gorm:"default:false"`
	Hidden              bool            `json:"hidden"              gorm:"default:false"`
	SoldOut             bool            `json:"soldOut"             gorm:"default:false"`
	ApiExport           bool            `json:"apiExport"           gorm:"default:false"`
	Pos                 int             `json:"pos"`
	TotalStock          int             `json:"totalStock"          gorm:"default:0"`
	UnitsSold           int             `json:"unitsSold"           gorm:"default:0"`
	SoldOutRequestCount int             `json:"soldOutRequestCount" gorm:"default:0"`
	Guestlists          []Guestlist     `json:"guestlists"`
}

func (p Product) GrossPrice(decimalPlaces int32) decimal.Decimal {
	return p.NetPrice.Add(p.VATAmount(decimalPlaces)).Round(decimalPlaces)
}

func (p Product) VATAmount(decimalPlaces int32) decimal.Decimal {
	const hundred = 100

	return p.NetPrice.Mul(p.VATRate.Div(decimal.NewFromInt(hundred))).Round(decimalPlaces)
}
