package main

import (
	"embed"
	"log"
	"os"
	"time"

	config "github.com/potibm/kasseapparat/internal/app/config"
	handlerHttp "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/monitor"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	"github.com/potibm/kasseapparat/internal/app/utils"
)

//go:embed assets
var staticFiles embed.FS

var (
	version = "0.0.0"
)

func main() {
	cfg := config.Load()
	cfg.SetVersion(version)
	cfg.OutputVersion()

	db := utils.ConnectToDatabase()

	initializer.InitializeSentry(cfg.SentryConfig)
	initializer.InitializeSumup(cfg.SumupConfig)

	sqliteRepository := sqliteRepo.NewRepository(db, int32(cfg.FormatConfig.FractionDigitsMax))
	sumupRepository := sumupRepo.NewRepository(initializer.GetSumupService())
	mailer := initializer.InitializeMailer(cfg.MailerConfig)
	jwtMiddleware := initializer.InitializeJwtMiddleware(sqliteRepository, cfg.JwtConfig)

	purchaseService := purchaseService.NewPurchaseService(
		sqliteRepository,
		sumupRepository,
		&mailer,
		int32(cfg.FormatConfig.FractionDigitsMax),
	)

	websocketHandler := websocket.NewHandler(
		sqliteRepository,
		sumupRepository,
		purchaseService,
		jwtMiddleware,
		&cfg.CorsAllowOrigins,
	)
	publisher := &websocket.WebsocketPublisher{}
	poller := monitor.NewPoller(sumupRepository, sqliteRepository, purchaseService, publisher)

	httpHandlerConfig := handlerHttp.HandlerConfig{
		Repo:            sqliteRepository,
		SumupRepository: sumupRepository,
		PurchaseService: purchaseService,
		Monitor:         poller,
		StatusPublisher: publisher,
		Mailer:          mailer,
		AppConfig:       cfg,
	}
	httpHandler := handlerHttp.NewHandler(httpHandlerConfig)

	router := initializer.InitializeHttpServer(
		*httpHandler,
		websocketHandler,
		*sqliteRepository,
		staticFiles,
		jwtMiddleware,
		cfg,
	)

	startPollerForPendingPurchases(poller, sqliteRepository)
	startCleanupForWebsocketConnections()

	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	log.Println("Listening on " + port + "...")

	err := router.Run(port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}

func startCleanupForWebsocketConnections() {
	const cleanupInterval = 5 * time.Minute

	websocket.StartCleanupRoutine(cleanupInterval)
}

func startPollerForPendingPurchases(poller monitor.Poller, sqliteRepository *sqliteRepo.Repository) {
	hasClientTransactionID := true

	filters := sqliteRepo.PurchaseFilters{
		PaymentMethods:         []models.PaymentMethod{models.PaymentMethodSumUp},
		StatusList:             &models.PurchaseStatusList{models.PurchaseStatusPending},
		HasClientTransactionID: &hasClientTransactionID,
	}

	const plentyOfTransactions = 1000

	activeTransactions, err := sqliteRepository.GetPurchases(plentyOfTransactions, 0, "createdAt", "ASC", filters)
	if err != nil {
		log.Printf("[Error] failed to get active purchases: %v", err)

		return
	}

	for _, tx := range activeTransactions {
		log.Println("Starting poller for active transaction:", tx.ID)
		poller.Start(tx.ID)
	}
}
