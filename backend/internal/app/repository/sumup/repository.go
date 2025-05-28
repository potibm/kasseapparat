package sumup

import (
	sumupService "github.com/potibm/kasseapparat/internal/app/service/sumup"
)

type Repository struct {
	service *sumupService.Service
}

func NewRepository(service *sumupService.Service) *Repository {
	return &Repository{
		service: service,
	}
}
