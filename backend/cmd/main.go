package main

import (
	"embed"
	"flag"
	"log/slog"
	"os"
	"strconv"
	"time"

	config "github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/exitcode"
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

const defaultPort = 3000

func main() {
	logLevel := flag.String("log-level", "info", "Set the log level (debug, info, warn, error)")
	port := flag.Int("port", defaultPort, "Set the port number for the server to listen on")

	flag.Parse()

	logger := initializer.InitJsonLogger(*logLevel)

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

	router, err := initializer.InitializeHttpServer(
		*httpHandler,
		websocketHandler,
		*sqliteRepository,
		staticFiles,
		jwtMiddleware,
		cfg,
		logger,
	)
	if err != nil {
		logger.Error("Failed to initialize HTTP server", "error", err)
		os.Exit(int(exitcode.Software))
	}

	startPollerForPendingPurchases(poller, sqliteRepository)
	startCleanupForWebsocketConnections()

	portStr := ":" + strconv.Itoa(*port)
	logger.Info("HTTP server listening", slog.Int("port", *port))

	err = router.Run(portStr)
	if err != nil {
		logger.Error("Failed to start server", "error", err.Error())
		os.Exit(int(exitcode.Software))
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
		slog.Error("Failed to get active purchases", "error", err)

		return
	}

	for _, tx := range activeTransactions {
		slog.Debug("Starting poller for active transaction", "transaction_id", tx.ID)
		poller.Start(tx.ID)
	}
}
