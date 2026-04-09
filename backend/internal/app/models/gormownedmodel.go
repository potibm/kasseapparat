package models

type GormOwnedModel struct {
	GormModel

	CreatedByID *int  `json:"createdById"`
	CreatedBy   *User `json:"createdBy"`
	UpdatedByID *int  `json:"updatedById"`
	UpdatedBy   *User `json:"updatedBy"`
	DeletedByID *int  `json:"deletedById"`
	DeletedBy   *User `json:"deletedBy"`
} // @name models.gormOwnedModel
