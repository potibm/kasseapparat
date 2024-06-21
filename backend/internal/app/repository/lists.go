package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

func (repo *Repository) GetLists(limit int, offset int, sort string, order string) ([]models.List, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getListsValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	var lists []models.List
	if err := repo.db.Order(sort + " " + order + ", Id ASC").Limit(limit).Offset(offset).Find(&lists).Error; err != nil {
		return nil, errors.New("Lists not found")
	}

	return lists, nil
}

func getListsValidFieldName(input string) (string, error) {
	switch input {
	case "id":
		return "ID", nil
	case "name":
		return "Name", nil
	}

	return "", errors.New("Invalid field name")
}


func (repo *Repository) GetTotalLists() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.List{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetListByID(id int) (*models.List, error) {
	var list models.List
	if err := repo.db.First(&list, id).Error; err != nil {
		return nil, errors.New("List not found")
	}

	return &list, nil
}

func (repo *Repository) UpdateListByID(id int, updatedList models.List) (*models.List, error) {
	var list models.List
	if err := repo.db.First(&list, id).Error; err != nil {
		return nil, errors.New("List not found")
	}

	list.Name = updatedList.Name
	list.TypeCode = updatedList.TypeCode
	list.UpdatedByID = updatedList.UpdatedByID

	if err := repo.db.Save(&list).Error; err != nil {
		return nil, errors.New("Failed to update list")
	}

	return &list, nil
}

func (repo *Repository) CreateList(list models.List) (models.List, error) {
	result := repo.db.Create(&list)

	return list, result.Error
}

func (repo *Repository) DeleteList(list models.List, deletedBy models.User) {
	repo.db.Model(&models.List{}).Where("id = ?", list.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&list)
}
