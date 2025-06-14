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
	ID              uuid.UUID
	TransactionCode string
	Amount          decimal.Decimal
	Currency        string
	CreatedAt       time.Time
	Status          string
}
