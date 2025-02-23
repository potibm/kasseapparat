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

	assert.True(t, decimal.NewFromFloat(1.31).Equal(purchaseItem.VATAmount(2)))
	assert.True(t, decimal.NewFromFloat(20).Equal(purchaseItem.GrossPrice(2)))
	assert.True(t, decimal.NewFromFloat(6.55).Equal(purchaseItem.TotalVATAmount(2)))
	assert.True(t, decimal.NewFromFloat(100).Equal(purchaseItem.TotalGrossPrice(2)))
}
