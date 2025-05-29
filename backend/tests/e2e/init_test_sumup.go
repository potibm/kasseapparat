package tests_e2e

import "github.com/potibm/kasseapparat/internal/app/repository/sumup"

type MockSumUpRepository struct {
	GetReadersFunc   func() ([]sumup.Reader, error)
	GetReaderFunc    func(readerId string) (*sumup.Reader, error)
	CreateReaderFunc func(pairingCode string, readerName string) (*sumup.Reader, error)
	DeleteReaderFunc func(readerId string) error
}

func (m *MockSumUpRepository) GetReaders() ([]sumup.Reader, error) {
	return m.GetReadersFunc()
}

func (m *MockSumUpRepository) GetReader(readerId string) (*sumup.Reader, error) {
	return m.GetReaderFunc(readerId)
}

func (m *MockSumUpRepository) CreateReader(pairingCode string, readerName string) (*sumup.Reader, error) {
	return m.CreateReaderFunc(pairingCode, readerName)
}

func (m *MockSumUpRepository) DeleteReader(readerId string) error {
	return m.DeleteReaderFunc(readerId)
}
