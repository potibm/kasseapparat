package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

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
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	PasswordChangeRequired bool `json:"-" gorm:"default:false"`
	Admin    bool   `json:"admin"`
}

func (u *User) Role() string {
	if u.Admin {
		return "admin"
	}	
	return "user"
}

func (u *User) GravatarURL() string {
	hasher := sha256.Sum256([]byte(strings.TrimSpace(u.Email)))
    hash := hex.EncodeToString(hasher[:])
	
	return "https://www.gravatar.com/avatar/" + hash
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

