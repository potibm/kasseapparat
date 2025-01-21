package models

import "github.com/shopspring/decimal"

// Product represents a product model
type Product struct {
	GormOwnedModel
	Name                string          `json:"name"`
	Price               decimal.Decimal `gorm:"type:TEXT" json:"price"`
	WrapAfter           bool            `gorm:"default:false" json:"wrapAfter"`
	Hidden              bool            `gorm:"default:false" json:"hidden"`
	SoldOut             bool            `gorm:"default:false" json:"soldOut"`
	ApiExport           bool            `gorm:"default:false" json:"apiExport"`
	Pos                 int             `json:"pos"`
	TotalStock          int             `gorm:"default:0" json:"totalStock"`
	UnitsSold           int             `gorm:"default:0" json:"unitsSold"`
	SoldOutRequestCount int             `gorm:"default:0" json:"soldOutRequestCount"`
	Guestlists          []Guestlist     `json:"guestlists"`
}

type ProductWithSalesAndInterrest struct {
	Product
	UnitsSold           int `json:"unitsSold"`
	SoldOutRequestCount int `json:"soldOutRequestCount"`
}

type ProductStats struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	SoldItems  int             `json:"soldItems"`
	TotalPrice decimal.Decimal `json:"totalPrice"`
}
