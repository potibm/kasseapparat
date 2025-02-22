package models

import (
	"os"
	"strconv"
)

func getFractionDigitsMax() int32 {
	fractionDigitsMax := 2 // Default value

	if value, exists := os.LookupEnv("FRACTION_DIGITS_MAX"); exists {
		if parsedValue, err := strconv.Atoi(value); err == nil {
			fractionDigitsMax = parsedValue
		}
	}

	return int32(fractionDigitsMax)
}
