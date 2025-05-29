package initializer

import (
	"os"
	"strconv"
	"sync"

	sumupSevice "github.com/potibm/kasseapparat/internal/app/service/sumup"
	"github.com/sumup/sumup-go"
	"github.com/sumup/sumup-go/client"
)

var (
	instance *sumupSevice.Service
	once     sync.Once
)

func InitializeSumup() *sumupSevice.Service {
	once.Do(func() {
		apiKey := getEnv("SUMUP_API_KEY", "")
		merchantCode := getEnv("SUMUP_MERCHANT_CODE", "")
		paymentCurrency := getEnv("CURRENCY_CODE", "DKK")
		paymentMinorUnit := getEnvAsInt("FRACTION_DIGITS_MAX", 2)

		options := client.New()
		clientOptions := options.WithAPIKey(apiKey)

		client := sumup.NewClient(clientOptions)

		instance = sumupSevice.NewService(client, merchantCode, paymentCurrency, uint(paymentMinorUnit))
	})

	return instance
}

func GetSumupService() *sumupSevice.Service {
	if instance == nil {
		panic("SumUp service is not initialized. Call InitializeSumup first.")
	}

	return instance
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
