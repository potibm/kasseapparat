package handler

import (
	"encoding/csv"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

func (handler *Handler) ExportPurchases(c *gin.Context) {
	filters := repository.PurchaseFilters{}

	purchases, err := handler.repo.GetFilteredPurchases(filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	filename := fmt.Sprintf("purchases_%s.csv", purchases[0].CreatedAt.Format("2006-01-02"))

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	_ = writer.Write([]string{"Time", "Gross Price", "Net Price", "VAT", "Payment Method"})

	for _, p := range purchases {
		Vat := p.TotalGrossPrice.Sub(p.TotalNetPrice)

		_ = writer.Write([]string{
			fmt.Sprint(p.CreatedAt.Format("2006-01-02 15:04")),
			p.TotalGrossPrice.StringFixed(handler.decimalPlaces),
			p.TotalNetPrice.StringFixed(handler.decimalPlaces),
			Vat.StringFixed(handler.decimalPlaces),
			p.PaymentMethod,
		})
	}
}
