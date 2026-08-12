package initializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildWebhookURL(t *testing.T) {
	tests := []struct {
		name      string
		publicURL string
		wantURL   *string
	}{
		{
			name:      "empty URL returns nil",
			publicURL: "",
			wantURL:   nil,
		},
		{
			name:      "URL without trailing slash",
			publicURL: "https://example.com",
			wantURL:   strPtr("https://example.com/api/" + APIVersion + "/sumup/webhook"),
		},
		{
			name:      "URL with trailing slash",
			publicURL: "https://example.com/",
			wantURL:   strPtr("https://example.com/api/" + APIVersion + "/sumup/webhook"),
		},
		{
			name:      "URL with path",
			publicURL: "https://example.com/app",
			wantURL:   strPtr("https://example.com/app/api/" + APIVersion + "/sumup/webhook"),
		},
		{
			name:      "URL with path and trailing slash",
			publicURL: "https://example.com/app/",
			wantURL:   strPtr("https://example.com/app/api/" + APIVersion + "/sumup/webhook"),
		},
		{
			name:      "URL with port",
			publicURL: "https://localhost:8443",
			wantURL:   strPtr("https://localhost:8443/api/" + APIVersion + "/sumup/webhook"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildWebhookURL(tt.publicURL)

			if tt.wantURL == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, *tt.wantURL, *got)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
