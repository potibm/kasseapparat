package repository_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/potibm/die-kassa/internal/app/handler"
	"github.com/potibm/die-kassa/internal/app/repository"
	"github.com/potibm/die-kassa/internal/app/utils"
)

var e *httpexpect.Expect

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

func TestGetProducts(t *testing.T) {
	myhandler := handler.NewHandler(repository.NewLocalRepository())

	gin.SetMode(gin.TestMode)
	e := setupTest(t, "/example", myhandler.GetProducts)

	obj := e.GET("/example").
		Expect().
		Status(http.StatusOK).JSON().Array()

	obj.Length().Ge(1)

	for i := 0; i < len(obj.Iter()); i++ {
		product := obj.Value(i).Object()
		product.Value("ID").Number().Gt(0)
		product.Value("Name").String().NotEmpty()
		product.Value("Price").Number().Ge(0)
		product.Value("WrapAfter").Boolean()
		product.Value("Pos").Number().Ge(0)
		product.Value("ApiExport").Boolean()
	}
}

func TestGetProduct(t *testing.T) {
	myhandler := handler.NewHandler(repository.NewLocalRepository())
	gin.SetMode(gin.TestMode)

	e := setupTest(t, "/example/:id", myhandler.GetProductByID)

	obj := e.GET("/example/1").
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Value("ID").Number().Gt(0)
	obj.Value("Name").String().NotEmpty()
	obj.Value("Price").Number().Ge(0)
	obj.Value("WrapAfter").Boolean()
	obj.Value("Pos").Number().Ge(0)
	obj.Value("ApiExport").Boolean()
}
