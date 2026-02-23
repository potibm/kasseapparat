package monitor

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/models"
)

func (n *transactionPoller) Start(transactionID uuid.UUID) {
	slog.Debug("Starting polling for transaction", "transaction_id", transactionID)

	if !registerPoller(transactionID) {
		slog.Info("Polling already running for transaction", "transaction_id", transactionID)

		return // already running
	}

	go func() {
		defer unregisterPoller(transactionID)

		slog.Debug("Polling started for transaction", "transaction_id", transactionID)

		const pollingInterval = 5 * time.Second

		ticker := time.NewTicker(pollingInterval)
		defer ticker.Stop()

		for range ticker.C {
			if done := n.handleTransactionPolling(transactionID); done {
				slog.Debug("Polling ended for transaction", "transaction_id", transactionID)

				return
			}
		}
	}()
}

func (n *transactionPoller) handleTransactionPolling(transactionID uuid.UUID) bool {
	ctx := context.Background()

	purchase, err := n.SqliteRepository.GetPurchaseByID(transactionID)
	if err != nil {
		slog.Error("Database error for transaction", "transaction_id", transactionID, "error", err)

		return false
	}

	if purchase.PaymentMethod != models.PaymentMethodSumUp {
		slog.Info("Skipping polling for transaction, not a SumUp transaction", "transaction_id", transactionID)

		return true
	}

	if isFinal(string(purchase.Status)) {
		n.StatusPublisher.PushUpdate(transactionID, purchase.Status)

		slog.Info("Polling ended for transaction", "transaction_id", transactionID)

		return true
	}

	if purchase.SumupClientTransactionID == nil {
		slog.Info("No SumUp client transaction ID for transaction", "transaction_id", transactionID)

		return true
	}

	// Fetch current status from SumUp
	transaction, err := n.SumupRepository.GetTransactionByClientTransactionId(*purchase.SumupClientTransactionID)
	if err != nil {
		slog.Error(
			"Error fetching transaction from SumUp",
			"transaction_id",
			purchase.SumupClientTransactionID,
			"error",
			err,
		)

		if strings.Contains(err.Error(), "NOT_FOUND") {
			slog.Info(
				"Transaction not found in SumUp, stopping polling",
				"transaction_id",
				purchase.SumupClientTransactionID,
			)

			_, err := n.PurchaseService.FailPurchase(ctx, transactionID)
			if err != nil {
				slog.Error("Error setting purchase to failed", "transaction_id", transactionID, "error", err)
			}

			n.StatusPublisher.PushUpdate(transactionID, models.PurchaseStatusFailed)

			return true
		}

		return false
	}

	if purchase.SumupTransactionID == nil {
		purchase, err = n.SqliteRepository.UpdatePurchaseSumupTransactionIDByID(
			transactionID,
			transaction.TransactionID,
		)
		if err != nil {
			slog.Error(
				"Error updating purchase with SumUp transaction ID",
				"transaction_id",
				transactionID,
				"error",
				err,
			)
		}
	}

	slog.Info("Transaction status update", "transaction_id", transactionID, "status", transaction.Status)

	return n.handleStatusUpdate(ctx, transactionID, transaction.Status, purchase)
}

func (n *transactionPoller) handleStatusUpdate(
	ctx context.Context,
	transactionID uuid.UUID,
	status string,
	purchase *models.Purchase,
) bool {
	var (
		updatedPurchase *models.Purchase
		err             error
	)

	switch status {
	case "PENDING":
		slog.Info("Transaction is still pending, continuing to poll", "transaction_id", purchase.SumupTransactionID)
		websocket.PushUpdate(transactionID, purchase.Status)

		return false
	case "SUCCESSFUL":
		updatedPurchase, err = n.PurchaseService.FinalizePurchase(ctx, transactionID)
	case "FAILED":
		updatedPurchase, err = n.PurchaseService.FailPurchase(ctx, transactionID)
	case "CANCELED":
		updatedPurchase, err = n.PurchaseService.CancelPurchase(ctx, transactionID)
	default:
		slog.Warn("Unknown transaction status", "status", status, "transaction_id", transactionID)

		return false
	}

	if err != nil {
		slog.Error("Error updating purchase status", "transaction_id", transactionID, "error", err)

		return false
	}

	n.StatusPublisher.PushUpdate(transactionID, updatedPurchase.Status)

	return isFinal(string(updatedPurchase.Status))
}
