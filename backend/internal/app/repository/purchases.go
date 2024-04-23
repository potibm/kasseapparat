package repository

import (
	"github.com/potibm/die-kassa/internal/app/models"
)

func (repo *Repository) StorePurchases(purchase models.Purchase) (models.Purchase, error) {
	result := repo.db.Create(&purchase)

	return purchase, result.Error
}

func (repo *Repository) DeletePurchases(id int) {
	repo.db.Delete(&models.Purchase{}, id)

	repo.db.Where("purchase_id = ?", id).Delete(&models.PurchaseItem{})
}
