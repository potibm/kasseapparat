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

func TestConfig_DropSumupWithoutCredentials(t *testing.T) {
	withSumup := func() PaymentMethods {
		return PaymentMethods{
			{Code: models.PaymentMethodCash, Name: "Cash"},
			{Code: models.PaymentMethodSumUp, Name: "Card"},
		}
	}

	tests := []struct {
		name        string
		sumup       SumupConfig
		methods     PaymentMethods
		wantDropped bool
		wantCodes   []models.PaymentMethod
	}{
		{
			name:        "keeps SumUp with credentials",
			sumup:       SumupConfig{APIKey: "key", MerchantCode: "M1"},
			methods:     withSumup(),
			wantDropped: false,
			wantCodes:   []models.PaymentMethod{models.PaymentMethodCash, models.PaymentMethodSumUp},
		},
		{
			name:        "drops SumUp without an API key",
			sumup:       SumupConfig{MerchantCode: "M1"},
			methods:     withSumup(),
			wantDropped: true,
			wantCodes:   []models.PaymentMethod{models.PaymentMethodCash},
		},
		{
			name:        "drops SumUp without a merchant code",
			sumup:       SumupConfig{APIKey: "key"},
			methods:     withSumup(),
			wantDropped: true,
			wantCodes:   []models.PaymentMethod{models.PaymentMethodCash},
		},
		{
			name:        "leaves other methods alone",
			methods:     PaymentMethods{{Code: models.PaymentMethodCash, Name: "Cash"}},
			wantDropped: false,
			wantCodes:   []models.PaymentMethod{models.PaymentMethodCash},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Sumup: tt.sumup, PaymentMethods: tt.methods}

			assert.Equal(t, tt.wantDropped, cfg.DropSumupWithoutCredentials())

			codes := make([]models.PaymentMethod, 0, len(cfg.PaymentMethods))

			for _, method := range cfg.PaymentMethods {
				codes = append(codes, method.Code)
			}

			assert.Equal(t, tt.wantCodes, codes)
		})
	}
}
