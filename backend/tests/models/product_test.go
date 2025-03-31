package tests_models

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProductVATAmount(t *testing.T) {
	netPrice := decimal.NewFromFloat(100.0)
	vatRate := decimal.NewFromFloat(20) // 20% VAT
	product := models.Product{
		NetPrice: netPrice,
		VATRate:  vatRate,
	}

	expectedVATAmount := decimal.NewFromFloat(20.0)
	actualVATAmount := product.VATAmount(2)

	assert.True(t, expectedVATAmount.Equal(actualVATAmount), "Expected VAT amount to be %s, but got %s", expectedVATAmount, actualVATAmount)
}

func TestProductGrossPrice(t *testing.T) {
	netPrice := decimal.NewFromFloat(100.0)
	vatRate := decimal.NewFromFloat(20) // 20% VAT
	product := models.Product{
		NetPrice: netPrice,
		VATRate:  vatRate,
	}

	expectedGrossPrice := decimal.NewFromFloat(120.0)
	actualGrossPrice := product.GrossPrice(2)

	assert.True(t, expectedGrossPrice.Equal(actualGrossPrice), "Expected gross price to be %s, but got %s", expectedGrossPrice, actualGrossPrice)
}
