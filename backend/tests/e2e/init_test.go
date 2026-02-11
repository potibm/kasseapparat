package tests_e2e

import (
	"embed"
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
	db = utils.ConnectToLocalDatabase()
	utils.PurgeDatabase(db)
	utils.MigrateDatabase(db)
	utils.SeedDatabase(db, true)
}

func setupTestEnvironment(t *testing.T) (*httptest.Server, func()) {
	cfg := config.Config{
		AppConfig: config.AppConfig{
			Version: "0.1.2",
			GinMode: "test",
		},
		FormatConfig: config.FormatConfig{
			FractionDigitsMax: 2,
			CurrencyCode:      "DKK",
			CurrencyLocale:    "dk-DK",
			DateLocale:        "dk-DK",
			DateOptions:       config.DefaultDateOptions,
			FractionDigitsMin: 0,
		},
		JwtConfig: config.JwtConfig{
			Realm:  "",
			Secret: "test",
		},
		VATRates:           config.DefaultVatRates,
		EnvironmentMessage: "Test environment",
		CorsAllowOrigins:   []string{"http://localhost:3000"},
		PaymentMethods: config.PaymentMethods{
			{Code: models.PaymentMethodCash, Name: "Cash"},
			{Code: models.PaymentMethodCC, Name: "Creditcard"},
			{Code: models.PaymentMethodSumUp, Name: "SumUp"},
		},
	}

	sqliteRepo := sqliteRepo.NewRepository(db, int32(cfg.FormatConfig.FractionDigitsMax))
	sumupRepo := NewMockSumUpRepository()
	mailer := mailer.NewMailer("smtp://127.0.0.1:1025")
	mailer.SetDisabled(true)

	jwtMiddleware := initializer.InitializeJwtMiddleware(sqliteRepo, cfg.JwtConfig)

	purchaseService := purchaseService.NewPurchaseService(sqliteRepo, sumupRepo, mailer, int32(cfg.FormatConfig.FractionDigitsMax))

	statusPublisher := MockStatusPublisher{}
	poller := monitor.NewPoller(sumupRepo, sqliteRepo, purchaseService, &statusPublisher)

	httpHandlerConfig := handlerHttp.HandlerConfig{
		Repo:            sqliteRepo,
		SumupRepository: sumupRepo,
		PurchaseService: purchaseService,
		Monitor:         poller,
		Mailer:          *mailer,
		AppConfig:       cfg,
	}
	handlerHttp := handlerHttp.NewHandler(httpHandlerConfig)
	websocketHandler := websocket.NewHandler(sqliteRepo, sumupRepo, purchaseService, jwtMiddleware, &cfg.CorsAllowOrigins)

	router := initializer.InitializeHttpServer(*handlerHttp, websocketHandler, *sqliteRepo, embed.FS{}, jwtMiddleware, cfg)

	ts := httptest.NewServer(router)

	e = httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(router),
			Jar:       httpexpect.NewCookieJar(),
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

func testAuthenticationForEntityEndpoints(t *testing.T, baseUrl string, urlWithId string) {
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
