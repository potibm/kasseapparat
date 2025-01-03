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

func (s *GuestlistService) FindAllWithParams(ctx context.Context, queryOptions storage.QueryOptions, filters storage.GuestListFilters) ([]*guestlist.Guestlist, error) {
	return s.guestListRepository.FindAllWithParams(ctx, queryOptions, filters)
}

func (s *GuestlistService) GetTotalCount(ctx context.Context) (int64, error) {
	return s.guestListRepository.GetTotalCount(ctx)
}

func NewGuestlistService(
	guestlistRepository storage.GuestlistRepository,
) *GuestlistService {
	return &GuestlistService{
		guestListRepository: guestlistRepository,
	}
}
