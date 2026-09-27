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

func TestAutoFlagForBehavior(t *testing.T) {
	assert.Empty(t, autoFlagForBehavior(config.OutOfStockIgnore))
	assert.Empty(t, autoFlagForBehavior(config.OutOfStockFail))
	assert.Empty(t, autoFlagForBehavior(""))
	assert.Equal(t, "soldOut", autoFlagForBehavior(config.OutOfStockAutoSoldOut))
	assert.Equal(t, "hidden", autoFlagForBehavior(config.OutOfStockAutoHide))
}

func TestInsufficientStockError(t *testing.T) {
	err := error(insufficientStockError{ProductName: "T-Shirt", Available: 2})
	assert.ErrorIs(t, err, ErrInsufficientStock)
	assert.Equal(t, "only 2 left in stock for T-Shirt", err.Error())

	soldOut := error(insufficientStockError{ProductName: "T-Shirt", Available: 0})
	assert.ErrorIs(t, soldOut, ErrInsufficientStock)
	assert.Equal(t, "T-Shirt is sold out", soldOut.Error())
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

func stockPurchaseInput(quantity uint) PurchaseInput {
	return PurchaseInput{
		Cart:            []PurchaseCartItem{{ID: 7, Quantity: quantity}},
		TotalNetPrice:   decimal.NewFromInt(0),
		TotalGrossPrice: decimal.NewFromInt(0),
		PaymentMethod:   models.PaymentMethodCash,
	}
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

// The pending SumUp path reserves stock, so a second purchase of the last unit is
// rejected while the first is still awaiting the terminal.
func TestCreatePendingPurchase_ReservesStock(t *testing.T) {
	product := stockTestProduct(2)
	service := newStockService(config.OutOfStockFail, product, 0)

	input := purchaseInputFor(product, 2)
	input.PaymentMethod = models.PaymentMethodSumUp

	pending, err := service.CreatePendingPurchase(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, pending)
	assert.Equal(t, models.PurchaseStatusPending, pending.Status)

	// The pending purchase holds both units, so nothing is left.
	input2 := purchaseInputFor(product, 1)
	input2.PaymentMethod = models.PaymentMethodSumUp

	_, err = service.CreatePendingPurchase(context.Background(), input2)
	assert.ErrorIs(t, err, ErrInsufficientStock)

	// Once the pending purchase is confirmed the units are still held.
	_, err = service.FinalizePurchase(context.Background(), pending.ID)
	require.NoError(t, err)

	_, err = service.CreatePendingPurchase(context.Background(), input2)
	assert.ErrorIs(t, err, ErrInsufficientStock)
}

// A failed purchase releases its units again, so the stock becomes sellable.
func TestFailPurchase_ReleasesReservedStock(t *testing.T) {
	product := stockTestProduct(2)
	service := newStockService(config.OutOfStockFail, product, 0)

	input := purchaseInputFor(product, 2)
	input.PaymentMethod = models.PaymentMethodSumUp

	pending, err := service.CreatePendingPurchase(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, 2, unitsSold(t, service, product.ID), "a pending purchase must reserve its units")

	_, err = service.FailPurchase(context.Background(), pending.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, unitsSold(t, service, product.ID), "a failed purchase must free its units")
}

// Refunding frees the units, and the automatically set flag must be cleared with them,
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

// The automatic flag must survive while the stock is still gone: clearing it on every
// status change would put a sold-out product back on sale.
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

// A refund re-enables a product whose stock came back. This also re-enables one an
// operator hid by hand, because the row does not record who set the flag; see the
// limitation noted on clearStockFlags.
func TestRefundPurchase_ReenablesProductWhenStockReturns(t *testing.T) {
	product := stockTestProduct(2)
	product.Hidden = true
	service := newStockService(config.OutOfStockAutoHide, product, 0)

	purchase, err := service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(product, 2),
	)
	require.NoError(t, err)

	_, err = service.RefundPurchase(context.Background(), purchase.ID)
	require.NoError(t, err)

	assert.False(t, product.Hidden, "the product is sellable again after the refund")

	_, err = service.CreateConfirmedPurchase(
		context.Background(),
		purchaseInputFor(product, 1),
	)
	assert.NoError(t, err)
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

func TestPurchaseStatusConsumesStock(t *testing.T) {
	assert.True(t, models.PurchaseStatusConfirmed.ConsumesStock())
	assert.True(t, models.PurchaseStatusPending.ConsumesStock())
	assert.False(t, models.PurchaseStatusRefunded.ConsumesStock())
	assert.False(t, models.PurchaseStatusFailed.ConsumesStock())
	assert.False(t, models.PurchaseStatusCancelled.ConsumesStock())
	assert.False(t, models.PurchaseStatus("").ConsumesStock())
}

func TestInsufficientStockError_WrapsThroughTransaction(t *testing.T) {
	// The stock error travels out of the transaction callback wrapped, which is why
	// the handler matches with errors.Is rather than by identity.
	inner := insufficientStockError{ProductName: "Mug", Available: 1}
	wrapped := fmt.Errorf("failed to store purchase: %w", error(inner))

	assert.True(t, errors.Is(wrapped, ErrInsufficientStock))
}
