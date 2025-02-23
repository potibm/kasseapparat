package initializer

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitializeDotenv() {
	_ = godotenv.Load()
}

func GetCurrencyDecimalPlaces() int32 {
	fractionDigitsMax := 2 // Default value

	if value, exists := os.LookupEnv("FRACTION_DIGITS_MAX"); exists {
		if parsedValue, err := strconv.Atoi(value); err == nil {
			fractionDigitsMax = parsedValue
		}
	}

	return int32(fractionDigitsMax)
}
