package config

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	vatRateZeroRateName = "Zero rate"
	vatRateStandardName = "Standard"
)

func TestDetermineVatRates(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected VatRatesConfig
	}{
		{
			name: "valid VAT rates",
			input: map[string]string{
				"1": "25:Standard",
				"2": "0:Zero rate",
			},
			expected: VatRatesConfig{
				{Rate: 0, Name: vatRateZeroRateName},
				{Rate: 25, Name: vatRateStandardName},
			},
		},
		{
			name: "invalid VAT rate format",
			input: map[string]string{
				"1": "invalid-format",
			},
			expected: VatRatesConfig{},
		},
		{
			name: "non-numeric VAT rate",
			input: map[string]string{
				"1": "not-a-number:Standard",
			},
			expected: VatRatesConfig{},
		},
		{
			name:     "empty input",
			input:    map[string]string{},
			expected: VatRatesConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineVatRates(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVatRatesConfigJson(t *testing.T) {
	vatRates := VatRatesConfig{
		{Rate: 25, Name: vatRateStandardName},
		{Rate: 0, Name: vatRateZeroRateName},
	}

	expectedJson := `[{"rate":25,"name":"Standard"},{"rate":0,"name":"Zero rate"}]`
	assert.Equal(t, expectedJson, vatRates.Json())

	// invalid config should return empty JSON array
	invalidConfig := VatRatesConfig{
		{Rate: math.NaN(), Name: "Invalid"},
	}

	result := invalidConfig.Json()

	assert.Equal(t, "[]", result)
}

func TestLoadVatRates(t *testing.T) {
	// Test with no environment variables set
	t.Setenv("VAT_RATES_1", "")
	t.Setenv("VAT_RATES_2", "")

	expected := DefaultVatRates
	result := loadVATRates()
	assert.Equal(t, expected, result)

	// Test with valid environment variables
	t.Setenv("VAT_RATES_1", "12:Standard")
	t.Setenv("VAT_RATES_2", "0:Zero rate")

	expected = VatRatesConfig{
		{Rate: 0, Name: vatRateZeroRateName},
		{Rate: 12, Name: vatRateStandardName},
	}
	result = loadVATRates()
	assert.Equal(t, expected, result)

	// Test with invalid environment variable format (will fallback to defaults)
	t.Setenv("VAT_RATES_1", "invalid-format")
	t.Setenv("VAT_RATES_2", "")

	expected = DefaultVatRates
	result = loadVATRates()
	assert.Equal(t, expected, result)
}
