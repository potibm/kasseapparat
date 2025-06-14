package sumup

import (
	"time"

	"github.com/google/uuid"
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
	CreateReaderCheckout(readerId string, amount decimal.Decimal, description string, affiliateTransactionId string, returnUrl string) (*uuid.UUID, error)
	CreateReaderTerminateAction(readerId string) error
	GetTransactions(oldestTime *time.Time) ([]Transaction, error)
	GetTransactionById(transactionId uuid.UUID) (*Transaction, error)
	GetTransactionByClientTransactionId(clientTransactionId uuid.UUID) (*Transaction, error)
	RefundTransaction(transactionId uuid.UUID) error
}

func NewRepository(service *sumupService.Service) RepositoryInterface {
	return &Repository{
		service: service,
	}
}
