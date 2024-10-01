package guestlist

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
	"github.com/potibm/kasseapparat/internal/app/storage"
)

type GuestlistService struct {
	guestListRepository storage.GuestlistRepository
}

func (s *GuestlistService) FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error) {
	return s.guestListRepository.FindByID(ctx, id)
}

func NewGuestlistService(
	guestlistRepository storage.GuestlistRepository,
) *GuestlistService {
	return &GuestlistService{
		guestListRepository: guestlistRepository,
	}
}
