package sqlite

import (
	"errors"
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/models"
)

var (
	ErrSectionNotFound = errors.New("section not found")
	// ErrSectionNotEmpty is returned when a section that still holds products is deleted.
	// Moving the products elsewhere would change the POS layout as a side effect of a
	// delete, so the caller has to empty the section first (ADR 006).
	ErrSectionNotEmpty = errors.New("section still contains products")
)

var sectionSortFieldMappings = map[string]string{
	"id":   "ID",
	"name": "LOWER(Name)",
	"pos":  "Pos",
}

func getSectionsValidSortFieldName(input string) (string, error) {
	if field, exists := sectionSortFieldMappings[input]; exists {
		return field, nil
	}

	return "", errors.New("invalid sort field name")
}

func (repo *Repository) GetSections(
	limit int,
	offset int,
	sort string,
	order string,
	ids []int,
) ([]models.Section, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sortField, err := getSectionsValidSortFieldName(sort)
	if err != nil {
		return nil, err
	}

	query := repo.db.Order(sortField + " " + order + ", Id ASC").Limit(limit).Offset(offset)

	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}

	var sections []models.Section
	if err := query.Find(&sections).Error; err != nil {
		return nil, fmt.Errorf("failed to load sections: %w", err)
	}

	return sections, nil
}

func (repo *Repository) GetTotalSections() (int64, error) {
	var totalRows int64
	if err := repo.db.Model(&models.Section{}).Count(&totalRows).Error; err != nil {
		return 0, err
	}

	return totalRows, nil
}

func (repo *Repository) GetSectionByID(id int) (*models.Section, error) {
	var section models.Section
	if err := repo.db.First(&section, id).Error; err != nil {
		return nil, ErrSectionNotFound
	}

	return &section, nil
}

// GetOrCreateDefaultSection returns the first section by position, and creates an
// unnamed one when there is none yet. Every product must belong to a section, so a
// product created without one goes here.
func (repo *Repository) GetOrCreateDefaultSection() (*models.Section, error) {
	var section models.Section

	result := repo.db.Order("Pos ASC, Id ASC").Limit(1).Find(&section)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to load the default section: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		return &section, nil
	}

	section = models.Section{}
	if err := repo.db.Create(&section).Error; err != nil {
		return nil, fmt.Errorf("failed to create the default section: %w", err)
	}

	return &section, nil
}

func (repo *Repository) CreateSection(section models.Section) (models.Section, error) {
	result := repo.db.Create(&section)

	return section, result.Error
}

func (repo *Repository) UpdateSectionByID(id int, updatedSection models.Section) (*models.Section, error) {
	var section models.Section
	if err := repo.db.First(&section, id).Error; err != nil {
		return nil, ErrSectionNotFound
	}

	section.Name = updatedSection.Name
	section.Pos = updatedSection.Pos

	if err := repo.db.Save(&section).Error; err != nil {
		return nil, errors.New("failed to update section")
	}

	return &section, nil
}

// DeleteSection deletes an empty section and refuses one that still holds products.
func (repo *Repository) DeleteSection(section models.Section) error {
	var productCount int64
	if err := repo.db.Model(&models.Product{}).
		Where("section_id = ?", section.ID).
		Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to count the products of the section: %w", err)
	}

	if productCount > 0 {
		return ErrSectionNotEmpty
	}

	if err := repo.db.Delete(&section).Error; err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}

	return nil
}
