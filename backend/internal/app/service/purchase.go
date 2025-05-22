package service

import (
	"context"
	"errors"

	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Repository interface {
	GetProductByID(id int) (*models.Product, error)
	GetFullGuestByID(id int) (*models.Guest, error)
	StorePurchasesTx(tx *gorm.DB, purchase models.Purchase) (models.Purchase, error)
	UpdateGuestByIDTx(tx *gorm.DB, id int, updatedGuest models.Guest) (*models.Guest, error)
}

type Mailer interface {
	SendNotificationOnArrival(email string, name string) error
}

type PurchaseService struct {
	Repo          Repository
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

func NewPurchaseService(repo *repository.Repository, mailer *mailer.Mailer, decimalPlaces int32) *PurchaseService {
	return &PurchaseService{
		Repo:          repo,
		DB:            repo.GetDB(),
		Mailer:        mailer,
		DecimalPlaces: decimalPlaces,
	}
}

func (s *PurchaseService) ValidateAndCalculatePrices(input PurchaseInput) (decimal.Decimal, decimal.Decimal, error) {
	totalNet := decimal.NewFromInt(0)
	totalGross := decimal.NewFromInt(0)

	for _, item := range input.Cart {
		product, err := s.Repo.GetProductByID(item.ID)
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
	guest, err := s.Repo.GetFullGuestByID(listInput.ID)
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

func (s *PurchaseService) CreatePurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error) {
	net, gross, err := s.ValidateAndCalculatePrices(input)
	if err != nil {
		return nil, err
	}

	guests, err := s.ValidateAndPrepareGuests(input)
	if err != nil {
		return nil, err
	}

	var savedPurchase *models.Purchase

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		purchase := &models.Purchase{
			TotalNetPrice:   net,
			TotalGrossPrice: gross,
			PaymentMethod:   input.PaymentMethod,
		}
		purchase.CreatedByID = uintPtr(uint(userID))

		for _, item := range input.Cart {
			product, err := s.Repo.GetProductByID(item.ID)
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

		stored, err := s.Repo.StorePurchasesTx(tx, *purchase)
		if err != nil {
			return err
		}

		savedPurchase = &stored

		for _, guest := range guests {
			guest.PurchaseID = &stored.ID
			if _, err := s.Repo.UpdateGuestByIDTx(tx, int(guest.ID), guest); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.notifyGuests(guests)

	return savedPurchase, nil
}
