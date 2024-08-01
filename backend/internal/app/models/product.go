package models

// Product represents a product model
type Product struct {
	GormOwnedModel
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	WrapAfter bool    `gorm:"default:false" json:"wrapAfter"`
	Hidden    bool    `gorm:"default:false" json:"hidden"`
	Pos       int     `json:"pos"`
	ApiExport bool    `gorm:"default:false" json:"apiExport"`
	Lists     []List  `json:"lists"`
}

type ProductStats struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	SoldItems  int     `json:"soldItems"`
	TotalPrice float64 `json:"totalPrice"`
}
