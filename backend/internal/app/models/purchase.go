package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PurchaseStatus string
type PurchaseStatusList []PurchaseStatus

const (
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusConfirmed PurchaseStatus = "confirmed"
	PurchaseStatusFailed    PurchaseStatus = "failed"
	PurchaseStatusCancelled PurchaseStatus = "cancelled"
	PurchaseStatusRefunded  PurchaseStatus = "refunded"
)

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
