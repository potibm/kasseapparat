package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

const maxProductsForStats = 1000

func (handler *Handler) GetHourlyRevenueStats(c *gin.Context) {
	stats, err := handler.repo.GetHourlyRevenueStats()
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	paymentMethodNames := make(map[string]string)
	for _, pm := range handler.config.PaymentMethods {
		paymentMethodNames[string(pm.Code)] = pm.Name
	}

	for i := range stats {
		if name, ok := paymentMethodNames[stats[i].PaymentMethod]; ok {
			stats[i].Name = name
		} else {
			stats[i].Name = stats[i].PaymentMethod
		}
	}

	c.Header("X-Total-Count", strconv.Itoa(len(stats)))
	c.JSON(http.StatusOK, stats)
}

func (handler *Handler) GetHourlyQuantityStats(c *gin.Context) {
	stats, err := handler.repo.GetHourlyQuantityStats()
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	productNames := make(map[int]string)

	products, _ := handler.repo.GetProducts(maxProductsForStats, 0, "name", "ASC", sqliteRepo.ProductFilters{})
	for _, p := range products {
		productNames[p.ID] = p.Name
	}

	for i := range stats {
		if name, ok := productNames[stats[i].ProductID]; ok {
			stats[i].ProductName = name
		} else {
			stats[i].ProductName = "Unknown"
		}
	}

	c.Header("X-Total-Count", strconv.Itoa(len(stats)))
	c.JSON(http.StatusOK, stats)
}
