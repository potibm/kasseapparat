package guestlist

import (
	"context"
	"errors"

	"github.com/potibm/kasseapparat/internal/app/entities/guestlist"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/storage"
	"github.com/potibm/kasseapparat/internal/app/storage/sqlite/product"
	"gorm.io/gorm"
)

const ErrGuestlistNotFound = "Guestlist not found"

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

func (r *GuestlistRepository) FindAllWithParams(ctx context.Context, queryOptions storage.QueryOptions, filters storage.GuestListFilters) ([]*guestlist.Guestlist, error) {
	order := "DESC"
	if queryOptions.SortAsc {
		order = "ASC"
	}
	sort, err := getValidFieldName(queryOptions.SortBy)
	if err != nil {
		return nil, err
	}

	query := r.db.WithContext(ctx).Preload("Product").Order(sort + " " + order + ", Id ASC").Limit(queryOptions.Limit).Offset(queryOptions.Offset)

	if len(filters.IDs) > 0 {
		query = query.Where("id IN ?", filters.IDs)
	}
	if filters.Query != "" {
		query = query.Where("lists.Name LIKE ?", "%"+filters.Query+"%")
	}

	var lists []GuestlistModel
	if err := query.Find(&lists).Error; err != nil {
		return nil, errors.New("Guestlists not found")
	}

	var resultList []*guestlist.Guestlist
	for _, list := range lists {
		resultList = append(resultList, list.CreateEntity())
	}

	return resultList, nil
}

func (r *GuestlistRepository) FindByID(ctx context.Context, id int) (*guestlist.Guestlist, error) {
	var list GuestlistModel
	if err := r.db.WithContext(ctx).First(&list, id).Error; err != nil {
		return nil, errors.New(ErrGuestlistNotFound)
	}

	return list.CreateEntity(), nil
}

func (r *GuestlistRepository) GetTotalCount(ctx context.Context) (int64, error) {
	var totalRows int64
	r.db.WithContext(ctx).Model(&GuestlistModel{}).Count(&totalRows)

	return totalRows, nil
}

func (r *GuestlistRepository) Save(ctx context.Context, guestlist *guestlist.Guestlist) (*guestlist.Guestlist, error) {
	var guestlistModel GuestlistModel

	guestlistModel.Name = guestlist.Name
	guestlistModel.TypeCode = guestlist.TypeCode
	guestlistModel.ProductID = guestlist.Product.ID

	result := r.db.WithContext(ctx).Create(&guestlistModel)
	if result.Error != nil {
		return nil, result.Error
	}

	guestlist.ID = guestlistModel.ID

	return guestlist, nil
}

func (r *GuestlistRepository) Update(ctx context.Context, guestlist *guestlist.Guestlist) (*guestlist.Guestlist, error) {
	var dbGuestlist GuestlistModel
	if err := r.db.WithContext(ctx).First(&dbGuestlist, guestlist.ID).Error; err != nil {
		return nil, errors.New(ErrGuestlistNotFound)
	}

	dbGuestlist.Name = guestlist.Name
	dbGuestlist.TypeCode = guestlist.TypeCode
	dbGuestlist.ProductID = guestlist.Product.ID

	result := r.db.WithContext(ctx).Save(&dbGuestlist)
	if result.Error != nil {
		return nil, result.Error
	}
	return dbGuestlist.CreateEntity(), nil
}

func (r *GuestlistRepository) Delete(ctx context.Context, guestlistID int, deletedByID int) error {
	var dbGuestlist GuestlistModel
	if err := r.db.WithContext(ctx).First(&dbGuestlist, guestlistID).Error; err != nil {
		return errors.New(ErrGuestlistNotFound)
	}

	if guestlistID != 0 {
		r.db.WithContext(ctx).Model(GuestlistModel{}).Where("id = ?", guestlistID).Update("DeletedByID", deletedByID)
	}

	/*
		This should probably go into a separate repository
		r.db.Model(&models.ListEntry{}).Where("list_id = ?", list.ID).Update("DeletedByID", deletedBy.ID)
		r.db.Delete(&models.ListEntry{}, "list_id = ?", list.ID)
	*/

	r.db.WithContext(ctx).Delete(&dbGuestlist)

	return nil
}

func NewGuestlistRepository(db *gorm.DB) *GuestlistRepository {
	return &GuestlistRepository{
		db: db,
	}
}
