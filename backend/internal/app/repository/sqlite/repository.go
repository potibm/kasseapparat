package sqlite

import (
	"gorm.io/gorm"
)

const whereIDEquals = "id = ?"

type Repository struct {
	db            *gorm.DB
	decimalPlaces int32
}

func NewRepository(db *gorm.DB, decimalPlaces int32) *Repository {
	return &Repository{db: db, decimalPlaces: decimalPlaces}
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}
