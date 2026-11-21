package utils

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The daily seeded test instance must show the POS in sections (ADR 006), so every
// seeded product has to land in one.
func TestSeed_AssignsEveryProductToASection(t *testing.T) {
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase("seedsections")
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	require.NoError(t, SeedDatabase(db, false))

	var sections []models.Section
	require.NoError(t, db.Order("pos ASC").Find(&sections).Error)
	require.Len(t, sections, 3)
	assert.Equal(t, "Tickets", sections[0].Name)
	assert.Equal(t, "T-Shirts", sections[1].Name)
	assert.Equal(t, "Merch", sections[2].Name)

	var withoutSection int64
	require.NoError(t, db.Model(&models.Product{}).Where("section_id = 0").Count(&withoutSection).Error)
	assert.Zero(t, withoutSection)

	sectionOf := func(name string) int {
		var product models.Product
		require.NoError(t, db.Where("name = ?", name).First(&product).Error)

		return product.SectionID
	}

	assert.Equal(t, sections[0].ID, sectionOf("🎟️ Regular"))
	assert.Equal(t, sections[0].ID, sectionOf("🎟️ Prepaid"))
	assert.Equal(t, sections[1].ID, sectionOf("👕 Male S"))
	assert.Equal(t, sections[1].ID, sectionOf("👕 Female XXL"))
	assert.Equal(t, sections[2].ID, sectionOf("☕ Coffee Mug"))
}
