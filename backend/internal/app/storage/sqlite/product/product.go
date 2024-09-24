package product

import (
	"github.com/potibm/kasseapparat/internal/app/entities/product"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
)

type Product struct {
	models.GormOwnedModel
	Name                 string          ``
	Price                decimal.Decimal `sql:"type:decimal(20,8);"`
	WrapAfter            bool            `gorm:"default:false" `
	Hidden               bool            `gorm:"default:false" `
	SoldOut              bool            `gorm:"default:false" `
	ApiExport            bool            `gorm:"default:false" `
	Pos                  uint            ``
	TotalStock           uint            `gorm:"default:0"`
	UnitsSold            uint            `gorm:"default:0" `
	SoldOutInterestCount uint            `gorm:"default:0"`
}

func (m Product) CreateEntity() *product.Product {
	return &product.Product{
		Name:                 m.Name,
		Price:                m.Price,
		WrapAfter:            m.WrapAfter,
		Hidden:               m.Hidden,
		SoldOut:              m.SoldOut,
		ApiExport:            m.ApiExport,
		Pos:                  m.Pos,
		TotalStock:           m.TotalStock,
		UnitsSold:            m.UnitsSold,
		SoldOutInterestCount: m.SoldOutInterestCount,
	}
}
