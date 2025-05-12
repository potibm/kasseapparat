package handler

import (
	"encoding/csv"
	"fmt"
	"strconv"

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

	_ = writer.Write([]string{"Time", "Purchase ID", "Quantity", "Product Name", "VAT Rate", "Gross Price", "Net Price", "VAT Amount", "Total Gross Price", "Total Net Price", "Total VAT Amount",
		"Purchase Gross Price", "Purchase Net Price", "Purchase VAT", "Payment Method"})

	for _, p := range purchases {
		Vat := p.Purchase.TotalGrossPrice.Sub(p.Purchase.TotalNetPrice)

		_ = writer.Write([]string{
			fmt.Sprint(p.CreatedAt.Format("2006-01-02 15:04:05")),
			strconv.Itoa(int(p.Purchase.ID)),
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
			p.Purchase.PaymentMethod,
		})
	}
}
