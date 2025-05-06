package initializer_test

import (
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/initializer"

	"github.com/stretchr/testify/assert"
)

func TestGetEnabledPaymentMethods_DefaultsToCash(t *testing.T) {
	os.Unsetenv("PAYMENT_METHODS") // no value set

	methods := initializer.GetEnabledPaymentMethods()

	assert.Len(t, methods, 1)
	assert.Contains(t, methods, "CASH")
}

func TestGetEnabledPaymentMethods_WithValidEnv(t *testing.T) {
	os.Setenv("PAYMENT_METHODS", "CASH,CC")

	methods := initializer.GetEnabledPaymentMethods()

	assert.Len(t, methods, 2)
	assert.Contains(t, methods, "CASH")
	assert.Contains(t, methods, "CC")
}

func TestGetEnabledPaymentMethods_IgnoresUnknown(t *testing.T) {
	os.Setenv("PAYMENT_METHODS", "CASH,BITCOIN,FOO")

	methods := initializer.GetEnabledPaymentMethods()

	assert.Len(t, methods, 1)
	assert.Contains(t, methods, "CASH")
	assert.NotContains(t, methods, "BITCOIN")
	assert.NotContains(t, methods, "FOO")
}
