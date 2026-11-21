package models

// ProductGroup bundles the variants of one article, for example the sizes of a
// T-shirt, into one row of the POS. A group lives inside one section; its products
// take their section from the group. See ADR 006.
type ProductGroup struct {
	GormOwnedModel

	SectionID int    `json:"sectionId" gorm:"index;default:0"`
	Name      string `json:"name"      gorm:""`
	Pos       int    `json:"pos"       gorm:"default:0"`
}
