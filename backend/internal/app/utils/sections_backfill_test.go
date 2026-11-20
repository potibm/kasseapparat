package utils

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupBackfillDatabase(t *testing.T, name string) *gorm.DB {
	t.Helper()
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase(name)
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	return db
}

func sectionIDOf(t *testing.T, db *gorm.DB, productID int) int {
	t.Helper()

	var product models.Product
	require.NoError(t, db.First(&product, productID).Error)

	return product.SectionID
}

func TestMigrate_AssignsProductsWithoutSectionToANewDefaultSection(t *testing.T) {
	db := setupBackfillDatabase(t, "backfillnew")

	// a product stored before sections existed
	product := models.Product{Name: "Mate"}
	require.NoError(t, db.Create(&product).Error)
	require.Zero(t, product.SectionID)

	require.NoError(t, MigrateDatabase(db))

	var sections []models.Section
	require.NoError(t, db.Find(&sections).Error)
	require.Len(t, sections, 1)
	assert.Empty(t, sections[0].Name, "the default section is unnamed")
	assert.Equal(t, sections[0].ID, sectionIDOf(t, db, product.ID))
}

func TestMigrate_AssignsProductsWithoutSectionToTheFirstSection(t *testing.T) {
	db := setupBackfillDatabase(t, "backfillexisting")

	merch := models.Section{Name: "Merch", Pos: 2}
	tickets := models.Section{Name: "Tickets", Pos: 1}
	require.NoError(t, db.Create(&merch).Error)
	require.NoError(t, db.Create(&tickets).Error)

	withoutSection := models.Product{Name: "Mate"}
	withSection := models.Product{Name: "Mug", SectionID: merch.ID}
	require.NoError(t, db.Create(&withoutSection).Error)
	require.NoError(t, db.Create(&withSection).Error)

	require.NoError(t, MigrateDatabase(db))

	assert.Equal(t, tickets.ID, sectionIDOf(t, db, withoutSection.ID))
	assert.Equal(t, merch.ID, sectionIDOf(t, db, withSection.ID), "assigned products are left alone")

	var count int64
	require.NoError(t, db.Model(&models.Section{}).Count(&count).Error)
	assert.Equal(t, int64(2), count, "no extra section is created")
}

func TestMigrate_CreatesNoSectionWhenEveryProductHasOne(t *testing.T) {
	db := setupBackfillDatabase(t, "backfillnoop")

	require.NoError(t, MigrateDatabase(db))

	var count int64
	require.NoError(t, db.Model(&models.Section{}).Count(&count).Error)
	assert.Zero(t, count)
}
