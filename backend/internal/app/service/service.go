package service

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
	"github.com/potibm/kasseapparat/internal/app/storage"
)

type GuestlistService interface {
	// FindByID retrieves the guest list by its ID.
	FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error)
	// FindAllWithParams returns guest lists using queryOptions and filters.
	FindAllWithParams(ctx context.Context, queryOptions storage.QueryOptions, filters storage.GuestListFilters) ([]*guestlist.Guestlist, error)
	// GetTotalCount returns the number of guest lists in storage.
	GetTotalCount(ctx context.Context) (int64, error)
	// Save creates a new guest list.
	Save(ctx context.Context, guestlist *guestlist.Guestlist) (*guestlist.Guestlist, error)
	// Update updates an existing guest list.
	Update(ctx context.Context, guestlist *guestlist.Guestlist) (*guestlist.Guestlist, error)	
	// Delete removes a guest list from storage.
	Delete(ctx context.Context, guestlistID int, deletedByID int) error
}

type Service struct {
	GuestlistService GuestlistService
}
