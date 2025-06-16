package initializer_test

import (
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/models"

	"github.com/stretchr/testify/assert"
)

func TestGetEnabledPaymentMethodsDefaultsToCash(t *testing.T) {
	os.Unsetenv("PAYMENT_METHODS") // no value set
	t.Cleanup(func() {
		os.Unsetenv("PAYMENT_METHODS")
	})

	methods := initializer.GetEnabledPaymentMethods()

	assert.Len(t, methods, 1)
	assert.Contains(t, methods, models.PaymentMethodCash)
}

func TestGetEnabledPaymentMethodsWithValidEnv(t *testing.T) {
	os.Setenv("PAYMENT_METHODS", "CASH,CC")
	t.Cleanup(func() {
		os.Unsetenv("PAYMENT_METHODS")
	})

	methods := initializer.GetEnabledPaymentMethods()

	assert.Len(t, methods, 2)
	assert.Contains(t, methods, models.PaymentMethodCash)
	assert.Contains(t, methods, models.PaymentMethodCC)
	// Verify labels match allAvailablePaymentMethods
	assert.Equal(t, "💶 Cash", methods[models.PaymentMethodCash])
	assert.Equal(t, "💳 Creditcard", methods[models.PaymentMethodCC])
}

func TestGetEnabledPaymentMethodsIgnoresUnknown(t *testing.T) {
	os.Setenv("PAYMENT_METHODS", "CASH,BITCOIN,FOO")
	t.Cleanup(func() {
		os.Unsetenv("PAYMENT_METHODS")
	})

	methods := initializer.GetEnabledPaymentMethods()

	assert.Len(t, methods, 1)
	assert.Contains(t, methods, models.PaymentMethodCash)
	assert.NotContains(t, methods, "BITCOIN")
	assert.NotContains(t, methods, "FOO")
}
