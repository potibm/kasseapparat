package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDeprecations_AutoHideBecomesAutoSoldOut(t *testing.T) {
	cfg := defaultTestConfig
	cfg.App.OutOfStockBehavior = "auto_hide"

	warnings := cfg.ApplyDeprecations()

	assert.Equal(t, OutOfStockAutoSoldOut, cfg.App.OutOfStockBehavior)
	assert.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], "auto_hide")
	assert.NoError(t, cfg.Validate())
}

func TestApplyDeprecations_LeavesCurrentValuesAlone(t *testing.T) {
	cfg := defaultTestConfig
	cfg.App.OutOfStockBehavior = OutOfStockFail

	assert.Empty(t, cfg.ApplyDeprecations())
	assert.Equal(t, OutOfStockFail, cfg.App.OutOfStockBehavior)
}
