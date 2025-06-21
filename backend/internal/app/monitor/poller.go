package monitor

import (
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
)

// List of final states where we can stop polling
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

type transactionPoller struct {
	SumupRepository  sumupRepo.RepositoryInterface
	SqliteRepository sqliteRepo.RepositoryInterface
	PurchaseService  purchaseService.Service
	StatusPublisher  StatusPublisher
	active           map[string]struct{}
}

func NewPoller(sumupRepo sumupRepo.RepositoryInterface, sqliteRepo sqliteRepo.RepositoryInterface, purchaseService purchaseService.Service, statusPublisher StatusPublisher) Poller {
	return &transactionPoller{
		SumupRepository:  sumupRepo,
		SqliteRepository: sqliteRepo,
		PurchaseService:  purchaseService,
		StatusPublisher:  statusPublisher,
		active:           make(map[string]struct{}),
	}
}
