package guestlist

import (
	"time"

	"github.com/potibm/kasseapparat/internal/app/entities/order"
)

type Guest struct {
	ID               uint
	GuestList        Guestlist
	Name             string
	Code             *string
	AdditionalGuests uint
	AttendedGuests   uint
	ArrivedAt        *time.Time
	ArrivalNote      *string
	Purchase         *order.Order
}
