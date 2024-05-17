package models

// Product represents a product model
type Product struct {
	GormModel
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	WrapAfter bool    `gorm:"default:false" json:"wrapAfter"`
	Pos       int     `json:"pos"`
	ApiExport bool    `gorm:"default:false" json:"apiExport"`
}
