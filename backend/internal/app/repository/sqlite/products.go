package sqlite

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

var ErrProductNotFound = errors.New("product not found")

var productSortFieldMappings = map[string]string{
	"id":         "ID",
	"name":       "Name",
	"vatRate":    "Vat_Rate",
	"grossPrice": "Net_Price * (1+ (Vat_Rate / 100))",
	"pos":        "Pos",
}

func (repo *Repository) GetProducts(limit int, offset int, sort string, order string, ids []int) ([]models.Product, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sortField, err := getProductsValidSortFieldName(sort)
	if err != nil {
		return nil, err
	}

	var products []models.Product

	query := repo.db.Table("Products").Preload("Guestlists").Order(sortField + " " + order + ", Pos ASC, Id ASC").Limit(limit).Offset(offset)

	if len(ids) > 0 {
		query = query.Where("Id IN ?", ids)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, errors.New("products not found")
	}

	return products, nil
}

func getProductsValidSortFieldName(input string) (string, error) {
	if field, exists := productSortFieldMappings[input]; exists {
		return field, nil
	}

	return "", errors.New("invalid sort field name")
}

func (repo *Repository) GetTotalProducts() (int64, error) {
	var totalRows int64

	repo.db.Model(&models.Product{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetProductByID(id int) (*models.Product, error) {
	var product models.Product
	if err := repo.db.Table("Products").First(&product, id).Error; err != nil {
		return nil, ErrProductNotFound
	}

	return &product, nil
}

func (repo *Repository) UpdateProductByID(id int, updatedProduct models.Product) (*models.Product, error) {
	var product models.Product
	if err := repo.db.First(&product, id).Error; err != nil {
		return nil, ErrProductNotFound
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
		return nil, errors.New("failed to update product")
	}

	return &product, nil
}

func (repo *Repository) CreateProduct(product models.Product) (models.Product, error) {
	result := repo.db.Create(&product)

	return product, result.Error
}

func (repo *Repository) DeleteProduct(product models.Product, deletedBy models.User) {
	repo.db.Model(&models.Product{}).Where(whereIDEquals, product.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&product)
}
