package repository

import (
	"github.com/potibm/kasseapparat/internal/app/models"
)

type ProductPurchaseStats struct {
	ProductID int
	Quantity  int
	Name      string
}

func (repo *Repository) StorePurchases(purchase models.Purchase) (models.Purchase, error) {
	result := repo.db.Create(&purchase)

	return purchase, result.Error
}

func (repo *Repository) DeletePurchases(id int) {
	repo.db.Delete(&models.Purchase{}, id)

	repo.db.Where("purchase_id = ?", id).Delete(&models.PurchaseItem{})
}

func (repo *Repository) GetLastPurchases(limit int, offset int) ([]models.Purchase, error) {
	var purchases []models.Purchase
	result := repo.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&purchases)

	return purchases, result.Error
}

func (repo *Repository) GetTotalPurchases() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.Purchase{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetPurchaseStats() ([]ProductPurchaseStats, error) {
	rows, err := repo.db.Raw("SELECT pu.product_id, SUM(pu.quantity) as quantity, p.name " +
		"FROM purchase_items AS pu " +
		"JOIN products as p ON p.id = pu.product_id " +
		"WHERE pu.deleted_at IS NULL AND p.api_export = 1 " +
		"GROUP BY pu.product_id").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []ProductPurchaseStats
	for rows.Next() {
		var purchase ProductPurchaseStats
		if err := rows.Scan(&purchase.ProductID, &purchase.Quantity, &purchase.Name); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}

	return purchases, nil
}
