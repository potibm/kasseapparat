package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestListEndpoints_RejectAnUnknownSortFieldWith400(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	endpoints := map[string]func(c *gin.Context){
		"products":   handler.GetProducts,
		"guests":     handler.GetGuests,
		"guestlists": handler.GetGuestlists,
		"purchases":  handler.GetPurchases,
		"sections":   handler.GetSections,
	}

	for name, list := range endpoints {
		t.Run(name, func(t *testing.T) {
			c, _ := ginContextWithParam(t, http.MethodGet, "/?_sort=bogus")

			list(c)

			assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
		})
	}
}

func TestListQueryError_OtherErrorsAre500(t *testing.T) {
	c, _ := ginContextWithParam(t, http.MethodGet, "/")

	listQueryError(c, errors.New("database is gone"))

	assert.Equal(t, http.StatusInternalServerError, lastErrorStatus(t, c))
}
