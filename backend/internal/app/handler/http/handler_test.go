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

	// Erstelle gültige Zahlungsmethoden
	validPaymentMethods := config.PaymentMethods{
		{Code: models.PaymentMethodCash, Name: "Cash"},
		{Code: models.PaymentMethodCC, Name: "Credit Card"},
	}

	// Erstelle einen Request mit passenden Query-Parametern
	req, _ := http.NewRequest(http.MethodGet, "/?paymentMethods=CASH,CC,INVALID", nil)
	w := httptest.NewRecorder()
	engine := gin.New()
	c := gin.CreateTestContextOnly(w, engine)
	c.Request = req

	// Teste die Funktion
	result := queryPaymentMethods(c, "paymentMethods", validPaymentMethods)

	assert.ElementsMatch(t, []models.PaymentMethod{models.PaymentMethodCash, models.PaymentMethodCC}, result)
}
