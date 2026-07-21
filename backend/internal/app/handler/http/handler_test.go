package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryPaymentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create valid payment methods for testing
	validPaymentMethods := config.PaymentMethods{
		{Code: models.PaymentMethodCash, Name: "Cash"},
		{Code: models.PaymentMethodCC, Name: "Credit Card"},
	}

	// Create a request with matching query parameters
	req, _ := http.NewRequest(http.MethodGet, "/?paymentMethods=CASH,CC,INVALID", http.NoBody)
	w := httptest.NewRecorder()
	engine := gin.New()
	c := gin.CreateTestContextOnly(w, engine)
	c.Request = req

	result := queryPaymentMethods(c, "paymentMethods", validPaymentMethods)

	assert.ElementsMatch(t, []models.PaymentMethod{models.PaymentMethodCash, models.PaymentMethodCC}, result)
}

func TestNewHandler(t *testing.T) {
	cfg := HandlerConfig{
		AppConfig: config.Config{
			Format: config.FormatConfig{
				Currency: config.CurrencyFormatConfig{
					FractionDigitsMax: 2,
				},
			},
		},
	}

	handler := NewHandler(cfg)
	require.NotNil(t, handler)
	assert.Equal(t, int32(2), handler.decimalPlaces)
}

func TestGetOIDCHandler(t *testing.T) {
	oidcHandler := &OIDCAuthHandler{}

	handler := &Handler{
		oidcHandler: oidcHandler,
	}

	result := handler.GetOIDCHandler()
	assert.Equal(t, oidcHandler, result)
}

func TestGetOIDCHandler_Nil(t *testing.T) {
	handler := &Handler{
		oidcHandler: nil,
	}

	result := handler.GetOIDCHandler()
	assert.Nil(t, result)
}

func TestGetMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/auth/me", http.NoBody)

	authUser := models.AuthUser{
		Username: "testuser",
		Role:     "user",
	}
	c.Set(middleware.AuthUserKey, authUser)

	handler.GetMe(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMe_NoAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/auth/me", http.NoBody)

	handler.GetMe(c)

	assert.True(t, len(c.Errors) > 0)
}

func TestGetUsernameFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	c.Set(middleware.IdentityKey, "testuser")

	username, err := handler.getUsernameFromContext(c)
	require.NoError(t, err)
	assert.Equal(t, "testuser", username)
}

func TestGetUsernameFromContext_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	username, err := handler.getUsernameFromContext(c)
	assert.Error(t, err)
	assert.Empty(t, username)
}

func TestGetUsernameFromContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	c.Set(middleware.IdentityKey, 12345)

	username, err := handler.getUsernameFromContext(c)
	assert.Error(t, err)
	assert.Empty(t, username)
}

func TestContextWithUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	c.Set(middleware.IdentityKey, "testuser")

	result := handler.contextWithUser(c)
	assert.NotNil(t, result)
}

func TestContextWithUser_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	result := handler.contextWithUser(c)
	assert.NotNil(t, result)
}
