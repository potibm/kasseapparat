package models

// List represents a guestlist
type List struct {
	GormOwnedModel
	Name      string  `json:"name"`
	TypeCode	bool 	`gorm:"default:false" json:"typeCode"`
}