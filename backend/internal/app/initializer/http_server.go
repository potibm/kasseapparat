//nolint:gocritic // unnecessary blocks are more readable with the grouping of routes
package initializer

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	httpHandler "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var r *gin.Engine

const APIVersion = "v2"

func InitializeHTTPServer(
	httpHdlr httpHandler.Handler,
	websocketHdlr websocket.TransactionWebSocketHandler,
	repository sqliteRepo.Repository,
	staticFiles embed.FS,
	jwtMiddleware *jwt.GinJWTMiddleware,
	cfg config.Config,
	logger *slog.Logger,
) (*gin.Engine, error) {
	gin.SetMode(cfg.App.GinMode)

	r = gin.New()
	r.Use(
		middleware.ErrorHandlingMiddleware(),
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic: false,
		}),
		sloggin.New(logger),
		otelgin.Middleware("kasseapparat-backend"),
	)

	r.GET("/api/"+APIVersion+"/purchases/stats", httpHdlr.GetPurchaseStats)

	r.Use(CreateCorsMiddleware(cfg.App.CorsAllowOrigins))

	folder, err := static.EmbedFolder(staticFiles, "assets")
	if err != nil {
		return nil, fmt.Errorf("create embedded folder: %w", err)
	}

	r.Use(static.Serve("/", folder))

	registerAuthMiddleware(jwtMiddleware)
	registerAPIRoutes(httpHdlr, websocketHdlr, jwtMiddleware)

	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") && !strings.Contains(c.Request.RequestURI, ".") {
			file, _ := staticFiles.ReadFile("assets/index.html")
			c.Data(
				http.StatusOK,
				"text/html; charset=utf-8",
				file,
			)
		}
	})

	return r, nil
}

func CreateCorsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = allowedOrigins
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization", "Credentials")
	corsConfig.AddExposeHeaders("X-Total-Count", "Content-Disposition")

	return cors.New(corsConfig)
}

func SlogUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(middleware.IdentityKey)
		if exists {
			if user, ok := user.(*models.User); ok {
				sloggin.AddCustomAttributes(c,
					slog.Int("user_id", user.ID),
				)
			}
		}

		c.Next()
	}
}

func registerAuthMiddleware(authMiddleware *jwt.GinJWTMiddleware) {
	r.Use(middleware.HandlerMiddleWare(authMiddleware))

	versionedGroup := r.Group("/api/" + APIVersion)

	middleware.RegisterRoute(versionedGroup, authMiddleware)
}

func SentryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(middleware.IdentityKey)
		if exists {
			if user, ok := user.(*models.User); ok {
				sentry.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetUser(sentry.User{
						ID: strconv.Itoa(int(user.ID)),
					})
				})
			}
		}

		c.Next()
	}
}

func registerAPIRoutes(
	httpHdlr httpHandler.Handler,
	websocketHdlr websocket.TransactionWebSocketHandler,
	authMiddleware *jwt.GinJWTMiddleware,
) {
	protectedAPIRouter := r.Group("/api/" + APIVersion)
	protectedAPIRouter.Use(authMiddleware.MiddlewareFunc(), SentryMiddleware(), SlogUserID())
	{
		registerProductRoutes(protectedAPIRouter, httpHdlr)
		registerProductInterestRoutes(protectedAPIRouter, httpHdlr)
		protectedAPIRouter.GET("/productStats", httpHdlr.GetProductStats)

		registerGuestlistRoutes(protectedAPIRouter, httpHdlr)
		registerGuestRoutes(protectedAPIRouter, httpHdlr)
		protectedAPIRouter.POST("/guestsUpload", httpHdlr.ImportGuestsFromDeineTicketsCsv)

		registerPurchaseRoutes(protectedAPIRouter, httpHdlr)
		registerUserRoutes(protectedAPIRouter, httpHdlr)

		registerSumupReadersRoutes(protectedAPIRouter, httpHdlr)
		registerSumupTransactionRoutes(protectedAPIRouter, httpHdlr)
	}

	// unprotected routes
	unprotectedAPIRouter := r.Group("/api/" + APIVersion)
	{
		unprotectedAPIRouter.GET("/config", httpHdlr.GetConfig)

		unprotectedAPIRouter.POST("/auth/changePasswordToken", httpHdlr.RequestChangePasswordToken)
		unprotectedAPIRouter.POST("/auth/changePassword", httpHdlr.UpdateUserPassword)

		unprotectedAPIRouter.POST("/sumup/webhook", httpHdlr.GetSumupTransactionWebhook)

		unprotectedAPIRouter.GET("/purchases/:id/ws", websocketHdlr.HandleTransactionWebSocket)
	}
}

func registerProductRoutes(rg *gin.RouterGroup, handler httpHandler.Handler) {
	products := rg.Group("/products")
	{
		products.GET("", handler.GetProducts)
		products.GET("/:id", handler.GetProductByID)
		products.GET("/:id/guests", handler.GetGuestsByProductID)
		products.PUT("/:id", handler.UpdateProductByID)
		products.DELETE("/:id", handler.DeleteProductByID)
		products.POST("", handler.CreateProduct)
	}
}

func registerGuestlistRoutes(rg *gin.RouterGroup, handler httpHandler.Handler) {
	guestlist := rg.Group("/guestlists")
	{
		guestlist.GET("", handler.GetGuestlists)
		guestlist.GET("/:id", handler.GetGuestlistByID)
		guestlist.PUT("/:id", handler.UpdateGuestlistByID)
		guestlist.DELETE("/:id", handler.DeleteGuestlistByID)
		guestlist.POST("", handler.CreateGuestlist)
	}
}

func registerGuestRoutes(rg *gin.RouterGroup, handler httpHandler.Handler) {
	guests := rg.Group("/guests")
	{
		guests.GET("", handler.GetGuests)
		guests.GET("/:id", handler.GetGuestByID)
		guests.PUT("/:id", handler.UpdateGuestByID)
		guests.DELETE("/:id", handler.DeleteGuestByID)
		guests.POST("", handler.CreateGuest)
	}
}

func registerPurchaseRoutes(
	rg *gin.RouterGroup,
	handler httpHandler.Handler,
) {
	purchases := rg.Group("/purchases")
	{
		purchases.GET("", handler.GetPurchases)
		purchases.GET("/:id", handler.GetPurchaseByID)
		purchases.POST("", handler.PostPurchases)
		purchases.DELETE("/:id", handler.DeletePurchase)
		purchases.GET("/export", handler.ExportPurchases)
		purchases.POST("/:id/refund", handler.RefundPurchase)
	}
}

func registerUserRoutes(rg *gin.RouterGroup, handler httpHandler.Handler) {
	users := rg.Group("/users")
	{
		users.GET("", handler.GetUsers)
		users.GET("/:id", handler.GetUserByID)
		users.PUT("/:id", handler.UpdateUserByID)
		users.DELETE("/:id", handler.DeleteUserByID)
		users.POST("", handler.CreateUser)
	}
}

func registerProductInterestRoutes(rg *gin.RouterGroup, handler httpHandler.Handler) {
	productInterests := rg.Group("/productInterests")
	{
		productInterests.GET("", handler.GetProductInterests)
		productInterests.DELETE("/:id", handler.DeleteProductInterestByID)
		productInterests.POST("", handler.CreateProductInterest)
	}
}

func registerSumupReadersRoutes(rg *gin.RouterGroup, handler httpHandler.Handler) {
	sumupReaders := rg.Group("/sumup/readers")
	{
		sumupReaders.GET("", handler.GetSumupReaders)
		sumupReaders.GET("/:id", handler.GetSumupReaderByID)
		sumupReaders.DELETE("/:id", handler.DeleteSumupReader)
		sumupReaders.POST("", handler.CreateSumupReader)
	}
}

func registerSumupTransactionRoutes(rg *gin.RouterGroup, handler httpHandler.Handler) {
	sumupTransactions := rg.Group("/sumup/transactions")
	{
		sumupTransactions.GET("", handler.GetSumupTransactions)
		sumupTransactions.GET("/:id", handler.GetSumupTransactionByID)
	}
}
