package http

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

func queryArrayInt(c *gin.Context, field string) []int {
	idStrings := c.QueryArray(field)

	ids := make([]int, 0, len(idStrings))

	for _, s := range idStrings {
		id, err := strconv.Atoi(s)
		if err != nil {
			slog.Warn("Error converting string to int", "value", s, "error", err)

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

func queryPaymentMethods(
	c *gin.Context,
	field string,
	validPaymentMethods config.PaymentMethods,
) []models.PaymentMethod {
	paymentMethods := c.DefaultQuery(field, "")

	result := make([]models.PaymentMethod, 0)

	paymentMethodsArray := strings.SplitSeq(paymentMethods, ",")
	for code := range paymentMethodsArray {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}

		if validPaymentMethods.Contains(models.PaymentMethod(code)) {
			result = append(result, models.PaymentMethod(code))
		}
	}

	return result
}

func queryPurchaseStatusList(c *gin.Context, field string) *models.PurchaseStatusList {
	status := c.QueryArray(field)

	if len(status) == 0 {
		return nil
	}

	statusMapper := map[string]models.PurchaseStatus{
		"pending":   models.PurchaseStatusPending,
		"confirmed": models.PurchaseStatusConfirmed,
		"failed":    models.PurchaseStatusFailed,
		"cancelled": models.PurchaseStatusCancelled,
		"refunded":  models.PurchaseStatusRefunded,
	}

	statusList := make(models.PurchaseStatusList, 0, len(status))

	for _, s := range status {
		if purchaseStatus, ok := statusMapper[strings.ToLower(s)]; ok {
			statusList = append(statusList, purchaseStatus)
		}
	}

	if len(statusList) == 0 {
		return nil
	}

	return &statusList
}

func (handler *Handler) IsValidPaymentMethod(code models.PaymentMethod) bool {
	// Check if the payment method code is valid
	return handler.config.PaymentMethods.Contains(code)
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
