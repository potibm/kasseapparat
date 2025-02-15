package repository

import (
	"database/sql"
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

type ProductPurchaseStats struct {
	ProductID int
	Quantity  int
	Name      string
}

var purchaseSortFieldMappings = map[string]string{
	"id":                 "purchases.ID",
	"createdAt":          "purchases.created_at",
	"totalGrossPrice":    "purchases.total_gross_price",
	"createdBy.username": "CreatedBy.username",
	"pos":                "Pos",
}

func (repo *Repository) StorePurchases(purchase models.Purchase) (models.Purchase, error) {
	result := repo.db.Create(&purchase)

	return purchase, result.Error
}

func (repo *Repository) DeletePurchaseByID(id int, deletedBy models.User) {
	// rollback list entries
	repo.db.Model(&models.Guest{}).Where("purchase_id = ?", id).Updates(map[string]interface{}{"purchase_id": nil, "attended_guests": 0, "arrived_at": nil})

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
	if err := repo.db.Joins("CreatedBy").Model(&models.Purchase{}).Preload("PurchaseItems").Preload("PurchaseItems.Product").Order(sort + " " + order + ", purchases.created_at DESC").Limit(limit).Offset(offset).Find(&purchases).Error; err != nil {
		return nil, errors.New("Purchases not found")
	}

	return purchases, nil
}

func getPurchasesValidFieldName(input string) (string, error) {
	if field, exists := purchaseSortFieldMappings[input]; exists {
		return field, nil
	}

	return "", errors.New("Invalid sort field name")
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

func (repo *Repository) GetPurchasedQuantitiesByProductID(productID uint) (int, error) {
	var sum sql.NullInt64

	err := repo.db.Table("purchase_items").
		Select("SUM(quantity)").
		Joins("JOIN purchases ON purchase_items.purchase_id = purchases.id").
		Where("purchase_items.product_id = ? AND purchase_items.deleted_at IS NULL AND purchases.deleted_at IS NULL", productID).
		Scan(&sum).Error

	if err != nil {
		return 0, err
	}

	if !sum.Valid {
		return 0, nil
	}

	return int(sum.Int64), nil
}
