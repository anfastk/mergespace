package valueobject

import "strings"

var commonPasswords = map[string]struct{}{
	"password":    {},
	"password123": {},
	"1234567890":   {},
}

func isCommonPassword(p string) bool {
	_, found := commonPasswords[strings.ToLower(p)]
	return found
}
