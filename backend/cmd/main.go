package main

import (
	"embed"
	"log"
	"os"

	handlerHttp "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
)

//go:embed assets
var staticFiles embed.FS

func main() {
	initializer.InitializeDotenv()
	initializer.InitializeVersion()
	initializer.InitializeSentry()
	initializer.OutputVersion()
	initializer.InitializeSumup()

	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	repository := sqliteRepo.NewRepository(initializer.GetCurrencyDecimalPlaces())
	sumupRepository := sumup.NewRepository(initializer.GetSumupService())

	mailer := initializer.InitializeMailer()
	myhandler := handlerHttp.NewHandler(repository, sumupRepository, mailer, initializer.GetVersion(), initializer.GetCurrencyDecimalPlaces(), initializer.GetEnabledPaymentMethods())

	router := initializer.InitializeHttpServer(*myhandler, *repository, staticFiles)

	log.Println("Listening on " + port + "...")

	err := router.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
