package monitor

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks basierend auf deinen Interfaces.
type MockSqlite struct{ mock.Mock }

func (m *MockSqlite) GetPurchaseByID(id uuid.UUID) (*models.Purchase, error) {
	args := m.Called(id)

	return args.Get(0).(*models.Purchase), args.Error(1)
}
func (m *MockSqlite) UpdatePurchaseSumupTransactionIDByID(id, sID uuid.UUID) (*models.Purchase, error) {
	args := m.Called(id, sID)

	return args.Get(0).(*models.Purchase), args.Error(1)
}

type MockSumup struct{ mock.Mock }

func (m *MockSumup) GetTransactionByClientTransactionId(id uuid.UUID) (*sumupRepo.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sumupRepo.Transaction), args.Error(1)
}

type MockService struct{ mock.Mock }

func (m *MockService) FinalizePurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*models.Purchase), args.Error(1)
}
func (m *MockService) FailPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*models.Purchase), args.Error(1)
}
func (m *MockService) CancelPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*models.Purchase), args.Error(1)
}

type MockPublisher struct{ mock.Mock }

func (m *MockPublisher) PushUpdate(id uuid.UUID, status models.PurchaseStatus) {
	m.Called(id, status)
}

func TestHandleTransactionPollingSuccess(t *testing.T) {
	tID := uuid.New()
	sClientID := uuid.New()
	sTransID := uuid.New()

	mSqlite := new(MockSqlite)
	mSumup := new(MockSumup)
	mService := new(MockService)
	mPub := new(MockPublisher)

	poller := &transactionPoller{
		SqliteRepository: mSqlite,
		SumupRepository:  mSumup,
		PurchaseService:  mService,
		StatusPublisher:  mPub,
	}

	// 1. Mock: DB liefert Purchase
	p := &models.Purchase{
		ID:                       tID,
		PaymentMethod:            models.PaymentMethodSumUp,
		SumupClientTransactionID: &sClientID,
	}
	mSqlite.On("GetPurchaseByID", tID).Return(p, nil)

	// 2. Mock: SumUp meldet SUCCESSFUL
	mSumup.On("GetTransactionByClientTransactionId", sClientID).Return(&sumupRepo.Transaction{
		TransactionID: sTransID,
		Status:        "SUCCESSFUL",
	}, nil)

	// 3. Mock: Update DB mit SumUp ID
	mSqlite.On("UpdatePurchaseSumupTransactionIDByID", tID, sTransID).Return(p, nil)

	// 4. Mock: Service beendet Kauf
	finalP := &models.Purchase{ID: tID, Status: models.PurchaseStatusConfirmed}
	mService.On("FinalizePurchase", mock.Anything, tID).Return(finalP, nil)
	mPub.On("PushUpdate", tID, models.PurchaseStatusConfirmed).Return()

	// Run
	shouldStop := poller.handleTransactionPolling(tID)

	assert.True(t, shouldStop)
	mService.AssertExpectations(t)
}

func TestHandleTransactionPollingNotFound(t *testing.T) {
	tID := uuid.New()
	sClientID := uuid.New()

	mSqlite := new(MockSqlite)
	mSumup := new(MockSumup)
	mService := new(MockService)
	mPub := new(MockPublisher)

	poller := &transactionPoller{
		SqliteRepository: mSqlite,
		SumupRepository:  mSumup,
		PurchaseService:  mService,
		StatusPublisher:  mPub,
	}

	p := &models.Purchase{ID: tID, PaymentMethod: models.PaymentMethodSumUp, SumupClientTransactionID: &sClientID}
	mSqlite.On("GetPurchaseByID", tID).Return(p, nil)

	// Simuliere NOT_FOUND
	mSumup.On("GetTransactionByClientTransactionId", sClientID).Return(nil, fmt.Errorf("SumUp error: NOT_FOUND"))

	mService.On("FailPurchase", mock.Anything, tID).Return(&models.Purchase{Status: models.PurchaseStatusFailed}, nil)
	mPub.On("PushUpdate", tID, models.PurchaseStatusFailed).Return()

	shouldStop := poller.handleTransactionPolling(tID)

	assert.True(t, shouldStop)
}
