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
	GetReadersFunc                  func() ([]sumup.Reader, error)
	GetReaderFunc                   func(readerId string) (*sumup.Reader, error)
	CreateReaderFunc                func(pairingCode string, readerName string) (*sumup.Reader, error)
	DeleteReaderFunc                func(readerId string) error
	CreateReaderCheckoutFunc        func(readerId string, amount decimal.Decimal, description string, affiliateTransactionId string, returnUrl string) (*uuid.UUID, error)
	CreateReaderTerminateActionFunc func(readerId string) error
	GetTransactionsFunc             func(oldestFrom *time.Time) ([]sumup.Transaction, error)
	GetTransactionByIdFunc          func(transactionId uuid.UUID) (*sumup.Transaction, error)
	RefundTransactionFunc           func(transactionId uuid.UUID) error
}

func (m *MockSumUpRepository) GetReaders() ([]sumup.Reader, error) {
	return m.GetReadersFunc()
}

func (m *MockSumUpRepository) GetReader(readerId string) (*sumup.Reader, error) {
	return m.GetReaderFunc(readerId)
}

func (m *MockSumUpRepository) CreateReader(pairingCode string, readerName string) (*sumup.Reader, error) {
	return m.CreateReaderFunc(pairingCode, readerName)
}

func (m *MockSumUpRepository) DeleteReader(readerId string) error {
	return m.DeleteReaderFunc(readerId)
}

func (m *MockSumUpRepository) CreateReaderCheckout(readerId string, amount decimal.Decimal, description string, affiliateTransactionId string, returnUrl string) (*uuid.UUID, error) {
	return m.CreateReaderCheckoutFunc(readerId, amount, description, affiliateTransactionId, returnUrl)
}

func (m *MockSumUpRepository) CreateReaderTerminateAction(readerId string) error {
	return m.CreateReaderTerminateActionFunc(readerId)
}

func (m *MockSumUpRepository) GetTransactions(oldestFrom *time.Time) ([]sumup.Transaction, error) {
	return m.GetTransactionsFunc(oldestFrom)
}

func (m *MockSumUpRepository) GetTransactionByClientTransactionId(transactionId uuid.UUID) (*sumup.Transaction, error) {
	return m.GetTransactionByIdFunc(transactionId)
}

func (m *MockSumUpRepository) GetTransactionById(transactionId uuid.UUID) (*sumup.Transaction, error) {
	return m.GetTransactionByIdFunc(transactionId)
}

func (m *MockSumUpRepository) RefundTransaction(transactionId uuid.UUID) error {
	return m.RefundTransactionFunc(transactionId)
}

func NewMockSumUpRepository() *MockSumUpRepository {
	const mockCheckoutUUID = "00000000-0000-4000-8000-000000000000"

	return &MockSumUpRepository{
		GetReadersFunc: func() ([]sumup.Reader, error) {
			return []sumup.Reader{
				{ID: "mock-1", Name: "Mock Reader 1"},
			}, nil
		},
		GetReaderFunc: func(readerId string) (*sumup.Reader, error) {
			return &sumup.Reader{ID: readerId, Name: "Mock Reader"}, nil
		},
		CreateReaderFunc: func(pairingCode string, readerName string) (*sumup.Reader, error) {
			return &sumup.Reader{ID: "created-1", Name: readerName}, nil
		},
		DeleteReaderFunc: func(readerId string) error {
			return nil
		},
		CreateReaderCheckoutFunc: func(readerId string, amount decimal.Decimal, description string, affiliateTransactionId string, returnUrl string) (*uuid.UUID, error) {
			checkoutId, _ := uuid.Parse(mockCheckoutUUID)
			return &checkoutId, nil
		},
		CreateReaderTerminateActionFunc: func(readerId string) error {
			return nil
		},
		GetTransactionsFunc: func(oldestFrom *time.Time) ([]sumup.Transaction, error) {
			return []sumup.Transaction{
				{ID: uuid.New().String(), Amount: decimal.NewFromFloat(10.00), Status: "COMPLETED"},
			}, nil
		},
		GetTransactionByIdFunc: func(transactionId uuid.UUID) (*sumup.Transaction, error) {
			if transactionId.String() == mockCheckoutUUID {
				return &sumup.Transaction{ID: transactionId.String(), Currency: "EUR", Amount: decimal.NewFromFloat(10.00), Status: "COMPLETED"}, nil
			}
			return nil, nil
		},
		RefundTransactionFunc: func(transactionId uuid.UUID) error {
			if transactionId.String() == mockCheckoutUUID {
				return nil
			}
			return fmt.Errorf("transaction not found")
		},
	}
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
