package sumup

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sumup/sumup-go/shared"
)

func getStringPtr(v url.Values, key string) *string {
	if s := v.Get(key); s != "" {
		return &s
	}

	return nil
}

// getStringSlicePtr returns a pointer to the slice of strings associated with key in v if the slice exists and has at least one element.
// It returns nil if the key is absent or the slice is empty.
func getStringSlicePtr(v url.Values, key string) *[]string {
	if list, ok := v[key]; ok && len(list) > 0 {
		return &list
	}

	return nil
}

// getStringSlice returns the slice of strings associated with the given key from v.
// If the key is absent or the slice is empty, it returns a non-nil empty slice.
func getStringSlice(v url.Values, key string) []string {
	if list, ok := v[key]; ok && len(list) > 0 {
		return list
	}

	return []string{}
}

// getIntPtr returns a pointer to the integer parsed from v for the given key.
// If the key is missing, the value is empty, or parsing fails, it returns nil.
func getIntPtr(v url.Values, key string) *int {
	if s := v.Get(key); s != "" {
		if i, err := strconv.Atoi(s); err == nil {
			return &i
		}
	}

	return nil
}

func getTimePtr(v url.Values, key string) *time.Time {
	if s := v.Get(key); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return &t
		}
	}

	return nil
}

func normalizeSumupError(err error) error {
	if err == nil {
		return nil
	}

	if apiErr, ok := err.(*shared.Error); ok {
		var code string

		var message string

		if apiErr.ErrorCode != nil {
			code = *apiErr.ErrorCode
		}

		if apiErr.Message != nil {
			message = *apiErr.Message
		}

		return fmt.Errorf("SumUp error %s: %s", code, message)
	}

	return fmt.Errorf("unexpected error: %v", err)
}