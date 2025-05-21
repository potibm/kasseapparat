package service

import (
	"context"
	"errors"
	"fmt"

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

type PurchaseService struct {
	Repo          Repository
	DB            *gorm.DB
	Mailer        mailer.Mailer
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
	Quantity  int
	ListItems []ListItemInput
}

var (
	ErrInvalidProductPrice = errors.New("invalid product price")
	ErrProductNotFound     = errors.New("product not found")
)

func uintPtr(v uint) *uint {
	return &v
}

func NewPurchaseService(repo *repository.Repository, mailer mailer.Mailer, decimalPlaces int32) *PurchaseService {
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

		net := product.NetPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		gross := product.GrossPrice(s.DecimalPlaces).Mul(decimal.NewFromInt(int64(item.Quantity)))

		totalNet = totalNet.Add(net)
		totalGross = totalGross.Add(gross)
	}

	if !totalNet.Equal(input.TotalNetPrice) {
		return totalNet, totalGross, ErrInvalidProductPrice
	}

	if !totalGross.Equal(input.TotalGrossPrice) {
		return totalNet, totalGross, ErrInvalidProductPrice
	}

	return totalNet, totalGross, nil
}

func (s *PurchaseService) ValidateAndPrepareGuests(input PurchaseInput) ([]models.Guest, error) {
	var updatedGuests []models.Guest

	for _, item := range input.Cart {
		for _, listInput := range item.ListItems {
			guest, err := s.Repo.GetFullGuestByID(listInput.ID)
			if err != nil || guest == nil {
				return nil, fmt.Errorf("list item %d not found", listInput.ID)
			}

			if guest.AttendedGuests != 0 {
				return nil, fmt.Errorf("list item %d has already been attended", guest.ID)
			}

			if guest.AdditionalGuests+1 < uint(listInput.AttendedGuests) {
				return nil, fmt.Errorf("too many additional guests for item %d", guest.ID)
			}

			if guest.Guestlist.ProductID != uint(item.ID) {
				return nil, fmt.Errorf("list item %d does not belong to product %d", guest.ID, item.ID)
			}

			guest.AttendedGuests = uint(listInput.AttendedGuests)
			guest.MarkAsArrived()

			updatedGuests = append(updatedGuests, *guest)
		}
	}

	return updatedGuests, nil
}

func (s *PurchaseService) CreatePurchase(ctx context.Context, input PurchaseInput, userID int) (*models.Purchase, error) {
	// Schritt 1: Preise validieren
	net, gross, err := s.ValidateAndCalculatePrices(input)
	if err != nil {
		return nil, err
	}

	// Schritt 2: Gästelisten validieren
	guests, err := s.ValidateAndPrepareGuests(input)
	if err != nil {
		return nil, err
	}

	// Schritt 3: Transaktion starten
	var savedPurchase *models.Purchase

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 3.1 Purchase-Objekt aufbauen
		purchase := &models.Purchase{
			TotalNetPrice:   net,
			TotalGrossPrice: gross,
			PaymentMethod:   input.PaymentMethod,
		}
		purchase.CreatedByID = uintPtr(uint(userID))

		// 3.2 PurchaseItems aufbauen
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
				//TotalPrice: product.GrossPrice(s.DecimalPlaces).Mul(decimal.NewFromInt(int64(item.Quantity))),
			}

			purchase.PurchaseItems = append(purchase.PurchaseItems, pi)
		}

		// 3.3 Speichern
		stored, err := s.Repo.StorePurchasesTx(tx, *purchase)
		if err != nil {
			return err
		}

		savedPurchase = &stored

		// 3.4 Gäste aktualisieren
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

	return savedPurchase, nil
}
