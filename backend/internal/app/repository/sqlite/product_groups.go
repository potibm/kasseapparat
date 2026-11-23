package sqlite

import (
	"errors"
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/gorm"
)

var (
	ErrProductGroupNotFound = errors.New("product group not found")
	// ErrProductGroupNotEmpty is returned when a group that still holds products is
	// deleted. Ungrouping them would change the POS layout as a side effect of a
	// delete (ADR 006).
	ErrProductGroupNotEmpty = errors.New("product group still contains products")
)

var productGroupSortFieldMappings = map[string]string{
	"id":        "ID",
	"name":      "LOWER(Name)",
	"pos":       "Pos",
	"sectionId": "Section_Id",
}

func getProductGroupsValidSortFieldName(input string) (string, error) {
	if field, exists := productGroupSortFieldMappings[input]; exists {
		return field, nil
	}

	return "", ErrInvalidSortField
}

func (repo *Repository) GetProductGroups(
	limit, offset int,
	sort, order string,
	filters ProductFilters,
) ([]models.ProductGroup, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sortField, err := getProductGroupsValidSortFieldName(sort)
	if err != nil {
		return nil, err
	}

	query := filters.apply(repo.db.Model(&models.ProductGroup{})).
		Order(sortField + " " + order + ", Id ASC").
		Limit(limit).
		Offset(offset)

	var groups []models.ProductGroup
	if err := query.Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("failed to load product groups: %w", err)
	}

	return groups, nil
}

func (repo *Repository) GetTotalProductGroups(filters ProductFilters) (int64, error) {
	var totalRows int64
	if err := filters.apply(repo.db.Model(&models.ProductGroup{})).Count(&totalRows).Error; err != nil {
		return 0, err
	}

	return totalRows, nil
}

func (repo *Repository) GetProductGroupByID(id int) (*models.ProductGroup, error) {
	var group models.ProductGroup
	if err := repo.db.First(&group, id).Error; err != nil {
		return nil, ErrProductGroupNotFound
	}

	return &group, nil
}

func (repo *Repository) CreateProductGroup(group models.ProductGroup) (models.ProductGroup, error) {
	result := repo.db.Create(&group)

	return group, result.Error
}

// UpdateProductGroupByID changes a group. When the group moves to another section,
// its products move with it, so a product never sits in a group of another section.
func (repo *Repository) UpdateProductGroupByID(
	id int,
	updatedGroup models.ProductGroup,
) (*models.ProductGroup, error) {
	var group models.ProductGroup

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&group, id).Error; err != nil {
			return ErrProductGroupNotFound
		}

		group.Name = updatedGroup.Name
		group.Pos = updatedGroup.Pos
		group.SectionID = updatedGroup.SectionID

		if err := tx.Save(&group).Error; err != nil {
			return fmt.Errorf("failed to update product group: %w", err)
		}

		if err := tx.Model(&models.Product{}).
			Where("product_group_id = ?", group.ID).
			Update("section_id", group.SectionID).Error; err != nil {
			return fmt.Errorf("failed to move the products of the group: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// DeleteProductGroup deletes an empty group and refuses one that still holds products.
func (repo *Repository) DeleteProductGroup(group models.ProductGroup) error {
	var productCount int64
	if err := repo.db.Model(&models.Product{}).
		Where("product_group_id = ?", group.ID).
		Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to count the products of the group: %w", err)
	}

	if productCount > 0 {
		return ErrProductGroupNotEmpty
	}

	if err := repo.db.Delete(&group).Error; err != nil {
		return fmt.Errorf("failed to delete product group: %w", err)
	}

	return nil
}
