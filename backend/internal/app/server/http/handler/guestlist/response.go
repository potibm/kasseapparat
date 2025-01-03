package guestlist

import "github.com/potibm/kasseapparat/internal/app/entities/guestlist"

type GuestlistResponse struct {
	ID       uint                      `json:"id"`
	Name     string                    `json:"name"`
	TypeCode bool                      `json:"typeCode"`
	Product  *GuestlistProductResponse `json:"product"`
}

type GuestlistProductResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func newGuestlistReponse(list guestlist.Guestlist) GuestlistResponse {
	response := GuestlistResponse{
		ID:       list.ID,
		Name:     list.Name,
		TypeCode: list.TypeCode,
	}

	if list.Product != nil && list.Product.ID != 0 {
		response.Product = &GuestlistProductResponse{
			Id:   list.Product.ID,
			Name: list.Product.Name,
		}
	}

	return response
}
