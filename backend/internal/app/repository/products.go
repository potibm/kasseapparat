package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

const ErrProductNotFound = "Product not found"

func (repo *Repository) GetProducts(limit int, offset int, sort string, order string, ids []int) ([]models.Product, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getProductsValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	var products []models.Product

	query := repo.db.Table("Products").Preload("Guestlists").Order(sort + " " + order + ", Pos ASC, Id ASC").Limit(limit).Offset(offset)

	if len(ids) > 0 {
		query = query.Where("Id IN ?", ids)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, errors.New("Products not found")
	}

	return products, nil
}

func getProductsValidFieldName(input string) (string, error) {
	switch input {
	case "id":
		return "ID", nil
	case "name":
		return "Name", nil
	case "vatRate":
		return "Vat_Rate", nil
	case "grossPrice":
		return "Net_Price * (1+ (Vat_Rate / 100))", nil
	case "pos":
		return "Pos", nil
	}

	return "", errors.New("Invalid field name")
}

func (repo *Repository) GetTotalProducts() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.Product{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetProductByID(id int) (*models.Product, error) {
	var product models.Product
	if err := repo.db.Table("Products").First(&product, id).Error; err != nil {
		return nil, errors.New(ErrProductNotFound)
	}

	return &product, nil
}

func (repo *Repository) UpdateProductByID(id int, updatedProduct models.Product) (*models.Product, error) {
	var product models.Product
	if err := repo.db.First(&product, id).Error; err != nil {
		return nil, errors.New(ErrProductNotFound)
	}

	// Update the product with the new values
	product.Name = updatedProduct.Name
	product.Pos = updatedProduct.Pos
	product.NetPrice = updatedProduct.NetPrice
	product.VATRate = updatedProduct.VATRate
	product.WrapAfter = updatedProduct.WrapAfter
	product.ApiExport = updatedProduct.ApiExport
	product.UpdatedByID = updatedProduct.UpdatedByID
	product.Hidden = updatedProduct.Hidden
	product.SoldOut = updatedProduct.SoldOut
	product.TotalStock = updatedProduct.TotalStock

	// Save the updated product to the database
	if err := repo.db.Save(&product).Error; err != nil {
		return nil, errors.New("Failed to update product")
	}

	return &product, nil
}

func (repo *Repository) CreateProduct(product models.Product) (models.Product, error) {
	result := repo.db.Create(&product)

	return product, result.Error
}

func (repo *Repository) DeleteProduct(product models.Product, deletedBy models.User) {
	repo.db.Model(&models.Product{}).Where("id = ?", product.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&product)
}
