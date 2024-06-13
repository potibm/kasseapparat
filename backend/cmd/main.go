package main

import (
	"embed"
	"log"
	"os"

	"github.com/potibm/kasseapparat/internal/app/handler"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

//go:embed assets
var staticFiles embed.FS

func main() {
	initializer.InitializeDotenv();
	initializer.InitializeSentry()

	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	repository := repository.NewRepository()
	myhandler := handler.NewHandler(repository)

	router := initializer.InitializeHttpServer(*myhandler, *repository, staticFiles)
	log.Println("Listening on " + port + "...")
	err := router.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
