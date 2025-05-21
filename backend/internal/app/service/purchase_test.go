package service

import (
	"errors"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

type MockRepository struct {
	Products map[int]*models.Product
	Guests   map[int]*models.Guest
}

func (m *MockRepository) GetProductByID(id int) (*models.Product, error) {
	p, ok := m.Products[id]
	if !ok {
		return nil, errors.New("not found")
	}

	return p, nil
}

func (m *MockRepository) GetFullGuestByID(id int) (*models.Guest, error) {
	g, ok := m.Guests[id]
	if !ok {
		return nil, errors.New("guest not found")
	}

	return g, nil
}

func TestValidateAndCalculatePricesWithSuccess(t *testing.T) {
	mockRepo := &MockRepository{
		Products: map[int]*models.Product{
			1: {
				NetPrice: decimal.NewFromFloat(10.00),
				VATRate:  decimal.NewFromInt(19),
			},
		},
	}

	service := &PurchaseService{
		Repo:          mockRepo,
		DecimalPlaces: 2,
	}

	input := PurchaseInput{
		Cart: []PurchaseCartItem{
			{ID: 1, Quantity: 2},
		},
		TotalNetPrice:   decimal.NewFromFloat(20.00),
		TotalGrossPrice: decimal.NewFromFloat(23.80),
	}

	net, gross, err := service.ValidateAndCalculatePrices(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !net.Equal(decimal.NewFromFloat(20.00)) {
		t.Errorf("unexpected net total: %s", net)
	}

	if !gross.Equal(decimal.NewFromFloat(23.80)) {
		t.Errorf("unexpected gross total: %s", gross)
	}
}

func TestValidateAndCalculatePricesWithProductNotFound(t *testing.T) {
	mockRepo := &MockRepository{
		Products: map[int]*models.Product{},
	}

	service := &PurchaseService{
		Repo:          mockRepo,
		DecimalPlaces: 2,
	}

	input := PurchaseInput{
		Cart: []PurchaseCartItem{
			{ID: 999, Quantity: 1},
		},
		TotalNetPrice:   decimal.NewFromFloat(10.00),
		TotalGrossPrice: decimal.NewFromFloat(11.90),
	}

	_, _, err := service.ValidateAndCalculatePrices(input)
	if err == nil || err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}

func TestValidateAndCalculatePricesWithPriceMismatch(t *testing.T) {
	mockRepo := &MockRepository{
		Products: map[int]*models.Product{
			1: {
				NetPrice: decimal.NewFromFloat(10.00),
				VATRate:  decimal.NewFromInt(19),
			},
		},
	}

	service := &PurchaseService{
		Repo:          mockRepo,
		DecimalPlaces: 2,
	}

	input := PurchaseInput{
		Cart: []PurchaseCartItem{
			{ID: 1, Quantity: 1},
		},
		TotalNetPrice:   decimal.NewFromFloat(15.00), // falsch!
		TotalGrossPrice: decimal.NewFromFloat(17.85), // falsch!
	}

	_, _, err := service.ValidateAndCalculatePrices(input)
	if err == nil || err != ErrInvalidProductPrice {
		t.Errorf("expected ErrInvalidProductPrice, got: %v", err)
	}
}

func TestValidateAndPrepareGuestsWithSuccess(t *testing.T) {
	mockRepo := &MockRepository{
		Guests: map[int]*models.Guest{
			42: {
				AttendedGuests:   0,
				AdditionalGuests: 1,
				Guestlist: models.Guestlist{
					ProductID: 1,
				},
			},
		},
	}

	service := &PurchaseService{
		Repo: mockRepo,
	}

	input := PurchaseInput{
		Cart: []PurchaseCartItem{
			{
				ID:       1,
				Quantity: 1,
				ListItems: []ListItemInput{
					{ID: 42, AttendedGuests: 1},
				},
			},
		},
	}

	guests, err := service.ValidateAndPrepareGuests(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(guests) != 1 {
		t.Errorf("expected 1 updated guest, got %d", len(guests))
	}

	if guests[0].AttendedGuests != 1 {
		t.Errorf("expected guest to have AttendedGuests = 1, got %d", guests[0].AttendedGuests)
	}
}

func TestValidateAndPrepareGuestsWithGuestNotFound(t *testing.T) {
	service := &PurchaseService{
		Repo: &MockRepository{Guests: map[int]*models.Guest{}}, // leer!
	}

	input := PurchaseInput{
		Cart: []PurchaseCartItem{
			{
				ID:       1,
				Quantity: 1,
				ListItems: []ListItemInput{
					{ID: 99, AttendedGuests: 1},
				},
			},
		},
	}

	_, err := service.ValidateAndPrepareGuests(input)
	if err == nil {
		t.Errorf("expected error for missing guest, got none")
	}
}

func TestValidateAndPrepareGuestsWithGuestAlreadyAttended(t *testing.T) {
	mockRepo := &MockRepository{
		Guests: map[int]*models.Guest{
			43: {
				AttendedGuests:   1,
				AdditionalGuests: 1,
				Guestlist: models.Guestlist{
					ProductID: 1,
				},
			},
		},
	}

	service := &PurchaseService{
		Repo: mockRepo,
	}

	input := PurchaseInput{
		Cart: []PurchaseCartItem{
			{
				ID:       1,
				Quantity: 1,
				ListItems: []ListItemInput{
					{ID: 43, AttendedGuests: 1},
				},
			},
		},
	}

	_, err := service.ValidateAndPrepareGuests(input)
	if err == nil {
		t.Errorf("expected error for already attended guest, got none")
	}
}
