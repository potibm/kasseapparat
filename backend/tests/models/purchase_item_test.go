package tests_models

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseItem_VATAmount(t *testing.T) {
	netPrice := decimal.NewFromFloat(18.69)
	vatRate := decimal.NewFromFloat(7)
	purchaseItem := models.PurchaseItem{
		NetPrice: netPrice,
		VATRate:  vatRate,
		Quantity: 5,
	}

	assert.Equal(t, "1.31", purchaseItem.VATAmount().String())
	assert.Equal(t, "20", purchaseItem.GrossPrice().String())
	assert.Equal(t, "6.55", purchaseItem.TotalVATAmount().String())
	assert.Equal(t, "100", purchaseItem.TotalGrossPrice().String())

}
