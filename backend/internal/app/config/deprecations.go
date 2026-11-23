package config

// deprecatedOutOfStockAutoHide hid depleted products. It was removed in v4 (ADR 007):
// a sold-out card is useful on the POS (it collects interest) and takes little room.
// Existing configs keep starting; the value now works like auto_sold_out.
const deprecatedOutOfStockAutoHide OutOfStockBehavior = "auto_hide"

// ApplyDeprecations rewrites settings that were removed to their replacement, so an
// old config keeps working. It returns one warning per rewrite for the log.
func (c *Config) ApplyDeprecations() []string {
	var warnings []string

	if c.App.OutOfStockBehavior == deprecatedOutOfStockAutoHide {
		c.App.OutOfStockBehavior = OutOfStockAutoSoldOut

		warnings = append(warnings,
			"out_of_stock_behavior \"auto_hide\" was removed and now works like \"auto_sold_out\"; "+
				"please update your config")
	}

	return warnings
}
