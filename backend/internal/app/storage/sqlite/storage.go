package sqlite

import (
	"github.com/potibm/kasseapparat/internal/app/storage"
	"github.com/potibm/kasseapparat/internal/app/storage/sqlite/guestlist"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *storage.Repository {
	return &storage.Repository{
		Guestlist: guestlist.NewGuestlistRepository(db),
	}
}
