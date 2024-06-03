package models

type Purchase struct {
	GormOwnedModel
	TotalPrice    float64        `json:"totalPrice"`
	PurchaseItems []PurchaseItem `gorm:"foreignKey:PurchaseID" json:"purchaseItems"`
}
