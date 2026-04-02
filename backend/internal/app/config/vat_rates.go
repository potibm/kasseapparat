package config

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

type VatRateConfig struct {
	Rate float64 `json:"rate"`
	Name string  `json:"name"`
}

type VatRatesConfig []VatRateConfig

const defaultVatRateStandard = 25.0
const defaultVatRateZero = 0.0

var DefaultVatRates = VatRatesConfig{
	{Rate: defaultVatRateStandard, Name: "Standard"},
	{Rate: defaultVatRateZero, Name: "Zero rate"},
}

func loadVATRates() VatRatesConfig {
	result := determineVatRates(getEnvsWithPrefix("VAT_RATES_"))

	if len(result) == 0 {
		return DefaultVatRates
	}

	return result
}

func determineVatRates(vatRateEnv map[string]string) VatRatesConfig {
	vatRates := VatRatesConfig{}

	const maxParts = 2

	for _, value := range vatRateEnv {
		parts := strings.SplitN(value, ":", maxParts)
		if len(parts) != maxParts {
			continue
		}

		rate, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}

		vatRates = append(vatRates, VatRateConfig{
			Rate: rate,
			Name: parts[1],
		})
	}

	// sort vatRates by Rate in ascending order
	sort.Slice(vatRates, func(i, j int) bool {
		return vatRates[i].Rate < vatRates[j].Rate
	})

	return vatRates
}

func (vr *VatRatesConfig) Json() string {
	jsonData, err := json.Marshal(*vr)
	if err != nil {
		return "[]"
	}

	return string(jsonData)
}
