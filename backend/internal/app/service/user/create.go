package user

import (
	"fmt"
	"time"

	"github.com/potibm/kasseapparat/internal/app/models"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
}

type Mailer interface {
	SendNewUserTokenMail(
		email string,
		userID uint,
		username string,
		token string,
	) error // Passe den Typ von userID an, falls nötig (z.B. int oder string)
}

type UserService struct {
	repo   UserRepository
	mailer Mailer
}

func NewUserService(repo UserRepository, mailer Mailer) *UserService {
	return &UserService{
		repo:   repo,
		mailer: mailer,
	}
}

func (s *UserService) CreateUser(username, email string, isAdmin bool) error {
	user := models.User{
		Username: username,
		Email:    email,
		Password: "",
		Admin:    isAdmin,
	}
	user.GenerateRandomPassword()

	const tokenValidityHours = 24

	validity := tokenValidityHours * time.Hour
	user.GenerateChangePasswordToken(&validity)

	// Aufruf über die Instanz-Variable "s.repo" (nicht mehr global!)
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Aufruf über die Instanz-Variable "s.mailer"
	err = s.mailer.SendNewUserTokenMail(
		createdUser.Email,
		createdUser.ID,
		createdUser.Username,
		*createdUser.ChangePasswordToken,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
