package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/potibm/kasseapparat/internal/app/exitcode"
	"github.com/potibm/kasseapparat/internal/app/models"
)

var allAvailablePaymentMethods = map[models.PaymentMethod]string{
	models.PaymentMethodCash:    "💶 Cash",
	models.PaymentMethodCC:      "💳 Creditcard",
	models.PaymentMethodVoucher: "🎟️ Voucher",
	models.PaymentMethodSumUp:   "💳 Sumup",
}

const defaultPaymentMethod = models.PaymentMethodCash

type PaymentMethodConfig struct {
	Code models.PaymentMethod
	Name string
}
type PaymentMethods []PaymentMethodConfig

func (pm PaymentMethods) Contains(code models.PaymentMethod) bool {
	for _, method := range pm {
		if method.Code == code {
			return true
		}
	}

	return false
}

func (pm PaymentMethods) GetName(code models.PaymentMethod) *string {
	for _, method := range pm {
		if method.Code == code {
			return &method.Name
		}
	}

	return nil
}

func loadPaymentMethods() PaymentMethods {
	raw := getEnv("PAYMENT_METHODS", "")
	codes := strings.Split(raw, ",")

	result := make(PaymentMethods, 0)
	seen := make(map[models.PaymentMethod]bool)

	for _, code := range codes {
		code = strings.TrimSpace(code)
		pm := models.PaymentMethod(code)

		if !isValidPaymentMethod(pm) || seen[pm] {
			continue
		}

		result = append(result, createPaymentMethodConfig(pm))
		seen[pm] = true
	}

	if len(result) == 0 {
		result = append(result, createPaymentMethodConfig(defaultPaymentMethod))
	}

	return result
}

func isValidPaymentMethod(code models.PaymentMethod) bool {
	_, exists := allAvailablePaymentMethods[code]

	return exists
}

func createPaymentMethodConfig(method models.PaymentMethod) PaymentMethodConfig {
	if name, exists := allAvailablePaymentMethods[method]; exists {
		return PaymentMethodConfig{
			Code: models.PaymentMethod(method),
			Name: name,
		}
	} else {
		slog.Error("Payment method is not supported", "method", method)
		os.Exit(int(exitcode.Config))

		return PaymentMethodConfig{}
	}
}
