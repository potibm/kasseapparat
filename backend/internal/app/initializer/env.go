package initializer

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/potibm/kasseapparat/internal/app/models"
)

var allAvailablePaymentMethods = map[models.PaymentMethod]string{
	models.PaymentMethodCash:    "💶 Cash",
	models.PaymentMethodCC:      "💳 Creditcard",
	models.PaymentMethodVoucher: "🎟️ Voucher",
	models.PaymentMethodSumUp:   "💳 Sumup",
}

const defaultPaymentMethod = models.PaymentMethodCash

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
func GetEnabledPaymentMethods() map[models.PaymentMethod]string {
	enabled := make(map[models.PaymentMethod]string)

	raw := os.Getenv("PAYMENT_METHODS")
	if raw == "" {
		// fallback: default to only CASH if not set
		raw = string(defaultPaymentMethod)
	}

	for _, code := range strings.Split(raw, ",") {
		pm := models.PaymentMethod(strings.TrimSpace(code))
		if label, ok := allAvailablePaymentMethods[pm]; ok {
			enabled[pm] = label
		}
	}

	if len(enabled) == 0 {
		enabled[defaultPaymentMethod] = allAvailablePaymentMethods[defaultPaymentMethod]
	}

	return enabled
}
