package sqlite

import (
	"errors"
	"strconv"

	"github.com/potibm/kasseapparat/internal/app/models"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"github.com/shopspring/decimal"
)

func (repo *Repository) GetHourlyRevenueStats() ([]response.HourlyRevenueStats, error) {
	stats := []response.HourlyRevenueStats{}

	query := repo.db.Table("purchases").
		Select(`strftime('%Y-%m-%d %H:', created_at) || printf('%02d', (strftime('%M', created_at) / 15) * 15) as time_bucket,
			payment_method,
			COALESCE(SUM(total_gross_price), 0) as total_gross_price`).
		Where("deleted_at IS NULL").
		Where("status = ?", string(models.PurchaseStatusConfirmed)).
		Where("created_at >= datetime('now', '-3 days')").
		Group("time_bucket, payment_method").
		Order("time_bucket ASC")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, errors.New("unable to retrieve hourly revenue stats")
	}

	for i := range stats {
		stats[i].ID = stats[i].TimeBucket + "_" + stats[i].PaymentMethod

		if stats[i].TotalGrossPrice.IsZero() {
			stats[i].TotalGrossPrice = decimal.NewFromFloat(0)
		}
	}

	return stats, nil
}

func (repo *Repository) GetHourlyQuantityStats() ([]response.HourlyQuantityStats, error) {
	stats := []response.HourlyQuantityStats{}

	query := repo.db.Table("purchase_items").
		Select(`strftime('%Y-%m-%d %H:', purchases.created_at) ||
			printf('%02d', (strftime('%M', purchases.created_at) / 15) * 15) as time_bucket,
			purchase_items.product_id,
			COALESCE(SUM(purchase_items.quantity), 0) as quantity`).
		Joins("JOIN purchases ON purchases.id = purchase_items.purchase_id").
		Where("purchases.deleted_at IS NULL").
		Where("purchases.status = ?", string(models.PurchaseStatusConfirmed)).
		Where("purchases.created_at >= datetime('now', '-3 days')").
		Where("purchase_items.deleted_at IS NULL").
		Group("time_bucket, purchase_items.product_id").
		Order("time_bucket ASC")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, errors.New("unable to retrieve hourly quantity stats")
	}

	for i := range stats {
		stats[i].ID = stats[i].TimeBucket + "_" + strconv.Itoa(stats[i].ProductID)
	}

	return stats, nil
}
