package utils

import "testing"

func TestCapitalizeFirstRune(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "Hello"},
		{" hello", "Hello"},
		{"äpfel", "Äpfel"},
		{"1 dog", "1 dog"},
	}

	for _, tt := range tests {
		result := CapitalizeFirstRune(tt.input)
		if result != tt.expected {
			t.Errorf("expected %q → %q, got %q", tt.input, tt.expected, result)
		}
	}
}
