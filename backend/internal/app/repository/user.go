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
