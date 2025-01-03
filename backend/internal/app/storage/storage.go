package storage

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
)

type GuestlistRepository interface {
	FindAll(ctx context.Context) ([]*guestlist.Guestlist, error)
	FindAllWithParams(ctx context.Context, queryOptions QueryOptions, filters GuestListFilters) ([]*guestlist.Guestlist, error)
	GetTotalCount(ctx context.Context) (int64, error)
	FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error)
	Save(ctx context.Context, guestlist *guestlist.Guestlist) error
	Update(ctx context.Context, guestlist *guestlist.Guestlist) error
}

type Repository struct {
	Guestlist GuestlistRepository
}
