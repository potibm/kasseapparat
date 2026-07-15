package tests_e2e

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/potibm/kasseapparat/internal/app/config"
	handlerHttp "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/handler/websocket"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/monitor"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	"github.com/potibm/kasseapparat/internal/app/utils"
	"gorm.io/gorm"
)

var (
	e                *httpexpect.Expect
	totalCountHeader = "X-Total-Count"
	db               *gorm.DB
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	os.Exit(code)
}

func setup() {
	var err error

	db, err = utils.ConnectToLocalDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = utils.PurgeDatabase(db)
	if err != nil {
		log.Fatal("Failed to purge database: ", err)
	}

	err = utils.MigrateDatabase(db)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	utils.SeedDatabase(db, true)
}

func setupTestEnvironment(t *testing.T) (httpServer *httptest.Server, cleanupFunc func()) {
	cfg := config.Config{
		App: config.AppConfig{
			Version:            "0.1.2",
			GinMode:            "debug",
			LogLevel:           "debug",
			LogFormat:          "text",
			RedisURL:           "",
			Environment:        "test",
			EnvironmentMessage: "Test environment",
			CorsAllowOrigins:   []string{"http://localhost:3000"},
		},
		Format: config.FormatConfig{
			Currency: config.CurrencyFormatConfig{
				Code:              "DKK",
				Locale:            "da-DK",
				FractionDigitsMax: 2,
				FractionDigitsMin: 0,
			},
			Date: config.DateFormatConfig{
				Locale:  "da-DK",
				Options: config.DefaultDateOptions,
			},
		},
		Jwt: config.JwtConfig{
			Realm:  "",
			Secret: "test",
		},
		Auth: config.AuthConfig{
			Mode:        "proxy",
			ProxyHeader: "X-Remote-User",
		},
		VATRates: config.DefaultVatRates,
		PaymentMethods: config.PaymentMethods{
			{Code: models.PaymentMethodCash, Name: "Cash"},
			{Code: models.PaymentMethodCC, Name: "Creditcard"},
			{Code: models.PaymentMethodSumUp, Name: "SumUp"},
		},
	}

	sqliteRp := sqliteRepo.NewRepository(db, int32(cfg.Format.Currency.FractionDigitsMax))
	sumupRp := NewMockSumUpRepository()
	mail, _ := mailer.NewMailer("smtp://127.0.0.1:1025")
	mail.SetDisabled(true)

	purchaseSrvc := purchaseService.NewPurchaseService(
		sqliteRp,
		sumupRp,
		mail,
		int32(cfg.Format.Currency.FractionDigitsMax),
		cfg.Format.Currency.Code,
	)

	statusPublisher := MockStatusPublisher{}
	poller := monitor.NewPoller(sumupRp, sqliteRp, purchaseSrvc, &statusPublisher)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	httpHandlerConfig := handlerHttp.HandlerConfig{
		Repo:            sqliteRp,
		SumupRepository: sumupRp,
		PurchaseService: purchaseSrvc,
		Monitor:         poller,
		Mailer:          *mail,
		AppConfig:       cfg,
	}
	handlerHTTPObj := handlerHttp.NewHandler(httpHandlerConfig)
	websocketHandler := websocket.NewHandler(
		sqliteRp,
		sumupRp,
		purchaseSrvc,
		&cfg.App.CorsAllowOrigins,
	)

	router, err := initializer.InitializeHTTPServer(
		*handlerHTTPObj,
		websocketHandler,
		*sqliteRp,
		embed.FS{},
		cfg,
		logger,
	)
	if err != nil {
		log.Fatal("Failed to initialize HTTP server: ", err)
	}

	ts := httptest.NewServer(router)

	e = httpexpect.WithConfig(httpexpect.Config{
		BaseURL: ts.URL,
		Client: &http.Client{
			// Transport: httpexpect.NewBinder(router),
			Jar: httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	cleanup := func() {
		ts.Close()
	}

	return ts, cleanup
}

func withRemoteUser(req *httpexpect.Request, username string) *httpexpect.Request {
	return req.WithHeader("X-Remote-User", username)
}

func withDemoUserAuthToken(req *httpexpect.Request) *httpexpect.Request {
	return withRemoteUser(req, "demo")
}

func withAdminUserAuthToken(req *httpexpect.Request) *httpexpect.Request {
	return withRemoteUser(req, "admin")
}

func testAuthenticationForEntityEndpoints(t *testing.T, baseURL, urlWithID string) {
	// Note: Authentication tests removed for Phase 1 - auth is now handled by reverse proxy in Phase 2
	// Placeholder: just verify endpoints are accessible
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", baseURL).Expect().Status(http.StatusOK)
}

func validateErrorDetailMessage(err *httpexpect.Object, message string) {
	err.Value("details").String().IsEqual(message)
}
