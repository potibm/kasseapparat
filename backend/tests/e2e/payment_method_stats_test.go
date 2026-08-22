package tests_e2e

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var paymentMethodStatsURL = "/api/v3/paymentMethodStats"

func TestPaymentMethodStats(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(paymentMethodStatsURL)).
		Expect().
		Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Ge(0)

	obj := res.JSON().Array()
	validatePaymentMethodStatsArray(obj)

	purchaseURL := createPurchase()

	res = withDemoUserAuthToken(e.GET(paymentMethodStatsURL)).
		Expect().
		Status(http.StatusOK)

	obj = res.JSON().Array()
	validatePaymentMethodStatsArray(obj)

	found := false

	for _, item := range obj.Iter() {
		paymentMethod := item.Object().Value("paymentMethod").String().Raw()
		if paymentMethod == "CASH" {
			found = true
			item.Object().Value("purchaseCount").Number().Ge(1)
			item.Object().Value("totalGrossPrice").String().NotEmpty()
			item.Object().Value("totalNetPrice").String().NotEmpty()
			item.Object().Value("name").String().NotEmpty()
		}
	}

	if !found {
		t.Error("Expected to find CASH payment method in stats")
	}

	deletePurchase(purchaseURL)
}

func validatePaymentMethodStatsArray(statsArray *httpexpect.Array) {
	for i := range len(statsArray.Iter()) {
		stats := statsArray.Value(i).Object()
		validatePaymentMethodStatsObject(stats)
	}
}

func validatePaymentMethodStatsObject(stats *httpexpect.Object) {
	stats.Value("paymentMethod").String().NotEmpty()
	stats.Value("name").String().NotEmpty()
	stats.Value("purchaseCount").Number().Ge(0)
	stats.Value("totalGrossPrice").String().NotEmpty()
	stats.Value("totalNetPrice").String().NotEmpty()
}
