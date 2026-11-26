package config

import (
	"github.com/potibm/kasseapparat/internal/app/models"
)

var allAvailablePaymentMethods = map[models.PaymentMethod]string{
	models.PaymentMethodCash:    "💶 Cash",
	models.PaymentMethodCC:      "💳 Creditcard",
	models.PaymentMethodVoucher: "🎟️ Voucher",
	models.PaymentMethodSumUp:   "💳 Sumup",
}

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

// UsesChangeCalculator reports whether the POS asks for the amount received with this
// payment method. Without an explicit setting only cash does.
func (pm PaymentMethodConfig) UsesChangeCalculator() bool {
	if pm.ChangeCalculator != nil {
		return *pm.ChangeCalculator
	}

	return pm.Code == models.PaymentMethodCash
}

func isValidPaymentMethod(code models.PaymentMethod) bool {
	_, exists := allAvailablePaymentMethods[code]

	return exists
}

// DropSumupWithoutCredentials removes SUMUP from the payment methods when the SumUp
// API key or merchant code is missing: the POS would offer a button that can never
// work. It reports whether SUMUP was removed, so the caller can warn about it.
func (c *Config) DropSumupWithoutCredentials() bool {
	if c.Sumup.APIKey != "" && c.Sumup.MerchantCode != "" {
		return false
	}

	if !c.PaymentMethods.Contains(models.PaymentMethodSumUp) {
		return false
	}

	kept := make(PaymentMethods, 0, len(c.PaymentMethods))

	for _, method := range c.PaymentMethods {
		if method.Code != models.PaymentMethodSumUp {
			kept = append(kept, method)
		}
	}

	c.PaymentMethods = kept

	return true
}
