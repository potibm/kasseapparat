package tests_e2e

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var purchaseExportBaseUrl = "/api/v2/purchases/export"

func TestPurchaseExportForEntityEndpoints(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", purchaseExportBaseUrl).Expect().Status(http.StatusUnauthorized)
}

func TestGetPurchaseExport(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(purchaseExportBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	res.Header("Content-Type").IsEqual("text/csv")
	res.Header("Content-Disposition").Contains("attachment; filename=\"purchases_")
	res.Header("Content-Disposition").Contains(".csv\"")

	raw := res.Body().Raw()

	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		if i == 0 || len(line) == 0 {
			continue
		}

		columns := strings.Split(line, ",")

		// assert that the number of columns is correct
		if len(columns) != 15 {
			t.Errorf("Expected 15 columns, got %d in line %d", len(columns), i)
		}

		// assert that the first column is a valid date
		if _, err := time.Parse("2006-01-02 15:04:05", columns[0]); err != nil {
			t.Errorf("Expected a valid date in the first column, got %s in line %d", columns[0], i)
		}

		// validate that column 2 is a valid uuid
		if _, err := uuid.Parse(columns[1]); err != nil {
			t.Errorf("Expected a valid integer in column 2 (id), got %s in line %d", columns[1], i)
		}

		// validate that column 3 is a valid integer
		if _, err := strconv.Atoi(columns[2]); err != nil {
			t.Errorf("Expected a valid integer in column 3 (quantity), got %s in line %d", columns[2], i)
		}

		// decimal columns should be in the format "0.00"
		for j := 5; j <= 13; j++ {
			// validate that the column is a valid decimal.Decimal value
			if _, err := decimal.NewFromString(columns[j]); err != nil {
				t.Errorf("Expected a valid decimal value in column %d, got %s in line %d", j, columns[j], i)
			}
		}
	}
}

func TestGetPurchaseExportFilterOnCASH(t *testing.T) {
	testGetPurchaseExportFilterOnPaymentMethod(t, "CASH")
}

func TestGetPurchaseExportFilterOnCC(t *testing.T) {
	testGetPurchaseExportFilterOnPaymentMethod(t, "CC")
}

func testGetPurchaseExportFilterOnPaymentMethod(t *testing.T, paymentMethod string) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(purchaseExportBaseUrl)).
		WithQuery("paymentMethods", paymentMethod).
		Expect()

	res.Status(http.StatusOK)

	raw := res.Body().Raw()

	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		if i == 0 || len(line) == 0 {
			continue
		}

		columns := strings.Split(line, ",")

		// assert that the number of columns is correct
		if len(columns) != 15 {
			t.Fatalf("Expected 15 columns, got %d in line %d", len(columns), i)
		}

		paymentMethodInCSV := columns[14]
		if paymentMethodInCSV != paymentMethod {
			t.Fatalf("Expected payment method to be %s, got %s in line %d", paymentMethod, paymentMethodInCSV, i)
		}
	}
}
