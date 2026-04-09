package sumup

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	sumup "github.com/sumup/sumup-go"
	sumupnullable "github.com/sumup/sumup-go/nullable"
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

func (r *Repository) GetReader(readerID string) (*Reader, error) {
	params := sumup.ReadersGetParams{}
	id := sumup.ReaderID(readerID)

	reader, err := r.service.Client.Readers.Get(context.Background(), r.service.MerchantCode, id, params)
	if err != nil {
		slog.Error("Error retrieving reader with ID", "reader_id", readerID, "error", err)

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

func (r *Repository) CreateReader(pairingCode, name string) (*Reader, error) {
	readerName := sumup.ReaderName(name)

	body := sumup.ReadersCreateParams{
		PairingCode: sumup.ReaderPairingCode(pairingCode),
		Name:        readerName,
	}

	createdReader, err := r.service.Client.Readers.Create(context.Background(), r.service.MerchantCode, body)
	if err != nil {
		return nil, err
	}

	return fromSDKReader(createdReader), nil
}

func (r *Repository) CreateReaderCheckout(
	readerID string,
	amount decimal.Decimal,
	description string,
	affiliateTransactionID string,
	returnURL *string,
) (*uuid.UUID, error) {
	amountStruct := sumup.CreateCheckoutRequestTotalAmount{
		Currency:  r.service.PaymentCurrency,
		Value:     getValueFromDecimal(amount, r.service.PaymentMinorUnit), // Example amount in cents (10.00 EUR)
		MinorUnit: int(r.service.PaymentMinorUnit),
	}

	var affiliateNullable *sumupnullable.Field[sumup.CreateCheckoutRequestAffiliate]

	if affiliateTransactionID != "" {
		affiliate := sumup.CreateCheckoutRequestAffiliate{
			AppID:                r.service.ApplicationID,
			Key:                  r.service.AffiliateKey,
			ForeignTransactionID: affiliateTransactionID,
		}
		affiliateNullable = sumupnullable.Value(affiliate)
	}

	body := sumup.ReadersCreateCheckoutParams{
		TotalAmount: amountStruct,
		Description: &description,
		Affiliate:   affiliateNullable,
		ReturnURL:   returnURL,
	}

	response, err := r.service.Client.Readers.CreateCheckout(
		context.Background(),
		r.service.MerchantCode,
		readerID,
		body,
	)
	if err != nil {
		slog.Error("Error creating SumUp reader checkout", "error", err)

		return nil, err
	}

	clientTransactionID, err := uuid.Parse(response.Data.ClientTransactionID)
	if err != nil {
		return nil, err
	}

	return &clientTransactionID, nil
}

func getValueFromDecimal(value decimal.Decimal, minorUnit int32) int {
	return int(value.Shift(minorUnit).IntPart())
}

func (r *Repository) CreateReaderTerminateAction(readerID string) error {
	return r.service.Client.Readers.TerminateCheckout(context.Background(), r.service.MerchantCode, readerID)
}

func (r *Repository) DeleteReader(readerID string) error {
	id := sumup.ReaderID(readerID)

	return r.service.Client.Readers.Delete(context.Background(), r.service.MerchantCode, id)
}

func fromSDKReader(sdkReader *sumup.Reader) *Reader {
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
