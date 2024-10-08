package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/gorm"
)

const ErrListEntryNotFound = "List Entry not found"

type ListEntryFilters struct {
	Query       string
	ListID      int
	ListGroupId int
	Present     bool
	NotPresent  bool
	IDs         []int
}

func (filters ListEntryFilters) AddWhere(query *gorm.DB) *gorm.DB {
	if len(filters.IDs) > 0 {
		query = query.Where("list_entries.ID IN ?", filters.IDs)
	}

	if filters.Query != "" {
		query = query.Where("list_entries.Name LIKE ? OR list_entries.Code LIKE ?", "%"+filters.Query+"%", filters.Query+"%")
	}
	if filters.ListID != 0 {
		query = query.Where("list_entries.list_id = ?", filters.ListID)
	}
	if filters.Present {
		query = query.Where("list_entries.attended_guests > 0")
	}
	if filters.NotPresent {
		query = query.Where("list_entries.attended_guests = 0")
	}

	return query
}

func (repo *Repository) GetListEntries(limit int, offset int, sort string, order string, filters ListEntryFilters) ([]models.ListEntry, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getListsEntriesValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	var listEntries []models.ListEntry
	query := repo.db.Joins("List").Order(sort + " " + order + ", list_entries.ID ASC").Limit(limit).Offset(offset)
	query = filters.AddWhere(query)

	if err := query.Find(&listEntries).Error; err != nil {
		return nil, errors.New("List Entries not found")
	}

	return listEntries, nil
}

func getListsEntriesValidFieldName(input string) (string, error) {
	switch input {
	case "id":
		return "list_entries.ID", nil
	case "name":
		return "list_entries.Name", nil
	case "list.name":
		return "List.Name", nil
	case "arrivedAt":
		return "list_entries.arrived_at", nil
	}

	return "", errors.New("Invalid field name")
}

func (repo *Repository) GetTotalListEntries(filters *ListEntryFilters) (int64, error) {
	var totalRows int64

	query := repo.db.Model(&models.ListEntry{}).Joins("List")
	if filters != nil {
		query = filters.AddWhere(query)
	}

	query.Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetUnattendedListEntriesByProductID(productId int, q string) (models.ListEntrySummarySlice, error) {
	var listEntries models.ListEntrySummarySlice
	query := repo.db.Model(&models.ListEntry{}).
		Select("list_entries.id, list_entries.name, list_entries.code, lists.name AS list_name, list_entries.additional_guests, list_entries.arrival_note").
		Joins("JOIN lists ON list_entries.list_id = lists.id").
		Joins("JOIN products ON lists.product_id = products.id").
		Where("products.id = ? AND list_entries.attended_guests = ?", productId, 0).
		Order("list_entries.name ASC")
	if q != "" {
		query = query.Where("list_entries.name LIKE ? OR code = ?", "%"+q+"%", q)
	}

	if err := query.Scan(&listEntries).Error; err != nil {
		return nil, errors.New("List Entries not found")
	}

	return listEntries, nil
}

func (repo *Repository) GetListEntryByID(id int) (*models.ListEntry, error) {
	var listEntry models.ListEntry
	if err := repo.db.First(&listEntry, id).Error; err != nil {
		return nil, errors.New(ErrListEntryNotFound)
	}

	return &listEntry, nil
}

func (repo *Repository) GetFullListEntryByID(id int) (*models.ListEntry, error) {
	var listEntry models.ListEntry
	if err := repo.db.Preload("List").Preload("List.Product").First(&listEntry, id).Error; err != nil {
		return nil, errors.New(ErrListEntryNotFound)
	}

	return &listEntry, nil
}

func (repo *Repository) UpdateListEntryByID(id int, updatedListEntry models.ListEntry) (*models.ListEntry, error) {
	var listEntry models.ListEntry
	if err := repo.db.First(&listEntry, id).Error; err != nil {
		return nil, errors.New(ErrListEntryNotFound)
	}

	updatedListEntry.ID = listEntry.ID

	if err := repo.db.Save(&updatedListEntry).Error; err != nil {
		return nil, errors.New("Failed to update List Entry")
	}

	return &updatedListEntry, nil
}

func (repo *Repository) CreateListEntry(listEntry models.ListEntry) (models.ListEntry, error) {
	result := repo.db.Create(&listEntry)

	return listEntry, result.Error
}

func (repo *Repository) DeleteListEntry(listEntry models.ListEntry, deletedBy models.User) {
	repo.db.Model(&models.ListEntry{}).Where("id = ?", listEntry.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&listEntry)
}

func (repo *Repository) GetListEntryByCode(code string) (*models.ListEntry, error) {
	var listEntry models.ListEntry
	if err := repo.db.Where("code = ?", code).First(&listEntry).Error; err != nil {
		return nil, errors.New(ErrListEntryNotFound)
	}

	return &listEntry, nil
}
