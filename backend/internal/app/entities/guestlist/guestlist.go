package guestlist

import "github.com/potibm/kasseapparat/internal/app/entities/product"

type Guestlist struct {
	ID       uint
	Name     string
	TypeCode bool
	Product  *product.Product
}
