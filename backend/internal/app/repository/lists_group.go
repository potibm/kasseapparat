package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

func (repo *Repository) GetListsGroups(limit int, offset int, sort string, order string, ids []int) ([]models.ListGroup, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getListGroupsValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	query := repo.db.Preload("List").Order(sort + " " + order + ", Id ASC").Limit(limit).Offset(offset);

	if (len(ids) > 0) {
		query = query.Where("id IN ?", ids)
	}

	var listGroups []models.ListGroup
	if err := query.Find(&listGroups).Error; err != nil {
		return nil, errors.New("List Groups not found")
	}

	return listGroups, nil
}

func getListGroupsValidFieldName(input string) (string, error) {
	switch input {
	case "id":
		return "ID", nil
	case "name":
		return "Name", nil
	}

	return "", errors.New("Invalid field name")
}


func (repo *Repository) GetTotalListGroups() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.ListGroup{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetListGroupByID(id int) (*models.ListGroup, error) {
	var listGroup models.ListGroup
	if err := repo.db.First(&listGroup, id).Error; err != nil {
		return nil, errors.New("List Group not found")
	}

	return &listGroup, nil
}

func (repo *Repository) UpdateListGroupByID(id int, updatedListGroup models.ListGroup) (*models.ListGroup, error) {
	var listGroup models.ListGroup
	if err := repo.db.First(&listGroup, id).Error; err != nil {
		return nil, errors.New("List Group not found")
	}

	updatedListGroup.ID = listGroup.ID

	if err := repo.db.Save(&updatedListGroup).Error; err != nil {
		return nil, errors.New("Failed to update List Group")
	}

	return &listGroup, nil
}

func (repo *Repository) CreateListGroup(listGroup models.ListGroup) (models.ListGroup, error) {
	result := repo.db.Create(&listGroup)

	return listGroup, result.Error
}

func (repo *Repository) DeleteListGroup(listGroup models.ListGroup, deletedBy models.User) {
	repo.db.Model(&models.ListGroup{}).Where("id = ?", listGroup.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&listGroup)
}
