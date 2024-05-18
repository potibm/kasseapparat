package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/potibm/kasseapparat/internal/app/handler"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

func main() {

	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	myhandler := handler.NewHandler(repository.NewRepository())

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddExposeHeaders("X-Total-Count")
	r.Use(cors.New(corsConfig))

	apiRouter := r.Group("/api/v1")
	{
		apiRouter.GET("/products", myhandler.GetProducts)
		apiRouter.GET("/products/:id", myhandler.GetProductByID)
		apiRouter.PUT("/products/:id", myhandler.UpdateProductByID)
		apiRouter.DELETE("/products/:id", myhandler.DeleteProductByID)
		apiRouter.POST("/products", myhandler.CreateProduct)

		apiRouter.OPTIONS("/purchases", myhandler.OptionsPurchases)
		apiRouter.GET("/purchases", myhandler.GetPurchases)
		apiRouter.GET("/purchases/:id", myhandler.GetPurchaseByID)
		apiRouter.POST("/purchases", myhandler.PostPurchases)
		apiRouter.DELETE("/purchases/:id", myhandler.DeletePurchases)

		apiRouter.GET("/purchases/stats", myhandler.GetPurchaseStats)
	}

	// Serve static files from the "public" directory for all other requests
	r.StaticFile("/", "./public/index.html")
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") && !strings.Contains(c.Request.RequestURI, ".") {
			c.File("./public/index.html")
		}
		//default 404 page not found
	})
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.Static("/static", "./public/static")

	log.Println("Listening on " + port + "...")
	err := r.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
