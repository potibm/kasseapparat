package http

import (
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/models"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	"github.com/shopspring/decimal"
)

type PurchaseListItemRequest struct {
	ID             int  `form:"ID"             binding:"required"`
	AttendedGuests uint `form:"attendedGuests" binding:"required"`
}

type PurchaseCartRequest struct {
	ID        int                       `form:"ID"        binding:"required"`
	Quantity  int                       `form:"quantity"  binding:"required"`
	NetPrice  decimal.Decimal           `form:"netPrice"  binding:"required"`
	ListItems []PurchaseListItemRequest `form:"listItems" binding:"required,dive"`
}

type PurchaseRequest struct {
	TotalNetPrice   decimal.Decimal       `form:"totalNetPrice"   binding:"required"`
	TotalGrossPrice decimal.Decimal       `form:"totalGrossPrice" binding:"required"`
	Cart            []PurchaseCartRequest `form:"cart"            binding:"required,dive"`
	PaymentMethod   models.PaymentMethod  `form:"paymentMethod"   binding:"required"`
	SumupReaderID   string                `form:"sumupReaderId"   binding:"omitempty"`
}

func (req PurchaseRequest) Validate() error {
	if req.TotalNetPrice.IsNegative() || req.TotalGrossPrice.IsNegative() {
		return fmt.Errorf("total price must not be negative")
	}

	if len(req.Cart) == 0 {
		return fmt.Errorf("cart must not be empty")
	}

	seen := make(map[int]struct{})
	for _, cart := range req.Cart {
		if err := validateCartItem(cart, seen); err != nil {
			return err
		}
	}

	return nil
}

func validateCartItem(cart PurchaseCartRequest, seen map[int]struct{}) error {
	if cart.Quantity < 1 {
		return fmt.Errorf("quantity must be at least 1")
	}

	if _, ok := seen[cart.ID]; ok {
		return fmt.Errorf("duplicate product ID: %d", cart.ID)
	}

	seen[cart.ID] = struct{}{}

	for _, li := range cart.ListItems {
		if err := validateListItem(li); err != nil {
			return err
		}
	}

	return nil
}

func validateListItem(li PurchaseListItemRequest) error {
	if li.ID <= 0 {
		return fmt.Errorf("list item has invalid ID")
	}

	if li.AttendedGuests < 1 {
		return fmt.Errorf("attendedGuests must be at least 1")
	}

	return nil
}

func (req PurchaseRequest) ToInput() purchaseService.PurchaseInput {
	input := purchaseService.PurchaseInput{
		PaymentMethod:   req.PaymentMethod,
		TotalNetPrice:   req.TotalNetPrice,
		TotalGrossPrice: req.TotalGrossPrice,
	}

	for _, cart := range req.Cart {
		item := purchaseService.PurchaseCartItem{
			ID:       cart.ID,
			Quantity: cart.Quantity,
			NetPrice: cart.NetPrice,
		}
		for _, li := range cart.ListItems {
			item.ListItems = append(item.ListItems, purchaseService.ListItemInput{
				ID:             li.ID,
				AttendedGuests: int(li.AttendedGuests),
			})
		}

		input.Cart = append(input.Cart, item)
	}

	return input
}
