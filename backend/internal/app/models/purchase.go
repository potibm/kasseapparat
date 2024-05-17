package models

type Purchase struct {
	GormModel
	TotalPrice    float64        `json:"totalPrice"`
	PurchaseItems []PurchaseItem `gorm:"foreignKey:PurchaseID" json:"purchaseItems"`
}
