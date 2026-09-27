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
	"github.com/stretchr/testify/require"
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

		registered := make(map[string]bool, len(routes))

		for _, r := range routes {
			registered[r.Method+" "+r.Path] = true
		}

		assert.True(t, registered[http.MethodGet+" /api/v3/config"], "The route /api/v3/config should be registered")
		assert.True(t, registered[http.MethodGet+" /health"], "The route /health should be registered")
		assert.True(t, registered[http.MethodGet+" /ready"], "The route /ready should be registered")
	})
}

// The statistics display is reached from arbitrary origins, so it needs both a bearer
// token and a preflight answer. The preflight part depends on the route being
// registered before the app-wide CORS middleware, which is easy to break by moving
// the registration, so pin the whole behaviour here rather than per middleware.
func TestPublicPurchaseStats_RequiresTokenAndAnswersPreflight(t *testing.T) {
	const token = "s3cr3t-public-token"

	gin.SetMode(gin.TestMode)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	engine, err := InitializeHTTPServer(
		httpHandler.Handler{},
		&stubWSHandler{},
		sqliteRepo.Repository{},
		testFS,
		config.Config{
			App: config.AppConfig{
				GinMode:          gin.TestMode,
				CorsAllowOrigins: []string{"http://localhost:8080"},
			},
			Auth: config.AuthConfig{PublicEndpointToken: token},
		},
		logger,
	)
	require.NoError(t, err)

	const statsPath = "/api/v3/purchases/stats"

	t.Run("rejects a request without a token", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, statsPath, http.NoBody))

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("rejects a request with a wrong token", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, statsPath, http.NoBody)
		request.Header.Set("Authorization", "Bearer nope")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("accepts a request with the configured token", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, statsPath, http.NoBody)
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Origin", "https://some-other-site.example")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, request)

		// The handler is a zero-value Handler here, so a panic-free 500 is proof the
		// token check let the request through to it.
		assert.NotEqual(t, http.StatusUnauthorized, recorder.Code)
		assert.NotEqual(t, http.StatusForbidden, recorder.Code)
		assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("answers a cross-origin preflight from an origin outside the app allow list", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodOptions, statsPath, http.NoBody)
		request.Header.Set("Origin", "https://some-other-site.example")
		request.Header.Set("Access-Control-Request-Method", "GET")
		request.Header.Set("Access-Control-Request-Headers", "Authorization")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
		assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, recorder.Header().Get("Access-Control-Allow-Headers"), "Authorization")
		assert.Contains(t, recorder.Header().Get("Access-Control-Allow-Methods"), "GET")
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

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := httpHandler.Handler{}

	tests := []struct {
		name          string
		registerFunc  func(rg *gin.RouterGroup, handler httpHandler.Handler)
		expectedPaths []string
	}{
		{
			name: "register product routes",
			registerFunc: func(rg *gin.RouterGroup, handler httpHandler.Handler) {
				registerProductRoutes(rg, handler)
			},
			expectedPaths: []string{
				"/products",
				"/products/:id",
				"/products/:id/guests",
				"/products/:id",
				"/products/:id",
				"/products",
			},
		},
		{
			name: "register guestlist routes",
			registerFunc: func(rg *gin.RouterGroup, handler httpHandler.Handler) {
				registerGuestlistRoutes(rg, handler)
			},
			expectedPaths: []string{
				"/guestlists",
				"/guestlists/:id",
				"/guestlists/:id",
				"/guestlists/:id",
				"/guestlists",
			},
		},
		{
			name: "register guest routes",
			registerFunc: func(rg *gin.RouterGroup, handler httpHandler.Handler) {
				registerGuestRoutes(rg, handler)
			},
			expectedPaths: []string{
				"/guests",
				"/guests/:id",
				"/guests/:id",
				"/guests/:id",
				"/guests",
			},
		},
		{
			name: "register purchase routes",
			registerFunc: func(rg *gin.RouterGroup, handler httpHandler.Handler) {
				registerPurchaseRoutes(rg, handler)
			},
			expectedPaths: []string{
				"/purchases",
				"/purchases/:id",
				"/purchases",
				"/purchases/:id",
				"/purchases/export",
				"/purchases/:id/refund",
			},
		},
		{
			name: "register product interest routes",
			registerFunc: func(rg *gin.RouterGroup, handler httpHandler.Handler) {
				registerProductInterestRoutes(rg, handler)
			},
			expectedPaths: []string{
				"/productInterests",
				"/productInterests/:id",
				"/productInterests",
			},
		},
		{
			name: "register sumup readers routes",
			registerFunc: func(rg *gin.RouterGroup, handler httpHandler.Handler) {
				registerSumupReadersRoutes(rg, handler)
			},
			expectedPaths: []string{
				"/sumup/readers",
				"/sumup/readers/:id",
				"/sumup/readers/:id",
				"/sumup/readers",
			},
		},
		{
			name: "register sumup transaction routes",
			registerFunc: func(rg *gin.RouterGroup, handler httpHandler.Handler) {
				registerSumupTransactionRoutes(rg, handler)
			},
			expectedPaths: []string{
				"/sumup/transactions",
				"/sumup/transactions/:id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			group := router.Group("/api/v3")

			tt.registerFunc(group, handler)

			routes := router.Routes()
			assert.NotEmpty(t, routes, "routes should be registered")

			for _, expectedPath := range tt.expectedPaths {
				found := false

				for _, route := range routes {
					if route.Path == "/api/v3"+expectedPath {
						found = true

						break
					}
				}

				assert.True(t, found, "expected path %s to be registered", expectedPath)
			}
		})
	}
}
