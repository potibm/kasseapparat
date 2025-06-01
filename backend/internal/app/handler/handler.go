package handler

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/shopspring/decimal"
)

type Handler struct {
	repo            *repository.Repository
	sumupRepository sumup.RepositoryInterface
	mailer          mailer.Mailer
	version         string
	decimalPlaces   int32
	paymentMethods  map[string]string
}

func NewHandler(repo *repository.Repository, sumupRepository sumup.RepositoryInterface, mailer mailer.Mailer, version string, decimalPlaces int32, paymentMethods map[string]string) *Handler {
	return &Handler{repo: repo, sumupRepository: sumupRepository, mailer: mailer, version: version, decimalPlaces: decimalPlaces, paymentMethods: paymentMethods}
}

func queryArrayInt(c *gin.Context, field string) []int {
	idStrings := c.QueryArray(field)

	var ids []int

	for _, s := range idStrings {
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("Error converting %s to int: %v", s, err)
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

func queryPaymentMethods(c *gin.Context, field string, validPaymentMethods map[string]string) []string {
	paymentMethods := c.DefaultQuery(field, "")

	result := make([]string, 0)

	paymentMethodsArray := strings.Split(paymentMethods, ",")
	for _, code := range paymentMethodsArray {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}

		if _, ok := validPaymentMethods[code]; ok {
			result = append(result, code)
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
	}

	if purchaseStatus, ok := statusMapper[strings.ToLower(status)]; ok {
		return &purchaseStatus
	}

	return nil
}

func (handler *Handler) IsValidPaymentMethod(code string) bool {
	// Check if the payment method code is valid
	if _, ok := handler.paymentMethods[code]; !ok {
		return false
	}

	return true
}

func (handler *Handler) ValidatePaymentMethodPayload(code string, sumupReaderId string) error {
	// Check if the payment method code is valid
	if !handler.IsValidPaymentMethod(code) {
		return errors.New("invalid payment method")
	}

	// If payment method is SUMUP, sumupReaderId must be provided
	if code == "SUMUP" && sumupReaderId == "" {
		return errors.New("the SumUp reader ID is required for SumUp payments")
	}

	return nil
}
