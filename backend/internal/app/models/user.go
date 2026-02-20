package models

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 14
const passwordLength = 32
const defaultChangePasswordTokenValidity = 15 * time.Minute

type User struct {
	GormModel

	ID                        uint    `json:"id"       gorm:"primarykey"`
	Username                  string  `json:"username" gorm:"unique"`
	Email                     string  `json:"email"    gorm:"unique"`
	Password                  string  `json:"-"`
	Admin                     bool    `json:"admin"`
	ChangePasswordToken       *string `json:"-"        gorm:"default:null"`
	ChangePasswordTokenExpiry *int64  `json:"-"        gorm:"default:null"`
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

func (u *User) SetPassword(password string) error {
	var err error

	u.Password, err = hashPassword(password)
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
		duration := defaultChangePasswordTokenValidity
		validity = &duration
	}

	token := randomString(passwordLength)
	u.ChangePasswordToken = &token
	currentTimestamp := time.Now().Unix()
	expiry := currentTimestamp + int64(validity.Seconds())
	u.ChangePasswordTokenExpiry = &expiry
}

func (u *User) GenerateRandomPassword() {
	u.Password = randomString(passwordLength)
}

func randomString(length int) string {
	const fractionOfDigits = 0.25

	numberOfDigitsInPassword := int(math.Round(float64(length) * fractionOfDigits))

	result, err := password.Generate(length, numberOfDigitsInPassword, 0, false, false)
	if err != nil {
		panic(err)
	}

	return result
}
