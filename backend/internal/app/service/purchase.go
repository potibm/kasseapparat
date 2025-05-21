package service

import (
	"errors"
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
	"github.com/shopspring/decimal"
)

type Repository interface {
	GetProductByID(id int) (*models.Product, error)
	GetFullGuestByID(id int) (*models.Guest, error)
}

type PurchaseService struct {
	Repo          Repository
	Mailer        mailer.Mailer
	DecimalPlaces int32
}

type PurchaseInput struct {
	Cart            []PurchaseCartItem
	TotalNetPrice   decimal.Decimal
	TotalGrossPrice decimal.Decimal
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

func NewPurchaseService(repo *repository.Repository, mailer mailer.Mailer, decimalPlaces int32) *PurchaseService {
	return &PurchaseService{
		Repo:          repo,
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
