package sumup

import (
	"net/url"
	"strconv"
	"time"
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
