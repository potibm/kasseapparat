package initializer

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/handler"
	"github.com/potibm/kasseapparat/internal/app/middleware"
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

func registerApiRoutes(myhandler handler.Handler, authMiddleware *jwt.GinJWTMiddleware) {

	apiRouter := r.Group("/api/v1")
	{
		apiRouter.GET("/products", myhandler.GetProducts)
		apiRouter.GET("/products/:id", myhandler.GetProductByID)
		apiRouter.GET("/products/:id/listEntries", myhandler.GetListEntriesByProductID)
		apiRouter.PUT("/products/:id", authMiddleware.MiddlewareFunc(), myhandler.UpdateProductByID)
		apiRouter.DELETE("/products/:id", authMiddleware.MiddlewareFunc(), myhandler.DeleteProductByID)
		apiRouter.POST("/products", authMiddleware.MiddlewareFunc(), myhandler.CreateProduct)

		apiRouter.GET("/productStats", myhandler.GetProductStats)
		
		apiRouter.GET("/lists", myhandler.GetLists)
		apiRouter.GET("/lists/:id", myhandler.GetListByID)
		apiRouter.PUT("/lists/:id", authMiddleware.MiddlewareFunc(), myhandler.UpdateListByID)
		apiRouter.DELETE("/lists/:id", authMiddleware.MiddlewareFunc(), myhandler.DeleteListByID)
		apiRouter.POST("/lists", authMiddleware.MiddlewareFunc(), myhandler.CreateList)

		apiRouter.GET("/listEntries", myhandler.GetListEntries)
		apiRouter.GET("/listEntries/:id", myhandler.GetListEntryByID)
		apiRouter.PUT("/listEntries/:id", authMiddleware.MiddlewareFunc(), myhandler.UpdateListEntryByID)
		apiRouter.DELETE("/listEntries/:id", authMiddleware.MiddlewareFunc(), myhandler.DeleteListEntryByID)
		apiRouter.POST("/listEntries", authMiddleware.MiddlewareFunc(), myhandler.CreateListEntry)

		apiRouter.OPTIONS("/purchases", myhandler.OptionsPurchases)
		apiRouter.GET("/purchases", myhandler.GetPurchases)
		apiRouter.GET("/purchases/:id", myhandler.GetPurchaseByID)
		apiRouter.POST("/purchases", authMiddleware.MiddlewareFunc(), myhandler.PostPurchases)
		apiRouter.DELETE("/purchases/:id", authMiddleware.MiddlewareFunc(), myhandler.DeletePurchase)

		apiRouter.GET("/users", myhandler.GetUsers)
		apiRouter.GET("/users/:id", myhandler.GetUserByID)
		apiRouter.PUT("/users/:id", authMiddleware.MiddlewareFunc(), myhandler.UpdateUserByID)
		apiRouter.DELETE("/users/:id", authMiddleware.MiddlewareFunc(), myhandler.DeleteUserByID)
		apiRouter.POST("/users", authMiddleware.MiddlewareFunc(), myhandler.CreateUser)
		
		apiRouter.POST("/auth/changePassword", myhandler.UpdateUserPassword)
		apiRouter.POST("/auth/changePasswordToken", myhandler.RequestChangePasswordToken)

		apiRouter.GET("/config", myhandler.GetConfig)
		apiRouter.GET("/purchases/stats", myhandler.GetPurchaseStats)
	}
}
