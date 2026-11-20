package models

// Section groups products on the POS under a heading. A section without a name is
// shown as a plain row break. See ADR 006.
type Section struct {
	GormOwnedModel

	Name string `json:"name" gorm:""`
	Pos  int    `json:"pos"  gorm:"default:0"`
}
