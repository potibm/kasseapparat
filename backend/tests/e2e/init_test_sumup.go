package tests_e2e

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/shopspring/decimal"
)

type MockSumUpRepository struct {
	GetReadersFunc           func() ([]sumup.Reader, error)
	GetReaderFunc            func(readerId string) (*sumup.Reader, error)
	CreateReaderFunc         func(pairingCode, readerName string) (*sumup.Reader, error)
	DeleteReaderFunc         func(readerId string) error
	CreateReaderCheckoutFunc func(readerId string, amount decimal.Decimal,
		description string, affiliateTransactionId string,
		returnUrl *string) (*uuid.UUID, error)
	CreateReaderTerminateActionFunc func(readerId string) error
	GetTransactionsFunc             func(oldestFrom *time.Time) ([]sumup.Transaction, error)
	GetTransactionByIDFunc          func(transactionId uuid.UUID) (*sumup.Transaction, error)
	RefundTransactionFunc           func(transactionId uuid.UUID) error
	GetWebhookURLFunc               func() *string
}

func NewMockSumUpRepository() *MockSumUpRepository {
	const mockCheckoutUUID = "00000000-0000-4000-8000-000000000000"

	return &MockSumUpRepository{
		GetReadersFunc: func() ([]sumup.Reader, error) {
			return []sumup.Reader{
				{ID: "mock-1", Name: "Mock Reader 1"},
			}, nil
		},
		GetReaderFunc: func(readerID string) (*sumup.Reader, error) {
			return &sumup.Reader{ID: readerID, Name: "Mock Reader"}, nil
		},
		CreateReaderFunc: func(pairingCode string, readerName string) (*sumup.Reader, error) {
			return &sumup.Reader{ID: "created-1", Name: readerName}, nil
		},
		DeleteReaderFunc: func(readerID string) error {
			return nil
		},
		CreateReaderCheckoutFunc: func(readerID string, amount decimal.Decimal,
			description string, affiliateTransactionID string,
			returnURL *string,
		) (*uuid.UUID, error) {
			checkoutID, _ := uuid.Parse(mockCheckoutUUID)

			return &checkoutID, nil
		},
		CreateReaderTerminateActionFunc: func(readerID string) error {
			return nil
		},
		GetTransactionsFunc: func(oldestFrom *time.Time) ([]sumup.Transaction, error) {
			return []sumup.Transaction{
				{ID: uuid.New().String(), Amount: decimal.NewFromFloat(10.00), Status: "COMPLETED"},
			}, nil
		},
		GetTransactionByIDFunc: func(transactionID uuid.UUID) (*sumup.Transaction, error) {
			if transactionID.String() == mockCheckoutUUID {
				return &sumup.Transaction{
					ID:       transactionID.String(),
					Currency: "EUR",
					Amount:   decimal.NewFromFloat(10.00),
					Status:   "COMPLETED",
				}, nil
			}

			return nil, nil
		},
		RefundTransactionFunc: func(transactionID uuid.UUID) error {
			if transactionID.String() == mockCheckoutUUID {
				return nil
			}

			return fmt.Errorf("transaction not found")
		},
		GetWebhookURLFunc: func() *string {
			url := "https://mock-webhook-url.example.com"

			return &url
		},
	}
}

func (m *MockSumUpRepository) GetReaders() ([]sumup.Reader, error) {
	return m.GetReadersFunc()
}

func (m *MockSumUpRepository) GetReader(readerID string) (*sumup.Reader, error) {
	return m.GetReaderFunc(readerID)
}

func (m *MockSumUpRepository) CreateReader(pairingCode, readerName string) (*sumup.Reader, error) {
	return m.CreateReaderFunc(pairingCode, readerName)
}

func (m *MockSumUpRepository) DeleteReader(readerID string) error {
	return m.DeleteReaderFunc(readerID)
}

func (m *MockSumUpRepository) CreateReaderCheckout(
	readerID string,
	amount decimal.Decimal,
	description string,
	affiliateTransactionID string,
	returnURL *string,
) (*uuid.UUID, error) {
	return m.CreateReaderCheckoutFunc(readerID, amount, description, affiliateTransactionID, returnURL)
}

func (m *MockSumUpRepository) CreateReaderTerminateAction(readerID string) error {
	return m.CreateReaderTerminateActionFunc(readerID)
}

func (m *MockSumUpRepository) GetTransactions(oldestFrom *time.Time) ([]sumup.Transaction, error) {
	return m.GetTransactionsFunc(oldestFrom)
}

func (m *MockSumUpRepository) GetTransactionByClientTransactionID(transactionID uuid.UUID) (*sumup.Transaction, error) {
	return m.GetTransactionByIDFunc(transactionID)
}

func (m *MockSumUpRepository) GetTransactionByID(transactionID uuid.UUID) (*sumup.Transaction, error) {
	return m.GetTransactionByIDFunc(transactionID)
}

func (m *MockSumUpRepository) RefundTransaction(transactionID uuid.UUID) error {
	return m.RefundTransactionFunc(transactionID)
}

func (m *MockSumUpRepository) GetWebhookURL() *string {
	return m.GetWebhookURLFunc()
}

type MockStatusPublisher struct {
	Calls []struct {
		PurchaseID uuid.UUID
		Status     string
	}
}

func (m *MockStatusPublisher) PushUpdate(purchaseID uuid.UUID, status models.PurchaseStatus) {
	m.Calls = append(m.Calls, struct {
		PurchaseID uuid.UUID
		Status     string
	}{purchaseID, string(status)})
}
