package sqlite

import (
	"testing"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupStockRepository drops and recreates the tables, so the sold quantity of a test
// cannot come from a purchase another test stored.
func setupStockRepository(t *testing.T) *Repository {
	t.Helper()

	repo := setupTestRepository(t)
	require.NoError(t, utils.PurgeDatabase(repo.db))
	require.NoError(t, utils.MigrateDatabase(repo.db))

	return repo
}

func stockTestProduct(name string, totalStock int) models.Product {
	return models.Product{
		Name:       name,
		NetPrice:   decimal.NewFromInt(10),
		VATRate:    decimal.NewFromInt(25),
		TotalStock: totalStock,
	}
}

// seedStockPurchase stores a purchase in a status, with one item of the given quantity.
func seedStockPurchase(
	t *testing.T,
	repo *Repository,
	product models.Product,
	quantity uint,
	status models.PurchaseStatus,
) {
	t.Helper()

	_, err := repo.StorePurchases(models.Purchase{
		ID:     uuid.New(),
		Status: status,
		PurchaseItems: []models.PurchaseItem{{
			ProductID: product.ID,
			Quantity:  quantity,
			NetPrice:  product.NetPrice,
			VATRate:   product.VATRate,
		}},
	})
	require.NoError(t, err)
}

// The reconciliation reads the sold quantities for all products in one query, so the
// query has to sum per product and per status.
func TestGetUnitsSoldByProductID(t *testing.T) {
	repo := setupStockRepository(t)

	shirts, err := repo.CreateProduct(stockTestProduct("Shirts", 3))
	require.NoError(t, err)

	mugs, err := repo.CreateProduct(stockTestProduct("Mugs", 3))
	require.NoError(t, err)

	unitsSold, err := repo.GetUnitsSoldByProductID()
	require.NoError(t, err)
	assert.Empty(t, unitsSold, "nothing is sold yet")

	seedStockPurchase(t, repo, shirts, 2, models.PurchaseStatusConfirmed)
	seedStockPurchase(t, repo, shirts, 2, models.PurchaseStatusConfirmed)
	seedStockPurchase(t, repo, shirts, 1, models.PurchaseStatusPending)
	seedStockPurchase(t, repo, mugs, 4, models.PurchaseStatusRefunded)

	unitsSold, err = repo.GetUnitsSoldByProductID()
	require.NoError(t, err)

	assert.Equal(t, 4, unitsSold[shirts.ID], "only confirmed purchases count")
	assert.NotContains(t, unitsSold, mugs.ID, "a refunded purchase released its units")
	assert.NotContains(t, unitsSold, mugs.ID)
}

// The reconciliation walks every product, so the query has to load them all and not page.
func TestGetProductsForStockReconciliation(t *testing.T) {
	repo := setupStockRepository(t)

	shirts, err := repo.CreateProduct(stockTestProduct("Shirts", 3))
	require.NoError(t, err)

	mugs, err := repo.CreateProduct(stockTestProduct("Mugs", 0))
	require.NoError(t, err)

	products, err := repo.GetProductsForStockReconciliation()
	require.NoError(t, err)

	require.Len(t, products, 2)
	assert.Equal(t, shirts.ID, products[0].ID, "ordered by id")
	assert.Equal(t, mugs.ID, products[1].ID)
	assert.Equal(t, 3, products[0].TotalStock)
	assert.Equal(t, 0, products[1].TotalStock)
	assert.False(t, products[0].SoldOut)
	assert.False(t, products[0].Hidden)
}

// The flags are written without reading the product first, so the other columns of a
// product the query selected only in part must survive.
func TestUpdateProductStockFlags(t *testing.T) {
	repo := setupStockRepository(t)

	product, err := repo.CreateProduct(stockTestProduct("Shirts", 3))
	require.NoError(t, err)

	product.SoldOut = true

	stored, err := repo.UpdateProductByID(product.ID, product)
	require.NoError(t, err)
	require.True(t, stored.SoldOut)

	err = repo.UpdateProductStockFlags(product.ID, false, true, true)
	require.NoError(t, err)

	reloaded, err := repo.GetProductByID(product.ID)
	require.NoError(t, err)

	assert.False(t, reloaded.SoldOut)
	assert.True(t, reloaded.Hidden)
	assert.True(t, reloaded.StockFlagsManaged)
	assert.Equal(t, 3, reloaded.TotalStock, "the stock must survive the flag update")
	assert.Equal(t, product.NetPrice.String(), reloaded.NetPrice.String(),
		"the price must survive the flag update")
}

func TestUpdateProductStockFlags_UnknownProduct(t *testing.T) {
	repo := setupStockRepository(t)

	err := repo.UpdateProductStockFlags(4242, true, false, false)

	assert.Error(t, err)
}
