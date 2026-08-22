package sqlite

import (
	"errors"

	"github.com/potibm/kasseapparat/internal/app/models"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"github.com/shopspring/decimal"
)

func (repo *Repository) GetPaymentMethodStats() ([]response.PaymentMethodStats, error) {
	stats := []response.PaymentMethodStats{}

	query := repo.db.Table("purchases").
		Select("payment_method as id, payment_method, COUNT(*) as purchase_count, COALESCE(SUM(total_net_price), 0) as total_net_price, COALESCE(SUM(total_gross_price), 0) as total_gross_price").
		Where("deleted_at IS NULL").
		Where("status = ?", string(models.PurchaseStatusConfirmed)).
		Group("payment_method")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, errors.New("unable to retrieve payment method stats")
	}

	for i := range stats {
		if stats[i].TotalNetPrice.IsZero() {
			stats[i].TotalNetPrice = decimal.NewFromFloat(0)
		}

		if stats[i].TotalGrossPrice.IsZero() {
			stats[i].TotalGrossPrice = decimal.NewFromFloat(0)
		}
	}

	return stats, nil
}
