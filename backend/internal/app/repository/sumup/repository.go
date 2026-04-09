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
	ReaderRepository
	TransactionRepository

	GetWebhookURL() *string
}

type ReaderRepository interface {
	GetReaders() ([]Reader, error)
	GetReader(readerID string) (*Reader, error)
	CreateReader(pairingCode, readerName string) (*Reader, error)
	DeleteReader(readerID string) error
	CreateReaderCheckout(
		readerID string,
		amount decimal.Decimal,
		description string,
		affiliateTransactionID string,
		returnURL *string,
	) (*uuid.UUID, error)
	CreateReaderTerminateAction(readerID string) error
}

type TransactionRepository interface {
	GetTransactions(oldestTime *time.Time) ([]Transaction, error)
	GetTransactionByID(transactionID uuid.UUID) (*Transaction, error)
	GetTransactionByClientTransactionID(clientTransactionID uuid.UUID) (*Transaction, error)
	RefundTransaction(transactionID uuid.UUID) error
}

func NewRepository(service *sumupService.Service) RepositoryInterface {
	return &Repository{
		service: service,
	}
}

func (r *Repository) GetWebhookURL() *string {
	return r.service.WebhookURL
}
