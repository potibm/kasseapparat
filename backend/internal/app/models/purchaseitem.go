package models

import "gorm.io/gorm"

type PurchaseItem struct {
	gorm.Model
	PurchaseID uint
	ProductID  uint // Foreign key to Product
	Product    Product
	Quantity   int
	Price      float64
	TotalPrice float64
}
