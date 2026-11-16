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

		// Short-circuit: under a non-restricting behaviour there is nothing to derive
		// unless a previous mode left a managed flag behind.
		if !restrictsStock(s.OutOfStockBehavior) {
			hasManaged := false
			for _, product := range products {
				if product.StockFlagsManaged {
					hasManaged = true
					break
				}
			}
			if !hasManaged {
				return nil
			}
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

	switch s.OutOfStockBehavior {
	case config.OutOfStockAutoSoldOut:
		if unlimited {
			// A limited product that became unlimited releases a flag the
			// behaviour took over.
			if managed {
				soldOut = false
				managed = false
			}
		} else {
			// sold_out is system-owned for products with a stock limit.
			soldOut = depleted
			managed = depleted
		}
	case config.OutOfStockAutoHide:
		if unlimited {
			if managed {
				hidden = false
				managed = false
			}
		} else {
			// hidden is system-owned for products with a stock limit.
			hidden = depleted
			managed = depleted
		}
	default:
		// ignore / fail never own a flag. A flag a previous mode took over is
		// released, which is what makes a mode switch not strand products.
		if managed {
			soldOut = false
			hidden = false
			managed = false
		}
	}

	if soldOut == product.SoldOut && hidden == product.Hidden && managed == product.StockFlagsManaged {
		return nil
	}

	if err := repo.UpdateProductStockFlags(product.ID, soldOut, hidden, managed); err != nil {
		return fmt.Errorf("failed to update stock flags on product %q: %w", product.Name, err)
	}

	return nil
}
