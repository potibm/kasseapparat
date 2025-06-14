package purchase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service interface {
	CreateConfirmedPurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error)
	CreatePendingPurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error)
	FinalizePurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	CancelPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	FailPurchase(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
}

var _ Service = (*PurchaseService)(nil)

type sqliteRepository interface {
	GetProductByID(id int) (*models.Product, error)
	GetFullGuestByID(id int) (*models.Guest, error)
	UpdateGuestByIDTx(tx *gorm.DB, id int, updatedGuest models.Guest) (*models.Guest, error)
	StorePurchasesTx(tx *gorm.DB, purchase models.Purchase) (models.Purchase, error)
	GetPurchaseByIDTx(tx *gorm.DB, id uuid.UUID) (*models.Purchase, error)
	UpdatePurchaseStatusByIDTx(tx *gorm.DB, id uuid.UUID, status models.PurchaseStatus) (*models.Purchase, error)
	GetDB() *gorm.DB
}

type sumupRepository interface {
	RefundTransaction(purchaseId uuid.UUID) error
}

type Mailer interface {
	SendNotificationOnArrival(email string, name string) error
}

type PurchaseService struct {
	sqliteRepo    sqliteRepository
	sumupRepo     sumupRepository
	DB            *gorm.DB
	Mailer        Mailer
	DecimalPlaces int32
}

type PurchaseInput struct {
	Cart            []PurchaseCartItem
	TotalNetPrice   decimal.Decimal
	TotalGrossPrice decimal.Decimal
	PaymentMethod   string
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

func NewPurchaseService(sqliteRepo sqliteRepository, sumupRepo sumupRepository, mailer Mailer, decimalPlaces int32) *PurchaseService {
	return &PurchaseService{
		sqliteRepo:    sqliteRepo,
		DB:            sqliteRepo.GetDB(),
		sumupRepo:     sumupRepo,
		Mailer:        mailer,
		DecimalPlaces: decimalPlaces,
	}
}

func (s *PurchaseService) ValidateAndCalculatePrices(input PurchaseInput) (decimal.Decimal, decimal.Decimal, error) {
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
		return
	}

	for _, guest := range guests {
		if guest.NotifyOnArrivalEmail != nil {
			_ = s.Mailer.SendNotificationOnArrival(*guest.NotifyOnArrivalEmail, guest.Name)
		}
	}
}

func (s *PurchaseService) CreateConfirmedPurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error) {
	savedPurchase, guests, err := s.createPurchaseWithStatus(ctx, input, userID, models.PurchaseStatusConfirmed)

	if err != nil {
		return nil, err
	}

	s.notifyGuests(guests)

	return savedPurchase, nil
}

func (s *PurchaseService) CreatePendingPurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error) {
	savedPurchase, _, err := s.createPurchaseWithStatus(ctx, input, userID, models.PurchaseStatusPending)

	return savedPurchase, err
}

func (s *PurchaseService) createPurchaseWithStatus(ctx context.Context, input PurchaseInput, userID int, status models.PurchaseStatus) (*models.Purchase, []models.Guest, error) {
	net, gross, err := s.ValidateAndCalculatePrices(input)
	if err != nil {
		return nil, nil, err
	}

	guests, err := s.ValidateAndPrepareGuests(input)
	if err != nil {
		return nil, nil, err
	}

	var savedPurchase *models.Purchase

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		purchase := &models.Purchase{
			TotalNetPrice:   net,
			TotalGrossPrice: gross,
			PaymentMethod:   input.PaymentMethod,
			Status:          status,
		}
		purchase.CreatedByID = uintPtr(uint(userID))

		for _, item := range input.Cart {
			product, err := s.sqliteRepo.GetProductByID(item.ID)
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

		stored, err := s.sqliteRepo.StorePurchasesTx(tx, *purchase)
		if err != nil {
			return err
		}

		savedPurchase = &stored

		for _, guest := range guests {
			guest.PurchaseID = &stored.ID
			if _, err := s.sqliteRepo.UpdateGuestByIDTx(tx, int(guest.ID), guest); err != nil {
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

func (s *PurchaseService) FinalizePurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error) {
	// update status of purchase to confirmed
	purchase, err := s.setPurchaseStatus(ctx, purchaseId, models.PurchaseStatusConfirmed)
	if err != nil {
		return nil, errors.New("failed to finalize purchase: " + err.Error())
	}

	// @TODO: notify about arrival of guests
	// 	s.notifyGuests(guests)

	return purchase, nil
}

func (s *PurchaseService) CancelPurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.rollbackPurchase(ctx, purchaseId, models.PurchaseStatusCancelled)

	if err != nil {
		return nil, errors.New("failed to cancel purchase: " + err.Error())
	}

	// @TODO: refund the purchase (for sumup), perhaps by flag

	return purchase, nil
}

func (s *PurchaseService) FailPurchase(ctx context.Context, purchaseId uuid.UUID) (*models.Purchase, error) {
	purchase, err := s.rollbackPurchase(ctx, purchaseId, models.PurchaseStatusFailed)

	if err != nil {
		return nil, errors.New("failed to set the purchase to failed: " + err.Error())
	}

	// @TODO: refund the purchase (for sumup), perhaps by flag

	return purchase, nil
}

func (s *PurchaseService) rollbackPurchase(ctx context.Context, purchaseId uuid.UUID, status models.PurchaseStatus) (*models.Purchase, error) {
	// update status of purchase
	purchase, err := s.setPurchaseStatus(ctx, purchaseId, status)

	// @TODO: rollback all guests to not attended

	return purchase, err
}

func (s *PurchaseService) setPurchaseStatus(ctx context.Context, purchaseId uuid.UUID, status models.PurchaseStatus) (*models.Purchase, error) {
	var purchase *models.Purchase

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		p, err := s.sqliteRepo.UpdatePurchaseStatusByIDTx(tx, purchaseId, status)

		if err != nil {
			return err
		}

		purchase = p

		return nil
	})

	return purchase, err
}
