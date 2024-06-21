package models

// Product represents a product model
type Product struct {
	GormOwnedModel
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	WrapAfter bool    `gorm:"default:false" json:"wrapAfter"`
	Hidden	  bool    `gorm:"default:false" json:"hidden"`
	Pos       int     `json:"pos"`
	ApiExport bool    `gorm:"default:false" json:"apiExport"`
	AssociatedListID	uint `json:"associatedListId"` 
	AssociatedList 	*List `json:"associatedList"`
}
