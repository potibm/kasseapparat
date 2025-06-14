package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PurchaseStatus string

const (
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusConfirmed PurchaseStatus = "confirmed"
	PurchaseStatusFailed    PurchaseStatus = "failed"
	PurchaseStatusCancelled PurchaseStatus = "cancelled"
)

type Purchase struct {
	ID        uuid.UUID `gorm:"type:text;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	GormOwnedModel
	TotalNetPrice            decimal.Decimal `gorm:"type:TEXT"             json:"totalNetPrice"`
	TotalGrossPrice          decimal.Decimal `gorm:"type:TEXT"             json:"totalGrossPrice"`
	PurchaseItems            []PurchaseItem  `gorm:"foreignKey:PurchaseID" json:"purchaseItems"`
	PaymentMethod            string          `gorm:"type:TEXT"             json:"paymentMethod"`
	SumupTransactionID       *uuid.UUID      `gorm:"type:TEXT"             json:"sumupTransactionId"`
	SumupClientTransactionID *uuid.UUID      `gorm:"type:TEXT"             json:"sumupClientTransactionId"`
	Status                   PurchaseStatus  `gorm:"type:TEXT;default:'confirmed'"             json:"status"`
}

func (p *Purchase) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return
}
