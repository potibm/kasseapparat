package purchase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailableStock(t *testing.T) {
	tests := []struct {
		name          string
		totalStock    int
		unitsSold     int
		wantAvailable int
		wantUnlimited bool
	}{
		{name: "no total stock is unlimited", totalStock: 0, unitsSold: 0, wantUnlimited: true},
		{name: "no total stock stays unlimited when oversold", totalStock: 0, unitsSold: 99, wantUnlimited: true},
		{name: "negative total stock is unlimited", totalStock: -5, unitsSold: 0, wantUnlimited: true},
		{name: "nothing sold", totalStock: 10, unitsSold: 0, wantAvailable: 10},
		{name: "partly sold", totalStock: 10, unitsSold: 4, wantAvailable: 6},
		{name: "exactly sold out", totalStock: 10, unitsSold: 10, wantAvailable: 0},
		{name: "oversold clamps to zero", totalStock: 10, unitsSold: 13, wantAvailable: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			available, unlimited := AvailableStock(
				models.Product{TotalStock: tc.totalStock},
				tc.unitsSold,
			)

			assert.Equal(t, tc.wantUnlimited, unlimited)
			assert.Equal(t, tc.wantAvailable, available)
		})
	}
}

func TestRestrictsStock(t *testing.T) {
	// An unset or unknown value must behave like the documented default and never
	// start rejecting sales.
	assert.False(t, restrictsStock(config.OutOfStockIgnore))
	assert.False(t, restrictsStock(""))
	assert.False(t, restrictsStock("nonsense"))
	assert.True(t, restrictsStock(config.OutOfStockFail))
	assert.True(t, restrictsStock(config.OutOfStockAutoSoldOut))
	assert.True(t, restrictsStock(config.OutOfStockAutoHide))
}

func TestInsufficientStockError(t *testing.T) {
	err := error(insufficientStockError{ProductName: "T-Shirt", Available: 2})
	assert.ErrorIs(t, err, ErrInsufficientStock)
	assert.Equal(t, "only 2 left in stock for T-Shirt", err.Error())

	soldOut := error(insufficientStockError{ProductName: "T-Shirt", Available: 0})
	assert.ErrorIs(t, soldOut, ErrInsufficientStock)
	assert.Equal(t, "T-Shirt is sold out", soldOut.Error())
}

func TestInsufficientStockError_WrapsThroughTransaction(t *testing.T) {
	// The stock error travels out of the transaction callback wrapped, which is why
	// the handler matches with errors.Is rather than by identity.
	inner := insufficientStockError{ProductName: "Mug", Available: 1}
	wrapped := fmt.Errorf("failed to store purchase: %w", error(inner))

	assert.True(t, errors.Is(wrapped, ErrInsufficientStock))
}

// stockTestProduct is a limited-stock product with nothing sold yet.
func stockTestProduct(totalStock int) *models.Product {
	return &models.Product{
		ID:         7,
		Name:       "T-Shirt",
		NetPrice:   decimal.NewFromInt(10),
		VATRate:    decimal.NewFromInt(25),
		TotalStock: totalStock,
	}
}

func newStockService(
	behavior config.OutOfStockBehavior,
	product *models.Product,
	unitsSold int,
) *PurchaseService {
	repo := &MockRepository{
		Products:      map[int]*models.Product{product.ID: product},
		Guests:        map[int]*models.Guest{},
		UpdatedGuests: map[int]*models.Guest{},
	}
	repo.seedUnitsSold(product, unitsSold)

	return NewPurchaseService(repo, nil, nil, 2, "EUR", behavior)
}

// unitsSold reads the sold quantity back through the repository interface, the same
// way the service does.
func unitsSold(t *testing.T, service *PurchaseService, productID int) int {
	t.Helper()

	sold, err := service.sqliteRepo.GetPurchasedQuantitiesByProductID(productID)
	require.NoError(t, err)

	return sold
}

// purchaseInputFor builds a cart whose prices and totals match the product, since the
// service validates them before it ever looks at stock.
func purchaseInputFor(product *models.Product, quantity uint) PurchaseInput {
	net := product.NetPrice.Mul(decimal.NewFromUint64(uint64(quantity)))

	return PurchaseInput{
		Cart: []PurchaseCartItem{{
			ID:       product.ID,
			NetPrice: product.NetPrice,
			Quantity: quantity,
		}},
		TotalNetPrice:   net,
		TotalGrossPrice: net.Add(net.Mul(product.VATRate.Div(decimal.NewFromInt(100)))),
		PaymentMethod:   models.PaymentMethodCash,
	}
}

// The reconciliation derives both flags and the ownership from the stock, so it is the
// single decision the behaviour makes.
func TestReconcileStockFlags(t *testing.T) {
	tests := []struct {
		name string

		behavior   config.OutOfStockBehavior
		totalStock int
		unitsSold  int

		soldOut  bool
		hidden   bool
		managed  bool
		wantTick bool
	}{
		{
			// ignore forbids nothing, so it never touches a product.
			name:     "ignore leaves a depleted product alone",
			behavior: config.OutOfStockIgnore, totalStock: 3, unitsSold: 3,
		},
		{
			// fail rejects the sale but does not manage the product.
			name:     "fail leaves a depleted product alone",
			behavior: config.OutOfStockFail, totalStock: 3, unitsSold: 3,
		},
		{
			name:     "auto_sold_out flags a depleted product",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 3, unitsSold: 3,
			wantTick: true,
		},
		{
			// The comparison is >=, so an already oversold product still converges.
			name:     "auto_sold_out flags an oversold product",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 2, unitsSold: 5,
			wantTick: true,
		},
		{
			name:     "auto_sold_out leaves a product with stock alone",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 3, unitsSold: 2,
		},
		{
			name:     "auto_hide hides a depleted product",
			behavior: config.OutOfStockAutoHide, totalStock: 3, unitsSold: 3,
			wantTick: true,
		},
		{
			// The stock is back through a refund, a failed or a deleted purchase.
			name:     "the stock is back, so the flag is released",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 3, unitsSold: 1,
			soldOut: true, managed: true, wantTick: true,
		},
		{
			// The mode changed since the flag was set, so a startup reconciliation
			// re-evaluates it instead of leaving it behind.
			name:     "the mode changed, so the flag is released",
			behavior: config.OutOfStockFail, totalStock: 3, unitsSold: 3,
			soldOut: true, managed: true, wantTick: true,
		},
		{
			name:     "a product the behaviour already flagged correctly is not rewritten",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 3, unitsSold: 3,
			soldOut: true, managed: true,
		},
		{
			// An operator hid or sold out the product by hand, so the behaviour does
			// not take it over and does not clear it either.
			name:     "a manual sold-out flag is left alone",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 3, unitsSold: 3,
			soldOut: true,
		},
		{
			name:     "a manual hidden flag is left alone",
			behavior: config.OutOfStockAutoHide, totalStock: 3, unitsSold: 3,
			hidden: true,
		},
		{
			name:     "an unlimited product is never flagged",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 0, unitsSold: 9,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			product := stockTestProduct(tc.totalStock)
			product.SoldOut = tc.soldOut
			product.Hidden = tc.hidden
			product.StockFlagsManaged = tc.managed

			service := newStockService(tc.behavior, product, tc.unitsSold)
			repo := service.sqliteRepo.(*MockRepository)

			err := service.ReconcileStockFlags(context.Background())
			require.NoError(t, err)

			wantSoldOut, wantHidden, wantManaged := wantedStockFlags(
				tc.behavior, tc.totalStock, tc.unitsSold, tc.soldOut, tc.hidden, tc.managed,
			)

			assert.Equal(t, wantSoldOut, product.SoldOut, "sold out flag")
			assert.Equal(t, wantHidden, product.Hidden, "hidden flag")
			assert.Equal(t, wantManaged, product.StockFlagsManaged, "ownership")

			if tc.wantTick {
				assert.Equal(t, []int{product.ID}, repo.UpdatedProducts, "expected exactly one product write")
			} else {
				assert.Empty(t, repo.UpdatedProducts, "nothing to do, so nothing may be written")
			}
		})
	}
}

// wantedStockFlags is the state the behaviour should leave a product in.
func wantedStockFlags(
	behavior config.OutOfStockBehavior,
	totalStock, unitsSold int,
	soldOut, hidden, managed bool,
) (wantSoldOut, wantHidden, wantManaged bool) {
	depleted := totalStock > 0 && unitsSold >= totalStock

	wantSoldOut = depleted && behavior == config.OutOfStockAutoSoldOut
	wantHidden = depleted && behavior == config.OutOfStockAutoHide
	wantManaged = wantSoldOut || wantHidden

	// A flag the behaviour does not own keeps its state, and so does a product whose
	// stock is still gone while the flag is the behaviour's.
	if !managed && (soldOut || hidden) {
		wantSoldOut, wantHidden, wantManaged = soldOut, hidden, false
	} else if !depleted && managed {
		wantSoldOut, wantHidden, wantManaged = false, false, false
	}

	return wantSoldOut, wantHidden, wantManaged
}

// A manual flag must survive a refund that frees the stock again, otherwise one refund
// would put a product an operator hid by hand back on sale.
func TestReconcileStockFlags_LeavesManualFlagAfterRefund(t *testing.T) {
	product := stockTestProduct(2)
	product.Hidden = true
	service := newStockService(config.OutOfStockAutoHide, product, 0)

	purchase, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(product, 2),
	)
	require.NoError(t, err)
	require.True(t, product.Hidden, "the operator hid the product")

	_, err = service.RefundPurchase(context.Background(), purchase.ID)
	require.NoError(t, err)

	assert.True(t, product.Hidden, "a manual flag is not the behaviour's to clear")
	assert.Equal(t, 0, unitsSold(t, service, product.ID))
}

// The reconciliation walks every product, so it also fixes one the purchase flow never
// touches: an operator raised the stock of a product that was flagged while it was
// depleted, and no purchase will ever release that flag.
func TestReconcileStockFlags_CoversProductsWithoutAPurchase(t *testing.T) {
	stillDepleted := stockTestProduct(2)
	stillDepleted.SoldOut = true
	stillDepleted.StockFlagsManaged = true

	restocked := stockTestProduct(5)
	restocked.ID = 8
	restocked.Name = "Mug"
	restocked.SoldOut = true
	restocked.StockFlagsManaged = true

	service := newStockService(config.OutOfStockAutoSoldOut, stillDepleted, 2)
	repo := service.sqliteRepo.(*MockRepository)
	repo.Products[restocked.ID] = restocked

	require.NoError(t, service.ReconcileStockFlags(context.Background()))

	assert.Equal(t, []int{restocked.ID}, repo.UpdatedProducts,
		"only the product whose flag changed is written")
	assert.True(t, stillDepleted.SoldOut, "the stock is still gone, so the flag stays")
	assert.False(t, restocked.SoldOut, "the raised stock releases the flag")
	assert.False(t, restocked.StockFlagsManaged)
}

func TestCreateConfirmedPurchase_StockBehaviours(t *testing.T) {
	const productID = 7

	tests := []struct {
		name            string
		behavior        config.OutOfStockBehavior
		totalStock      int
		unitsSold       int
		quantity        uint
		wantErr         bool
		wantUnitsSold   int
		wantSoldOutFlag bool
		wantHiddenFlag  bool
	}{
		{
			name: "ignore allows going beyond stock", behavior: config.OutOfStockIgnore,
			totalStock: 2, unitsSold: 1, quantity: 5,
			wantErr: false, wantUnitsSold: 6,
		},
		{
			name: "ignore leaves an unlimited product alone", behavior: config.OutOfStockIgnore,
			totalStock: 0, unitsSold: 0, quantity: 100,
			wantErr: false, wantUnitsSold: 100,
		},
		{
			name: "fail allows exactly the remaining stock", behavior: config.OutOfStockFail,
			totalStock: 5, unitsSold: 3, quantity: 2,
			wantErr: false, wantUnitsSold: 5,
		},
		{
			name: "fail rejects beyond the remaining stock", behavior: config.OutOfStockFail,
			totalStock: 5, unitsSold: 3, quantity: 3,
			wantErr: true, wantUnitsSold: 3,
		},
		{
			name: "fail never touches the product flags", behavior: config.OutOfStockFail,
			totalStock: 2, unitsSold: 2, quantity: 1,
			wantErr: true, wantUnitsSold: 2,
		},
		{
			name: "fail ignores unlimited products", behavior: config.OutOfStockFail,
			totalStock: 0, unitsSold: 0, quantity: 50,
			wantErr: false, wantUnitsSold: 50,
		},
		{
			name:     "auto_sold_out sells the last unit and flags the product",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 3, unitsSold: 2, quantity: 1,
			wantErr: false, wantUnitsSold: 3, wantSoldOutFlag: true,
		},
		{
			name:     "auto_sold_out rejects the unit after the last",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 3, unitsSold: 3, quantity: 1,
			wantErr: true, wantUnitsSold: 3,
		},
		{
			name:     "auto_sold_out does not flag while stock remains",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 10, unitsSold: 2, quantity: 1,
			wantErr: false, wantUnitsSold: 3, wantSoldOutFlag: false,
		},
		{
			name:     "auto_hide sells the last unit and hides the product",
			behavior: config.OutOfStockAutoHide, totalStock: 3, unitsSold: 2, quantity: 1,
			wantErr: false, wantUnitsSold: 3, wantHiddenFlag: true,
		},
		{
			name:     "auto_hide rejects the unit after the last",
			behavior: config.OutOfStockAutoHide, totalStock: 3, unitsSold: 3, quantity: 1,
			wantErr: true, wantUnitsSold: 3,
		},
		{
			name:     "auto_hide ignores unlimited products",
			behavior: config.OutOfStockAutoHide, totalStock: 0, unitsSold: 0, quantity: 9,
			wantErr: false, wantUnitsSold: 9, wantHiddenFlag: false,
		},
		{
			// The comparison is >=, so an already oversold product still converges
			// instead of never being flagged.
			name:     "auto_sold_out flags an already oversold product",
			behavior: config.OutOfStockAutoSoldOut, totalStock: 2, unitsSold: 5, quantity: 1,
			wantErr: true, wantUnitsSold: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			product := stockTestProduct(tc.totalStock)
			service := newStockService(tc.behavior, product, tc.unitsSold)

			input := purchaseInputFor(product, tc.quantity)

			purchase, err := service.CreateConfirmedPurchase(context.Background(), input)

			if tc.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInsufficientStock)
				assert.Nil(t, purchase)
			} else {
				require.NoError(t, err)
				require.NotNil(t, purchase)
			}

			assert.Equal(t, tc.wantUnitsSold, unitsSold(t, service, productID),
				"a rejected purchase must not change the sold quantity")
			assert.Equal(t, tc.wantSoldOutFlag, product.SoldOut)
			assert.Equal(t, tc.wantHiddenFlag, product.Hidden)
		})
	}
}

// A rejected purchase must leave no trace: no stored purchase and no product change.
func TestCreateConfirmedPurchase_RejectionStoresNothing(t *testing.T) {
	product := stockTestProduct(1)
	service := newStockService(config.OutOfStockFail, product, 1)
	repo := service.sqliteRepo.(*MockRepository)

	// The seeded purchase is already stored, so compare against its count.
	storedBefore := len(repo.storedPurchases)

	_, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(product, 1),
	)

	require.ErrorIs(t, err, ErrInsufficientStock)
	assert.Len(t, repo.storedPurchases, storedBefore, "a rejected purchase must not be stored")
	assert.False(t, product.SoldOut)
	assert.False(t, product.Hidden)
}

// TotalStock 0 is excluded from every mode, so a stock-less product always sells.
func TestCreateConfirmedPurchase_UnlimitedProductNeverBlocked(t *testing.T) {
	for _, behavior := range []config.OutOfStockBehavior{
		config.OutOfStockIgnore,
		config.OutOfStockFail,
		config.OutOfStockAutoSoldOut,
		config.OutOfStockAutoHide,
	} {
		t.Run(string(behavior), func(t *testing.T) {
			product := stockTestProduct(0)
			service := newStockService(behavior, product, 250)

			_, err := service.CreateConfirmedPurchase(
				context.Background(),
				purchaseInputFor(product, 25),
			)

			require.NoError(t, err)
			assert.False(t, product.SoldOut, "an unlimited product must never be flagged")
			assert.False(t, product.Hidden, "an unlimited product must never be hidden")
		})
	}
}

// Several lines of the same product in one cart must be checked as a whole, otherwise
// each line would be validated against the same availability and oversell.
func TestCreateConfirmedPurchase_ChecksTotalCartQuantityPerProduct(t *testing.T) {
	product := stockTestProduct(3)
	service := newStockService(config.OutOfStockFail, product, 1)

	net := product.NetPrice.Mul(decimal.NewFromInt(4))
	input := PurchaseInput{
		Cart: []PurchaseCartItem{
			{ID: product.ID, NetPrice: product.NetPrice, Quantity: 2},
			{ID: product.ID, NetPrice: product.NetPrice, Quantity: 2},
		},
		TotalNetPrice:   net,
		TotalGrossPrice: net.Add(net.Mul(product.VATRate.Div(decimal.NewFromInt(100)))),
		PaymentMethod:   models.PaymentMethodCash,
	}

	_, err := service.CreateConfirmedPurchase(context.Background(), input)

	require.ErrorIs(t, err, ErrInsufficientStock)
	assert.Equal(t, 1, unitsSold(t, service, product.ID))
}

// Only a confirmed purchase holds its units, so a pending SumUp payment does not block
// another client while it is still in flight, and does not flag the product either.
func TestCreatePendingPurchase_DoesNotHoldUnits(t *testing.T) {
	product := stockTestProduct(2)
	service := newStockService(config.OutOfStockAutoSoldOut, product, 0)

	pending := purchaseInputFor(product, 2)
	pending.PaymentMethod = models.PaymentMethodSumUp

	_, err := service.CreatePendingPurchase(context.Background(), pending)
	require.NoError(t, err)
	assert.Equal(t, 0, unitsSold(t, service, product.ID),
		"a pending purchase is not sold yet")
	assert.False(t, product.SoldOut, "a pending purchase does not deplete the stock")

	second := purchaseInputFor(product, 2)
	second.PaymentMethod = models.PaymentMethodSumUp

	_, err = service.CreatePendingPurchase(context.Background(), second)
	assert.NoError(t, err, "a pending purchase does not reserve the stock either")
}

// The units only become sold once the payment is confirmed, which is also when the
// product flag is set.
func TestFinalizePurchase_SellsTheReservedUnits(t *testing.T) {
	product := stockTestProduct(2)
	service := newStockService(config.OutOfStockAutoSoldOut, product, 0)

	pending := purchaseInputFor(product, 2)
	pending.PaymentMethod = models.PaymentMethodSumUp

	created, err := service.CreatePendingPurchase(context.Background(), pending)
	require.NoError(t, err)

	_, err = service.FinalizePurchase(context.Background(), created.ID)
	require.NoError(t, err)

	assert.Equal(t, 2, unitsSold(t, service, product.ID))
	assert.True(t, product.SoldOut, "the last unit is sold once the payment is confirmed")
	assert.True(t, product.StockFlagsManaged)
}

// A cancelled payment leaves the stock untouched, so an in-flight payment that never
// completes cannot remove a product from sale.
func TestCancelPurchase_KeepsStockAvailable(t *testing.T) {
	product := stockTestProduct(2)
	service := newStockService(config.OutOfStockAutoHide, product, 0)

	pending := purchaseInputFor(product, 2)
	pending.PaymentMethod = models.PaymentMethodSumUp

	created, err := service.CreatePendingPurchase(context.Background(), pending)
	require.NoError(t, err)

	_, err = service.CancelPurchase(context.Background(), created.ID)
	require.NoError(t, err)

	assert.Equal(t, 0, unitsSold(t, service, product.ID))
	assert.False(t, product.Hidden)

	_, err = service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(product, 2),
	)
	assert.NoError(t, err, "the stock is still sellable")
}

// Refunding frees the units, and the flag the behaviour set must be cleared with them,
// otherwise one refund keeps the product off sale for the rest of the event.
func TestRefundPurchase_ClearsAutomaticFlag(t *testing.T) {
	for _, tc := range []struct {
		name     string
		behavior config.OutOfStockBehavior
	}{
		{name: "auto_sold_out", behavior: config.OutOfStockAutoSoldOut},
		{name: "auto_hide", behavior: config.OutOfStockAutoHide},
	} {
		t.Run(tc.name, func(t *testing.T) {
			product := stockTestProduct(2)
			service := newStockService(tc.behavior, product, 0)

			// Take the whole stock, which sets the automatic flag.
			purchase, err := service.CreateConfirmedPurchase(
				context.Background(),
				purchaseInputFor(product, 2),
			)
			require.NoError(t, err)

			if tc.behavior == config.OutOfStockAutoSoldOut {
				require.True(t, product.SoldOut, "selling the last unit should flag the product")
			} else {
				require.True(t, product.Hidden, "selling the last unit should hide the product")
			}

			// Further sales are rejected while the stock is gone.
			_, err = service.CreateConfirmedPurchase(
				context.Background(),
				purchaseInputFor(product, 1),
			)
			require.ErrorIs(t, err, ErrInsufficientStock)

			_, err = service.RefundPurchase(context.Background(), purchase.ID)
			require.NoError(t, err)

			assert.False(t, product.SoldOut, "refunding should clear an automatic sold-out flag")
			assert.False(t, product.Hidden, "refunding should clear an automatic hidden flag")
			assert.False(t, product.StockFlagsManaged)
			assert.Equal(t, 0, unitsSold(t, service, product.ID))

			// And the product sells again.
			_, err = service.CreateConfirmedPurchase(
				context.Background(),
				purchaseInputFor(product, 1),
			)
			assert.NoError(t, err)
		})
	}
}

// The flag must survive while the stock is still gone, so a status change of an
// unrelated purchase cannot put a sold-out product back on sale.
func TestCancelPurchase_KeepsFlagWhileStockIsGone(t *testing.T) {
	product := stockTestProduct(1)
	service := newStockService(config.OutOfStockAutoSoldOut, product, 0)

	first, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(product, 1),
	)
	require.NoError(t, err)
	require.NotNil(t, first)
	require.True(t, product.SoldOut)

	// A second product's purchase is cancelled; the sold-out product keeps its flag.
	other := &models.Product{
		ID: 8, Name: "Mug", NetPrice: decimal.NewFromInt(5),
		VATRate: decimal.NewFromInt(25), TotalStock: 5,
	}
	repo := service.sqliteRepo.(*MockRepository)
	repo.Products[other.ID] = other

	second, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(other, 1),
	)
	require.NoError(t, err)

	_, err = service.CancelPurchase(context.Background(), second.ID)
	require.NoError(t, err)

	assert.True(t, product.SoldOut, "stock is still gone, so the flag must stay")
}

// The flag must survive while the stock is still gone, so a refund of one product
// cannot put a different sold-out product back on sale.
func TestRefundPurchase_KeepsFlagOnOtherProducts(t *testing.T) {
	soldOut := stockTestProduct(1)
	other := &models.Product{
		ID: 8, Name: "Mug", NetPrice: decimal.NewFromInt(5),
		VATRate: decimal.NewFromInt(25), TotalStock: 5,
	}
	service := newStockService(config.OutOfStockAutoSoldOut, soldOut, 0)
	repo := service.sqliteRepo.(*MockRepository)
	repo.Products[other.ID] = other

	// Taking the last unit of one product flags it.
	_, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(soldOut, 1),
	)
	require.NoError(t, err)
	require.True(t, soldOut.SoldOut)

	// Refunding an unrelated purchase of another product leaves it flagged.
	mugPurchase, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(other, 1),
	)
	require.NoError(t, err)

	_, err = service.RefundPurchase(context.Background(), mugPurchase.ID)
	require.NoError(t, err)

	assert.True(t, soldOut.SoldOut, "its stock is still gone, so the flag must stay")
}

// A failing reconciliation must not fail the purchase that succeeded: the sale went
// through, and the operator would be told otherwise.
func TestCreateConfirmedPurchase_ReconciliationFailureIsLogged(t *testing.T) {
	product := stockTestProduct(2)
	service := newStockService(config.OutOfStockAutoSoldOut, product, 0)
	repo := service.sqliteRepo.(*MockRepository)
	repo.reconcileError = errors.New("reconciliation is down")

	purchase, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(product, 1),
	)

	require.NoError(t, err)
	require.NotNil(t, purchase)
	assert.Equal(t, 1, unitsSold(t, service, product.ID))
}
