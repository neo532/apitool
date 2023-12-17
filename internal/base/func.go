package base

import "strings"

func UpperFirstChar(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
