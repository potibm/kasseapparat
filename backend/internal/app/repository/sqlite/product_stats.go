package sqlite

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"github.com/shopspring/decimal"
)

func (repo *Repository) GetProductStats() ([]response.ProductStats, error) {
	var products = []response.ProductStats{}

	query := repo.db.Table("products").
		Select("products.id, products.name, 0 as sold_items, 0 as total_net_price, 0 as total_gross_price").
		Where("products.deleted_at IS NULL").
		Group("products.id").
		Order("products.pos ASC")

	if err := query.Scan(&products).Error; err != nil {
		return nil, errors.New("unable to retrieve the products")
	}

	for i := range products {
		var purchaseItems []models.PurchaseItem

		purchaseQuery := repo.db.Table("purchase_items").
			Select("purchase_items.quantity, purchase_items.net_price, purchase_items.vat_rate").
			Joins("JOIN purchases ON purchases.id = purchase_items.purchase_id").
			Where("purchase_items.product_id = ?", products[i].ID).
			Where("purchases.deleted_at IS NULL").
			Where("purchases.status = ?", string(models.PurchaseStatusConfirmed)).
			Where("purchase_items.deleted_at IS NULL")

		if err := purchaseQuery.Scan(&purchaseItems).Error; err != nil {
			return nil, errors.New("unable to retrieve the purchases for this product")
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
