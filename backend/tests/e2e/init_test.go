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
	demoJwt          string
	adminJwt         string
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
				Locale:  "dk-DK",
				Options: config.DefaultDateOptions,
			},
		},
		Jwt: config.JwtConfig{
			Realm:  "",
			Secret: "test",
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

	jwtMiddleware := initializer.InitializeJwtMiddleware(sqliteRp, cfg.Jwt, nil)

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
	handlerHttpObj := handlerHttp.NewHandler(httpHandlerConfig)
	websocketHandler := websocket.NewHandler(
		sqliteRp,
		sumupRp,
		purchaseSrvc,
		jwtMiddleware,
		&cfg.App.CorsAllowOrigins,
	)

	router, err := initializer.InitializeHttpServer(
		*handlerHttpObj,
		websocketHandler,
		*sqliteRp,
		embed.FS{},
		jwtMiddleware,
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
			//Transport: httpexpect.NewBinder(router),
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

func getJwtForUser(username, password string) string {
	// Perform login request
	login := e.POST("/api/v2/auth/login").
		WithJSON(map[string]string{
			"login":    username,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Read the JWT token from the response
	jwt := login.Value("access_token").String().Raw()

	return jwt
}

func getJwtForDemoUser() string {
	if demoJwt == "" {
		demoJwt = getJwtForUser("demo", "demo")
	}

	return demoJwt
}

func getJwtForAdminUser() string {
	if adminJwt == "" {
		adminJwt = getJwtForUser("admin", "admin")
	}

	return adminJwt
}

func withAuthToken(req *httpexpect.Request, token string) *httpexpect.Request {
	return req.WithHeader("Authorization", "Bearer "+token)
}

func withDemoUserAuthToken(req *httpexpect.Request) *httpexpect.Request {
	return withAuthToken(req, getJwtForDemoUser())
}

func withAdminUserAuthToken(req *httpexpect.Request) *httpexpect.Request {
	return withAuthToken(req, getJwtForAdminUser())
}

func testAuthenticationForEntityEndpoints(t *testing.T, baseUrl, urlWithId string) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", baseUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("GET", urlWithId).Expect().Status(http.StatusUnauthorized)
	e.Request("POST", baseUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("PUT", urlWithId).Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", urlWithId).Expect().Status(http.StatusUnauthorized)
}

func validateErrorDetailMessage(err *httpexpect.Object, message string) {
	err.Value("details").String().IsEqual(message)
}
