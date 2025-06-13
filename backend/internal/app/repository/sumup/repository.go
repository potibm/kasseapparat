package sumup

import (
	"time"

	sumupService "github.com/potibm/kasseapparat/internal/app/service/sumup"
	"github.com/shopspring/decimal"
)

type Repository struct {
	service *sumupService.Service
}

var _ RepositoryInterface = (*Repository)(nil)

type RepositoryInterface interface {
	GetReaders() ([]Reader, error)
	GetReader(readerId string) (*Reader, error)
	CreateReader(pairingCode string, readerName string) (*Reader, error)
	DeleteReader(readerId string) error
	CreateReaderCheckout(readerId string, amount decimal.Decimal, description string, affiliateTransactionId string, returnUrl string) (*string, error)
	CreateReaderTerminateAction(readerId string) error
	GetCheckouts() ([]Checkout, error)
	GetCheckout(id string) (*Checkout, error)
	GetTransactions(oldestTime *time.Time) ([]Transaction, error)
	GetTransactionById(transactionId string) (*Transaction, error)
	GetTransactionByClientTransactionId(clientTransactionId string) (*Transaction, error)
	RefundTransaction(transactionId string) error
}

func NewRepository(service *sumupService.Service) RepositoryInterface {
	return &Repository{
		service: service,
	}
}
