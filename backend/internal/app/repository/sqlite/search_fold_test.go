package sqlite

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFoldForSearch(t *testing.T) {
	tests := map[string]string{
		"Müller":       "muller",
		"MÜLLER":       "muller",
		"Straße":       "strasse",
		"Ångström":     "angstrom",
		"Zoë Çelik":    "zoe celik",
		"Søren Łukasz": "soren lukasz",
		"plain":        "plain",
	}

	for input, want := range tests {
		assert.Equal(t, want, FoldForSearch(input), input)
	}
}

func TestGuestMatchesQuery(t *testing.T) {
	code := "ABC123"

	assert.True(t, guestMatchesQuery("Jana Müller", nil, "mul"), "diacritics in the name")
	assert.True(t, guestMatchesQuery("Jana Müller", nil, "MÜL"), "diacritics in the query")
	assert.True(t, guestMatchesQuery("Hans Straße", nil, "strasse"), "ß")
	assert.True(t, guestMatchesQuery("Anyone", &code, "abc"), "code prefix, any case")
	assert.True(t, guestMatchesQuery("Anyone", nil, "  "), "blank query")
	assert.False(t, guestMatchesQuery("Anyone", &code, "123"), "code only matches as prefix")
	assert.False(t, guestMatchesQuery("Jana Müller", nil, "meier"))
}

func TestGetUnattendedGuestsByProductID_IgnoresDiacritics(t *testing.T) {
	repo := setupSectionRepository(t)

	product, err := repo.CreateProduct(models.Product{Name: "Free"})
	require.NoError(t, err)

	list := &models.Guestlist{Name: "Crew", ProductID: product.ID}
	require.NoError(t, repo.db.Create(list).Error)

	code := "XYZ987"

	for _, guest := range []*models.Guest{
		{Name: "Jana Müller", GuestlistID: list.ID},
		{Name: "Hans Straße", GuestlistID: list.ID},
		{Name: "Kim Kasse", GuestlistID: list.ID, Code: &code},
		{Name: "Max Mulder", GuestlistID: list.ID, AttendedGuests: 1},
	} {
		require.NoError(t, repo.db.Create(guest).Error)
	}

	names := func(q string) []string {
		guests, err := repo.GetUnattendedGuestsByProductID(product.ID, q)
		require.NoError(t, err)

		result := make([]string, 0, len(guests))

		for _, guest := range guests {
			result = append(result, guest.Name)
		}

		return result
	}

	assert.Equal(t, []string{"Jana Müller"}, names("mul"), "attended guests are left out")
	assert.Equal(t, []string{"Hans Straße"}, names("STRASSE"))
	assert.Equal(t, []string{"Kim Kasse"}, names("xyz"))
	assert.Len(t, names(""), 3)
	assert.Empty(t, names("nobody"))
}
