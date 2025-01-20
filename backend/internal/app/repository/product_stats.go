package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

func (repo *Repository) GetProductStats() ([]models.ProductStats, error) {
	var products []models.ProductStats

	query := repo.db.Table("products").
		Select("products.id, products.name, 0 as sold_items, 0 as total_price").
		Where("products.deleted_at IS NULL").
		Group("products.id").
		Order("products.pos ASC")

	if err := query.Scan(&products).Error; err != nil {
		return nil, errors.New("Unable to retrieve the products")
	}

	for i := range products {
		var purchaseItems []models.PurchaseItem
		purchaseQuery := repo.db.Table("purchase_items").
			Select("quantity, total_price").
			Where("product_id = ?", products[i].ID).
			Where("deleted_at IS NULL")

		if err := purchaseQuery.Scan(&purchaseItems).Error; err != nil {
			return nil, errors.New("Unable to retrieve the purchases for this product")
		}

		products[i].TotalPrice = decimal.NewFromFloat(0)
		for j := range purchaseItems {
			products[i].SoldItems += purchaseItems[j].Quantity
			products[i].TotalPrice = products[i].TotalPrice.Add(purchaseItems[j].TotalPrice)
		}
	}

	return products, nil
}
