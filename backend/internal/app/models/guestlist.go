package models

// Guestlist represents a list of guests
type Guestlist struct {
	GormOwnedModel
	Name      string  `json:"name"`
	TypeCode  bool    `gorm:"default:false" json:"typeCode"`
	ProductID uint    `json:"productId"`
	Product   Product `gorm:"" json:"product"`
}
