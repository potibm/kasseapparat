package http

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestQueryPurchaseStatusList(t *testing.T) {
	testCases := []struct {
		name           string
		queryValues    []string
		expectedStatus *models.PurchaseStatusList
	}{
		{
			name:           "no status provided",
			queryValues:    []string{},
			expectedStatus: nil,
		},
		{
			name:           "single valid status",
			queryValues:    []string{"pending"},
			expectedStatus: &models.PurchaseStatusList{models.PurchaseStatusPending},
		},
		{
			name:           "multiple valid statuses",
			queryValues:    []string{"pending", "confirmed"},
			expectedStatus: &models.PurchaseStatusList{models.PurchaseStatusPending, models.PurchaseStatusConfirmed},
		},
		{
			name:           "invalid status is ignored",
			queryValues:    []string{"pending", "invalid"},
			expectedStatus: &models.PurchaseStatusList{models.PurchaseStatusPending},
		},
		{
			name:           "all invalid statuses",
			queryValues:    []string{"invalid1", "invalid2"},
			expectedStatus: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &gin.Context{}
			c.Request = &http.Request{URL: &url.URL{}}

			q := c.Request.URL.Query()
			for _, value := range tc.queryValues {
				q.Add("status", value)
			}

			c.Request.URL.RawQuery = q.Encode()

			result := queryPurchaseStatusList(c, "status")

			assert.Equal(t, tc.expectedStatus, result)
		})
	}
}
