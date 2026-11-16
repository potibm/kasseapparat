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

// hasManagedStockFlags reports whether any product carries a flag the behaviour set
// itself, which is the only reason a non-restricting mode still has work to do.
func hasManagedStockFlags(products []models.Product) bool {
	for _, product := range products {
		if product.StockFlagsManaged {
			return true
		}
	}

	return false
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

		// Short-circuit: under a non-restricting behaviour there is nothing to derive
		// unless a previous mode left a managed flag behind.
		if !restrictsStock(s.OutOfStockBehavior) && !hasManagedStockFlags(products) {
			return nil
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

// reconcileStockFlag updates the flag the behaviour owns for one product and leaves
// the other one untouched. A product an operator flagged by hand keeps that flag,
// because the behaviour only manages the column its mode owns.
func (s *PurchaseService) reconcileStockFlag(
	repo sqlite.RepositoryInterface,
	product models.Product,
	unitsSold int,
) error {
	unlimited := product.TotalStock <= 0
	depleted := !unlimited && unitsSold >= product.TotalStock

	soldOut := product.SoldOut
	hidden := product.Hidden
	managed := product.StockFlagsManaged

	// A product an operator flagged by hand keeps its flags: the behaviour only
	// manages products whose flags it set itself.
	if !managed && (soldOut || hidden) {
		return nil
	}

	switch s.OutOfStockBehavior {
	case config.OutOfStockAutoSoldOut:
		soldOut, hidden, managed = s.applyAutoSoldOut(unlimited, depleted, managed, soldOut, hidden)
	case config.OutOfStockAutoHide:
		soldOut, hidden, managed = s.applyAutoHide(unlimited, depleted, managed, soldOut, hidden)
	default:
		soldOut, hidden, managed = s.releaseManagedFlags(managed, soldOut, hidden)
	}

	if soldOut == product.SoldOut && hidden == product.Hidden && managed == product.StockFlagsManaged {
		return nil
	}

	if err := repo.UpdateProductStockFlags(product.ID, soldOut, hidden, managed); err != nil {
		return fmt.Errorf("failed to update stock flags on product %q: %w", product.Name, err)
	}

	return nil
}

// applyAutoSoldOut derives the sold_out flag the auto_sold_out mode owns. A limited
// product that became unlimited releases a flag the behaviour took over.
func (s *PurchaseService) applyAutoSoldOut(
	unlimited, depleted, managed, soldOut, hidden bool,
) (newSoldOut, newHidden, newManaged bool) {
	if unlimited {
		if managed {
			return false, hidden, false
		}

		return soldOut, hidden, managed
	}

	return depleted, hidden, depleted
}

// applyAutoHide derives the hidden flag the auto_hide mode owns. A limited product
// that became unlimited releases a flag the behaviour took over.
func (s *PurchaseService) applyAutoHide(
	unlimited, depleted, managed, soldOut, hidden bool,
) (newSoldOut, newHidden, newManaged bool) {
	if unlimited {
		if managed {
			return soldOut, false, false
		}

		return soldOut, hidden, managed
	}

	return soldOut, depleted, depleted
}

// releaseManagedFlags drops any flag a previous mode took over, which is what makes
// a mode switch not strand products. A hand-set flag is left untouched.
func (s *PurchaseService) releaseManagedFlags(
	managed, soldOut, hidden bool,
) (newSoldOut, newHidden, newManaged bool) {
	if managed {
		return false, false, false
	}

	return soldOut, hidden, managed
}
