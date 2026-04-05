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

func isValidPaymentMethod(code models.PaymentMethod) bool {
	_, exists := allAvailablePaymentMethods[code]

	return exists
}
