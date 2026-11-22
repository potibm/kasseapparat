package sqlite

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGuestlistStats(t *testing.T) {
	repo := setupSectionRepository(t)

	product, err := repo.CreateProduct(models.Product{Name: "Free"})
	require.NoError(t, err)

	crew := &models.Guestlist{Name: "Crew", ProductID: product.ID}
	require.NoError(t, repo.db.Create(crew).Error)

	empty := &models.Guestlist{Name: "Empty", ProductID: product.ID}
	require.NoError(t, repo.db.Create(empty).Error)

	for _, guest := range []*models.Guest{
		{Name: "A", GuestlistID: crew.ID, AdditionalGuests: 2, AttendedGuests: 3},
		{Name: "B", GuestlistID: crew.ID, AdditionalGuests: 1},
		{Name: "C", GuestlistID: crew.ID},
	} {
		require.NoError(t, repo.db.Create(guest).Error)
	}

	stats, err := repo.GetGuestlistStats([]int{crew.ID, empty.ID})
	require.NoError(t, err)

	assert.Equal(t, GuestlistStats{Entries: 3, ExpectedGuests: 6, ArrivedGuests: 3}, stats[crew.ID])

	_, found := stats[empty.ID]
	assert.False(t, found, "a list without entries has no stats")

	none, err := repo.GetGuestlistStats(nil)
	require.NoError(t, err)
	assert.Empty(t, none)
}
