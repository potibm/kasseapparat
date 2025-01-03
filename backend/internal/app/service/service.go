package service

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
	"github.com/potibm/kasseapparat/internal/app/storage"
)

type GuestlistService interface {
	FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error)
	FindAllWithParams(ctx context.Context, queryOptions storage.QueryOptions, filters storage.GuestListFilters) ([]*guestlist.Guestlist, error)
	GetTotalCount(ctx context.Context) (int64, error)
}

type Service struct {
	GuestlistService GuestlistService
}
