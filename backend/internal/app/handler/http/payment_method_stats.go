package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)

func (handler *Handler) GetPaymentMethodStats(c *gin.Context) {
	stats, err := handler.repo.GetPaymentMethodStats()
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

func (handler *Handler) getPaymentMethodName(code models.PaymentMethod) string {
	for _, pm := range handler.config.PaymentMethods {
		if pm.Code == code {
			return pm.Name
		}
	}

	return string(code)
}
