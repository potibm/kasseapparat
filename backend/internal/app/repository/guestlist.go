package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

const ErrGuestlistNotFound = "Guestlist not found"

type GuestlistFilters = struct {
	Query string
	IDs   []int
}

func (repo *Repository) GetGuestlists(limit int, offset int, sort string, order string, filters GuestlistFilters) ([]models.Guestlist, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getListsValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	query := repo.db.Preload("Product").Order(sort + " " + order + ", Id ASC").Limit(limit).Offset(offset)

	if len(filters.IDs) > 0 {
		query = query.Where("id IN ?", filters.IDs)
	}
	if filters.Query != "" {
		query = query.Where("guestlists.Name LIKE ?", "%"+filters.Query+"%")
	}

	var guestlists []models.Guestlist
	if err := query.Find(&guestlists).Error; err != nil {
		return nil, errors.New("Guestlists not found")
	}

	return guestlists, nil
}

func getListsValidFieldName(input string) (string, error) {
	switch input {
	case "id":
		return "ID", nil
	case "name":
		return "LOWER(Name)", nil
	}

	return "", errors.New("Invalid field name")
}

func (repo *Repository) GetTotalGuestlists() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.Guestlist{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetGuestlistByID(id int) (*models.Guestlist, error) {
	var guestlist models.Guestlist
	if err := repo.db.First(&guestlist, id).Error; err != nil {
		return nil, errors.New(ErrGuestlistNotFound)
	}

	return &guestlist, nil
}

func (repo *Repository) GetGuestlistWithTypeCode() (*models.Guestlist, error) {
	var guestlist models.Guestlist
	if err := repo.db.Where("type_code = ?", "1").First(&guestlist).Error; err != nil {
		return nil, errors.New(ErrGuestlistNotFound)
	}
	return &guestlist, nil
}

func (repo *Repository) UpdateGuestlistByID(id int, updatedGuestlist models.Guestlist) (*models.Guestlist, error) {
	var guestlist models.Guestlist
	if err := repo.db.First(&guestlist, id).Error; err != nil {
		return nil, errors.New(ErrGuestlistNotFound)
	}

	guestlist.Name = updatedGuestlist.Name
	guestlist.TypeCode = updatedGuestlist.TypeCode
	guestlist.ProductID = updatedGuestlist.ProductID
	guestlist.UpdatedByID = updatedGuestlist.UpdatedByID

	if err := repo.db.Save(&guestlist).Error; err != nil {
		return nil, errors.New("Failed to update guestlist")
	}

	return &guestlist, nil
}

func (repo *Repository) CreateGuestlist(guestlist models.Guestlist) (models.Guestlist, error) {
	result := repo.db.Create(&guestlist)

	return guestlist, result.Error
}

func (repo *Repository) DeleteGuestlist(guestlist models.Guestlist, deletedBy models.User) {
	repo.db.Model(&models.Guestlist{}).Where("id = ?", guestlist.ID).Update("DeletedByID", deletedBy.ID)
	repo.db.Model(&models.ListEntry{}).Where("guestlist_id = ?", guestlist.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&models.ListEntry{}, "guestlist_id = ?", guestlist.ID)
	repo.db.Delete(&guestlist)
}
