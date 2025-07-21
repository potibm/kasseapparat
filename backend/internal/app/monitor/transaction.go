package monitor

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/models"
)

// Starts a polling loop for a given transaction ID.
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
		n.StatusPublisher.PushUpdate(transactionID, purchase.Status)

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
		log.Printf("Error fetching transaction %s from SumUp: %v", purchase.SumupClientTransactionID, err)

		if strings.Contains(err.Error(), "NOT_FOUND") {
			log.Printf("Transaction %s not found in SumUp, stopping polling", purchase.SumupClientTransactionID)

			_, err := n.PurchaseService.FailPurchase(ctx, transactionID)
			if err != nil {
				log.Printf("Error setting purchase to failed %s: %v", transactionID, err)
			}

			n.StatusPublisher.PushUpdate(transactionID, models.PurchaseStatusFailed)

			return true
		}

		return false
	}

	if purchase.SumupTransactionID == nil {
		purchase, err = n.SqliteRepository.UpdatePurchaseSumupTransactionIDByID(transactionID, transaction.TransactionID)
		if err != nil {
			log.Printf("Error updating purchase %s with SumUp transaction ID: %v", transactionID, err)
		}
	}

	log.Printf("Transaction %s status: %s", transactionID, transaction.Status)

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
		websocket.PushUpdate(transactionID, purchase.Status)

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

	n.StatusPublisher.PushUpdate(transactionID, updatedPurchase.Status)

	return isFinal(string(updatedPurchase.Status))
}
