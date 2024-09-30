package storage

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
)

type GuestlistRepository interface {
	FindAll(ctx context.Context) ([]*guestlist.Guestlist, error)
	FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error)
	Save(ctx context.Context, guestlist *guestlist.Guestlist) error
	Update(ctx context.Context, guestlist *guestlist.Guestlist) error
}

type Repository struct {
	Guestlist GuestlistRepository
}