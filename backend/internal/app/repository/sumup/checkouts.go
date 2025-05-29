package sumup

import (
	"context"
	"log"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/sumup/sumup-go/checkouts"
)

func (r *Repository) GetCheckouts() ([]Checkout, error) {
	result := []Checkout{}
	params := checkouts.ListCheckoutsParams{}

	checkoutsResp, err := r.service.Client.Checkouts.List(context.Background(), params)
	if err != nil {
		return nil, err
	}

	// iterate over the checkouts and convert them to our Checkout type
	for _, checkout := range *checkoutsResp {
		result = append(result, *fromSDKCheckout(&checkout))
	}

	return result, nil
}

func (r *Repository) GetCheckout(id string) (*Checkout, error) {
	checkoutResp, err := r.service.Client.Checkouts.Get(context.Background(), id)
	if err != nil {
		log.Printf("Error retrieving checkout with ID %s: %v", id, err)

		return nil, err
	}

	return fromSDKCheckout(checkoutResp), nil
}

func fromSDKCheckout(sdkCheckout *checkouts.CheckoutSuccess) *Checkout {
	return &Checkout{
		ID:              utils.StrPtr(sdkCheckout.Id),
		Amount:          utils.F64PtrToDecimal(sdkCheckout.Amount),
		Currency:        string(*sdkCheckout.Currency),
		Description:     utils.StrPtr(sdkCheckout.Description),
		Status:          string(*sdkCheckout.Status),
		TransactionCode: utils.StrPtr(sdkCheckout.TransactionCode),
		CreatedAt:       utils.TimePtr(sdkCheckout.Date),
	}
}
