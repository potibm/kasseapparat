package tests_e2e

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var productStatsUrl = "/api/v2/productStats"

func TestProductStats(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(productStatsUrl)).
		Expect().
		Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Ge(1)

	obj := res.JSON().Array()
	validateProductStatsArray(obj)

	purchaseUrl := createPurchase()

	res = withDemoUserAuthToken(e.GET(productStatsUrl)).
		Expect().
		Status(http.StatusOK)

	obj = res.JSON().Array()
	validateProductStatsArray(obj)

	item := obj.Value(0).Object()
	item.Value("soldItems").Number().IsEqual(1)
	item.Value("totalGrossPrice").String().IsEqual("40")
	item.Value("totalNetPrice").String().IsEqual("37.38")

	deletePurchase(purchaseUrl)
}

func validateProductStatsArray(productStatsArray *httpexpect.Array) {
	productStatsArray.Length().Gt(0)

	for i := range len(productStatsArray.Iter()) {
		productStats := productStatsArray.Value(i).Object()
		validateProductStatsObject(productStats)
	}
}

func validateProductStatsObject(productStats *httpexpect.Object) {
	productStats.Value("id").Number().Gt(0)
	productStats.Value("name").String().NotEmpty()
	productStats.Value("soldItems").Number().Ge(0)
	productStats.Value("totalGrossPrice").String().NotEmpty()
	productStats.Value("totalNetPrice").String().NotEmpty()
}

func TestProductStatsAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", productStatsUrl).Expect().Status(http.StatusUnauthorized)
}
