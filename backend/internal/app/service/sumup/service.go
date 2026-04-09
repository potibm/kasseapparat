package sumup

import "github.com/sumup/sumup-go"

type Service struct {
	Client           *sumup.Client
	MerchantCode     string
	ApplicationID    string
	AffiliateKey     string
	PaymentCurrency  string
	PaymentMinorUnit int32
	WebhookURL       *string
}

func NewService(
	client *sumup.Client,
	merchantCode string,
	applicationID string,
	affiliateKey string,
	paymentCurrency string,
	paymentMinorUnit int32,
	webhookURL *string,
) *Service {
	return &Service{
		Client:           client,
		MerchantCode:     merchantCode,
		ApplicationID:    applicationID,
		AffiliateKey:     affiliateKey,
		PaymentCurrency:  paymentCurrency,
		PaymentMinorUnit: paymentMinorUnit,
		WebhookURL:       webhookURL,
	}
}
