package models

import "gorm.io/gorm"

// Product represents a product model
type Product struct {
	gorm.Model
	Name      string
	Price     float64
	WrapAfter bool `gorm:"default:false"`
	Pos       int
	ApiExport bool `gorm:"default:false"`
}
