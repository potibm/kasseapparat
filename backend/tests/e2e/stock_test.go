package tests_e2e

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stockProduct creates a product and then sets its stock. The create endpoint takes no
// totalStock, so it has to be set with an update.
func stockProduct(t *testing.T, name string, totalStock int) int {
	t.Helper()

	res := withDemoUserAuthToken(e.POST(productBaseURL)).
		WithJSON(map[string]any{
			"name":     name,
			"netPrice": "10.00",
			"vatRate":  "25",
			"pos":      90,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	productID := int(res.Value("id").Number().Raw())

	withDemoUserAuthToken(e.PUT(productBaseURL + "/" + strconv.Itoa(productID))).
		WithJSON(map[string]any{
			"name":       name,
			"netPrice":   "10.00",
			"vatRate":    "25",
			"pos":        90,
			"totalStock": totalStock,
		}).
		Expect().
		Status(http.StatusOK)

	require.Equal(t, totalStock, reloadProduct(t, productID).TotalStock,
		"the product must carry the requested stock")

	return productID
}

// unlimitedProduct creates a product that is never stock-restricted.
func unlimitedProduct(t *testing.T, name string) int {
	t.Helper()

	res := withDemoUserAuthToken(e.POST(productBaseURL)).
		WithJSON(map[string]any{
			"name":     name,
			"netPrice": "10.00",
			"vatRate":  "25",
			"pos":      91,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	return int(res.Value("id").Number().Raw())
}

// reloadProduct reads the product back through the API so the assertions cover what
// the POS actually sees.
func reloadProduct(t *testing.T, productID int) models.Product {
	t.Helper()

	var product models.Product

	require.NoError(t, db.First(&product, productID).Error)

	return product
}

func purchasePayload(productID, quantity int) map[string]any {
	// The seeded test product is 10.00 net at 25% VAT, so gross is 12.50 per unit.
	unitNet := 10.0
	unitGross := unitNet * 1.25

	return map[string]any{
		"paymentMethod":   "CASH",
		"totalNetPrice":   fmt.Sprintf("%.2f", unitNet*float64(quantity)),
		"totalGrossPrice": fmt.Sprintf("%.2f", unitGross*float64(quantity)),
		"cart": []map[string]any{
			{
				"ID":        productID,
				"quantity":  quantity,
				"netPrice":  "10.00",
				"listItems": []map[string]any{},
			},
		},
	}
}

func TestOutOfStockBehavior_Ignore(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockIgnore)
	defer cleanup()

	productID := stockProduct(t, "Ignore Mugs", 2)

	// Buying beyond the stock is allowed and recorded.
	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 5)).
		Expect().
		Status(http.StatusCreated)

	product := reloadProduct(t, productID)
	assert.False(t, product.SoldOut, "ignore must not flag the product")
	assert.False(t, product.Hidden, "ignore must not hide the product")
}

func TestOutOfStockBehavior_Fail(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockFail)
	defer cleanup()

	productID := stockProduct(t, "Fail Mugs", 2)

	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 2)).
		Expect().
		Status(http.StatusCreated)

	// The stock is gone, so the next purchase is rejected with a clear message.
	res := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 1)).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object()

	assert.Contains(t, res.Value("details").String().Raw(), "sold out")

	product := reloadProduct(t, productID)
	assert.False(t, product.SoldOut, "fail must not flag the product")
	assert.False(t, product.Hidden, "fail must not hide the product")
}

func TestOutOfStockBehavior_FailReportsRemainingStock(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockFail)
	defer cleanup()

	productID := stockProduct(t, "Partial Mugs", 5)

	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 3)).
		Expect().
		Status(http.StatusCreated)

	res := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 3)).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object()

	// The handler capitalises the first rune of the detail.
	assert.Contains(t, res.Value("details").String().Raw(), "Only 2 left in stock")
}

func TestOutOfStockBehavior_AutoSoldOut(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockAutoSoldOut)
	defer cleanup()

	productID := stockProduct(t, "Auto Sold Out Shirts", 2)

	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 2)).
		Expect().
		Status(http.StatusCreated)

	product := reloadProduct(t, productID)
	assert.True(t, product.SoldOut, "selling the last unit should mark the product sold out")
	assert.False(t, product.Hidden, "auto_sold_out must not hide the product")

	// Further sales are rejected.
	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 1)).
		Expect().
		Status(http.StatusBadRequest)
}

func TestOutOfStockBehavior_AutoHide(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockAutoHide)
	defer cleanup()

	productID := stockProduct(t, "Auto Hide Shirts", 2)

	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 2)).
		Expect().
		Status(http.StatusCreated)

	product := reloadProduct(t, productID)
	assert.True(t, product.Hidden, "selling the last unit should hide the product")
	assert.False(t, product.SoldOut, "auto_hide must not mark the product sold out")

	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 1)).
		Expect().
		Status(http.StatusBadRequest)
}

// A product with no total stock is unlimited and must never be restricted, whatever
// the behaviour is.
func TestOutOfStockBehavior_UnlimitedProductNeverBlocked(t *testing.T) {
	for _, behavior := range []config.OutOfStockBehavior{
		config.OutOfStockFail,
		config.OutOfStockAutoSoldOut,
		config.OutOfStockAutoHide,
	} {
		t.Run(string(behavior), func(t *testing.T) {
			_, cleanup := setupTestEnvironmentWithBehavior(t, behavior)

			defer cleanup()

			productID := unlimitedProduct(t, "Drinks "+string(behavior))

			withDemoUserAuthToken(e.POST(purchaseBaseURL)).
				WithJSON(purchasePayload(productID, 50)).
				Expect().
				Status(http.StatusCreated)

			product := reloadProduct(t, productID)
			assert.False(t, product.SoldOut, "an unlimited product must never be flagged")
			assert.False(t, product.Hidden, "an unlimited product must never be hidden")
		})
	}
}

// The stock check runs inside the purchase transaction, so the sold quantity it reads
// must reflect confirmed purchases only.
func TestOutOfStockBehavior_CountsConfirmedPurchases(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockFail)
	defer cleanup()

	productID := stockProduct(t, "Counted Mugs", 3)

	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 2)).
		Expect().
		Status(http.StatusCreated)

	// One unit left, so a single unit still fits.
	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 1)).
		Expect().
		Status(http.StatusCreated)

	// Nothing left.
	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 1)).
		Expect().
		Status(http.StatusBadRequest)
}

func TestOutOfStockBehavior_RefundRestoresStock(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockAutoSoldOut)
	defer cleanup()

	productID := stockProduct(t, "Refunded Mugs", 2)

	res := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 2)).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	purchaseID := res.Value("id").String().Raw()

	require.True(t, reloadProduct(t, productID).SoldOut)

	// Refunding frees the units, so the product is sellable again and unflagged.
	withDemoUserAuthToken(e.POST(purchaseBaseURL + "/" + purchaseID + "/refund")).
		Expect().
		Status(http.StatusOK)

	product := reloadProduct(t, productID)
	assert.False(t, product.SoldOut, "refunding should clear the sold-out flag")

	withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(purchasePayload(productID, 2)).
		Expect().
		Status(http.StatusCreated)
}

// The behaviour must reach the frontend so it can guide the cart.
func TestOutOfStockBehavior_ExposedOnConfig(t *testing.T) {
	_, cleanup := setupTestEnvironmentWithBehavior(t, config.OutOfStockAutoHide)
	defer cleanup()

	e.GET(configURL).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("outOfStockBehavior").
		String().
		IsEqual("auto_hide")
}
