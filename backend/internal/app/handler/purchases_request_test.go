package handler

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestToInputWithInvalidQuantity(t *testing.T) {
	req := PurchaseRequest{
		Cart: []PurchaseCartRequest{
			{ID: 1, Quantity: 0},
		},
	}

	err := req.Validate()
	if err == nil || !strings.Contains(err.Error(), "quantity") {
		t.Errorf("expected quantity error, got %v", err)
	}
}

func TestToInputWithNegativeTotalNetPrice(t *testing.T) {
	req := PurchaseRequest{
		TotalNetPrice: decimal.NewFromInt(-1),
		Cart:          []PurchaseCartRequest{},
	}

	err := req.Validate()
	if err == nil || !strings.Contains(err.Error(), "total price must not be negative") {
		t.Errorf("expected total price not negative error, got %v", err)
	}
}

func TestToInputWithNegativeTotalGrossPrice(t *testing.T) {
	req := PurchaseRequest{
		TotalGrossPrice: decimal.NewFromInt(-1),
		Cart:            []PurchaseCartRequest{},
	}

	err := req.Validate()
	if err == nil || !strings.Contains(err.Error(), "total price must not be negative") {
		t.Errorf("expected total price not negative error, got %v", err)
	}
}

func TestToInputWithDuplicateProductId(t *testing.T) {
	req := PurchaseRequest{
		Cart: []PurchaseCartRequest{
			{ID: 1, Quantity: 1},
			{ID: 2, Quantity: 1},
			{ID: 3, Quantity: 1},
			{ID: 2, Quantity: 1},
		},
	}

	err := req.Validate()
	if err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("expected duplicate product id error, got %v", err)
	}
}

func TestToInputWithEmptyCart(t *testing.T) {
	req := PurchaseRequest{
		Cart: []PurchaseCartRequest{},
	}

	err := req.Validate()
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Errorf("expected empty cart error, got %v", err)
	}
}

func TestToInputWithInvalidAttendedGuestsValue(t *testing.T) {
	req := PurchaseRequest{
		Cart: []PurchaseCartRequest{
			{ID: 1, Quantity: 1, ListItems: []PurchaseListItemRequest{{ID: 1, AttendedGuests: 0}}},
		},
	}

	err := req.Validate()
	if err == nil || !strings.Contains(err.Error(), "attendedGuests") {
		t.Errorf("expected attendedGuests must at least be 1 error, got %v", err)
	}
}
