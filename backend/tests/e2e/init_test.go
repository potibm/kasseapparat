package tests_e2e

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	handlerHttp "github.com/potibm/kasseapparat/internal/app/handler/http"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"github.com/potibm/kasseapparat/internal/app/utils"
)

var (
	e                *httpexpect.Expect
	demoJwt          string
	adminJwt         string
	totalCountHeader = "X-Total-Count"
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	os.Exit(code)
}

func setup() {
	db := utils.ConnectToLocalDatabase()
	utils.PurgeDatabase(db)
	utils.MigrateDatabase(db)
	utils.SeedDatabase(db, true)
}

func setupTestEnvironment(t *testing.T) (*httptest.Server, func()) {
	t.Setenv("CORS_ALLOW_ORIGINS", "http://localhost:3000")
	t.Setenv("JWT_SECRET", "test")

	currencyDecimalPlaces := int32(2)
	paymentMethods := map[string]string{
		"CASH":  "💶 Cash",
		"CC":    "💳 Creditcard",
		"SUMUP": "💳 SumUp",
	}

	repo := sqliteRepo.NewLocalRepository(currencyDecimalPlaces)
	sumupRepo := NewMockSumUpRepository()
	mailer := mailer.NewMailer("smtp://127.0.0.1:1025")
	mailer.SetDisabled(true)
	handlerHttp := handlerHttp.NewHandler(repo, sumupRepo, *mailer, "v1", currencyDecimalPlaces, paymentMethods)

	router := initializer.InitializeHttpServer(*handlerHttp, *repo, embed.FS{})

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
	login := e.POST("/login").
		WithJSON(map[string]string{
			"login":    username,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Read the JWT token from the response
	jwt := login.Value("token").String().Raw()

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
