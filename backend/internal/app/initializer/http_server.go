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

	r.Use(createCorsMiddleware())
	r.Use(sentrygin.New(sentrygin.Options{}))

	r.Use(static.Serve("/", static.EmbedFolder(staticFiles, "assets")))

	authMiddleware := registerAuthMiddleware(repository)
	//r.Use(SentryMiddleware())
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
	corsConfig.AddAllowHeaders("Authorization","Credentials")
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
						ID:    strconv.Itoa(int(user.ID)),
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
		protectedApiRouter.GET("/products", myhandler.GetProducts)
		protectedApiRouter.GET("/products/:id", myhandler.GetProductByID)
		protectedApiRouter.GET("/products/:id/listEntries", myhandler.GetListEntriesByProductID)
		protectedApiRouter.PUT("/products/:id", myhandler.UpdateProductByID)
		protectedApiRouter.DELETE("/products/:id", myhandler.DeleteProductByID)
		protectedApiRouter.POST("/products", myhandler.CreateProduct)

		protectedApiRouter.GET("/productStats", myhandler.GetProductStats)
		
		protectedApiRouter.GET("/lists", myhandler.GetLists)
		protectedApiRouter.GET("/lists/:id", myhandler.GetListByID)
		protectedApiRouter.PUT("/lists/:id", myhandler.UpdateListByID)
		protectedApiRouter.DELETE("/lists/:id", myhandler.DeleteListByID)
		protectedApiRouter.POST("/lists", myhandler.CreateList)

		protectedApiRouter.GET("/listEntries", myhandler.GetListEntries)
		protectedApiRouter.GET("/listEntries/:id", myhandler.GetListEntryByID)
		protectedApiRouter.PUT("/listEntries/:id", myhandler.UpdateListEntryByID)
		protectedApiRouter.DELETE("/listEntries/:id", myhandler.DeleteListEntryByID)
		protectedApiRouter.POST("/listEntries", myhandler.CreateListEntry)

		protectedApiRouter.OPTIONS("/purchases", myhandler.OptionsPurchases)
		protectedApiRouter.GET("/purchases", myhandler.GetPurchases)
		protectedApiRouter.GET("/purchases/:id", myhandler.GetPurchaseByID)
		protectedApiRouter.POST("/purchases", myhandler.PostPurchases)
		protectedApiRouter.DELETE("/purchases/:id", myhandler.DeletePurchase)

		protectedApiRouter.GET("/users", myhandler.GetUsers)
		protectedApiRouter.GET("/users/:id", myhandler.GetUserByID)
		protectedApiRouter.PUT("/users/:id", myhandler.UpdateUserByID)
		protectedApiRouter.DELETE("/users/:id", myhandler.DeleteUserByID)
		protectedApiRouter.POST("/users", myhandler.CreateUser)
		
		protectedApiRouter.POST("/auth/changePassword", myhandler.UpdateUserPassword)
		protectedApiRouter.POST("/auth/changePasswordToken", myhandler.RequestChangePasswordToken)
	}

	// unprotected routes
	unprotectedApiRouter := r.Group("/api/v1")
	{
		unprotectedApiRouter.GET("/config", myhandler.GetConfig)
		unprotectedApiRouter.GET("/purchases/stats", myhandler.GetPurchaseStats)
	}
}
