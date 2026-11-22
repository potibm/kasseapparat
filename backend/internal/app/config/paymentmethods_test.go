package config

import (
	"strings"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentMethodConfig_UsesChangeCalculator(t *testing.T) {
	on, off := true, false

	tests := []struct {
		name   string
		method PaymentMethodConfig
		want   bool
	}{
		{"cash by default", PaymentMethodConfig{Code: models.PaymentMethodCash}, true},
		{"card by default", PaymentMethodConfig{Code: models.PaymentMethodCC}, false},
		{"cash switched off", PaymentMethodConfig{Code: models.PaymentMethodCash, ChangeCalculator: &off}, false},
		{"voucher switched on", PaymentMethodConfig{Code: models.PaymentMethodVoucher, ChangeCalculator: &on}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.method.UsesChangeCalculator())
		})
	}
}

func TestPaymentMethodConfig_ReadsChangeCalculatorFromYAML(t *testing.T) {
	v := viper.New()
	v.SetConfigType("yaml")

	require.NoError(t, v.ReadConfig(strings.NewReader(`
payment_methods:
  - code: CASH
    name: Cash
    change_calculator: false
  - code: VOUCHER
    name: Voucher
    change_calculator: true
  - code: CC
    name: Card
`)))

	var methods PaymentMethods
	require.NoError(t, v.UnmarshalKey("payment_methods", &methods))
	require.Len(t, methods, 3)

	assert.False(t, methods[0].UsesChangeCalculator())
	assert.True(t, methods[1].UsesChangeCalculator())
	assert.False(t, methods[2].UsesChangeCalculator())
	assert.Nil(t, methods[2].ChangeCalculator)
}
