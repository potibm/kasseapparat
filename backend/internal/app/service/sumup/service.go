package sumup

import "github.com/sumup/sumup-go"

type Service struct {
	Client           *sumup.Client
	MerchantCode     string
	ApplicationId    string
	AffiliateKey     string
	PaymentCurrency  string
	PaymentMinorUnit uint
	ApiKey           string // remove this, once the GetPurchaseById method works with the updated SDK
}

func NewService(client *sumup.Client, merchantCode string, applicationId string, affiliateKey string, paymentCurrency string, paymentMinorUnit uint, apiKey string) *Service {
	return &Service{Client: client, MerchantCode: merchantCode, ApplicationId: applicationId, AffiliateKey: affiliateKey, PaymentCurrency: paymentCurrency, PaymentMinorUnit: paymentMinorUnit, ApiKey: apiKey}
}
