package repository

import (
	"errors"

	"github.com/potibm/die-kassa/internal/app/models"
)

func (repo *Repository) GetProducts() ([]models.Product, error) {
	var products []models.Product
	if err := repo.db.Find(&products).Error; err != nil {
		return nil, errors.New("Products not found")
	}

	return products, nil
}

func (repo *Repository) GetProductByID(id int) (*models.Product, error) {
	var product models.Product
	if err := repo.db.First(&product, id).Error; err != nil {
		return nil, errors.New("Product not found")
	}

	return &product, nil
}
