package utils

import (
	"strings"
	"unicode"
)

func CapitalizeFirstRune(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}
