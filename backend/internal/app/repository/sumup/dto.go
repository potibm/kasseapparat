package sumup

import (
	"time"

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

type Checkout struct {
	ID              string
	Amount          decimal.Decimal
	Currency        string
	Description     string
	Status          string
	TransactionCode string
	CreatedAt       time.Time
}

type Transaction struct {
	ID              string
	TransactionCode string
	Amount          decimal.Decimal
	Currency        string
	CreatedAt       time.Time
	Status          string
}
