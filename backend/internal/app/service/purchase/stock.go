package purchase

import (
	"context"
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

// AvailableStock reports how many units of a product can still be sold, and whether
// the product is unlimited. A product with no total stock is always unlimited and is
// excluded from every restriction, which is what makes stock-less products such as
// drinks work without configuration.
func AvailableStock(product models.Product, unitsSold int) (available int, unlimited bool) {
	if product.TotalStock <= 0 {
		return 0, true
	}

	available = product.TotalStock - unitsSold
	if available < 0 {
		available = 0
	}

	return available, false
}

// restrictsStock reports whether the behaviour stops a sale that exceeds the stock.
// An unset or unrecognised value falls back to ignoring, which is both the documented
// default and the behaviour from before the setting existed, so a missing config value
// can never start rejecting sales.
func restrictsStock(behavior config.OutOfStockBehavior) bool {
	switch behavior {
	case config.OutOfStockFail, config.OutOfStockAutoSoldOut, config.OutOfStockAutoHide:
		return true
	default:
		return false
	}
}

// ReconcileStockFlags brings the sold-out and hidden flags of every product in line
// with the stock that is actually left. It is the only place that maintains those
// flags, and it derives them from the purchases instead of remembering what an earlier
// call did, so it is idempotent and also repairs what an incremental update never saw:
// a raised total stock, a deleted purchase, a changed mode, or a restart.
//
// A product an operator flagged by hand is left alone entirely: the behaviour only
// manages products whose flags it set itself.
func (s *PurchaseService) ReconcileStockFlags(ctx context.Context) error {
	return s.sqliteRepo.WithTransaction(ctx, func(repo sqlite.RepositoryInterface) error {
		products, err := repo.GetProductsForStockReconciliation()
		if err != nil {
			return fmt.Errorf("failed to load products for stock reconciliation: %w", err)
		}

		unitsSoldByProduct, err := repo.GetUnitsSoldByProductID()
		if err != nil {
			return fmt.Errorf("failed to load sold quantities: %w", err)
		}

		for _, product := range products {
			if err := s.reconcileStockFlag(repo, product, unitsSoldByProduct[product.ID]); err != nil {
				return err
			}
		}

		return nil
	})
}

// reconcileStockFlag updates one product to the state the behaviour asks for, and
// leaves it untouched otherwise. A product that already carries that state is not
// written at all, which keeps a reconciliation that has nothing to do free of writes.
func (s *PurchaseService) reconcileStockFlag(
	repo sqlite.RepositoryInterface,
	product models.Product,
	unitsSold int,
) error {
	// A product without a stock limit is never restricted, so it never takes a flag. A
	// flag the behaviour took over is still released, which is what happens when an
	// operator makes a limited product unlimited again.
	depleted := product.TotalStock > 0 && unitsSold >= product.TotalStock

	soldOut := depleted && s.OutOfStockBehavior == config.OutOfStockAutoSoldOut
	hidden := depleted && s.OutOfStockBehavior == config.OutOfStockAutoHide

	if !product.StockFlagsManaged && (product.SoldOut || product.Hidden) {
		// An operator flagged this product by hand, so its flags are not ours to
		// change, neither to set nor to clear.
		return nil
	}

	managed := soldOut || hidden

	if soldOut == product.SoldOut && hidden == product.Hidden && managed == product.StockFlagsManaged {
		return nil
	}

	if err := repo.UpdateProductStockFlags(product.ID, soldOut, hidden, managed); err != nil {
		return fmt.Errorf("failed to update stock flags on product %q: %w", product.Name, err)
	}

	return nil
}
