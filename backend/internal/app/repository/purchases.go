package repository

import (
	"errors"

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

func (repo *Repository) DeletePurchaseByID(id int, deletedBy models.User) {
	repo.db.Model(&models.Purchase{}).Where("id = ?", id).Update("DeletedByID", deletedBy.ID)
	repo.db.Delete(&models.Purchase{}, id)

	repo.db.Where("purchase_id = ?", id).Delete(&models.PurchaseItem{})
}

func (repo *Repository) GetPurchaseByID(id int) (*models.Purchase, error) {
	var purchase models.Purchase
	if err := repo.db.Model(&models.Purchase{}).Preload("PurchaseItems").Preload("PurchaseItems.Product").First(&purchase, id).Error; err != nil {
		return nil, errors.New("Purchase not found")
	}

	return &purchase, nil
}

func (repo *Repository) GetPurchases(limit int, offset int, sort string, order string) ([]models.Purchase, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getPurchasesValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	var purchases []models.Purchase
	if err := repo.db.Model(&models.Purchase{}).Preload("CreatedBy").Order(sort + " " + order + ", created_at DESC").Limit(limit).Offset(offset).Find(&purchases).Error; err != nil {
		return nil, errors.New("Purchases not found")
	}

	return purchases, nil
}

func getPurchasesValidFieldName(input string) (string, error) {
	switch input {
	case "id":
		return "ID", nil
	case "createdAt":
		return "created_at", nil
	case "totalPrice":
		return "total_price", nil
	}

	return "", errors.New("Invalid field name")
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
