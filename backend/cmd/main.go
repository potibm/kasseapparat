package main

import (
	"embed"
	"log"
	"os"

	handlerHttp "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/monitor"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	"github.com/potibm/kasseapparat/internal/app/utils"
)

//go:embed assets
var staticFiles embed.FS

func main() {
	db := utils.ConnectToDatabase()

	initializer.InitializeDotenv()
	initializer.InitializeVersion()
	initializer.InitializeSentry()
	initializer.OutputVersion()
	initializer.InitializeSumup()

	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	sqliteRepository := sqliteRepo.NewRepository(db, initializer.GetCurrencyDecimalPlaces())
	sumupRepository := sumupRepo.NewRepository(initializer.GetSumupService())
	mailer := initializer.InitializeMailer()

	purchaseService := purchaseService.NewPurchaseService(sqliteRepository, sumupRepository, &mailer, initializer.GetCurrencyDecimalPlaces())

	poller := monitor.NewPoller(sumupRepository, sqliteRepository, purchaseService)

	httpHandler := handlerHttp.NewHandler(sqliteRepository, sumupRepository, purchaseService, poller, mailer, initializer.GetVersion(), initializer.GetCurrencyDecimalPlaces(), initializer.GetEnabledPaymentMethods())
	websocketHandler := websocket.NewHandler(sqliteRepository, sumupRepository, purchaseService)

	router := initializer.InitializeHttpServer(*httpHandler, websocketHandler, *sqliteRepository, staticFiles)

	// @TODO we should restart the poller for active transactions

	log.Println("Listening on " + port + "...")

	err := router.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
