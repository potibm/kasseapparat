package sumup

import "github.com/sumup/sumup-go"

type Service struct {
	Client       *sumup.Client
	MerchantCode string
}

func NewService(client *sumup.Client, merchantCode string) *Service {
	return &Service{Client: client, MerchantCode: merchantCode}
}
