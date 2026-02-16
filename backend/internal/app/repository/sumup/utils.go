package sumup

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	sumup "github.com/sumup/sumup-go"
)

func getStringPtr(v url.Values, key string) *string {
	if s := v.Get(key); s != "" {
		return &s
	}

	return nil
}

func getStringSlicePtr(v url.Values, key string) *[]string {
	if list, ok := v[key]; ok && len(list) > 0 {
		return &list
	}

	return nil
}

func getStringSlice(v url.Values, key string) []string {
	if list, ok := v[key]; ok && len(list) > 0 {
		return list
	}

	return []string{}
}

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

	if apiErr, ok := err.(*sumup.Error); ok {
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
