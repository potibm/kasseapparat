package monitor

import (
	"context"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
)

// List of final states where we can stop polling.
var finalStates = map[string]bool{
	"confirmed": true,
	"canceled":  true,
	"failed":    true,
	"timeout":   true,
}

func isFinal(status string) bool {
	return finalStates[status]
}

type Poller interface {
	Start(transactionID uuid.UUID)
}

type StatusPublisher interface {
	PushUpdate(purchaseID uuid.UUID, status models.PurchaseStatus)
}

type SumupTransactionReader interface {
	GetTransactionByClientTransactionId(clientTransactionId uuid.UUID) (*sumupRepo.Transaction, error)
}

type PurchaseRepository interface {
	GetPurchaseByID(id uuid.UUID) (*models.Purchase, error)
	UpdatePurchaseSumupTransactionIDByID(id uuid.UUID, sumupTransactionID uuid.UUID) (*models.Purchase, error)
}

type PurchaseStatusService interface {
	FinalizePurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	CancelPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	FailPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
}

type transactionPoller struct {
	SumupRepository  SumupTransactionReader
	SqliteRepository PurchaseRepository
	PurchaseService  PurchaseStatusService
	StatusPublisher  StatusPublisher
	active           map[string]struct{}
}

func NewPoller(
	sumupRp SumupTransactionReader,
	sqliteRp PurchaseRepository,
	purchaseSrvc PurchaseStatusService,
	statusPblshr StatusPublisher,
) Poller {
	return &transactionPoller{
		SumupRepository:  sumupRp,
		SqliteRepository: sqliteRp,
		PurchaseService:  purchaseSrvc,
		StatusPublisher:  statusPblshr,
		active:           make(map[string]struct{}),
	}
}
