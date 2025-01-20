package tests_e2e

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	productStatsUrl = "/api/v1/productStats"
)

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
	item.Value("totalPrice").String().IsEqual("40")

	deletePurchase(purchaseUrl)
}

func validateProductStatsArray(productStatsArray *httpexpect.Array) {
	productStatsArray.Length().Gt(0)

	for i := 0; i < len(productStatsArray.Iter()); i++ {
		productStats := productStatsArray.Value(i).Object()
		validateProductStatsObject(productStats)
	}
}

func validateProductStatsObject(productStats *httpexpect.Object) {
	productStats.Value("id").Number().Gt(0)
	productStats.Value("name").String().NotEmpty()
	productStats.Value("soldItems").Number().Ge(0)
	productStats.Value("totalPrice").String().NotEmpty()
}

func TestProductStatsAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", productStatsUrl).Expect().Status(http.StatusUnauthorized)

}
