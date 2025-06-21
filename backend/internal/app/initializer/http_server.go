package initializer

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	httpHandler "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

var (
	r              *gin.Engine
	authMiddleware *jwt.GinJWTMiddleware
)

func InitializeHttpServer(httpHandler httpHandler.Handler, websocketHandler websocket.HandlerInterface, repository sqliteRepo.Repository, staticFiles embed.FS) *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	r = gin.Default()
	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(middleware.ErrorHandlingMiddleware())

	r.GET("/api/v2/purchases/stats", httpHandler.GetPurchaseStats)

	r.Use(CreateCorsMiddleware())

	folder, err := static.EmbedFolder(staticFiles, "assets")
	if err != nil {
		log.Fatalf("Failed to create embedded folder: %v", err)
	}

	r.Use(static.Serve("/", folder))

	authMiddleware := registerAuthMiddleware(repository)
	registerApiRoutes(httpHandler, websocketHandler, authMiddleware)

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

	return r
}

func CreateCorsMiddleware() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	corsAllowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if corsAllowOrigins == "" {
		log.Fatalf("CORS_ALLOW_ORIGINS is not set in env")
	}

	corsConfig.AllowOrigins = strings.Split(corsAllowOrigins, ",")
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization", "Credentials")
	corsConfig.AddExposeHeaders("X-Total-Count", "Content-Disposition")

	return cors.New(corsConfig)
}

func registerAuthMiddleware(repository sqliteRepo.Repository) *jwt.GinJWTMiddleware {
	authMiddleware, _ = jwt.New(middleware.InitParams(repository, os.Getenv("JWT_REALM"), os.Getenv("JWT_SECRET"), 10))
	r.Use(middleware.HandlerMiddleWare(authMiddleware))

	middleware.RegisterRoute(r, authMiddleware)

	return authMiddleware
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

func registerApiRoutes(httpHandler httpHandler.Handler, websockeHandler websocket.HandlerInterface, authMiddleware *jwt.GinJWTMiddleware) {
	protectedApiRouter := r.Group("/api/v2")
	protectedApiRouter.Use(authMiddleware.MiddlewareFunc(), SentryMiddleware())
	{
		registerProductRoutes(protectedApiRouter, httpHandler)
		registerProductInterestRoutes(protectedApiRouter, httpHandler)
		protectedApiRouter.GET("/productStats", httpHandler.GetProductStats)

		registerGuestlistRoutes(protectedApiRouter, httpHandler)
		registerGuestRoutes(protectedApiRouter, httpHandler)
		protectedApiRouter.POST("/guestsUpload", httpHandler.ImportGuestsFromDeineTicketsCsv)

		registerPurchaseRoutes(protectedApiRouter, httpHandler, websockeHandler)
		registerUserRoutes(protectedApiRouter, httpHandler)

		registerSumupReadersRoutes(protectedApiRouter, httpHandler)
		registerSumupTransactionRoutes(protectedApiRouter, httpHandler)
	}

	// unprotected routes
	unprotectedApiRouter := r.Group("/api/v2")
	{
		unprotectedApiRouter.GET("/config", httpHandler.GetConfig)

		unprotectedApiRouter.POST("/auth/changePasswordToken", httpHandler.RequestChangePasswordToken)
		unprotectedApiRouter.POST("/auth/changePassword", httpHandler.UpdateUserPassword)

		unprotectedApiRouter.POST("/sumup/webhook", httpHandler.GetSumupTransactionWebhook)
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

func registerPurchaseRoutes(rg *gin.RouterGroup, handler httpHandler.Handler, websockeHandler websocket.HandlerInterface) {
	purchases := rg.Group("/purchases")
	{
		purchases.GET("", handler.GetPurchases)
		purchases.GET(":id", handler.GetPurchaseByID)
		purchases.POST("", handler.PostPurchases)
		purchases.DELETE(":id", handler.DeletePurchase)
		purchases.GET("export", handler.ExportPurchases)
		purchases.POST(":id/refund", handler.RefundPurchase)
		purchases.GET(":id/ws", websockeHandler.HandleTransactionWebSocket)
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
