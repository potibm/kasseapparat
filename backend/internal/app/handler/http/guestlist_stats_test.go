package http

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGuestlists_IncludesStats(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	product, err := handler.repo.CreateProduct(models.Product{Name: "Free"})
	require.NoError(t, err)

	crew, err := handler.repo.CreateGuestlist(models.Guestlist{Name: "Crew", ProductID: product.ID})
	require.NoError(t, err)

	arrived := time.Now()

	_, err = handler.repo.CreateGuest(models.Guest{
		Name: "A", GuestlistID: crew.ID, AdditionalGuests: 1, AttendedGuests: 2, ArrivedAt: &arrived,
	})
	require.NoError(t, err)

	c, w := ginContextWithParam(t, http.MethodGet, "/?_start=0&_end=10")

	handler.GetGuestlists(c)

	require.Empty(t, c.Errors)

	var lists []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &lists))

	var found map[string]any

	for _, list := range lists {
		if list["name"] == "Crew" {
			found = list
		}
	}

	require.NotNil(t, found, "the crew list is in the response")
	assert.InDelta(t, 1, found["entries"], 0)
	assert.InDelta(t, 1, found["arrivedEntries"], 0)
	assert.InDelta(t, 2, found["arrivedGuests"], 0)
}
