package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

func (repo *Repository) GetProductInterests(limit int, offset int, ids []int) ([]models.ProductInterest, error) {

	query := repo.db.Preload("Product").Order("created_at DESC").Limit(limit).Offset(offset)

	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}

	var productInterests []models.ProductInterest
	if err := query.Find(&productInterests).Error; err != nil {
		return nil, errors.New("ProductInterests not found")
	}

	return productInterests, nil
}

func (repo *Repository) GetTotalroductInterests() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.ProductInterest{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetProductInterestByID(id int) (*models.ProductInterest, error) {
	var productInterest models.ProductInterest
	if err := repo.db.Preload("Product").First(&productInterest, id).Error; err != nil {
		return nil, errors.New("ProductInterest not found")
	}

	return &productInterest, nil
}

func (repo *Repository) DeleteProductInterest(productInterest models.ProductInterest, deletedBy models.User) {
	repo.db.Model(&models.ProductInterest{}).Where("id = ?", productInterest.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&productInterest)
}

func (repo *Repository) CreateProductInterest(productInterest models.ProductInterest, createdBy models.User) (models.ProductInterest, error) {
	productInterest.CreatedByID = &createdBy.ID

	result := repo.db.Create(&productInterest)

	return productInterest, result.Error
}
