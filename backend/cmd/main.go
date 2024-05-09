package main

import (
	"log"
	"os"

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

	/*
		r.Use(cors.New(cors.Config{
			AllowAllOrigins:  false,
			AllowMethods:     []string{"POST", "DELETE", "PUT", "GET", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
			//MaxAge: 12 * time.Hour,
		}))*/
	r.Use(cors.Default())

	apiRouter := r.Group("/api/v1")
	{
		apiRouter.GET("/products", myhandler.GetProducts)
		apiRouter.GET("/products/:id", myhandler.GetProductByID)
		apiRouter.OPTIONS("/purchases", myhandler.OptionsPurchases)
		apiRouter.GET("/purchases", myhandler.GetLastPurchases)
		apiRouter.POST("/purchases", myhandler.PostPurchases)
		apiRouter.DELETE("/purchases/:id", myhandler.DeletePurchases)
		apiRouter.GET("/purchases/stats", myhandler.GetPurchaseStats)
	}

	// Serve static files from the "public" directory for all other requests
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.Static("/static", "./public/static")

	log.Println("Listening on " + port + "...")
	err := r.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
