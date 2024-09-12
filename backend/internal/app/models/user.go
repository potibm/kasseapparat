package models

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	bcryptCost = 14
)

type User struct {
	ID uint `gorm:"primarykey" json:"id"`
	GormModel
	Username                  string  `json:"username" gorm:"unique"`
	Email                     string  `json:"email" gorm:"unique"`
	Password                  string  `json:"-"`
	Admin                     bool    `json:"admin"`
	ChangePasswordToken       *string `json:"-" gorm:"default:null"`
	ChangePasswordTokenExpiry *int64  `json:"-" gorm:"default:null"`
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
	return u.hashAndSetPassword()
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return u.hashAndSetPassword()
}

func (u *User) hashAndSetPassword() (err error) {
	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) ChangePasswordTokenIsValid(token string) bool {
	return u.ChangePasswordToken != nil && *u.ChangePasswordToken == token &&
		!u.ChangePasswordTokenIsExpired()
}

func (u *User) ChangePasswordTokenIsExpired() bool {
	currentTimestamp := time.Now().Unix()
	return u.ChangePasswordTokenExpiry != nil && *u.ChangePasswordTokenExpiry < currentTimestamp
}

func (u *User) GenerateChangePasswordToken(validity *time.Duration) {
	if validity == nil {
		duration := 15 * time.Minute
		validity = &duration
	}

	token := randomString(32)
	u.ChangePasswordToken = &token
	currentTimestamp := time.Now().Unix()
	expiry := currentTimestamp + int64(validity.Seconds())
	u.ChangePasswordTokenExpiry = &expiry
}

func (u *User) GenerateRandomPassword() {
	u.Password = randomString(32)
}

func randomString(length int) string {
	result, err := password.Generate(length, int(math.Round(float64(length)/4)), 0, false, false)
	if err != nil {
		panic(err)
	}
	return result
}
