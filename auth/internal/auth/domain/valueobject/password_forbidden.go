package valueobject

import "strings"

func containsForbiddenSubstring(p string, forbidden []string) bool {
	p = strings.ToLower(p)

	for _, f := range forbidden {
		f = strings.ToLower(strings.TrimSpace(f))
		if len(f) < 3 {
			continue
		}
		if strings.Contains(p, f) {
			return true
		}
	}

	return false
}
