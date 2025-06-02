package sqlite

import (
	"github.com/potibm/kasseapparat/internal/app/utils"
	"gorm.io/gorm"
)

const whereIDEquals = "id = ?"

type Repository struct {
	db            *gorm.DB
	decimalPlaces int32
}

func NewRepository(decimalPlaces int32) *Repository {
	db := utils.ConnectToDatabase()

	return &Repository{db: db, decimalPlaces: decimalPlaces}
}

func NewLocalRepository(decimalPlaces int32) *Repository {
	db := utils.ConnectToLocalDatabase()

	return &Repository{db: db, decimalPlaces: decimalPlaces}
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}
