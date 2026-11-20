package utils

import (
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/gorm"
)

const productsWithoutSection = "section_id = 0 OR section_id IS NULL"

// assignProductsWithoutSection puts products that have no section yet into the default
// section: the first by position, or a new unnamed one if there is none. It covers
// databases created before sections existed (ADR 006), whose products would otherwise
// show up in no section at all. On a fresh installation there is nothing to do.
func assignProductsWithoutSection(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Product{}).Where(productsWithoutSection).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count products without a section: %w", err)
	}

	if count == 0 {
		return nil
	}

	var section models.Section

	result := db.Order("pos ASC, id ASC").Limit(1).Find(&section)
	if result.Error != nil {
		return fmt.Errorf("failed to load the default section: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		if err := db.Create(&section).Error; err != nil {
			return fmt.Errorf("failed to create the default section: %w", err)
		}
	}

	if err := db.Model(&models.Product{}).
		Where(productsWithoutSection).
		Update("section_id", section.ID).Error; err != nil {
		return fmt.Errorf("failed to assign products to the default section: %w", err)
	}

	return nil
}
