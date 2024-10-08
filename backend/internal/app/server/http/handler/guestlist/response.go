package guestlist

import "github.com/potibm/kasseapparat/internal/app/entities/guestlist"

type GuestlistResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	TypeCode  bool   `json:"typeCode"`
	ProductId uint   `json:"productId"`
}

func newGuestlistReponse(list guestlist.Guestlist) GuestlistResponse {
	return GuestlistResponse{
		ID:        list.ID,
		Name:      list.Name,
		TypeCode:  list.TypeCode,
		ProductId: list.Product.ID,
	}
}
