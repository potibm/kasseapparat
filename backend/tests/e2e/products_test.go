package tests_e2e

import (
	"net/http"
	"testing"
)

func TestGetProducts(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := e.GET("/api/v1/products").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeader := res.Header("X-Total-Count").AsNumber()
	totalCountHeader.Gt(10)

	obj := res.JSON().Array()
	obj.Length().IsEqual(10)

	for i := 0; i < len(obj.Iter()); i++ {
		product := obj.Value(i).Object()
		product.Value("id").Number().Gt(0)
		product.Value("name").String().NotEmpty()
		product.Value("price").Number().Ge(0)
		product.Value("wrapAfter").Boolean()
		product.Value("pos").Number().Ge(0)
		product.Value("apiExport").Boolean()
		product.Value("totalStock").Number().Ge(0)
		product.Value("unitsSold").Number().Ge(0)
		product.Value("soldOutRequestCount").Number().Ge(0)
		product.Value("lists").Array()
	}

	product := obj.Value(0).Object()
	product.Value("id").Number().IsEqual(1)
	product.Value("name").String().Contains("Regular")
	product.Value("price").Number().IsEqual(40)
	product.Value("wrapAfter").Boolean().IsFalse()
	product.Value("hidden").Boolean().IsFalse()
	product.Value("pos").Number().IsEqual(1)
	product.Value("totalStock").Number().IsEqual(0)
	product.Value("unitsSold").Number().IsEqual(0)
	product.Value("soldOutRequestCount").Number().IsEqual(0)
	product.Value("apiExport").Boolean().IsTrue()
	product.Value("lists").Array().IsEmpty()

	productWithList := obj.Value(1).Object()
	productWithList.Value("name").String().Contains("Reduced")
	productWithList.Value("lists").Array().Length().IsEqual(2)
}

func TestGetProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	product := e.GET("/api/v1/products/1").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	product.Value("id").Number().IsEqual(1)
	product.Value("name").String().Contains("Regular")
	product.Value("price").Number().IsEqual(40)
	product.Value("wrapAfter").Boolean().IsFalse()
	product.Value("hidden").Boolean().IsFalse()
	product.Value("pos").Number().IsEqual(1)
	product.Value("totalStock").Number().IsEqual(0)
	product.Value("unitsSold").Number().IsEqual(0)
	product.Value("soldOutRequestCount").Number().IsEqual(0)
	product.Value("apiExport").Boolean().IsTrue()
	product.Value("lists").IsNull()
}
