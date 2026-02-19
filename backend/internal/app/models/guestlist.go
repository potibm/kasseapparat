package models

// Guestlist represents a list of guests.
type Guestlist struct {
	GormOwnedModel

	Name      string  `json:"name"`
	TypeCode  bool    `json:"typeCode"  gorm:"default:false"`
	ProductID uint    `json:"productId"`
	Product   Product `json:"product"   gorm:""`
}
