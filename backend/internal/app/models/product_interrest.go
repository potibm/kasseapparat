package models

type ProductInterrest struct {
	GormOwnedModel
	ProductID uint    `json:"productID"`
	Product   Product `gorm:"" json:"product"`
}
