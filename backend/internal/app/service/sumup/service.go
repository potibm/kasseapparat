package sumup

import "github.com/sumup/sumup-go"

type Service struct {
	Client           *sumup.Client
	MerchantCode     string
	PaymentCurrency  string
	PaymentMinorUnit uint
}

func NewService(client *sumup.Client, merchantCode string, paymentCurrency string, paymentMinorUnit uint) *Service {
	return &Service{Client: client, MerchantCode: merchantCode, PaymentCurrency: paymentCurrency, PaymentMinorUnit: paymentMinorUnit}
}
