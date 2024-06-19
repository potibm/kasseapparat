package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	bcryptCost = 14
)

type User struct {
	ID uint `gorm:"primarykey" json:"id"`
	GormModel
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"`
	Admin    bool   `json:"admin"`
}

func (u *User) Role() string {
	if u.Admin {
		return "admin"
	}	
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
	u.Password = string(bytes)

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
	u.Password = string(bytes)

	return
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

