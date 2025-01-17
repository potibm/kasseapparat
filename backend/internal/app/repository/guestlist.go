package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

const ErrGuestlistNotFound = "Guestlist not found"

type ListFilters = struct {
	Query string
	IDs   []int
}

func (repo *Repository) GetGuestlists(limit int, offset int, sort string, order string, filters ListFilters) ([]models.Guestlist, error) {
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
		query = query.Where("lists.Name LIKE ?", "%"+filters.Query+"%")
	}

	var guestlists []models.Guestlist
	if err := query.Find(&guestlists).Error; err != nil {
		return nil, errors.New("Lists not found")
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

func (repo *Repository) UpdateListByID(id int, updatedList models.Guestlist) (*models.Guestlist, error) {
	var guestlist models.Guestlist
	if err := repo.db.First(&guestlist, id).Error; err != nil {
		return nil, errors.New(ErrGuestlistNotFound)
	}

	guestlist.Name = updatedList.Name
	guestlist.TypeCode = updatedList.TypeCode
	guestlist.ProductID = updatedList.ProductID
	guestlist.UpdatedByID = updatedList.UpdatedByID

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
	repo.db.Model(&models.ListEntry{}).Where("list_id = ?", guestlist.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&models.ListEntry{}, "list_id = ?", guestlist.ID)
	repo.db.Delete(&guestlist)
}
