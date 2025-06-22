package http

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

func queryArrayInt(c *gin.Context, field string) []int {
	idStrings := c.QueryArray(field)

	var ids []int

	for _, s := range idStrings {
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("Error converting %s to int: %v", s, err)
			continue // skip invalid integers
		}

		ids = append(ids, id)
	}

	return ids
}

func queryDecimal(c *gin.Context, field string) *decimal.Decimal {
	value := c.DefaultQuery(field, "none")

	if value == "none" {
		return nil
	} else {
		decimalValue, err := decimal.NewFromString(value)
		if err != nil {
			// ignore silently
			return nil
		}

		return &decimalValue
	}
}

func queryTime(c *gin.Context, field string, defaultValue *time.Time) *time.Time {
	timeString := c.DefaultQuery(field, "")

	if timeString == "" {
		return defaultValue
	} else {
		t, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			return defaultValue
		}

		return &t
	}
}

func queryPaymentMethods(c *gin.Context, field string, validPaymentMethods map[models.PaymentMethod]string) []models.PaymentMethod {
	paymentMethods := c.DefaultQuery(field, "")

	result := make([]models.PaymentMethod, 0)

	paymentMethodsArray := strings.Split(paymentMethods, ",")
	for _, code := range paymentMethodsArray {
		code = strings.TrimSpace(string(code))
		if code == "" {
			continue
		}

		if _, ok := validPaymentMethods[models.PaymentMethod(code)]; ok {
			result = append(result, models.PaymentMethod(code))
		}
	}

	return result
}

func queryPurchaseStatus(c *gin.Context, field string) *models.PurchaseStatus {
	status := c.DefaultQuery(field, "")

	if status == "" {
		return nil
	}

	statusMapper := map[string]models.PurchaseStatus{
		"pending":   models.PurchaseStatusPending,
		"confirmed": models.PurchaseStatusConfirmed,
		"failed":    models.PurchaseStatusFailed,
		"cancelled": models.PurchaseStatusCancelled,
		"refunded":  models.PurchaseStatusRefunded,
	}

	if purchaseStatus, ok := statusMapper[strings.ToLower(status)]; ok {
		return &purchaseStatus
	}

	return nil
}

func (handler *Handler) IsValidPaymentMethod(code models.PaymentMethod) bool {
	// Check if the payment method code is valid
	if _, ok := handler.paymentMethods[code]; !ok {
		return false
	}

	return true
}

func (handler *Handler) ValidatePaymentMethodPayload(code models.PaymentMethod, sumupReaderId string) error {
	// Check if the payment method code is valid
	if !handler.IsValidPaymentMethod(code) {
		return errors.New("invalid payment method")
	}

	// If payment method is SUMUP, sumupReaderId must be provided
	if code == models.PaymentMethodSumUp && sumupReaderId == "" {
		return errors.New("the SumUp reader ID is required for SumUp payments")
	}

	return nil
}
