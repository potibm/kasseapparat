package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

func (repo *Repository) GetProductInterrests(limit int, offset int, ids []int) ([]models.ProductInterrest, error) {

	query := repo.db.Preload("Product").Order("created_at DESC").Limit(limit).Offset(offset)

	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}

	var productInterrests []models.ProductInterrest
	if err := query.Find(&productInterrests).Error; err != nil {
		return nil, errors.New("ProductInterrests not found")
	}

	return productInterrests, nil
}

func (repo *Repository) GetTotalroductInterrests() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.ProductInterrest{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetProductInterrestByID(id int) (*models.ProductInterrest, error) {
	var productInterrest models.ProductInterrest
	if err := repo.db.Preload("Product").First(&productInterrest, id).Error; err != nil {
		return nil, errors.New("ProductInterrest not found")
	}

	return &productInterrest, nil
}

func (repo *Repository) DeleteProductInterrest(productInterrest models.ProductInterrest, deletedBy models.User) {
	repo.db.Model(&models.ProductInterrest{}).Where("id = ?", productInterrest.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&productInterrest)
}

func (repo *Repository) CreateProductInterrest(productInterrest models.ProductInterrest, createdBy models.User) (models.ProductInterrest, error) {
	productInterrest.CreatedByID = &createdBy.ID

	result := repo.db.Create(&productInterrest)

	return productInterrest, result.Error
}
