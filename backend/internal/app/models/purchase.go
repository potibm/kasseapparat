package models

import "gorm.io/gorm"

type Purchase struct {
	gorm.Model
	TotalPrice    float64
	PurchaseItems []PurchaseItem `gorm:"foreignKey:PurchaseID"`
}
