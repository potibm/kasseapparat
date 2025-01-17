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
	"github.com/potibm/kasseapparat/internal/app/handler"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

var r *gin.Engine
var authMiddleware *jwt.GinJWTMiddleware

func InitializeHttpServer(myhandler handler.Handler, repository repository.Repository, staticFiles embed.FS) *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	r = gin.Default()
	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(middleware.ErrorHandlingMiddleware())

	r.GET("/api/v1/purchases/stats", myhandler.GetPurchaseStats)

	r.Use(createCorsMiddleware())

	r.Use(static.Serve("/", static.EmbedFolder(staticFiles, "assets")))

	authMiddleware := registerAuthMiddleware(repository)
	registerApiRoutes(myhandler, authMiddleware)

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

func createCorsMiddleware() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsAllowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if corsAllowOrigins == "" {
		log.Fatalf("CORS_ALLOW_ORIGINS is not set in env")
	}
	corsConfig.AllowOrigins = strings.Split(corsAllowOrigins, ",")
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization", "Credentials")
	corsConfig.AddExposeHeaders("X-Total-Count")

	return cors.New(corsConfig)
}

func registerAuthMiddleware(repository repository.Repository) *jwt.GinJWTMiddleware {
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

func registerApiRoutes(myhandler handler.Handler, authMiddleware *jwt.GinJWTMiddleware) {

	protectedApiRouter := r.Group("/api/v1")
	protectedApiRouter.Use(authMiddleware.MiddlewareFunc(), SentryMiddleware())
	{
		registerProductRoutes(protectedApiRouter, myhandler)
		registerProductInterestRoutes(protectedApiRouter, myhandler)
		protectedApiRouter.GET("/productStats", myhandler.GetProductStats)

		registerListRoutes(protectedApiRouter, myhandler)
		registerListEntryRoutes(protectedApiRouter, myhandler)
		protectedApiRouter.POST("/listEntriesUpload", myhandler.ImportListEntriesFromDeineTicketsCsv)

		registerPurchaseRoutes(protectedApiRouter, myhandler)
		registerUserRoutes(protectedApiRouter, myhandler)
	}

	// unprotected routes
	unprotectedApiRouter := r.Group("/api/v1")
	{
		unprotectedApiRouter.GET("/config", myhandler.GetConfig)

		unprotectedApiRouter.POST("/auth/changePasswordToken", myhandler.RequestChangePasswordToken)
		unprotectedApiRouter.POST("/auth/changePassword", myhandler.UpdateUserPassword)
	}
}

func registerProductRoutes(rg *gin.RouterGroup, handler handler.Handler) {
	products := rg.Group("/products")
	{
		products.GET("", handler.GetProducts)
		products.GET("/:id", handler.GetProductByID)
		products.GET("/:id/listEntries", handler.GetListEntriesByProductID)
		products.PUT("/:id", handler.UpdateProductByID)
		products.DELETE("/:id", handler.DeleteProductByID)
		products.POST("", handler.CreateProduct)
	}
}

func registerListRoutes(rg *gin.RouterGroup, handler handler.Handler) {
	lists := rg.Group("/guestlists")
	{
		lists.GET("", handler.GetGuestlists)
		lists.GET("/:id", handler.GetGuestlistByID)
		lists.PUT("/:id", handler.UpdateGuestlistByID)
		lists.DELETE(":id", handler.DeleteGuestlistByID)
		lists.POST("", handler.CreateGuestlist)
	}
}

func registerListEntryRoutes(rg *gin.RouterGroup, handler handler.Handler) {
	listEntries := rg.Group("/listEntries")
	{
		listEntries.GET("", handler.GetListEntries)
		listEntries.GET(":id", handler.GetListEntryByID)
		listEntries.PUT("/:id", handler.UpdateListEntryByID)
		listEntries.DELETE(":id", handler.DeleteListEntryByID)
		listEntries.POST("", handler.CreateListEntry)
	}
}

func registerPurchaseRoutes(rg *gin.RouterGroup, handler handler.Handler) {
	purchases := rg.Group("/purchases")
	{
		purchases.GET("", handler.GetPurchases)
		purchases.GET("/:id", handler.GetPurchaseByID)
		purchases.POST("", handler.PostPurchases)
		purchases.DELETE("/:id", handler.DeletePurchase)
	}
}

func registerUserRoutes(rg *gin.RouterGroup, handler handler.Handler) {
	users := rg.Group("/users")
	{
		users.GET("", handler.GetUsers)
		users.GET("/:id", handler.GetUserByID)
		users.PUT("/:id", handler.UpdateUserByID)
		users.DELETE("/:id", handler.DeleteUserByID)
		users.POST("", handler.CreateUser)
	}
}

func registerProductInterestRoutes(rg *gin.RouterGroup, handler handler.Handler) {
	productInterests := rg.Group("/productInterests")
	{
		productInterests.GET("", handler.GetProductInterests)
		productInterests.DELETE("/:id", handler.DeleteProductInterestByID)
		productInterests.POST("", handler.CreateProductInterest)
	}
}
