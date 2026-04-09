package initializer

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	httpHandler "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
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

	jwtMiddleware := &jwt.GinJWTMiddleware{
		Key: []byte("secret_test_key"),
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("should initialize server and register routes successfully", func(t *testing.T) {
		engine, err := InitializeHTTPServer(
			emptyHTTPHandler,
			mockWs,
			emptyRepo,
			testFS,
			jwtMiddleware,
			cfg,
			logger,
		)

		assert.NoError(t, err)
		assert.NotNil(t, engine)

		routes := engine.Routes()

		var foundConfigRoute bool

		for _, r := range routes {
			if r.Path == "/api/v2/config" && r.Method == http.MethodGet {
				foundConfigRoute = true

				break
			}
		}

		assert.True(t, foundConfigRoute, "The route /api/v2/config should be registered")
	})
}
