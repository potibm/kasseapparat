package http

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

func (handler *Handler) ExportPurchases(c *gin.Context) {
	confirmed := models.PurchaseStatusConfirmed

	filters := sqliteRepo.PurchaseFilters{}
	filters.PaymentMethods = queryPaymentMethods(c, "paymentMethods", handler.paymentMethods)
	filters.Status = &confirmed
	status := models.PurchaseStatusConfirmed
	filters.Status = &status

	purchases, err := handler.repo.GetFilteredPurchases(filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	paymentMethodsString := "all"

	if len(filters.PaymentMethods) > 0 {
		methods := make([]string, len(filters.PaymentMethods))
		for i, m := range filters.PaymentMethods {
			methods[i] = string(m)
		}

		paymentMethodsString = strings.ToLower(strings.Join(methods, "_"))
	}

	time := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("purchases_%s_%s.csv", time, paymentMethodsString)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	err = writer.Write([]string{"Time", "Purchase ID", "Quantity", "Product Name", "VAT Rate", "Gross Price", "Net Price", "VAT Amount", "Total Gross Price", "Total Net Price", "Total VAT Amount",
		"Purchase Gross Price", "Purchase Net Price", "Purchase VAT", "Payment Method"})

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Failed to write CSV header: "+err.Error()))
		return
	}

	for _, p := range purchases {
		Vat := p.Purchase.TotalGrossPrice.Sub(p.Purchase.TotalNetPrice)

		err = writer.Write([]string{
			fmt.Sprint(p.CreatedAt.Format("2006-01-02 15:04:05")),
			p.Purchase.ID.String(),
			strconv.Itoa(p.Quantity),
			p.Product.Name,
			p.VATRate.String() + "%",
			p.GrossPrice(handler.decimalPlaces).StringFixed(handler.decimalPlaces),
			p.NetPrice.StringFixed(handler.decimalPlaces),
			p.VATAmount(handler.decimalPlaces).StringFixed(handler.decimalPlaces),
			p.TotalGrossPrice(handler.decimalPlaces).StringFixed(handler.decimalPlaces),
			p.TotalNetPrice(handler.decimalPlaces).StringFixed(handler.decimalPlaces),
			p.TotalVATAmount(handler.decimalPlaces).StringFixed(handler.decimalPlaces),
			p.Purchase.TotalGrossPrice.StringFixed(handler.decimalPlaces),
			p.Purchase.TotalNetPrice.StringFixed(handler.decimalPlaces),
			Vat.StringFixed(handler.decimalPlaces),
			string(p.Purchase.PaymentMethod),
		})

		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Failed to write CSV header: "+err.Error()))
			return
		}
	}
}
