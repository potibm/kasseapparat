package config

import (
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadPaymentMethodsWithValidAndInvalid(t *testing.T) {
	_ = os.Setenv("PAYMENT_METHODS", "CASH,SUMUP,INVALID")

	result := loadPaymentMethods()

	assert.Len(t, result, 2)
	assert.Equal(t, models.PaymentMethod("CASH"), result[0].Code)
	assert.Contains(t, result[0].Name, "Cash")
	assert.Equal(t, models.PaymentMethod("SUMUP"), result[1].Code)
	assert.Contains(t, result[1].Name, "Sumup")
}

func TestLoadPaymentMethodsWithOrderAndDuplicates(t *testing.T) {
	_ = os.Setenv("PAYMENT_METHODS", "VOUCHER,CASH,SUMUP,CASH")

	result := loadPaymentMethods()

	assert.Len(t, result, 3)
	assert.Equal(t, models.PaymentMethod("VOUCHER"), result[0].Code)
	assert.Equal(t, models.PaymentMethod("CASH"), result[1].Code)
	assert.Equal(t, models.PaymentMethod("SUMUP"), result[2].Code)
}

func TestLoadPaymentMethodsWithEmptyEnvUsesDefault(t *testing.T) {
	_ = os.Unsetenv("PAYMENT_METHODS") // or Set to empty

	result := loadPaymentMethods()

	assert.Len(t, result, 1)
	assert.Equal(t, models.PaymentMethod("CASH"), result[0].Code)
}

func TestLoadPaymentMethodsWithInvalidEnvUsesDefault(t *testing.T) {
	_ = os.Setenv("PAYMENT_METHODS", "WURST,HANS")

	result := loadPaymentMethods()

	assert.Len(t, result, 1)
	assert.Equal(t, models.PaymentMethod("CASH"), result[0].Code)
}

func TestPaymentMethodsContains(t *testing.T) {
	paymentMethods := PaymentMethods{
		{Code: models.PaymentMethodCash, Name: "Cash"},
		{Code: models.PaymentMethodCC, Name: "Credit Card"},
	}

	assert.True(t, paymentMethods.Contains(models.PaymentMethodCash))
	assert.False(t, paymentMethods.Contains(models.PaymentMethodVoucher))
}

func TestPaymentMethodsGetName(t *testing.T) {
	paymentMethods := PaymentMethods{
		{Code: models.PaymentMethodCash, Name: "Cash"},
		{Code: models.PaymentMethodCC, Name: "Credit Card"},
	}

	name := paymentMethods.GetName(models.PaymentMethodCash)
	assert.NotNil(t, name)
	assert.Equal(t, "Cash", *name)

	name = paymentMethods.GetName(models.PaymentMethodVoucher)
	assert.Nil(t, name)
}
