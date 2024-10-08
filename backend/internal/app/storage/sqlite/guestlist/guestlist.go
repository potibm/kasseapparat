package guestlist

import (
	"context"
	"errors"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/storage/sqlite"
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


type GuestListFilters struct {
	Query string
	IDs []uint     
}


func (m GuestlistModel) CreateEntity() *guestlist.Guestlist {
	return &guestlist.Guestlist{
		ID:       m.ID,
		Name:     m.Name,
		TypeCode: m.TypeCode,
		Product:  m.Product.CreateEntity(),
	}
}

func getValidFieldName(SortBy string) (string, error) {
	switch SortBy {
	case "id":
		return "ID", nil
	case "name":
		return "LOWER(Name)", nil
	}

	return "", errors.New("Invalid field name")
}

type GuestlistRepository struct {
	db *gorm.DB
}

func (r *GuestlistRepository) FindAll(ctx context.Context) ([]*guestlist.Guestlist, error) {
	return nil, nil
}

func (r *GuestlistRepository) FindAllWithParams(ctx context.Context, queryOptions sqlite.QueryOptions, filters GuestListFilters) ([]*guestlist.Guestlist, error) {
	order := "DESC"
	if queryOptions.SortAsc{
		order = "ASC"
	}
	sort, err := getValidFieldName(queryOptions.SortBy)
	if err != nil {
		return nil, err
	}

	query := r.db.Preload("Product").Order(sort + " " + order + ", Id ASC").Limit(queryOptions.Limit).Offset(queryOptions.Offset)

	if len(filters.IDs) > 0 {
		query = query.Where("id IN ?", filters.IDs)
	}
	if filters.Query != "" {
		query = query.Where("lists.Name LIKE ?", "%"+filters.Query+"%")
	}

	var lists []GuestlistModel
	if err := query.Find(&lists).Error; err != nil {
		return nil, errors.New("Lists not found")
	}

	var resultList []*guestlist.Guestlist
	for _, list := range lists {
		resultList = append(resultList, list.CreateEntity())
	}

	return resultList, nil


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
