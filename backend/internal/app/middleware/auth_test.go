package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ginjwtCore "github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock for the User Repository.
type MockAuthRepo struct{ mock.Mock }

func (m *MockAuthRepo) GetUserByLoginAndPassword(l, p string) (*models.User, error) {
	args := m.Called(l, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

func TestInitParams(t *testing.T) {
	t.Run("Should return valid GinJWTMiddleware with correct parameters", func(t *testing.T) {
		mockRepo := new(MockAuthRepo)
		realm := "test-realm"
		secret := "super-secret-key-that-is-long-enough"
		timeout := 10
		secureCookie := true

		middleware := InitParams(mockRepo, realm, secret, timeout, secureCookie, nil)

		assert.NotNil(t, middleware)
		assert.Equal(t, realm, middleware.Realm)
		assert.Equal(t, []byte(secret), middleware.Key)
		assert.Equal(t, time.Minute*time.Duration(timeout), middleware.Timeout)
		assert.Equal(t, secureCookie, middleware.SecureCookie)
	})

	t.Run("Should return valid GinJWTMiddleware when secret is empty", func(t *testing.T) {
		mockRepo := new(MockAuthRepo)
		realm := "test-realm"
		timeout := 10
		secureCookie := true

		middleware := InitParams(mockRepo, realm, "", timeout, secureCookie, nil)

		assert.NotNil(t, middleware)
		assert.Equal(t, []byte("secret"), middleware.Key)
	})
}

func TestAuthMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("IdentityHandler: Extraction from Claims", func(t *testing.T) {
		f := identityHandler()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// We simulate what gin-jwt would put in the context
		c.Set("JWT_PAYLOAD", jwt.MapClaims{
			IdentityKey: float64(456),
		})

		identity := f(c)
		user, ok := identity.(*models.User)

		assert.True(t, ok, "Identity should be a *models.User")
		assert.Equal(t, uint(456), user.ID)
	})

	t.Run("Authorizer: Grant and Deny", func(t *testing.T) {
		f := authorizer()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		assert.True(t, f(c, &models.User{}), "Should allow if data is a user")
		assert.False(t, f(c, "not a user"), "Should deny otherwise")
	})

	t.Run("LoginResponse: DTO and Metrics", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", http.NoBody)

		user := &models.User{ID: 789, Username: "test-admin"}
		c.Set(IdentityKey, user)

		token := &ginjwtCore.Token{
			AccessToken: "secret-token",
		}

		loginResponse(c, token)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp loginResponseDTO

		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "secret-token", resp.AccessToken)

		if assert.NotNil(t, resp.Id) {
			assert.Equal(t, uint(789), *resp.Id)
		}
	})

	t.Run("Unauthorized: Correct Event Types", func(t *testing.T) {
		f := unauthorized()

		paths := []struct {
			url      string
			expected string
		}{
			{loginEndpoint, "login"},
			{refreshEndpoint, "refresh"},
			{"/api/v1/something", "request"},
		}

		for _, p := range paths {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, p.url, http.NoBody)

			f(c, http.StatusUnauthorized, "fail")

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		}
	})
}
