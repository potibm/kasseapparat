package sumup

import "github.com/sumup/sumup-go"

type Service struct {
	Client           *sumup.Client
	MerchantCode     string
	ApplicationId    string
	AffiliateKey     string
	PaymentCurrency  string
	PaymentMinorUnit uint
}

func NewService(client *sumup.Client, merchantCode string, applicationId string, affiliateKey string, paymentCurrency string, paymentMinorUnit uint) *Service {
	return &Service{Client: client, MerchantCode: merchantCode, ApplicationId: applicationId, AffiliateKey: affiliateKey, PaymentCurrency: paymentCurrency, PaymentMinorUnit: paymentMinorUnit}
}
