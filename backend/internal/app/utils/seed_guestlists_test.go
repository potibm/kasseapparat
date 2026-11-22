package utils

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// crewGuestCount seeds a fresh database and counts the guests on the crew list of the
// free ticket; it fails if that list is missing.
func crewGuestCount(t *testing.T, name string, withTestData bool) int64 {
	t.Helper()
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase(name)
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	require.NoError(t, SeedDatabase(db, withTestData))

	var free models.Product
	require.NoError(t, db.Where("name = ?", "🎟️ Free").First(&free).Error)

	var crew models.Guestlist
	require.NoError(t, db.Where("product_id = ? AND name = ?", free.ID, "Crew").First(&crew).Error)

	var guests int64
	require.NoError(t, db.Model(&models.Guest{}).Where("guestlist_id = ?", crew.ID).Count(&guests).Error)

	return guests
}

// The free ticket is a guest list product on the daily seeded instance too, which
// runs without test data.
func TestSeed_FreeTicketHasACrewGuestlist(t *testing.T) {
	t.Run("without test data", func(t *testing.T) {
		assert.Zero(t, crewGuestCount(t, "seedcrew", false))
	})

	t.Run("with test data", func(t *testing.T) {
		assert.Equal(t, int64(5), crewGuestCount(t, "seedcrewtest", true))
	})
}
