package initializer

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var allAvailablePaymentMethods = map[string]string{
	"CASH":    "💶 Cash",
	"CC":      "💳 Creditcard",
	"VOUCHER": "🎟️ Voucher",
}

const defaultPaymentMethod = "CASH"

func InitializeDotenv() {
	_ = godotenv.Load()
}

func GetCurrencyDecimalPlaces() int32 {
	fractionDigitsMax := 2 // Default value

	if value, exists := os.LookupEnv("FRACTION_DIGITS_MAX"); exists {
		if parsedValue, err := strconv.Atoi(value); err == nil {
			fractionDigitsMax = parsedValue
		}
	}

	return int32(fractionDigitsMax)
}
func GetEnabledPaymentMethods() map[string]string {
	enabled := make(map[string]string)

	raw := os.Getenv("PAYMENT_METHODS")
	if raw == "" {
		// fallback: default to only CASH if not set
		raw = defaultPaymentMethod
	}

	for _, code := range strings.Split(raw, ",") {
		code = strings.TrimSpace(code)
		if label, ok := allAvailablePaymentMethods[code]; ok {
			enabled[code] = label
		}
	}

	if len(enabled) == 0 {
		enabled[defaultPaymentMethod] = allAvailablePaymentMethods[defaultPaymentMethod]
	}

	return enabled
}
