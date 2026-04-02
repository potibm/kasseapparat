package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestQueryPaymentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create valid payment methods for testing
	validPaymentMethods := config.PaymentMethods{
		{Code: models.PaymentMethodCash, Name: "Cash"},
		{Code: models.PaymentMethodCC, Name: "Credit Card"},
	}

	// Create a request with matching query parameters
	req, _ := http.NewRequest(http.MethodGet, "/?paymentMethods=CASH,CC,INVALID", http.NoBody)
	w := httptest.NewRecorder()
	engine := gin.New()
	c := gin.CreateTestContextOnly(w, engine)
	c.Request = req

	result := queryPaymentMethods(c, "paymentMethods", validPaymentMethods)

	assert.ElementsMatch(t, []models.PaymentMethod{models.PaymentMethodCash, models.PaymentMethodCC}, result)
}
