package utils

import "strings"

func BearerToken(header string) (token string, ok bool) {
	return strings.CutPrefix(header, "Bearer ")
}
