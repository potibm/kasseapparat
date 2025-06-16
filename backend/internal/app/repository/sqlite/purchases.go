package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductPurchaseStats struct {
	ProductID int
	Quantity  int
	Name      string
}

type PurchaseFilters struct {
	CreatedByID        int
	PaymentMethods     []models.PaymentMethod
	Status             *models.PurchaseStatus
	TotalGrossPriceLte *decimal.Decimal
	TotalGrossPriceGte *decimal.Decimal
	IDs                []int
}

func (filters PurchaseFilters) AddWhere(query *gorm.DB) *gorm.DB {
	if len(filters.IDs) > 0 {
		query = query.Where("purchases.ID IN ?", filters.IDs)
	}

	if filters.CreatedByID != 0 {
		query = query.Where("purchases.created_by_id = ?", filters.CreatedByID)
	}

	if len(filters.PaymentMethods) > 0 {
		query = query.Where("purchases.payment_method IN ?", filters.PaymentMethods)
	}

	if filters.TotalGrossPriceLte != nil {
		query = query.Where("purchases.total_gross_price <= ?", filters.TotalGrossPriceLte)
	}

	if filters.TotalGrossPriceGte != nil {
		query = query.Where("purchases.total_gross_price >= ?", filters.TotalGrossPriceGte)
	}

	if filters.Status != nil {
		query = query.Where("purchases.status = ?", *filters.Status)
	}

	return query
}

var purchaseSortFieldMappings = map[string]string{
	"id":                 "purchases.ID",
	"createdAt":          "purchases.created_at",
	"totalGrossPrice":    "purchases.total_gross_price",
	"createdBy.username": "CreatedBy.username",
	"paymentMethod":      "purchases.payment_method",
	"status":             "purchases.status",
	"pos":                "Pos",
}

func (repo *Repository) StorePurchases(purchase models.Purchase) (models.Purchase, error) {
	return repo.StorePurchasesTx(repo.db, purchase)
}

func (repo *Repository) StorePurchasesTx(tx *gorm.DB, purchase models.Purchase) (models.Purchase, error) {
	result := tx.Create(&purchase)
	return purchase, result.Error
}

func (repo *Repository) DeletePurchaseByID(id uuid.UUID, deletedBy models.User) {
	repo.db.Model(&models.Purchase{}).Where(whereIDEquals, id).Update("DeletedByID", deletedBy.ID)
	repo.db.Where(whereIDEquals, id).Delete(&models.Purchase{})

	repo.db.Where("purchase_id = ?", id).Delete(&models.PurchaseItem{})
}

func (repo *Repository) GetPurchaseByID(id uuid.UUID) (*models.Purchase, error) {
	return repo.GetPurchaseByIDTx(repo.db, id)
}

func (repo *Repository) GetPurchaseByIDTx(tx *gorm.DB, id uuid.UUID) (*models.Purchase, error) {
	var purchase models.Purchase
	if err := tx.Model(&models.Purchase{}).
		Preload("PurchaseItems").
		Preload("PurchaseItems.Product").
		Where(whereIDEquals, id.String()).
		First(&purchase).
		Error; err != nil {
		return nil, errors.New("purchase not found")
	}

	return &purchase, nil
}

func (repo *Repository) UpdatePurchaseStatusByIDTx(tx *gorm.DB, id uuid.UUID, status models.PurchaseStatus) (*models.Purchase, error) {
	return repo.updatePurchaseFieldByIDTx(tx, id, map[string]interface{}{
		"status": string(status),
	})
}

func (repo *Repository) UpdatePurchaseSumupTransactionIDByIDTx(tx *gorm.DB, id uuid.UUID, sumupTransactionID uuid.UUID) (*models.Purchase, error) {
	return repo.updatePurchaseFieldByIDTx(tx, id, map[string]interface{}{
		"sumup_transaction_id": sumupTransactionID.String(),
	})
}

func (repo *Repository) UpdatePurchaseSumupClientTransactionIDByIDTx(tx *gorm.DB, id uuid.UUID, sumupClientTransactionID uuid.UUID) (*models.Purchase, error) {
	return repo.updatePurchaseFieldByIDTx(tx, id, map[string]interface{}{
		"sumup_client_transaction_id": sumupClientTransactionID.String(),
	})
}

func (repo *Repository) updatePurchaseFieldByIDTx(tx *gorm.DB, id uuid.UUID, fields map[string]interface{}) (*models.Purchase, error) {
	var purchase models.Purchase
	if err := tx.Model(&purchase).
		Where(whereIDEquals, id.String()).
		Updates(fields).
		Error; err != nil {
		fmt.Printf("failed to update purchase with ID %s: %v", id.String(), err)
		return nil, fmt.Errorf("failed to update purchase fields: %v", fields)
	}

	return repo.GetPurchaseByIDTx(tx, id)
}

func (repo *Repository) GetPurchases(limit int, offset int, sort string, order string, filters PurchaseFilters) ([]models.Purchase, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getPurchasesValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	var purchases []models.Purchase

	query := repo.db.Joins("CreatedBy").Model(&models.Purchase{}).Preload("PurchaseItems").Preload("PurchaseItems.Product").Order(sort + " " + order + ", purchases.created_at DESC").Limit(limit).Offset(offset)
	query = filters.AddWhere(query)

	if err := query.Find(&purchases).Error; err != nil {
		return nil, errors.New("purchases not found")
	}

	return purchases, nil
}

func (repo *Repository) GetFilteredPurchases(filters PurchaseFilters) ([]models.PurchaseItem, error) {
	var purchaseItems []models.PurchaseItem

	query := repo.db.
		Model(&models.PurchaseItem{}).
		Joins("JOIN purchases ON purchases.id = purchase_items.purchase_id").
		Preload("Product").
		Preload("Purchase")

	query = filters.AddWhere(query)

	if err := query.Find(&purchaseItems).Error; err != nil {
		return nil, errors.New("purchases not found")
	}

	return purchaseItems, nil
}

func getPurchasesValidFieldName(input string) (string, error) {
	if field, exists := purchaseSortFieldMappings[input]; exists {
		return field, nil
	}

	return "", errors.New("invalid sort field name")
}

func (repo *Repository) GetTotalPurchases(filters PurchaseFilters) (int64, error) {
	var totalRows int64

	query := repo.db.Model(&models.Purchase{})
	query = filters.AddWhere(query)
	query.Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetPurchaseStats() ([]ProductPurchaseStats, error) {
	var purchases []ProductPurchaseStats

	err := repo.db.
		Model(&models.PurchaseItem{}).
		Select("purchase_items.product_id, SUM(purchase_items.quantity) AS quantity, products.name").
		Joins("JOIN purchases ON purchases.id = purchase_items.purchase_id AND purchases.status = ? ", models.PurchaseStatusConfirmed).
		Joins("JOIN products ON products.id = purchase_items.product_id AND products.api_export = ?", 1).
		Where("purchase_items.deleted_at IS NULL").
		Group("purchase_items.product_id, products.name").
		Scan(&purchases).Error

	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func (repo *Repository) GetPurchasedQuantitiesByProductID(productID uint) (int, error) {
	var sum sql.NullInt64

	err := repo.db.Table("purchase_items").
		Select("SUM(quantity)").
		Joins("JOIN purchases ON (purchase_items.purchase_id = purchases.id AND purchase_items.deleted_at IS NULL)").
		Where("purchase_items.product_id = ? AND purchases.deleted_at IS NULL AND purchases.status = ?", productID, models.PurchaseStatusConfirmed).
		Scan(&sum).Error
	if err != nil {
		return 0, err
	}

	if !sum.Valid {
		return 0, nil
	}

	return int(sum.Int64), nil
}
