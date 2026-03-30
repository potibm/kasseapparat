package purchase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("kasseapparat")
var (
	salesOrdersCounter, _ = meter.Int64Counter("kasseapparat_sales_orders_total",
		metric.WithDescription("Total number of processed orders or refunds"))

	salesAmountCounter, _ = meter.Int64Counter("kasseapparat_sales_amount_total",
		metric.WithDescription("Total monetary value in cents"),
		metric.WithUnit("ct"))
)

type Service interface {
	CreateConfirmedPurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error)
	CreatePendingPurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error)
	FinalizePurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	CancelPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	FailPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	RefundPurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error)
}

var _ Service = (*PurchaseService)(nil)

var _ sqlite.RepositoryInterface = (*sqlite.Repository)(nil)

type sumupRepository interface {
	RefundTransaction(purchaseId uuid.UUID) error
}

type Mailer interface {
	SendNotificationOnArrival(email string, name string) error
}

type PurchaseService struct {
	sqliteRepo    sqlite.RepositoryInterface
	sumupRepo     sumupRepository
	Mailer        Mailer
	DecimalPlaces int32
	CurrencyCode  string
}

type PurchaseInput struct {
	Cart            []PurchaseCartItem
	TotalNetPrice   decimal.Decimal
	TotalGrossPrice decimal.Decimal
	PaymentMethod   models.PaymentMethod
}

type ListItemInput struct {
	ID             int
	AttendedGuests int
}

type PurchaseCartItem struct {
	ID        int
	NetPrice  decimal.Decimal
	Quantity  int
	ListItems []ListItemInput
}

var (
	ErrInvalidTotalGrossPrice  = errors.New("total gross price does not match")
	ErrInvalidTotalNetPrice    = errors.New("total net price does not match")
	ErrInvalidProductPrice     = errors.New("invalid product price")
	ErrProductNotFound         = errors.New("product not found")
	ErrGuestNotFound           = errors.New("guest not found")
	ErrGuestAlreadyAttended    = errors.New("guest already attended")
	ErrTooManyAdditionalGuests = errors.New("additional guests exceed available guests")
	ErrListItemWrongProduct    = errors.New("list item does not belong to product")
)

func uintPtr(v uint) *uint {
	return &v
}

func NewPurchaseService(
	sqliteRepo sqlite.RepositoryInterface,
	sumupRepo sumupRepository,
	mailer Mailer,
	decimalPlaces int32,
	currencyCode string,
) *PurchaseService {
	return &PurchaseService{
		sqliteRepo:    sqliteRepo,
		sumupRepo:     sumupRepo,
		Mailer:        mailer,
		DecimalPlaces: decimalPlaces,
		CurrencyCode:  currencyCode,
	}
}

func (s *PurchaseService) ValidateAndCalculatePrices(
	input PurchaseInput,
) (totalNetResult, totalGrossResult decimal.Decimal, err error) {
	totalNet := decimal.NewFromInt(0)
	totalGross := decimal.NewFromInt(0)

	for _, item := range input.Cart {
		product, err := s.sqliteRepo.GetProductByID(item.ID)
		if err != nil || product == nil {
			return decimal.Zero, decimal.Zero, ErrProductNotFound
		}

		if !product.NetPrice.Round(s.DecimalPlaces).Equal(item.NetPrice.Round(s.DecimalPlaces)) {
			return decimal.Zero, decimal.Zero, ErrInvalidProductPrice
		}

		net := product.NetPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		gross := product.GrossPrice(s.DecimalPlaces).Mul(decimal.NewFromInt(int64(item.Quantity)))

		totalNet = totalNet.Add(net)
		totalGross = totalGross.Add(gross)
	}

	if !totalNet.Equal(input.TotalNetPrice) {
		return totalNet, totalGross, ErrInvalidTotalNetPrice
	}

	if !totalGross.Equal(input.TotalGrossPrice) {
		return totalNet, totalGross, ErrInvalidTotalGrossPrice
	}

	return totalNet, totalGross, nil
}

func (s *PurchaseService) ValidateAndPrepareGuests(input PurchaseInput) ([]models.Guest, error) {
	var updatedGuests []models.Guest

	for _, item := range input.Cart {
		for _, listInput := range item.ListItems {
			guest, err := s.validateGuest(listInput, item.ID)
			if err != nil {
				return nil, err
			}

			updatedGuests = append(updatedGuests, *guest)
		}
	}

	return updatedGuests, nil
}

func (s *PurchaseService) CreateConfirmedPurchase(
	ctx context.Context,
	input PurchaseInput,
	userID int,
) (*models.Purchase, error) {
	savedPurchase, guests, err := s.createPurchaseWithStatus(ctx, input, userID, models.PurchaseStatusConfirmed)
	if err != nil {
		return nil, err
	}

	s.notifyGuests(guests)

	s.recordTransactionMetrics(
		ctx,
		savedPurchase.TotalGrossPrice,
		savedPurchase.TotalNetPrice,
		string(savedPurchase.PaymentMethod),
		false,
	)

	return savedPurchase, nil
}

func (s *PurchaseService) CreatePendingPurchase(
	ctx context.Context,
	input PurchaseInput,
	userID int,
) (*models.Purchase, error) {
	savedPurchase, _, err := s.createPurchaseWithStatus(ctx, input, userID, models.PurchaseStatusPending)

	return savedPurchase, err
}

func (s *PurchaseService) FinalizePurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error) {
	// update status of purchase to confirmed
	purchase, err := s.setPurchaseStatus(ctx, purchaseId, models.PurchaseStatusConfirmed, false)
	if err != nil {
		return nil, errors.New("failed to finalize purchase: " + err.Error())
	}

	// notify guests
	guests, err := s.sqliteRepo.GetGuestsByPurchaseID(purchaseId)
	if guests == nil || err != nil {
		args := []any{"purchase_id", purchaseId}
		if err != nil {
			args = append(args, "error", err)
		}

		slog.Warn("No guests found for purchase, skipping notification", args...)
	} else {
		s.notifyGuests(guests)
	}

	s.recordTransactionMetrics(
		ctx,
		purchase.TotalGrossPrice,
		purchase.TotalNetPrice,
		string(purchase.PaymentMethod),
		false,
	)

	return purchase, nil
}

func (s *PurchaseService) recordTransactionMetrics(
	ctx context.Context,
	gross, net decimal.Decimal,
	method string,
	isRefund bool,
) {
	precision := s.DecimalPlaces
	multiplier := decimal.New(1, int32(precision))

	direction := int64(1)
	entryType := "purchase"

	if isRefund {
		direction = -1
		entryType = "refund"
	}

	grossSubUnits := gross.Mul(multiplier).IntPart() * direction
	netSubUnits := net.Mul(multiplier).IntPart() * direction

	commonAttrs := []attribute.KeyValue{
		attribute.String("type", entryType),
		attribute.String("currency", s.CurrencyCode),
		attribute.String("payment_method", method),
	}

	// Gross
	salesAmountCounter.Add(ctx, grossSubUnits, metric.WithAttributes(
		append(commonAttrs, attribute.String("tax_status", "gross"))...,
	))

	// Net
	salesAmountCounter.Add(ctx, netSubUnits, metric.WithAttributes(
		append(commonAttrs, attribute.String("tax_status", "net"))...,
	))

	salesOrdersCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("payment_method", method),
		attribute.String("type", entryType),
	))
}

func (s *PurchaseService) CancelPurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.rollbackPurchase(ctx, purchaseId, models.PurchaseStatusCancelled)
	if err != nil {
		return nil, errors.New("failed to cancel purchase: " + err.Error())
	}

	return purchase, nil
}

func (s *PurchaseService) FailPurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.rollbackPurchase(ctx, purchaseId, models.PurchaseStatusFailed)
	if err != nil {
		return nil, errors.New("failed to set the purchase to failed: " + err.Error())
	}

	return purchase, nil
}

func (s *PurchaseService) RefundPurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.sqliteRepo.GetPurchaseByID(purchaseId)
	if err != nil {
		return nil, errors.New("failed to get purchase by ID: " + err.Error())
	}

	// Validate current status
	if purchase.Status != models.PurchaseStatusConfirmed {
		return nil, fmt.Errorf("cannot refund purchase with status: %s", purchase.Status)
	}

	// refund the purchase via sumup
	if purchase.PaymentMethod == models.PaymentMethodSumUp && purchase.SumupTransactionID != nil {
		slog.Debug("Refunding transaction via SumUp for transaction", "transaction_id", *purchase.SumupTransactionID)

		if err := s.sumupRepo.RefundTransaction(*purchase.SumupTransactionID); err != nil {
			return nil, errors.New("failed to refund purchase via sumup: " + err.Error())
		}
	}

	// update status of purchase to refunded
	purchase, err = s.rollbackPurchase(ctx, purchaseId, models.PurchaseStatusRefunded)
	if err != nil {
		return nil, errors.New("failed to set the purchase to refunded: " + err.Error())
	}

	s.recordTransactionMetrics(
		ctx,
		purchase.TotalGrossPrice,
		purchase.TotalNetPrice,
		string(purchase.PaymentMethod),
		true,
	)

	return purchase, nil
}

func (s *PurchaseService) rollbackPurchase(
	ctx context.Context,
	purchaseId uuid.UUID,
	status models.PurchaseStatus,
) (*models.Purchase, error) {
	purchase, err := s.setPurchaseStatus(ctx, purchaseId, status, true)
	if err != nil {
		return nil, errors.New("failed to rollback purchase: " + err.Error())
	}

	return purchase, err
}

func (s *PurchaseService) setPurchaseStatus(
	ctx context.Context,
	purchaseId uuid.UUID,
	status models.PurchaseStatus,
	rollbackGuests bool,
) (*models.Purchase, error) {
	var purchase *models.Purchase

	err := s.sqliteRepo.WithTransaction(ctx, func(txRepo sqlite.RepositoryInterface) error {
		p, err := txRepo.UpdatePurchaseStatusByID(purchaseId, status)
		if err != nil {
			return err
		}

		purchase = p

		if rollbackGuests {
			if err := txRepo.RollbackVisitedGuestsByPurchaseID(purchaseId); err != nil {
				return fmt.Errorf("failed to rollback visited guests: %w", err)
			}
		}

		return nil
	})

	return purchase, err
}

func (s *PurchaseService) validateGuest(listInput ListItemInput, productID int) (*models.Guest, error) {
	guest, err := s.sqliteRepo.GetFullGuestByID(listInput.ID)
	if err != nil || guest == nil {
		return nil, ErrGuestNotFound
	}

	if guest.AttendedGuests != 0 {
		return nil, ErrGuestAlreadyAttended
	}

	if guest.AdditionalGuests+1 < uint(listInput.AttendedGuests) {
		return nil, ErrTooManyAdditionalGuests
	}

	if guest.Guestlist.ProductID != uint(productID) {
		return nil, ErrListItemWrongProduct
	}

	guest.AttendedGuests = uint(listInput.AttendedGuests)
	guest.MarkAsArrived()

	return guest, nil
}

func (s *PurchaseService) notifyGuests(guests []models.Guest) {
	if s.Mailer == nil {
		slog.Warn("Mailer is not configured, skipping guest notifications")

		return
	}

	for _, guest := range guests {
		if guest.NotifyOnArrivalEmail != nil {
			err := s.Mailer.SendNotificationOnArrival(*guest.NotifyOnArrivalEmail, guest.Name)
			if err != nil {
				slog.Error(
					"Failed to send notification email to guest",
					"guest_id",
					guest.ID,
					"error",
					err,
				)
			}
		}
	}
}

func (s *PurchaseService) createPurchaseWithStatus(
	ctx context.Context,
	input PurchaseInput,
	userID int,
	status models.PurchaseStatus,
) (*models.Purchase, []models.Guest, error) {
	net, gross, err := s.ValidateAndCalculatePrices(input)
	if err != nil {
		return nil, nil, err
	}

	guests, err := s.ValidateAndPrepareGuests(input)
	if err != nil {
		return nil, nil, err
	}

	var savedPurchase *models.Purchase

	err = s.sqliteRepo.WithTransaction(ctx, func(txRepo sqlite.RepositoryInterface) error {
		purchase := &models.Purchase{
			TotalNetPrice:   net,
			TotalGrossPrice: gross,
			PaymentMethod:   input.PaymentMethod,
			Status:          status,
		}
		purchase.CreatedByID = uintPtr(uint(userID))

		for _, item := range input.Cart {
			product, err := txRepo.GetProductByID(item.ID)
			if err != nil {
				return err
			}

			pi := models.PurchaseItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				NetPrice:  product.NetPrice,
				VATRate:   product.VATRate,
			}

			purchase.PurchaseItems = append(purchase.PurchaseItems, pi)
		}

		stored, err := txRepo.StorePurchases(*purchase)
		if err != nil {
			return err
		}

		savedPurchase = &stored

		for _, guest := range guests {
			guest.PurchaseID = &stored.ID
			if _, err := txRepo.UpdateGuestByID(int(guest.ID), guest); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return savedPurchase, guests, nil
}
