package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/potibm/kasseapparat/internal/app/handler"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

//go:embed assets
var staticFiles embed.FS

func main() {

	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	repository := repository.NewRepository()
	myhandler := handler.NewHandler(repository)

	r := gin.Default()

	authMiddleware, err := jwt.New(middleware.InitParams(*repository))
	r.Use(middleware.HandlerMiddleWare(authMiddleware))

	// register route
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddExposeHeaders("X-Total-Count")
	r.Use(cors.New(corsConfig))

	r.Use(static.Serve("/", static.EmbedFolder(staticFiles, "assets")))

	apiRouter := r.Group("/api/v1")
	{
		apiRouter.GET("/products", myhandler.GetProducts)
		apiRouter.GET("/products/:id", myhandler.GetProductByID)
		apiRouter.PUT("/products/:id", authMiddleware.MiddlewareFunc(), myhandler.UpdateProductByID)
		apiRouter.DELETE("/products/:id", authMiddleware.MiddlewareFunc(), myhandler.DeleteProductByID)
		apiRouter.POST("/products", authMiddleware.MiddlewareFunc(), myhandler.CreateProduct)

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

		apiRouter.GET("/purchases/stats", myhandler.GetPurchaseStats)
	}

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

	middleware.RegisterRoute(r, authMiddleware)

	log.Println("Listening on " + port + "...")
	err = r.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
