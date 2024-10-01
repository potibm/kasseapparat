package domain

import (
	"github.com/potibm/kasseapparat/internal/app/service"
	"github.com/potibm/kasseapparat/internal/app/service/domain/guestlist"
	"github.com/potibm/kasseapparat/internal/app/storage"
)

func NewService(repositories *storage.Repository) *service.Service {
	return &service.Service{
		GuestlistService: guestlist.NewGuestlistService(repositories.Guestlist),
	}
}
