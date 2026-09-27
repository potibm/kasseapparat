package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type (
	PurchaseStatus     string
	PurchaseStatusList []PurchaseStatus
)

const (
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusConfirmed PurchaseStatus = "confirmed"
	PurchaseStatusFailed    PurchaseStatus = "failed"
	PurchaseStatusCancelled PurchaseStatus = "cancelled"
	PurchaseStatusRefunded  PurchaseStatus = "refunded"
)

// StockConsumingStatuses are the purchase statuses whose units are still spoken for.
// Confirmed purchases are sold; pending ones are reserved so a second client cannot
// take the last units while a SumUp payment is in flight. The remaining statuses have
// released their units, which is what makes stock available again.
var StockConsumingStatuses = []PurchaseStatus{
	PurchaseStatusConfirmed,
	PurchaseStatusPending,
}

// ConsumesStock reports whether a purchase in this status still holds its units.
func (p PurchaseStatus) ConsumesStock() bool {
	for _, status := range StockConsumingStatuses {
		if p == status {
			return true
		}
	}

	return false
}

type PaymentMethod string

const (
	PaymentMethodCash    PaymentMethod = "CASH"
	PaymentMethodCC      PaymentMethod = "CC"
	PaymentMethodSumUp   PaymentMethod = "SUMUP"
	PaymentMethodVoucher PaymentMethod = "VOUCHER"
)

type Purchase struct {
	GormOwnedModel

	ID                       uuid.UUID       `json:"id"                       gorm:"type:text;primaryKey"`
	CreatedAt                time.Time       `json:"createdAt"                gorm:"index"`
	TotalNetPrice            decimal.Decimal `json:"totalNetPrice"            gorm:"type:TEXT"`
	TotalGrossPrice          decimal.Decimal `json:"totalGrossPrice"          gorm:"type:TEXT"`
	PurchaseItems            []PurchaseItem  `json:"purchaseItems"            gorm:"foreignKey:PurchaseID"`
	PaymentMethod            PaymentMethod   `json:"paymentMethod"            gorm:"type:TEXT"`
	SumupTransactionID       *uuid.UUID      `json:"sumupTransactionId"       gorm:"type:TEXT"`
	SumupClientTransactionID *uuid.UUID      `json:"sumupClientTransactionId" gorm:"type:TEXT"`
	Status                   PurchaseStatus  `json:"status"                   gorm:"type:TEXT;default:'confirmed'"`
}

func (p *Purchase) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return
}
