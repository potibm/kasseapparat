package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"github.com/shopspring/decimal"
)

func (repo *Repository) GetProductStats() ([]response.ProductStats, error) {
	var products []response.ProductStats

	query := repo.db.Table("products").
		Select("products.id, products.name, 0 as sold_items, 0 as total_net_price, 0 as total_gross_price").
		Where("products.deleted_at IS NULL").
		Group("products.id").
		Order("products.pos ASC")

	if err := query.Scan(&products).Error; err != nil {
		return nil, errors.New("Unable to retrieve the products")
	}

	for i := range products {
		var purchaseItems []models.PurchaseItem

		purchaseQuery := repo.db.Table("purchase_items").
			Select("quantity, net_price, vat_rate").
			Where("product_id = ?", products[i].ID).
			Where("deleted_at IS NULL")

		if err := purchaseQuery.Scan(&purchaseItems).Error; err != nil {
			return nil, errors.New("Unable to retrieve the purchases for this product")
		}

		products[i].TotalNetPrice = decimal.NewFromFloat(0)
		products[i].TotalGrossPrice = decimal.NewFromFloat(0)

		for j := range purchaseItems {
			products[i].SoldItems += purchaseItems[j].Quantity
			products[i].TotalNetPrice = products[i].TotalNetPrice.Add(purchaseItems[j].TotalNetPrice(repo.decimalPlaces))
			products[i].TotalGrossPrice = products[i].TotalGrossPrice.Add(purchaseItems[j].TotalGrossPrice(repo.decimalPlaces))
		}
	}

	return products, nil
}
