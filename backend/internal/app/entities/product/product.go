package product

import "github.com/shopspring/decimal"

type Product struct {
	ID                   uint
	Name                 string
	Price                decimal.Decimal
	WrapAfter            bool
	Hidden               bool
	SoldOut              bool
	ApiExport            bool
	Pos                  uint
	TotalStock           uint
	UnitsSold            uint
	SoldOutInterestCount uint
}
