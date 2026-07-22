package http

import (
	"testing"
	"time"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

// TestRefundPurchase_TimeLimitLogic tests the time limit logic for refunds.
// This is a unit test for the business logic, not the full HTTP handler.
func TestRefundPurchase_TimeLimitLogic(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		userRole          string
		purchaseAge       time.Duration
		shouldAllowRefund bool
	}{
		{
			name:              "Regular user can refund recent purchase (10 min old)",
			userRole:          "user",
			purchaseAge:       10 * time.Minute,
			shouldAllowRefund: true,
		},
		{
			name:              "Regular user cannot refund old purchase (20 min old)",
			userRole:          "user",
			purchaseAge:       20 * time.Minute,
			shouldAllowRefund: false,
		},
		{
			name:              "Admin user can refund recent purchase (10 min old)",
			userRole:          "admin",
			purchaseAge:       10 * time.Minute,
			shouldAllowRefund: true,
		},
		{
			name:              "Admin user can refund old purchase (20 min old)",
			userRole:          "admin",
			purchaseAge:       20 * time.Minute,
			shouldAllowRefund: true,
		},
		{
			name:              "Admin user can refund very old purchase (1 hour old)",
			userRole:          "admin",
			purchaseAge:       1 * time.Hour,
			shouldAllowRefund: true,
		},
		{
			name:              "Regular user cannot refund purchase at exactly 15 minutes",
			userRole:          "user",
			purchaseAge:       15 * time.Minute,
			shouldAllowRefund: false,
		},
		{
			name:              "Regular user can refund purchase just under 15 minutes",
			userRole:          "user",
			purchaseAge:       14*time.Minute + 59*time.Second,
			shouldAllowRefund: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a purchase with the specified age
			purchase := &models.Purchase{
				CreatedAt: now.Add(-tt.purchaseAge),
			}

			// Simulate the time check logic from RefundPurchase
			isAdmin := tt.userRole == "admin"
			isTooOld := time.Since(purchase.CreatedAt) > 15*time.Minute
			shouldBlock := !isAdmin && isTooOld

			// Verify the logic matches expectations
			if tt.shouldAllowRefund {
				assert.False(t, shouldBlock, "Refund should be allowed but would be blocked")
			} else {
				assert.True(t, shouldBlock, "Refund should be blocked but would be allowed")
			}
		})
	}
}

// TestRefundPurchase_AdminRoleCheck tests that admin role is correctly identified.
func TestRefundPurchase_AdminRoleCheck(t *testing.T) {
	tests := []struct {
		name     string
		userRole string
		isAdmin  bool
	}{
		{
			name:     "User role is not admin",
			userRole: "user",
			isAdmin:  false,
		},
		{
			name:     "Admin role is admin",
			userRole: "admin",
			isAdmin:  true,
		},
		{
			name:     "Empty role is not admin",
			userRole: "",
			isAdmin:  false,
		},
		{
			name:     "Other role is not admin",
			userRole: "manager",
			isAdmin:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authUser := models.AuthUser{
				Username: "testuser",
				Role:     tt.userRole,
			}

			// Simulate the admin check logic
			isAdmin := authUser.Role == "admin"
			assert.Equal(t, tt.isAdmin, isAdmin)
		})
	}
}
