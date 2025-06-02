package purchase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const errUnexpected = "unexpected error: %v"

type MockRepository struct {
	Products       map[int]*models.Product
	Guests         map[int]*models.Guest
	StoredPurchase *models.Purchase
	UpdatedGuests  map[int]*models.Guest
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

func (m *MockRepository) StorePurchasesTx(tx *gorm.DB, purchase models.Purchase) (models.Purchase, error) {
	purchase.ID = uuid.New()
	m.StoredPurchase = &purchase

	return purchase, nil
}

func (m *MockRepository) UpdateGuestByIDTx(tx *gorm.DB, id int, guest models.Guest) (*models.Guest, error) {
	g, ok := m.Guests[id]
	if !ok || g == nil {
		return nil, fmt.Errorf("guest %d not found in mock", id)
	}

	g.AttendedGuests = guest.AttendedGuests
	g.PurchaseID = guest.PurchaseID
	g.ArrivedAt = guest.ArrivedAt

	if m.UpdatedGuests == nil {
		m.UpdatedGuests = make(map[int]*models.Guest)
	}

	m.UpdatedGuests[id] = g

	return g, nil
}

func testDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

type MockMailer struct {
	Sent []string
}

func (m *MockMailer) SendNotificationOnArrival(email, name string) error {
	m.Sent = append(m.Sent, email+"|"+name)
	return nil
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
			{ID: 1, Quantity: 2, NetPrice: decimal.NewFromFloat(10.00)},
		},
		TotalNetPrice:   decimal.NewFromFloat(20.00),
		TotalGrossPrice: decimal.NewFromFloat(23.80),
	}

	net, gross, err := service.ValidateAndCalculatePrices(input)
	if err != nil {
		t.Errorf(errUnexpected, err)
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
			{ID: 999, Quantity: 1, NetPrice: decimal.NewFromFloat(10.00)},
		},
		TotalNetPrice:   decimal.NewFromFloat(10.00),
		TotalGrossPrice: decimal.NewFromFloat(11.90),
	}

	_, _, err := service.ValidateAndCalculatePrices(input)
	if err == nil || err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}

func TestValidateAndCalculatePricesWithTotalPriceMismatch(t *testing.T) {
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
			{ID: 1, Quantity: 1, NetPrice: decimal.NewFromFloat(10.00)},
		},
		TotalNetPrice:   decimal.NewFromFloat(15.00), // wrong
		TotalGrossPrice: decimal.NewFromFloat(17.85), // wrong
	}

	_, _, err := service.ValidateAndCalculatePrices(input)
	if err == nil || err != ErrInvalidTotalNetPrice {
		t.Errorf("expected ErrInvalidTotalNetPrice, got: %v", err)
	}
}

func TestValidateAndCalculatePricesWithProductPriceMismatch(t *testing.T) {
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
			{ID: 1, Quantity: 1, NetPrice: decimal.NewFromFloat(15.00)},
		},
		TotalNetPrice:   decimal.NewFromFloat(10.00),
		TotalGrossPrice: decimal.NewFromFloat(11.90),
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
		t.Errorf(errUnexpected, err)
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

func TestValidateGuestWithGuestNotFound(t *testing.T) {
	service := &PurchaseService{
		Repo: &MockRepository{
			Guests: map[int]*models.Guest{},
		},
	}

	_, err := service.validateGuest(ListItemInput{ID: 99, AttendedGuests: 1}, 1)
	if err != ErrGuestNotFound {
		t.Fatalf("expected ErrGuestNotFound, got %v", err)
	}
}

func TestValidateGuestWithGuestAlreadyAttended(t *testing.T) {
	service := &PurchaseService{
		Repo: &MockRepository{
			Guests: map[int]*models.Guest{
				42: {
					AttendedGuests: 1,
					Guestlist:      models.Guestlist{ProductID: 1},
				},
			},
		},
	}

	_, err := service.validateGuest(ListItemInput{ID: 42, AttendedGuests: 1}, 1)
	if err != ErrGuestAlreadyAttended {
		t.Fatalf("expected ErrGuestAlreadyAttended, got %v", err)
	}
}

func TestValidateGuestWithTooManyAdditionalGuests(t *testing.T) {
	service := &PurchaseService{
		Repo: &MockRepository{
			Guests: map[int]*models.Guest{
				42: {
					AttendedGuests:   0,
					AdditionalGuests: 1,
					Guestlist:        models.Guestlist{ProductID: 1},
				},
			},
		},
	}

	_, err := service.validateGuest(ListItemInput{ID: 42, AttendedGuests: 3}, 1)
	if err != ErrTooManyAdditionalGuests {
		t.Fatalf("expected ErrTooManyAdditionalGuests, got %v", err)
	}
}

func TestValidateGuestWithWrongProductId(t *testing.T) {
	service := &PurchaseService{
		Repo: &MockRepository{
			Guests: map[int]*models.Guest{
				42: {
					AttendedGuests:   0,
					AdditionalGuests: 1,
					Guestlist:        models.Guestlist{ProductID: 1},
				},
			},
		},
	}

	_, err := service.validateGuest(ListItemInput{ID: 42, AttendedGuests: 1}, 2)
	if err != ErrListItemWrongProduct {
		t.Fatalf("expected ErrListItemWrongProduct, got %v", err)
	}
}

func TestCreatePurchaseWithSuccess(t *testing.T) {
	ctx := context.Background()

	g := &models.Guest{
		AdditionalGuests: 1,
		AttendedGuests:   0,
		Guestlist: models.Guestlist{
			ProductID: 1,
		},
	}
	g.ID = 42

	p := &models.Product{
		NetPrice: decimal.NewFromFloat(10.00),
		VATRate:  decimal.NewFromFloat(19),
	}
	p.ID = 1

	mockRepo := &MockRepository{
		Products: map[int]*models.Product{
			1: p,
		},
		Guests: map[int]*models.Guest{
			42: g,
		},
	}

	service := &PurchaseService{
		Repo:          mockRepo,
		DB:            testDB(),
		DecimalPlaces: 2,
	}

	input := PurchaseInput{
		PaymentMethod:   "CASH",
		TotalNetPrice:   decimal.NewFromFloat(10.00),
		TotalGrossPrice: decimal.NewFromFloat(11.90),
		Cart: []PurchaseCartItem{
			{
				ID:       1,
				Quantity: 1,
				NetPrice: decimal.NewFromFloat(10.00),
				ListItems: []ListItemInput{
					{ID: 42, AttendedGuests: 1},
				},
			},
		},
	}

	purchase, err := service.CreatePurchase(ctx, input, 7)
	if err != nil {
		t.Fatalf(errUnexpected, err)
	}

	if purchase == nil {
		t.Fatal("expected purchase, got nil")
	}

	if purchase.PaymentMethod != "CASH" {
		t.Errorf("unexpected payment method: %s", purchase.PaymentMethod)
	}

	if len(mockRepo.UpdatedGuests) != 1 {
		t.Errorf("expected 1 updated guest, got %d", len(mockRepo.UpdatedGuests))
	}

	if g.AttendedGuests != 1 {
		t.Errorf("guest not marked as attended: %+v", g)
	}

	if g.PurchaseID == nil {
		t.Error("guest has no purchase ID")
	}
}

func TestNotifyGuestsWithSendsExpectedEmails(t *testing.T) {
	mailer := &MockMailer{}
	service := &PurchaseService{
		Mailer: mailer,
	}

	guests := []models.Guest{
		{
			Name:                 "Alice",
			NotifyOnArrivalEmail: ptr("alice@example.com"),
		},
		{
			Name: "Bob",
			// no email, should not trigger
		},
		{
			Name:                 "Eve",
			NotifyOnArrivalEmail: ptr("eve@example.com"),
		},
	}

	service.notifyGuests(guests)

	if len(mailer.Sent) != 2 {
		t.Fatalf("expected 2 mails, got %d", len(mailer.Sent))
	}

	want := []string{
		"alice@example.com|Alice",
		"eve@example.com|Eve",
	}
	for i, expected := range want {
		if mailer.Sent[i] != expected {
			t.Errorf("mail %d: expected %q, got %q", i, expected, mailer.Sent[i])
		}
	}
}

func ptr(s string) *string {
	return &s
}
