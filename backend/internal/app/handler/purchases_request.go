package handler

import (
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/service"
	"github.com/shopspring/decimal"
)

type PurchaseListItemRequest struct {
	ID             int  `binding:"required" form:"ID"`
	AttendedGuests uint `binding:"required" form:"attendedGuests"`
}

type PurchaseCartRequest struct {
	ID        int                       `binding:"required"      form:"ID"`
	Quantity  int                       `binding:"required"      form:"quantity"`
	NetPrice  decimal.Decimal           `binding:"required" form:"netPrice"`
	ListItems []PurchaseListItemRequest `binding:"required,dive" form:"listItems"`
}

type PurchaseRequest struct {
	TotalNetPrice   decimal.Decimal       `binding:"required"      form:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal       `binding:"required"      form:"totalGrossPrice"`
	Cart            []PurchaseCartRequest `binding:"required,dive" form:"cart"`
	PaymentMethod   string                `binding:"required"      form:"paymentMethod"`
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

func (req PurchaseRequest) ToInput() service.PurchaseInput {
	input := service.PurchaseInput{
		PaymentMethod:   req.PaymentMethod,
		TotalNetPrice:   req.TotalNetPrice,
		TotalGrossPrice: req.TotalGrossPrice,
	}

	for _, cart := range req.Cart {
		item := service.PurchaseCartItem{
			ID:       cart.ID,
			Quantity: cart.Quantity,
			NetPrice: cart.NetPrice,
		}
		for _, li := range cart.ListItems {
			item.ListItems = append(item.ListItems, service.ListItemInput{
				ID:             li.ID,
				AttendedGuests: int(li.AttendedGuests),
			})
		}

		input.Cart = append(input.Cart, item)
	}

	return input
}
