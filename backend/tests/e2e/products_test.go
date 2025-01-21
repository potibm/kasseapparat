package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	productBaseUrl   = "/api/v2/products"
	productUrlWithId = productBaseUrl + "/1"
)

func TestGetProducts(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(productBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Gt(10)

	obj := res.JSON().Array()
	obj.Length().IsEqual(10)

	for i := 0; i < len(obj.Iter()); i++ {
		product := obj.Value(i).Object()
		validateProduct(product)
		product.Value("guestlists").Array()
	}

	product := obj.Value(0).Object()
	validateProductOne(product)
	product.Value("guestlists").Array().IsEmpty()

	productWithList := obj.Value(1).Object()
	productWithList.Value("name").String().Contains("Reduced")
	productWithList.Value("guestlists").Array().Length().IsEqual(2)
}

func TestGetProductsWithSort(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// define an array of sort fields
	sortFields := []string{"id", "name", "price", "pos"}

	for _, sortField := range sortFields {

		withDemoUserAuthToken(e.GET(productBaseUrl)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestGetProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	product := withDemoUserAuthToken(e.GET(productUrlWithId)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	validateProductOne(product)
}

func TestCreateUpdateAndDeleteProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	var originalName = "Test Product"
	var changedName = "Test Product Updated"

	product := withDemoUserAuthToken(e.POST(productBaseUrl)).
		WithJSON(map[string]interface{}{
			"name":      originalName,
			"price":     10,
			"wrapAfter": false,
			"pos":       123,
			"hidden":    false,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	product.Value("id").Number().Gt(0)
	product.Value("name").String().Contains(originalName)

	productId := product.Value("id").Number().Raw()
	productUrl := productBaseUrl + "/" + strconv.FormatFloat(productId, 'f', -1, 64)

	product = withDemoUserAuthToken(e.GET(productUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	product.Value("id").Number().Gt(0)
	product.Value("name").String().Contains(originalName)

	withDemoUserAuthToken(e.PUT(productUrl)).
		WithJSON(map[string]interface{}{
			"name":      changedName,
			"price":     10,
			"wrapAfter": false,
			"pos":       123,
			"hidden":    false,
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	product = withDemoUserAuthToken(e.GET(productUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	product.Value("id").Number().Gt(0)
	product.Value("name").String().Contains(changedName)

	withAdminUserAuthToken(e.DELETE(productUrl)).
		Expect().
		Status(http.StatusOK)

	withDemoUserAuthToken(e.GET(productUrl)).
		Expect().
		Status(http.StatusNotFound)

}

func TestDemoUserIsNotAllowedToDeleteAProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	withDemoUserAuthToken(e.DELETE(productUrlWithId)).
		Expect().
		Status(http.StatusForbidden)
}

func TestProductAuthentication(t *testing.T) {
	testAuthenticationForEntityEndpoints(t, productBaseUrl, productUrlWithId)
}

func validateProduct(product *httpexpect.Object) {
	product.Value("id").Number().Gt(0)
	product.Value("name").String().NotEmpty()
	product.Value("price").String().NotEmpty()
	product.Value("wrapAfter").Boolean()
	product.Value("pos").Number().Ge(0)
	product.Value("apiExport").Boolean()
	product.Value("totalStock").Number().Ge(0)
	product.Value("unitsSold").Number().Ge(0)
	product.Value("soldOutRequestCount").Number().Ge(0)

}

func validateProductOne(product *httpexpect.Object) {
	product.Value("id").Number().IsEqual(1)
	product.Value("name").String().Contains("Regular")
	product.Value("price").String().IsEqual("40")
	product.Value("wrapAfter").Boolean().IsFalse()
	product.Value("hidden").Boolean().IsFalse()
	product.Value("pos").Number().IsEqual(1)
	product.Value("totalStock").Number().IsEqual(0)
	product.Value("unitsSold").Number().IsEqual(0)
	product.Value("soldOutRequestCount").Number().IsEqual(0)
	product.Value("apiExport").Boolean().IsTrue()
}
