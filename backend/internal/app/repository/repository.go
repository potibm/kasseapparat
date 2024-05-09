package repository

import (
	"github.com/potibm/die-kassa/internal/app/utils"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	db := utils.ConnectToDatabase()
	return &Repository{db: db}
}

func NewLocalRepository() *Repository {
	db := utils.ConnectToLocalDatabase()
	return &Repository{db: db}
}
