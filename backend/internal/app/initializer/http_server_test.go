package initializer

import (
	"context"
	"embed"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	httpHandler "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"github.com/stretchr/testify/assert"
)

// Important: create a folder named "assets" in the same directory as this test file and add
// an index.html file to it, otherwise the test will fail because the static files cannot be
// loaded. The content of index.html can be anything, e.g. "<p>Test</p>".

//go:embed assets/*
var testFS embed.FS

// --- STUB FOR THE WEBSOCKET INTERFACE ---.
type stubWSHandler struct{}

// HandleTransactionWebSocket implements the  TransactionWebSocketHandler interface.
func (s *stubWSHandler) HandleTransactionWebSocket(c *gin.Context) { /* mocked implementation */ }

// --- TEST ---.
func TestInitializeHttpServer(t *testing.T) {
	// Switch to test mode to avoid side effects on global Gin state
	gin.SetMode(gin.TestMode)

	emptyHTTPHandler := httpHandler.Handler{}
	emptyRepo := sqliteRepo.Repository{}

	var mockWs websocket.TransactionWebSocketHandler = &stubWSHandler{}

	cfg := config.Config{
		App: config.AppConfig{
			GinMode:          gin.TestMode,
			CorsAllowOrigins: []string{"http://localhost:8080"},
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("should initialize server and register routes successfully", func(t *testing.T) {
		engine, err := InitializeHTTPServer(
			emptyHTTPHandler,
			mockWs,
			emptyRepo,
			testFS,
			cfg,
			logger,
		)

		assert.NoError(t, err)
		assert.NotNil(t, engine)

		routes := engine.Routes()

		var foundConfigRoute bool

		for _, r := range routes {
			if r.Path == "/api/v3/config" && r.Method == http.MethodGet {
				foundConfigRoute = true

				break
			}
		}

		assert.True(t, foundConfigRoute, "The route /api/v3/config should be registered")
	})
}

func TestInitializeOIDCHandler_ProxyMode(t *testing.T) {
	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode: "proxy",
		},
	}

	handler, err := InitializeOIDCHandler(context.Background(), cfg)
	assert.NoError(t, err)
	assert.Nil(t, handler)
}

func TestInitializeOIDCHandler_OIDCMode_MissingConfig(t *testing.T) {
	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:             "oidc",
			OidcIssuer:       "",
			OidcClientID:     "",
			OidcClientSecret: "",
			OidcCallbackURL:  "",
			SessionSecret:    "test-secret-that-is-long-enough-for-testing",
		},
		App: config.AppConfig{
			FrontendURL: "http://localhost:3000",
		},
	}

	handler, err := InitializeOIDCHandler(context.Background(), cfg)
	assert.Error(t, err)
	assert.Nil(t, handler)
}

func TestCreateCorsMiddleware(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000", "http://localhost:8080"}

	corsMiddleware := CreateCorsMiddleware(allowedOrigins)
	assert.NotNil(t, corsMiddleware)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(corsMiddleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSlogUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	slogMiddleware := SlogUserID()
	assert.NotNil(t, slogMiddleware)

	router := gin.New()
	router.Use(slogMiddleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("X-Test", "value")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSlogUserID_WithUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	slogMiddleware := SlogUserID()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.IdentityKey, "testuser")
		c.Next()
	})
	router.Use(slogMiddleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSentryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sentryMiddleware := SentryMiddleware()
	assert.NotNil(t, sentryMiddleware)

	router := gin.New()
	router.Use(sentryMiddleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSentryMiddleware_WithUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sentryMiddleware := SentryMiddleware()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.IdentityKey, "testuser")
		c.Next()
	})
	router.Use(sentryMiddleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
