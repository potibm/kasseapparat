package repository

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
)

func (repo *Repository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).First(&user, id).Error; err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func (repo *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).Where("Username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func (repo *Repository) GetUserByUsernameAndPassword(username string, password string) (*models.User, error) {
	user, err := repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	err = user.ComparePassword(password)
	if err != nil {
		return nil, errors.New("Invalid password")
	}

	return user, nil
}

func (repo *Repository) GetUsers(limit int, offset int, sort string, order string) ([]models.User, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getUsersValidFieldName(sort)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := repo.db.Order(sort + " " + order + ", ID ASC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, errors.New("Users not found")
	}

	return users, nil
}

func getUsersValidFieldName(input string) (string, error) {
	switch input {
	case "id":
		return "ID", nil
	case "username":
		return "Username", nil
	}

	return "", errors.New("Invalid field name")
}

func (repo *Repository) GetTotalUsers() (int64, error) {
	var totalRows int64
	repo.db.Model(&models.User{}).Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) CreateUser(user models.User) (models.User, error) {
	result := repo.db.Create(&user)

	return user, result.Error
}

func (repo *Repository) DeleteUser(user models.User) {
	repo.db.Delete(&user)
}

func (repo *Repository) UpdateUserByID(id int, updatedUser models.User) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, errors.New("User not found")
	}

	// Update the product with the new values
	user.Username = updatedUser.Username
	user.Admin = updatedUser.Admin
	if updatedUser.Password != "" {
		user.Password = updatedUser.Password
	}

	// Save the updated product to the database
	if err := repo.db.Save(&user).Error; err != nil {
		return nil, errors.New("Failed to update user")
	}

	return &user, nil
}
