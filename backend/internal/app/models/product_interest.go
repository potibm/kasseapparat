package models

type ProductInterest struct {
	GormOwnedModel

	ProductID int     `json:"productID"`
	Product   Product `json:"product"   gorm:""`
}
