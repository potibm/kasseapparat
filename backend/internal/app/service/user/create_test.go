package user

import (
	"errors"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateUser(user models.User) (models.User, error) {
	args := m.Called(user)
	u := args.Get(0).(models.User)

	return u, args.Error(1)
}

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendNewUserTokenMail(email string, userID int, username, token string) error {
	args := m.Called(email, userID, username, token)

	return args.Error(0)
}

func TestCreateUserSuccess(t *testing.T) {
	repo := new(MockRepo)
	mailer := new(MockMailer)
	svc := NewUserService(repo, mailer)

	testEmail := "test@example.com"
	testUser := "testuser"

	repo.On("CreateUser", mock.MatchedBy(func(u models.User) bool {
		return u.Username == testUser && u.Email == testEmail
	})).Return(models.User{
		ID:                  1337,
		Username:            testUser,
		Email:               testEmail,
		ChangePasswordToken: stringPtr("valid-token"), // Helper siehe unten
	}, nil)

	mailer.On("SendNewUserTokenMail", testEmail, int(1337), testUser, "valid-token").Return(nil)

	err := svc.CreateUser(testUser, testEmail, true)

	// Assertions
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	mailer.AssertExpectations(t)
}

func TestCreateUserRepoError(t *testing.T) {
	repo := new(MockRepo)
	mailer := new(MockMailer)
	svc := NewUserService(repo, mailer)

	repo.On("CreateUser", mock.Anything).Return(models.User{}, errors.New("db crash"))

	err := svc.CreateUser("foo", "bar@baz.com", false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	mailer.AssertNotCalled(t, "SendNewUserTokenMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateUserMailerError(t *testing.T) {
	repo := new(MockRepo)
	mailer := new(MockMailer)
	svc := NewUserService(repo, mailer)

	repo.On("CreateUser", mock.Anything).Return(models.User{
		ID:                  1,
		Email:               "a@b.com",
		ChangePasswordToken: stringPtr("token"),
	}, nil)

	mailer.On("SendNewUserTokenMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("smtp down"))

	err := svc.CreateUser("foo", "a@b.com", false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

func stringPtr(s string) *string {
	return &s
}
