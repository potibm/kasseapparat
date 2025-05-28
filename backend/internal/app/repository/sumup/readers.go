package sumup

import (
	"context"

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
		return nil, err
	}

	return fromSDKReader(reader), nil
}

func (r *Repository) CreateReader(pairingCode string, readerName string) (*Reader, error) {
	name := readers.ReaderName(readerName)

	body := readers.CreateReaderBody{
		PairingCode: readers.ReaderPairingCode(pairingCode),
		Name:        &name,
	}
	createdReader, err := r.service.Client.Readers.Create(context.Background(), r.service.MerchantCode, body)
	if err != nil {
		return nil, err
	}

	return fromSDKReader(createdReader), nil
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
