package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestQueryPaymentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Erstelle gültige Zahlungsmethoden
	validPaymentMethods := map[string]string{
		"CASH": "Cash",
		"CC":   "Credit Card",
		"EC":   "EC Card",
	}

	// Erstelle einen Request mit passenden Query-Parametern
	req, _ := http.NewRequest("GET", "/?paymentMethods=CASH,CC,INVALID", nil)
	w := httptest.NewRecorder()
	engine := gin.New()
	c := gin.CreateTestContextOnly(w, engine)
	c.Request = req

	// Teste die Funktion
	result := queryPaymentMethods(c, "paymentMethods", validPaymentMethods)

	assert.ElementsMatch(t, []string{"CASH", "CC"}, result)
}
