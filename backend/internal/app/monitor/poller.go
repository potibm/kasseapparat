package monitor

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
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
	PushUpdate(purchaseID uuid.UUID, status string)
}

type transactionPoller struct {
	SumupRepository  sumupRepo.RepositoryInterface
	SqliteRepository *sqliteRepo.Repository
	PurchaseService  purchaseService.Service
	StatusPublisher  StatusPublisher
	active           map[string]struct{}
}

func NewPoller(sumupRepo sumupRepo.RepositoryInterface, sqliteRepo *sqliteRepo.Repository, purchaseService purchaseService.Service, statusPublisher StatusPublisher) Poller {
	return &transactionPoller{
		SumupRepository:  sumupRepo,
		SqliteRepository: sqliteRepo,
		PurchaseService:  purchaseService,
		StatusPublisher:  statusPublisher,
		active:           make(map[string]struct{}),
	}
}

// Starts a polling loop for a given transaction ID
func (n *transactionPoller) Start(transactionID uuid.UUID) {
	log.Println("Starting polling for transaction:", transactionID)

	if !registerPoller(transactionID) {
		log.Println("Polling already running for transaction:", transactionID)
		return // already running
	}

	go func() {
		defer unregisterPoller(transactionID)

		log.Printf("Polling started for %s\n", transactionID)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if done := n.handleTransactionPolling(transactionID); done {
				log.Printf("Polling ended for %s", transactionID)
				return
			}
		}
	}()
}

func (n *transactionPoller) handleTransactionPolling(transactionID uuid.UUID) bool {
	ctx := context.Background()

	purchase, err := n.SqliteRepository.GetPurchaseByID(transactionID)
	if err != nil {
		log.Printf("DB error for %s: %v", transactionID, err)
		return false
	}

	if purchase.PaymentMethod != models.PaymentMethodSumUp {
		log.Printf("Skipping polling for %s, not a SumUp transaction", transactionID)
		return true
	}

	if isFinal(string(purchase.Status)) {
		n.StatusPublisher.PushUpdate(transactionID, string(purchase.Status))

		log.Printf("Polling ended for %s", transactionID)

		return true
	}

	if purchase.SumupClientTransactionID == nil {
		log.Printf("No SumUp client transaction ID for %s, skipping polling", transactionID)
		return true
	}

	// Fetch current status from SumUp
	transaction, err := n.SumupRepository.GetTransactionByClientTransactionId(*purchase.SumupClientTransactionID)
	if err != nil {
		log.Printf("Error fetching transaction %s from SumUp: %v", purchase.SumupTransactionID, err)
		return false
	}

	if purchase.SumupTransactionID == nil {
		purchase, err = n.SqliteRepository.UpdatePurchaseSumupTransactionIDByID(transactionID, transaction.TransactionID)
		if err != nil {
			log.Printf("Error updating purchase %s with SumUp transaction ID: %v", transactionID, err)
		}
	}

	log.Printf("Transaction %s status: %s", purchase.SumupTransactionID, transaction.Status)

	return n.handleStatusUpdate(ctx, transactionID, transaction.Status, purchase)
}

func (n *transactionPoller) handleStatusUpdate(ctx context.Context, transactionID uuid.UUID, status string, purchase *models.Purchase) bool {
	var (
		updatedPurchase *models.Purchase
		err             error
	)

	switch status {
	case "PENDING":
		log.Printf("Transaction %s is still pending, continuing to poll", purchase.SumupTransactionID)
		websocket.PushUpdate(transactionID, string(purchase.Status))

		return false
	case "SUCCESSFUL":
		updatedPurchase, err = n.PurchaseService.FinalizePurchase(ctx, transactionID)
	case "FAILED":
		updatedPurchase, err = n.PurchaseService.FailPurchase(ctx, transactionID)
	case "CANCELED":
		updatedPurchase, err = n.PurchaseService.CancelPurchase(ctx, transactionID)
	default:
		log.Printf("Unknown transaction status %s for %s, skipping update", status, transactionID)
		return false
	}

	if err != nil {
		log.Printf("Error updating purchase %s status: %v", transactionID, err)
		return false
	}

	n.StatusPublisher.PushUpdate(transactionID, string(updatedPurchase.Status))

	return isFinal(string(updatedPurchase.Status))
}
