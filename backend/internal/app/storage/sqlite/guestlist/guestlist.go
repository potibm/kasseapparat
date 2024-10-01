package guestlist

import (
	"context"
	"errors"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/storage/sqlite/product"
	"gorm.io/gorm"
)

type GuestlistModel struct {
	models.GormOwnedModel
	Name      string          ``
	TypeCode  bool            `gorm:"default:false"`
	ProductID uint            ``
	Product   product.Product `gorm:"foreignKey:ProductID"`
}

func (m GuestlistModel) CreateEntity() *guestlist.Guestlist {
	return &guestlist.Guestlist{
		ID:       m.ID,
		Name:     m.Name,
		TypeCode: m.TypeCode,
		Product:  m.Product.CreateEntity(),
	}
}

type GuestlistRepository struct {
	db *gorm.DB
}

func (r *GuestlistRepository) FindAll(ctx context.Context) ([]*guestlist.Guestlist, error) {
	return nil, nil
}

func (r *GuestlistRepository) FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error) {
	var list GuestlistModel
	if err := r.db.First(&list, id).Error; err != nil {
		return nil, errors.New("Guestlist not found")
	}

	return list.CreateEntity(), nil
}

func (r *GuestlistRepository) Save(ctx context.Context, guestlist *guestlist.Guestlist) error {
	return nil
}

func (r *GuestlistRepository) Update(ctx context.Context, guestlist *guestlist.Guestlist) error {
	return nil
}

func NewGuestlistRepository(db *gorm.DB) *GuestlistRepository {
	return &GuestlistRepository{
		db: db,
	}
}
