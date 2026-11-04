package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupStockFlagHandler wires a handler with a real repository and purchase service, so a
// test can observe what a write does to the product flags.
func setupStockFlagHandler(t *testing.T, behavior config.OutOfStockBehavior) *Handler {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := utils.ConnectToLocalDatabase()
	require.NoError(t, err)
	require.NoError(t, utils.PurgeDatabase(db))
	require.NoError(t, utils.MigrateDatabase(db))

	repo := sqliteRepo.NewRepository(db, 2)
	service := purchaseService.NewPurchaseService(repo, nil, nil, 2, "EUR", behavior)

	return NewHandler(HandlerConfig{
		Repo:            repo,
		PurchaseService: service,
	})
}

// seedStockFlagProduct stores a limited product with a flag the behaviour owns.
func seedStockFlagProduct(t *testing.T, repo sqliteRepo.RepositoryInterface) models.Product {
	t.Helper()

	product, err := repo.CreateProduct(models.Product{
		Name:              "Shirts",
		NetPrice:          decimal.NewFromInt(10),
		VATRate:           decimal.NewFromInt(25),
		TotalStock:        1,
		SoldOut:           true,
		StockFlagsManaged: true,
	})
	require.NoError(t, err)

	return product
}

func ginContextWithParam(t *testing.T, method, target string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	req := httptest.NewRequest(method, target, http.NoBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

func jsonBody(t *testing.T, c *gin.Context, payload any) {
	t.Helper()

	encoded, err := json.Marshal(payload)
	require.NoError(t, err)

	c.Request = httptest.NewRequest(c.Request.Method, c.Request.URL.String(), bytes.NewReader(encoded))
	c.Request.Header.Set("Content-Type", "application/json")
}

// Editing the stock limit is one of the writes the purchase flow never sees, so the
// handler has to re-derive the flags itself.
func TestUpdateProductByID_ReconcilesStockFlags(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockAutoSoldOut)
	repo := handler.repo

	product := seedStockFlagProduct(t, repo)

	c, _ := ginContextWithParam(t, http.MethodPut, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(product.ID)}}
	jsonBody(t, c, ProductRequestUpdate{
		Name:       product.Name,
		NetPrice:   decimal.NewFromInt(10),
		VATRate:    decimal.NewFromInt(25),
		Pos:        1,
		TotalStock: 5,
	})

	handler.UpdateProductByID(c)

	require.Equal(t, http.StatusOK, c.Writer.Status())

	reloaded, err := repo.GetProductByID(product.ID)
	require.NoError(t, err)
	assert.False(t, reloaded.SoldOut, "the raised stock releases the flag")
	assert.False(t, reloaded.StockFlagsManaged)
}

// A deleted purchase released its units, so the flags of its products are stale.
func TestDeletePurchase_ReconcilesStockFlags(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockAutoSoldOut)
	repo := handler.repo

	product := seedStockFlagProduct(t, repo)

	purchase := models.Purchase{
		ID:     uuid.New(),
		Status: models.PurchaseStatusConfirmed,
		PurchaseItems: []models.PurchaseItem{{
			ProductID: product.ID,
			Quantity:  1,
			NetPrice:  decimal.NewFromInt(10),
			VATRate:   decimal.NewFromInt(25),
		}},
	}
	_, err := repo.StorePurchases(purchase)
	require.NoError(t, err)

	c, _ := ginContextWithParam(t, http.MethodDelete, "/")
	c.Params = gin.Params{{Key: "id", Value: purchase.ID.String()}}

	handler.DeletePurchase(c)

	require.Equal(t, http.StatusNoContent, c.Writer.Status())

	reloaded, err := repo.GetProductByID(product.ID)
	require.NoError(t, err)
	assert.False(t, reloaded.SoldOut, "the deleted purchase released its units")
	assert.False(t, reloaded.StockFlagsManaged)
}
