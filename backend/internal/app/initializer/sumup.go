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
		apiKey := sumupConfig.ApiKey
		merchantCode := sumupConfig.MerchantCode
		paymentCurrency := sumupConfig.CurrencyCode
		paymentMinorUnit := sumupConfig.CurrencyMinorUnit
		affiliateKey := sumupConfig.AffiliateKey
		applicationId := sumupConfig.ApplicationId

		var webhookUrl *string

		publicUrl := sumupConfig.PublicUrl
		if publicUrl != "" {
			webhookUrl = &publicUrl
			*webhookUrl += "/api/sumup/webhook"
		}

		clientOptions := client.WithAPIKey(apiKey)
		client := sumup.NewClient(clientOptions)

		instance = sumupService.NewService(
			client,
			merchantCode,
			applicationId,
			affiliateKey,
			paymentCurrency,
			uint(paymentMinorUnit),
			webhookUrl,
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
