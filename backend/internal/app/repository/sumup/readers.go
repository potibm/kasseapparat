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
	params := readers.GetParams{}
	id := readers.ReaderID(readerId)

	reader, err := r.service.Client.Readers.Get(context.Background(), r.service.MerchantCode, id, params)
	if err != nil {
		log.Printf("Error retrieving reader with ID %s: %v", readerId, err)

		if isReaderNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return fromSDKReader(reader), nil
}

func isReaderNotFoundError(err error) bool {
	return err != nil && err.Error() == "The requested Reader resource does not exists."
}

func (r *Repository) CreateReader(pairingCode string, name string) (*Reader, error) {
	readerName := readers.ReaderName(name)

	body := readers.Create{
		PairingCode: readers.ReaderPairingCode(pairingCode),
		Name:        readerName,
	}

	createdReader, err := r.service.Client.Readers.Create(context.Background(), r.service.MerchantCode, body)
	if err != nil {
		return nil, err
	}

	return fromSDKReader(createdReader), nil
}

func (r *Repository) CreateReaderCheckout(readerId string, amount decimal.Decimal, description string, affiliateTransactionId string, returnUrl *string) (*uuid.UUID, error) {
	amountStruct := readers.CreateCheckoutTotalAmount{
		Currency:  r.service.PaymentCurrency,
		Value:     getValueFromDecimal(amount, int(r.service.PaymentMinorUnit)), // Example amount in cents (10.00 EUR)
		MinorUnit: int(r.service.PaymentMinorUnit),
	}

	var affiliate *readers.CreateCheckoutAffiliate
	if affiliateTransactionId != "" {
		affiliate = &readers.CreateCheckoutAffiliate{
			AppID:                r.service.ApplicationId,
			Key:                  r.service.AffiliateKey,
			ForeignTransactionID: affiliateTransactionId,
		}
	}

	body := readers.CreateCheckout{
		TotalAmount: amountStruct,
		Description: &description,
		Affiliate:   affiliate,
		ReturnURL:   returnUrl,
	}

	response, err := r.service.Client.Readers.CreateCheckout(context.Background(), r.service.MerchantCode, readerId, body)
	if err != nil {
		log.Printf("Error creating SumUp reader checkout: %v", err)

		return nil, err
	}

	clientTransactionId, err := uuid.Parse(response.Data.ClientTransactionID)
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
	id := readers.ReaderID(readerId)

	return r.service.Client.Readers.Delete(context.Background(), r.service.MerchantCode, id)
}

func fromSDKReader(sdkReader *readers.Reader) *Reader {
	return &Reader{
		ID:               string(sdkReader.ID),
		Name:             string(sdkReader.Name),
		Status:           string(sdkReader.Status),
		DeviceIdentifier: string(sdkReader.Device.Identifier),
		DeviceModel:      string(sdkReader.Device.Model),
		CreatedAt:        sdkReader.CreatedAt,
		UpdatedAt:        sdkReader.UpdatedAt,
	}
}
