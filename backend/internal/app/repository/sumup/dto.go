package sumup

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Reader struct {
	ID               string
	Name             string
	Status           string
	DeviceIdentifier string
	DeviceModel      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Transaction struct {
	ID              string
	TransactionID   uuid.UUID
	TransactionCode string
	Amount          decimal.Decimal
	Currency        string
	CardType        string
	CreatedAt       time.Time
	Events          []TransactionEvent
	Status          string
}

type TransactionEvent struct {
	ID        int
	Timestamp time.Time
	Type      string
	Amount    decimal.Decimal
	Status    string
}
