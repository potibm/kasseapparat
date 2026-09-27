package purchase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/config"
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
	CreateConfirmedPurchase(ctx context.Context, input PurchaseInput) (*models.Purchase, error)
	CreatePendingPurchase(ctx context.Context, input PurchaseInput) (*models.Purchase, error)
	FinalizePurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	CancelPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	FailPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	RefundPurchase(ctx context.Context, purchaseID uuid.UUID) (*models.Purchase, error)
}

var _ Service = (*PurchaseService)(nil)

var _ sqlite.RepositoryInterface = (*sqlite.Repository)(nil)

type Refunder interface {
	RefundTransaction(purchaseID uuid.UUID) error
}

type Mailer interface {
	SendNotificationOnArrival(email, name string) error
}

type PurchaseService struct {
	sqliteRepo    sqlite.RepositoryInterface
	sumupRepo     Refunder
	Mailer        Mailer
	DecimalPlaces int32
	CurrencyCode  string

	OutOfStockBehavior config.OutOfStockBehavior
}

type PurchaseInput struct {
	Cart            []PurchaseCartItem
	TotalNetPrice   decimal.Decimal
	TotalGrossPrice decimal.Decimal
	PaymentMethod   models.PaymentMethod
}

type ListItemInput struct {
	ID             int
	AttendedGuests uint
}

type PurchaseCartItem struct {
	ID        int
	NetPrice  decimal.Decimal
	Quantity  uint
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

	// ErrInsufficientStock is the sentinel for a purchase that exceeds the stock
	// available under the configured out-of-stock behaviour. Wrap it with
	// insufficientStockError to tell the caller how many units are left.
	ErrInsufficientStock = errors.New("not enough stock left")
)

// insufficientStockError carries the remaining stock so the handler can tell the
// operator how many units are actually left.
type insufficientStockError struct {
	ProductName string
	Available   int
}

func (e insufficientStockError) Error() string {
	if e.Available <= 0 {
		return fmt.Sprintf("%s is sold out", e.ProductName)
	}

	return fmt.Sprintf("only %d left in stock for %s", e.Available, e.ProductName)
}

func (e insufficientStockError) Is(target error) bool {
	return target == ErrInsufficientStock
}

func (e insufficientStockError) Unwrap() error { return ErrInsufficientStock }

// AvailableStock reports how many units of a product can still be sold, and whether
// the product is unlimited. A product with no total stock is always unlimited and is
// excluded from every restriction, which is what makes stock-less products such as
// drinks work without configuration.
func AvailableStock(product models.Product, unitsSold int) (available int, unlimited bool) {
	if product.TotalStock <= 0 {
		return 0, true
	}

	available = product.TotalStock - unitsSold
	if available < 0 {
		available = 0
	}

	return available, false
}

// restrictsStock reports whether the behaviour stops a sale that exceeds the stock.
// An unset or unrecognised value falls back to ignoring, which is both the documented
// default and the behaviour from before the setting existed, so a missing config value
// can never start rejecting sales.
func restrictsStock(behavior config.OutOfStockBehavior) bool {
	switch behavior {
	case config.OutOfStockFail, config.OutOfStockAutoSoldOut, config.OutOfStockAutoHide:
		return true
	default:
		return false
	}
}

// autoFlagForBehavior returns the product field the behaviour should set once the
// stock is depleted, or an empty string when it should not touch the product.
func autoFlagForBehavior(behavior config.OutOfStockBehavior) string {
	switch behavior {
	case config.OutOfStockAutoSoldOut:
		return "soldOut"
	case config.OutOfStockAutoHide:
		return "hidden"
	default:
		return ""
	}
}

func NewPurchaseService(
	sqliteRepo sqlite.RepositoryInterface,
	sumupRepo Refunder,
	mailer Mailer,
	decimalPlaces int32,
	currencyCode string,
	outOfStockBehavior config.OutOfStockBehavior,
) *PurchaseService {
	return &PurchaseService{
		sqliteRepo:         sqliteRepo,
		sumupRepo:          sumupRepo,
		Mailer:             mailer,
		DecimalPlaces:      decimalPlaces,
		CurrencyCode:       currencyCode,
		OutOfStockBehavior: outOfStockBehavior,
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

		net := product.NetPrice.Mul(decimal.NewFromUint64(uint64(item.Quantity)))
		gross := product.GrossPrice(s.DecimalPlaces).Mul(decimal.NewFromUint64(uint64(item.Quantity)))

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
) (*models.Purchase, error) {
	savedPurchase, guests, err := s.createPurchaseWithStatus(ctx, input, models.PurchaseStatusConfirmed)
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
) (*models.Purchase, error) {
	savedPurchase, _, err := s.createPurchaseWithStatus(ctx, input, models.PurchaseStatusPending)

	return savedPurchase, err
}

func (s *PurchaseService) FinalizePurchase(ctx context.Context, purchaseID uuid.UUID) (*models.Purchase, error) {
	// update status of purchase to confirmed
	purchase, err := s.setPurchaseStatus(ctx, purchaseID, models.PurchaseStatusConfirmed, false)
	if err != nil {
		return nil, errors.New("failed to finalize purchase: " + err.Error())
	}

	// notify guests
	guests, err := s.sqliteRepo.GetGuestsByPurchaseID(purchaseID)
	if guests == nil || err != nil {
		args := []any{"purchase_id", purchaseID}
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

func (s *PurchaseService) CancelPurchase(ctx context.Context, purchaseID uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.rollbackPurchase(ctx, purchaseID, models.PurchaseStatusCancelled)
	if err != nil {
		return nil, errors.New("failed to cancel purchase: " + err.Error())
	}

	return purchase, nil
}

func (s *PurchaseService) FailPurchase(ctx context.Context, purchaseID uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.rollbackPurchase(ctx, purchaseID, models.PurchaseStatusFailed)
	if err != nil {
		return nil, errors.New("failed to set the purchase to failed: " + err.Error())
	}

	return purchase, nil
}

func (s *PurchaseService) RefundPurchase(ctx context.Context, purchaseID uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.sqliteRepo.GetPurchaseByID(purchaseID)
	if err != nil {
		return nil, errors.New("failed to get purchase by ID: " + err.Error())
	}

	// Validate current status
	if purchase.Status != models.PurchaseStatusConfirmed {
		return nil, fmt.Errorf("cannot refund purchase with status: %s", purchase.Status)
	}

	// refund the purchase via SumUp
	if purchase.PaymentMethod == models.PaymentMethodSumUp && purchase.SumupTransactionID != nil {
		slog.Debug("Refunding transaction via SumUp for transaction", "transaction_id", *purchase.SumupTransactionID)

		if err := s.sumupRepo.RefundTransaction(*purchase.SumupTransactionID); err != nil {
			return nil, errors.New("failed to refund purchase via sumup: " + err.Error())
		}
	}

	// update status of purchase to refunded
	purchase, err = s.rollbackPurchase(ctx, purchaseID, models.PurchaseStatusRefunded)
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
	purchaseID uuid.UUID,
	status models.PurchaseStatus,
) (*models.Purchase, error) {
	purchase, err := s.setPurchaseStatus(ctx, purchaseID, status, true)
	if err != nil {
		return nil, errors.New("failed to rollback purchase: " + err.Error())
	}

	return purchase, err
}

func (s *PurchaseService) setPurchaseStatus(
	ctx context.Context,
	purchaseID uuid.UUID,
	status models.PurchaseStatus,
	rollbackGuests bool,
) (*models.Purchase, error) {
	var purchase *models.Purchase

	err := s.sqliteRepo.WithTransaction(ctx, func(txRepo sqlite.RepositoryInterface) error {
		p, err := txRepo.UpdatePurchaseStatusByID(purchaseID, status)
		if err != nil {
			return err
		}

		purchase = p

		if rollbackGuests {
			if err := txRepo.RollbackVisitedGuestsByPurchaseID(purchaseID); err != nil {
				return fmt.Errorf("failed to rollback visited guests: %w", err)
			}
		}

		// Moving out of a stock-consuming status releases the units, so any flag the
		// behaviour set for them may no longer apply.
		if !status.ConsumesStock() {
			if err := s.clearStockFlags(txRepo, p); err != nil {
				return err
			}
		}

		return nil
	})

	return purchase, err
}

// clearStockFlags removes a sold-out or hidden flag once the stock is available again,
// which happens when a purchase is refunded, failed or cancelled. Without this a
// single refund would keep a product off sale for the rest of the event.
//
// Known limitation: the product row does not record whether the flag was set here or by
// an operator, so a product an operator hid by hand also reappears once a refund frees
// its stock. Distinguishing the two would need columns tracking auto-managed flags.
// The flag is only cleared when the stock is genuinely available again, so a product
// that is still sold out keeps its flag.
func (s *PurchaseService) clearStockFlags(
	repo sqlite.RepositoryInterface,
	purchase *models.Purchase,
) error {
	field := autoFlagForBehavior(s.OutOfStockBehavior)
	if field == "" || purchase == nil {
		return nil
	}

	for _, item := range purchase.PurchaseItems {
		if err := s.clearStockFlag(repo, item.ProductID, field); err != nil {
			return err
		}
	}

	return nil
}

// clearStockFlag clears one product's flag if the behaviour owns that flag and the
// stock has become available again.
func (s *PurchaseService) clearStockFlag(
	repo sqlite.RepositoryInterface,
	productID int,
	field string,
) error {
	// Read the product fresh rather than trusting a preloaded copy: the flag may have
	// been set after the purchase was loaded.
	product, err := repo.GetProductByID(productID)
	if err != nil {
		return fmt.Errorf("failed to load product %d for stock update: %w", productID, err)
	}

	isSet := product.SoldOut
	if field == "hidden" {
		isSet = product.Hidden
	}

	if !isSet {
		return nil
	}

	unitsSold, err := repo.GetPurchasedQuantitiesByProductID(product.ID)
	if err != nil {
		return fmt.Errorf("failed to check stock for product %q: %w", product.Name, err)
	}

	available, unlimited := AvailableStock(*product, unitsSold)
	if !unlimited && available <= 0 {
		return nil
	}

	updated := *product
	if field == "soldOut" {
		updated.SoldOut = false
	} else {
		updated.Hidden = false
	}

	if _, err := repo.UpdateProductByID(product.ID, updated); err != nil {
		return fmt.Errorf("failed to clear %s on product %q: %w", field, product.Name, err)
	}

	return nil
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

	if guest.Guestlist.ProductID != productID {
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
		items, err := s.buildPurchaseItems(txRepo, input.Cart)
		if err != nil {
			return err
		}

		purchase := &models.Purchase{
			TotalNetPrice:   net,
			TotalGrossPrice: gross,
			PaymentMethod:   input.PaymentMethod,
			Status:          status,
			PurchaseItems:   items,
		}

		stored, err := txRepo.StorePurchases(*purchase)
		if err != nil {
			return err
		}

		savedPurchase = &stored

		if err := s.applyStockBehavior(txRepo, input.Cart); err != nil {
			return err
		}

		for _, guest := range guests {
			guest.PurchaseID = &stored.ID
			if _, err := txRepo.UpdateGuestByID(guest.ID, guest); err != nil {
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

// buildPurchaseItems turns the cart into purchase items, checking the stock of every
// product first. It runs on the transaction-scoped repository so the availability it
// reads and the purchase it belongs to see the same state.
func (s *PurchaseService) buildPurchaseItems(
	repo sqlite.RepositoryInterface,
	cart []PurchaseCartItem,
) ([]models.PurchaseItem, error) {
	// Aggregate the requested quantity per product first. Checking each cart line
	// against the same availability would let two lines of the same product both pass
	// and oversell, because neither sees the other's quantity.
	requestedByProduct := make(map[int]uint, len(cart))

	for _, item := range cart {
		requestedByProduct[item.ID] += item.Quantity
	}

	stockChecked := make(map[int]bool, len(requestedByProduct))
	items := make([]models.PurchaseItem, 0, len(cart))

	for _, item := range cart {
		product, err := repo.GetProductByID(item.ID)
		if err != nil {
			return nil, err
		}

		if !stockChecked[item.ID] {
			stockChecked[item.ID] = true

			if err := s.checkStock(repo, *product, requestedByProduct[item.ID]); err != nil {
				return nil, err
			}
		}

		items = append(items, models.PurchaseItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			NetPrice:  product.NetPrice,
			VATRate:   product.VATRate,
		})
	}

	return items, nil
}

// checkStock rejects a cart line that would take more units than are available. It
// runs on the transaction-scoped repository so the reading and the purchase that
// follows it see the same state; reading through the outer repository here would both
// race and wait for a second connection the bounded pool cannot provide.
func (s *PurchaseService) checkStock(
	repo sqlite.RepositoryInterface,
	product models.Product,
	quantity uint,
) error {
	if !restrictsStock(s.OutOfStockBehavior) {
		return nil
	}

	unitsSold, err := repo.GetPurchasedQuantitiesByProductID(product.ID)
	if err != nil {
		return fmt.Errorf("failed to check stock for product %q: %w", product.Name, err)
	}

	available, unlimited := AvailableStock(product, unitsSold)
	if unlimited {
		return nil
	}

	if int(quantity) > available {
		return insufficientStockError{ProductName: product.Name, Available: available}
	}

	return nil
}

// applyStockBehavior sets the product flag the configured behaviour asks for once the
// stock is depleted. The purchase has already been stored, so unitsSold now includes
// this purchase and the comparison is against the post-purchase total. The comparison
// is >= rather than == so a product that was already oversold still converges.
func (s *PurchaseService) applyStockBehavior(
	repo sqlite.RepositoryInterface,
	cart []PurchaseCartItem,
) error {
	field := autoFlagForBehavior(s.OutOfStockBehavior)
	if field == "" {
		return nil
	}

	for _, item := range cart {
		product, err := repo.GetProductByID(item.ID)
		if err != nil {
			return fmt.Errorf("failed to load product %d for stock update: %w", item.ID, err)
		}

		if _, unlimited := AvailableStock(*product, 0); unlimited {
			continue
		}

		unitsSold, err := repo.GetPurchasedQuantitiesByProductID(product.ID)
		if err != nil {
			return fmt.Errorf("failed to check stock for product %q: %w", product.Name, err)
		}

		available, _ := AvailableStock(*product, unitsSold)
		if available > 0 {
			continue
		}

		updated := *product
		if field == "soldOut" {
			updated.SoldOut = true
		} else {
			updated.Hidden = true
		}

		if _, err := repo.UpdateProductByID(product.ID, updated); err != nil {
			return fmt.Errorf("failed to mark product %q as %s: %w", product.Name, field, err)
		}
	}

	return nil
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
