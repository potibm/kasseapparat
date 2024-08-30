package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/gorm"
)

const ErrUserNotFound = "User not found"

type UserFilters struct {
	Query   string
	IsAdmin bool
	IDs     []int
}

func (filters UserFilters) AddWhere(query *gorm.DB) *gorm.DB {
	if len(filters.IDs) > 0 {
		query = query.Where("list_entries.ID IN ?", filters.IDs)
	}
	if filters.Query != "" {
		query = query.Where("Username LIKE ? OR Email LIKE ?", "%"+filters.Query+"%", "%"+filters.Query+"%")
	}

	if filters.IsAdmin {
		query = query.Where("Admin = ?", true)
	}

	return query
}

func (repo *Repository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).First(&user, id).Error; err != nil {
		return nil, errors.New(ErrUserNotFound)
	}

	return &user, nil
}

func (repo *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).Where("LOWER(Username) = ?", strings.ToLower(username)).First(&user).Error; err != nil {
		return nil, errors.New(ErrUserNotFound + " by username")
	}

	return &user, nil
}

func (repo *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).Where("LOWER(Email) = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, errors.New(ErrUserNotFound + " by email")
	}

	return &user, nil
}

func (repo *Repository) GetUserByUserameOrEmail(login string) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).Where("LOWER(Username) = ? OR LOWER(Email) = ?", strings.ToLower(login), strings.ToLower(login)).First(&user).Error; err != nil {
		return nil, errors.New(ErrUserNotFound)
	}

	return &user, nil
}

func (repo *Repository) GetUserByLoginAndPassword(login string, password string) (*models.User, error) {
	user, err := repo.GetUserByUserameOrEmail(login)
	if err != nil {
		return nil, err
	}

	err = user.ComparePassword(password)
	if err != nil {
		return nil, errors.New("Invalid password")
	}

	return user, nil
}

func (repo *Repository) GetUsers(limit int, offset int, sort string, order string, filters UserFilters) ([]models.User, error) {
	sort, err := getUsersValidFieldName(sort)
	if err != nil {
		return nil, err
	}
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	query := repo.db.Order(sort + " " + order + ", ID ASC").Limit(limit).Offset(offset)
	query = filters.AddWhere(query)

	var users []models.User
	if err := query.Find(&users).Error; err != nil {
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

func (repo *Repository) GetTotalUsers(filters *UserFilters) (int64, error) {
	var totalRows int64
	query := repo.db.Model(&models.User{})
	if filters != nil {
		query = filters.AddWhere(query)
	}
	query.Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) CreateUser(user models.User) (models.User, error) {
	user.Username = strings.ToLower(user.Username)

	result := repo.db.Create(&user)

	return user, result.Error
}

func (repo *Repository) DeleteUser(user models.User) {
	// update the user to be deleted: postfix the username with "_deleted" and the current timestamp, and prefix the email with "deleted_" and the current timestamp
	now := time.Now().Format("20060102150405")
	user.Username = user.Username + "_deleted_" + now
	user.Email = "deleted_" + now + "_" + user.Email
	repo.db.Save(&user)

	repo.db.Delete(&user)
}

func (repo *Repository) UpdateUserByID(id int, updatedUser models.User) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, errors.New(ErrUserNotFound)
	}

	// Update the product with the new values
	user.Username = strings.ToLower(updatedUser.Username)
	user.Admin = updatedUser.Admin
	user.Email = updatedUser.Email
	user.ChangePasswordToken = updatedUser.ChangePasswordToken
	user.ChangePasswordTokenExpiry = updatedUser.ChangePasswordTokenExpiry
	if updatedUser.Password != "" {
		user.Password = updatedUser.Password
	}

	// Save the updated product to the database
	if err := repo.db.Save(&user).Error; err != nil {
		return nil, errors.New("Failed to update user")
	}

	return &user, nil
}
