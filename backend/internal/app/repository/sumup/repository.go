package sumup

import (
	sumupService "github.com/potibm/kasseapparat/internal/app/service/sumup"
)

type Repository struct {
	service *sumupService.Service
}

var _ RepositoryInterface = (*Repository)(nil)

type RepositoryInterface interface {
	GetReaders() ([]Reader, error)
	GetReader(readerId string) (*Reader, error)
	CreateReader(pairingCode string, readerName string) (*Reader, error)
	DeleteReader(readerId string) error
}

func NewRepository(service *sumupService.Service) RepositoryInterface {
	return &Repository{
		service: service,
	}
}
