package repository

import (
	"github.com/potibm/kasseapparat/internal/app/utils"
	"gorm.io/gorm"
)

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
