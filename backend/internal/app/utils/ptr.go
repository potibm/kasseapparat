package utils

import (
	"time"

	"github.com/shopspring/decimal"
)

func F64PtrToDecimal(p *float64) decimal.Decimal {
	if p != nil {
		return decimal.NewFromFloat(*p)
	}

	return decimal.Zero
}

func StrPtr(p *string) string {
	if p != nil {
		return *p
	}

	return ""
}

func TimePtr(p *time.Time) time.Time {
	if p != nil {
		return *p
	}

	return time.Time{}
}
