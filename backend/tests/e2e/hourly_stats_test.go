package tests_e2e

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	hourlyRevenueStatsURL  = "/api/v3/hourlyRevenueStats"
	hourlyQuantityStatsURL = "/api/v3/hourlyQuantityStats"
)

func TestHourlyRevenueStats(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(hourlyRevenueStatsURL)).
		Expect().
		Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Ge(0)

	obj := res.JSON().Array()
	validateHourlyRevenueStatsArray(obj)

	purchaseURL := createPurchase()

	res = withDemoUserAuthToken(e.GET(hourlyRevenueStatsURL)).
		Expect().
		Status(http.StatusOK)

	obj = res.JSON().Array()
	validateHourlyRevenueStatsArray(obj)

	if obj.Length().Raw() > 0 {
		found := false

		for _, item := range obj.Iter() {
			paymentMethod := item.Object().Value("paymentMethod").String().Raw()
			if paymentMethod == "CASH" {
				found = true

				item.Object().Value("totalGrossPrice").String().NotEmpty()
				item.Object().Value("name").String().NotEmpty()
				item.Object().Value("timeBucket").String().NotEmpty()
				item.Object().Value("id").String().NotEmpty()
			}
		}

		if !found {
			t.Error("Expected to find CASH payment method in hourly revenue stats")
		}
	}

	deletePurchase(purchaseURL)
}

func TestHourlyQuantityStats(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(hourlyQuantityStatsURL)).
		Expect().
		Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Ge(0)

	obj := res.JSON().Array()
	validateHourlyQuantityStatsArray(obj)

	purchaseURL := createPurchase()

	res = withDemoUserAuthToken(e.GET(hourlyQuantityStatsURL)).
		Expect().
		Status(http.StatusOK)

	obj = res.JSON().Array()
	validateHourlyQuantityStatsArray(obj)

	if obj.Length().Raw() > 0 {
		item := obj.Value(0).Object()
		item.Value("quantity").Number().Ge(0)
		item.Value("productName").String().NotEmpty()
		item.Value("timeBucket").String().NotEmpty()
		item.Value("id").String().NotEmpty()
	}

	deletePurchase(purchaseURL)
}

func validateHourlyRevenueStatsArray(statsArray *httpexpect.Array) {
	for i := range len(statsArray.Iter()) {
		stats := statsArray.Value(i).Object()
		validateHourlyRevenueStatsObject(stats)
	}
}

func validateHourlyRevenueStatsObject(stats *httpexpect.Object) {
	stats.Value("id").String().NotEmpty()
	stats.Value("timeBucket").String().NotEmpty()
	stats.Value("paymentMethod").String().NotEmpty()
	stats.Value("name").String().NotEmpty()
	stats.Value("totalGrossPrice").String().NotEmpty()
}

func validateHourlyQuantityStatsArray(statsArray *httpexpect.Array) {
	for i := range len(statsArray.Iter()) {
		stats := statsArray.Value(i).Object()
		validateHourlyQuantityStatsObject(stats)
	}
}

func validateHourlyQuantityStatsObject(stats *httpexpect.Object) {
	stats.Value("id").String().NotEmpty()
	stats.Value("timeBucket").String().NotEmpty()
	stats.Value("productId").Number()
	stats.Value("productName").String().NotEmpty()
	stats.Value("quantity").Number().Ge(0)
}
