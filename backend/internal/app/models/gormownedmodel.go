package models

type GormOwnedModel struct {
	GormModel

	CreatedByID *uint `json:"createdById"`
	CreatedBy   *User `json:"createdBy"`
	UpdatedByID *uint `json:"updatedById"`
	UpdatedBy   *User `json:"updatedBy"`
	DeletedByID *uint `json:"deletedById"`
	DeletedBy   *User `json:"deletedBy"`
} // @name models.gormOwnedModel
