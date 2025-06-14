package sumup

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sumup/sumup-go/readers"
)

func (r *Repository) GetReaders() ([]Reader, error) {
	result := []Reader{}

	readers, err := r.service.Client.Readers.List(context.Background(), r.service.MerchantCode)
	if err != nil {
		return nil, err
	}

	for _, sdkReader := range readers.Items {
		result = append(result, *fromSDKReader(&sdkReader))
	}

	return result, nil
}

func (r *Repository) GetReader(readerId string) (*Reader, error) {
	params := readers.GetReaderParams{}
	id := readers.ReaderId(readerId)

	reader, err := r.service.Client.Readers.Get(context.Background(), r.service.MerchantCode, id, params)
	if err != nil {
		log.Printf("Error retrieving reader with ID %s: %v", readerId, err)

		if err.Error() == "The requested Reader resource does not exists." {
			return nil, nil
		}

		return nil, err
	}

	return fromSDKReader(reader), nil
}

func (r *Repository) CreateReader(pairingCode string, name string) (*Reader, error) {
	readerName := readers.ReaderName(name)

	body := readers.CreateReaderBody{
		PairingCode: readers.ReaderPairingCode(pairingCode),
		Name:        &readerName,
	}

	createdReader, err := r.service.Client.Readers.Create(context.Background(), r.service.MerchantCode, body)
	if err != nil {
		return nil, err
	}

	return fromSDKReader(createdReader), nil
}

func (r *Repository) CreateReaderCheckout(readerId string, amount decimal.Decimal, description string, affiliateTransactionId string, returnUrl string) (*uuid.UUID, error) {
	amountStruct := readers.CreateReaderCheckoutAmount{
		Currency:  r.service.PaymentCurrency,
		Value:     getValueFromDecimal(amount, int(r.service.PaymentMinorUnit)), // Example amount in cents (10.00 EUR)
		MinorUnit: int(r.service.PaymentMinorUnit),
	}

	var affiliate *readers.Affiliate
	if affiliateTransactionId != "" {
		affiliate = &readers.Affiliate{
			AppId:                r.service.ApplicationId,
			Key:                  r.service.AffiliateKey,
			ForeignTransactionId: affiliateTransactionId,
		}
	}

	body := readers.CreateReaderCheckoutBody{
		TotalAmount: amountStruct,
		Description: &description,
		Affiliate:   affiliate,
		ReturnUrl:   &returnUrl,
	}

	response, err := r.service.Client.Readers.CreateCheckout(context.Background(), r.service.MerchantCode, readerId, body)
	if err != nil {
		return nil, err
	}

	clientTransactionId, err := uuid.Parse(*response.Data.ClientTransactionId)
	if err != nil {
		return nil, err
	}

	return &clientTransactionId, nil
}

func getValueFromDecimal(value decimal.Decimal, minorUnit int) int {
	return int(value.Shift(int32(minorUnit)).IntPart())
}

func (r *Repository) CreateReaderTerminateAction(readerId string) error {
	return r.service.Client.Readers.TerminateCheckout(context.Background(), r.service.MerchantCode, readerId)
}

func (r *Repository) DeleteReader(readerId string) error {
	id := readers.ReaderId(readerId)
	return r.service.Client.Readers.DeleteReader(context.Background(), r.service.MerchantCode, id)
}

func fromSDKReader(sdkReader *readers.Reader) *Reader {
	return &Reader{
		ID:               string(sdkReader.Id),
		Name:             string(sdkReader.Name),
		Status:           string(sdkReader.Status),
		DeviceIdentifier: string(sdkReader.Device.Identifier),
		DeviceModel:      string(sdkReader.Device.Model),
		CreatedAt:        sdkReader.CreatedAt,
		UpdatedAt:        sdkReader.UpdatedAt,
	}
}
