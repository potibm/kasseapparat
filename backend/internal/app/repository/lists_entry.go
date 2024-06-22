package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

type ListEntryFilters = struct {
	Query string;
	ListID int;
	ListGroupId int;
	Present bool;
	NotPresent bool;
}

func (repo *Repository) GetListEntries(limit int, offset int, sort string, order string, ids []int, filters ListEntryFilters) ([]models.ListEntry, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getListsEntriesValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	var listEntries []models.ListEntry
	query := repo.db.Joins("List").Joins("ListGroup").Order(sort + " " + order + ", list_entries.ID ASC").Limit(limit).Offset(offset);
	
	if (len(ids) > 0) {
		query = query.Where("list_entries.ID IN ?", ids)
	}

	if filters.Query != "" {
		query = query.Where("list_entries.Name LIKE ?", "%" + filters.Query + "%")
	}
	if filters.ListID != 0 {
		query = query.Where("list_entries.list_id = ?", filters.ListID)
	}
	if filters.ListGroupId != 0 {
		query = query.Where("list_entries.list_group_id = ?", filters.ListGroupId)
	}
	if filters.Present {
		query = query.Where("list_entries.attended_guests > 0")
	}
	if filters.NotPresent {
		query = query.Where("list_entries.attended_guests = 0")
	}

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
	case "listGroup.name":
		return "ListGroup.Name", nil
	}

	return "", errors.New("Invalid field name")
}


func (repo *Repository) GetTotalListEntries() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.ListEntry{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetListEntryByID(id int) (*models.ListEntry, error) {
	var listEntry models.ListEntry
	if err := repo.db.First(&listEntry, id).Error; err != nil {
		return nil, errors.New("List Entry not found")
	}

	return &listEntry, nil
}

func (repo *Repository) UpdateListEntryByID(id int, updatedListEntry models.ListEntry) (*models.ListEntry, error) {
	var listEntry models.ListEntry
	if err := repo.db.First(&listEntry, id).Error; err != nil {
		return nil, errors.New("List Entry not found")
	}

	updatedListEntry.ID = listEntry.ID
	if updatedListEntry.ListGroup != nil {
		updatedListEntry.ListID = updatedListEntry.ListGroup.ListID
	}

	if err := repo.db.Save(&updatedListEntry).Error; err != nil {
		return nil, errors.New("Failed to update List Entry")
	}

	return &listEntry, nil
}

func (repo *Repository) CreateListEntry(listEntry models.ListEntry) (models.ListEntry, error) {
	result := repo.db.Create(&listEntry)

	if listEntry.ListGroup != nil {
		listEntry.ListID = listEntry.ListGroup.ListID
	}

	return listEntry, result.Error
}

func (repo *Repository) DeleteListEntry(listEntry models.ListEntry, deletedBy models.User) {
	repo.db.Model(&models.ListEntry{}).Where("id = ?", listEntry.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&listEntry)
}
