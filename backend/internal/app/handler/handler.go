package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/repository"
	"github.com/shopspring/decimal"
)

type Handler struct {
	repo           *repository.Repository
	mailer         mailer.Mailer
	version        string
	decimalPlaces  int32
	paymentMethods map[string]string
}

func NewHandler(repo *repository.Repository, mailer mailer.Mailer, version string, decimalPlaces int32, paymentMethods map[string]string) *Handler {
	return &Handler{repo: repo, mailer: mailer, version: version, decimalPlaces: decimalPlaces, paymentMethods: paymentMethods}
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

func queryPaymentMethods(c *gin.Context, field string, validPaymentMethods  map[string]string) []string {
	paymentMethods := c.DefaultQuery(field,"")
	
	result := make([]string, 0)

	paymentMethodsArray := strings.Split(paymentMethods, ",")
	for _, code := range paymentMethodsArray {
		if code == "" {
			continue
		}
		if _, ok := validPaymentMethods[code]; ok {
			result = append(result, code)
		}
	}
	
	return result
}

func (handler *Handler) IsValidPaymentMethod(code string) bool {
	if _, ok := handler.paymentMethods[code]; ok {
		return true
	}

	return false
}
