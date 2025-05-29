package sumup

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestGetValueFromDecimal(t *testing.T) {
	dec := decimal.NewFromFloat(123.456)
	value := getValueFromDecimal(dec, 2) // EURO
	expected := 12345

	if value != expected {
		t.Errorf("Expected %d, got %d", expected, value)
	}

	value = getValueFromDecimal(dec, 0) // HUF
	expected = 123

	if value != expected {
		t.Errorf("Expected %d, got %d", expected, value)
	}
}
