package utils

import (
	"time"

	"github.com/shopspring/decimal"
)

// F64PtrToDecimal converts a *float64 to a decimal.Decimal, returning decimal.Zero when p is nil.
// It returns decimal.NewFromFloat(*p) if p is non-nil, or decimal.Zero otherwise.
func F64PtrToDecimal(p *float64) decimal.Decimal {
	if p != nil {
		return decimal.NewFromFloat(*p)
	}

	return decimal.Zero
}

// F32PtrToDecimal converts a *float32 to a decimal.Decimal.
// If p is nil it returns decimal.Zero.
func F32PtrToDecimal(p *float32) decimal.Decimal {
	if p != nil {
		return decimal.NewFromFloat(float64(*p))
	}

	return decimal.Zero
}

// StrPtr returns the string value pointed to by p, or the empty string if p is nil.
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