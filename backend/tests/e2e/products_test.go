package tests_e2e

import (
	"net/http"
	"testing"
)

//var e *httpexpect.Expect

/*
func TestMain(m *testing.M) {
	// Setup test environment
	setupTestDatabase()

	code := m.Run()

	os.Exit(code)
}

func setupTest(t *testing.T, route string, handlerFunc gin.HandlerFunc) *httpexpect.Expect {
	// Initialize Gin router
	engine := gin.New()

	// Add route to your handler
	engine.GET(route, handlerFunc)

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(engine),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	return e
}

func setupTestDatabase() {
	db := utils.ConnectToLocalDatabase()
	utils.PurgeDatabase(db)
	utils.MigrateDatabase(db)
	utils.SeedDatabase(db)
}

func setupMailer() mailer.Mailer {
	mailer := mailer.NewMailer("smtp://user:password@localhost:1025")

	return *mailer
}
*/

func TestGetProducts(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
    defer cleanup()

	/*
	// Setze Umgebungsvariablen
    t.Setenv("CORS_ALLOW_ORIGINS", "http://localhost:3000")

    // Initialisiere Repository, Mailer und Handler
    repo := repository.NewLocalRepository()
    mailer := mailer.NewMailer("smtp://127.0.0.1:1025")
    handler := handler.NewHandler(repo, *mailer, "v1")

    // Initialisiere den Gin-Router
    router := initializer.InitializeHttpServer(*handler, *repo, embed.FS{})

    // Erstelle den Testserver
    ts := httptest.NewServer(router)

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(router),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	defer ts.Close()
	*/
	
	res := e.GET("/api/v1/products").
		WithHeader("Authorization", "Bearer " + getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeader := res.Header("X-Total-Count").AsNumber()
    totalCountHeader.Gt(10)

    // Überprüfen der JSON-Antwort
    obj := res.JSON().Array()
    obj.Length().Ge(1)

	for i := 0; i < len(obj.Iter()); i++ {
		product := obj.Value(i).Object()
		product.Value("id").Number().Gt(0)
		product.Value("name").String().NotEmpty()
		product.Value("price").Number().Ge(0)
		product.Value("wrapAfter").Boolean()
		product.Value("pos").Number().Ge(0)
		product.Value("apiExport").Boolean()
	}
}

func TestGetProduct(t *testing.T) {
	/*
	myhandler := handler.NewHandler(repository.NewLocalRepository(), setupMailer(), "1.0.0")
	gin.SetMode(gin.TestMode)

	e := setupTest(t, "/example/:id", myhandler.GetProductByID)
	*/
	_, cleanup := setupTestEnvironment(t)
    defer cleanup()

	obj := e.GET("/api/v1/products/1").
		WithHeader("Authorization", "Bearer " + getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Value("id").Number().Gt(0)
	obj.Value("name").String().NotEmpty()
	obj.Value("price").Number().Ge(0)
	obj.Value("wrapAfter").Boolean()
	obj.Value("pos").Number().Ge(0)
	obj.Value("apiExport").Boolean()
}
