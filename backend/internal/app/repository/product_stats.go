package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

func (repo *Repository) GetProductStats() ([]models.ProductStats, error) {
	var products []models.ProductStats

	query := repo.db.Table("products").
		Select("products.id, products.name, SUM(purchase_items.quantity) as sold_items, SUM(purchase_items.total_price) as total_price").
		Joins("LEFT JOIN purchase_items ON products.id = purchase_items.product_id AND purchase_items.deleted_at IS NULL").
		Where("products.deleted_at IS NULL").
		Group("products.id").
		Order("products.pos ASC")

	if err := query.Scan(&products).Error; err != nil {
		return nil, errors.New("Products not found")
	}

	return products, nil
}
