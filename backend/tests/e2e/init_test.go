package tests_e2e

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/potibm/kasseapparat/internal/app/handler"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/repository"
	"github.com/potibm/kasseapparat/internal/app/utils"
)

var (
	e       *httpexpect.Expect
	demoJwt string
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
	utils.SeedDatabase(db)
}

func setupTestEnvironment(t *testing.T) (*httptest.Server, func()) {
	t.Setenv("CORS_ALLOW_ORIGINS", "http://localhost:3000")
	t.Setenv("JWT_SECRET", "test")

	repo := repository.NewLocalRepository()
	mailer := mailer.NewMailer("smtp://127.0.0.1:1025")
	handler := handler.NewHandler(repo, *mailer, "v1")

	router := initializer.InitializeHttpServer(*handler, *repo, embed.FS{})

	ts := httptest.NewServer(router)

	e = httpexpect.Default(t, ts.URL)

	cleanup := func() {
		ts.Close()
	}

	return ts, cleanup
}

func getJwtForUser(username, password string) string {
	// Login durchführen
	login := e.POST("/login").
		WithJSON(map[string]string{
			"login":    username,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// JWT auslesen
	jwt := login.Value("token").String().Raw()

	return jwt
}

func getJwtForDemoUser() string {
	if demoJwt == "" {
		demoJwt = getJwtForUser("demo", "demo")
	}
	return demoJwt
}
