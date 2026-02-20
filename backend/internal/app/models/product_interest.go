package models

type ProductInterest struct {
	GormOwnedModel

	ProductID uint    `json:"productID"`
	Product   Product `json:"product"   gorm:""`
}
