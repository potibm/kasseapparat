package cmd

import (
	"embed"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/spf13/cobra"

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
	port         int
	otelEndpoint string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the HTTP server for the Kasseapparat application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Context
		ctx := cmd.Context()

		// 2. Initialize Telemetry
		shutdownFn, err := initializer.InitTelemetry(ctx, otelEndpoint, Cfg.App.Version)
		if err != nil {
			return fmt.Errorf("failed to initialize telemetry: %w", err)
		}

		if shutdownFn != nil {
			defer shutdownFn()
		}

		// 3. Connect to Database
		db, err := utils.ConnectToDatabase(Cfg.App.DbFilename)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		// 4. Initialize external services (Sentry, SumUp, etc.)
		initializer.InitializeSentry(Cfg.Sentry)
		initializer.InitializeSumup(Cfg.Sumup)

		// 5. Dependency Injection (Repositories & Middleware)
		sqliteRepository := sqliteRepo.NewRepository(db, int32(Cfg.Format.Currency.FractionDigitsMax))
		sumupRepository := sumupRepo.NewRepository(initializer.GetSumupService())
		mailer := initializer.InitializeMailer(Cfg.Mailer)
		jwtMiddleware := initializer.InitializeJwtMiddleware(sqliteRepository, Cfg.Jwt, &Cfg.App.RedisURL)

		// 6. Services & Handler
		purchaseSvc := purchaseService.NewPurchaseService(
			sqliteRepository,
			sumupRepository,
			&mailer,
			int32(Cfg.Format.Currency.FractionDigitsMax),
			Cfg.Format.Currency.Code,
		)

		websocketHandler := websocket.NewHandler(
			sqliteRepository,
			sumupRepository,
			purchaseSvc,
			jwtMiddleware,
			&Cfg.App.CorsAllowOrigins,
		)
		publisher := &websocket.WebsocketPublisher{}
		poller := monitor.NewPoller(sumupRepository, sqliteRepository, purchaseSvc, publisher)

		httpHandlerConfig := handlerHttp.HandlerConfig{
			Repo:            sqliteRepository,
			SumupRepository: sumupRepository,
			PurchaseService: purchaseSvc,
			Monitor:         poller,
			StatusPublisher: publisher,
			Mailer:          mailer,
			AppConfig:       Cfg,
		}
		httpHandler := handlerHttp.NewHandler(httpHandlerConfig)

		// 7. Initialize HTTP Server
		router, err := initializer.InitializeHttpServer(
			*httpHandler,
			websocketHandler,
			*sqliteRepository,
			staticFiles,
			jwtMiddleware,
			Cfg,
			slog.Default(),
		)
		if err != nil {
			return fmt.Errorf("failed to initialize HTTP server: %w", err)
		}

		// 8. Start background tasks
		startPollerForPendingPurchases(poller, sqliteRepository)
		startCleanupForWebsocketConnections()

		// 9. Server hochfahren
		portStr := ":" + strconv.Itoa(port)
		slog.Info("HTTP server listening", slog.Int("port", port))

		if err := router.Run(portStr); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Flags spezifisch für den Server
	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "Set the port number for the server to listen on")
	serveCmd.Flags().
		StringVar(&otelEndpoint, "otel-endpoint", "", "Set the OpenTelemetry endpoint (e.g., localhost:4317)")
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
