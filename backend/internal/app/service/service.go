package service

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
)

type GuestlistService interface {
	FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error)
}

type Service struct {
	GuestlistService GuestlistService
}
