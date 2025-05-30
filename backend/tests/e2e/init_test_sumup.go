package tests_e2e

import (
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/shopspring/decimal"
)

type MockSumUpRepository struct {
	GetReadersFunc                  func() ([]sumup.Reader, error)
	GetReaderFunc                   func(readerId string) (*sumup.Reader, error)
	CreateReaderFunc                func(pairingCode string, readerName string) (*sumup.Reader, error)
	DeleteReaderFunc                func(readerId string) error
	CreateReaderCheckoutFunc        func(readerId string, readerName string, amount decimal.Decimal) (*string, error)
	CreateReaderTerminateActionFunc func(readerId string) error
	GetCheckoutsFunc                func() ([]sumup.Checkout, error)
	GetCheckoutFunc                 func(id string) (*sumup.Checkout, error)
	GetTransactionsFunc             func() ([]sumup.Transaction, error)
	GetTransactionByIdFunc          func(transactionId string) (*sumup.Transaction, error)
	RefundTransactionFunc           func(transactionId string) error
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

func (m *MockSumUpRepository) CreateReaderCheckout(readerId string, readerName string, amount decimal.Decimal) (*string, error) {
	return m.CreateReaderCheckoutFunc(readerId, readerName, amount)
}

func (m *MockSumUpRepository) CreateReaderTerminateAction(readerId string) error {
	return m.CreateReaderTerminateActionFunc(readerId)
}
func (m *MockSumUpRepository) GetCheckouts() ([]sumup.Checkout, error) {
	return m.GetCheckoutsFunc()
}
func (m *MockSumUpRepository) GetCheckout(id string) (*sumup.Checkout, error) {
	return m.GetCheckoutFunc(id)
}
func (m *MockSumUpRepository) GetTransactions() ([]sumup.Transaction, error) {
	return m.GetTransactionsFunc()
}
func (m *MockSumUpRepository) GetTransactionById(transactionId string) (*sumup.Transaction, error) {
	return m.GetTransactionByIdFunc(transactionId)
}
func (m *MockSumUpRepository) RefundTransaction(transactionId string) error {
	return m.RefundTransactionFunc(transactionId)
}

func NewMockSumUpRepository() *MockSumUpRepository {
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
		CreateReaderCheckoutFunc: func(readerId string, readerName string, amount decimal.Decimal) (*string, error) {
			checkoutId := "checkout-1"
			return &checkoutId, nil
		},
		CreateReaderTerminateActionFunc: func(readerId string) error {
			return nil
		},
		GetCheckoutsFunc: func() ([]sumup.Checkout, error) {
			return []sumup.Checkout{
				{ID: "checkout-1", Currency: "EUR", Amount: decimal.NewFromFloat(10.00), Status: "COMPLETED"},
			}, nil
		},
		GetCheckoutFunc: func(id string) (*sumup.Checkout, error) {
			if id == "checkout-1" {
				return &sumup.Checkout{ID: id, Currency: "EUR", TransactionCode: "transaction-1", Amount: decimal.NewFromFloat(10.00), Status: "COMPLETED"}, nil
			}
			return nil, nil
		},
		GetTransactionsFunc: func() ([]sumup.Transaction, error) {
			return []sumup.Transaction{
				{ID: "transaction-1", Amount: decimal.NewFromFloat(10.00), Status: "COMPLETED"},
			}, nil
		},
		GetTransactionByIdFunc: func(transactionId string) (*sumup.Transaction, error) {
			if transactionId == "transaction-1" {
				return &sumup.Transaction{ID: transactionId, Currency: "EUR", Amount: decimal.NewFromFloat(10.00), Status: "COMPLETED"}, nil
			}
			return nil, nil
		},
		RefundTransactionFunc: func(transactionId string) error {
			if transactionId == "transaction-1" {
				return nil
			}
			return fmt.Errorf("transaction not found")
		},
	}
}
