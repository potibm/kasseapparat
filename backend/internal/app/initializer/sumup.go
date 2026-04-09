package initializer

import (
	"sync"

	"github.com/potibm/kasseapparat/internal/app/config"
	sumupService "github.com/potibm/kasseapparat/internal/app/service/sumup"
	"github.com/sumup/sumup-go"
	"github.com/sumup/sumup-go/client"
)

var (
	instance *sumupService.Service
	once     sync.Once
)

func InitializeSumup(sumupConfig config.SumupConfig) *sumupService.Service {
	once.Do(func() {
		apiKey := sumupConfig.APIKey
		merchantCode := sumupConfig.MerchantCode
		paymentCurrency := sumupConfig.CurrencyCode
		paymentMinorUnit := sumupConfig.CurrencyMinorUnit
		affiliateKey := sumupConfig.AffiliateKey
		applicationID := sumupConfig.ApplicationID

		var webhookURL *string

		publicURL := sumupConfig.PublicURL
		if publicURL != "" {
			webhookURL = &publicURL
			*webhookURL += "/api/sumup/webhook"
		}

		clientOptions := client.WithAPIKey(apiKey)
		clnt := sumup.NewClient(clientOptions)

		instance = sumupService.NewService(
			clnt,
			merchantCode,
			applicationID,
			affiliateKey,
			paymentCurrency,
			paymentMinorUnit,
			webhookURL,
		)
	})

	return instance
}

func GetSumupService() *sumupService.Service {
	if instance == nil {
		panic("SumUp service is not initialized. Call InitializeSumup first.")
	}

	return instance
}
