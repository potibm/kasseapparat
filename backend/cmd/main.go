package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/potibm/die-kassa/internal/app/repository"
)

func main() {

	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	r := gin.Default()

	apiRouter := r.Group("/api/v1")
	{
		apiRouter.GET("/products", repository.GetProducts)
		apiRouter.GET("/products/:id", repository.GetProductByID)
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
