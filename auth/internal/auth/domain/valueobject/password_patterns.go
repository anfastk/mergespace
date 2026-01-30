package valueobject

import "strings"

func containsSequentialPatterns(p string) bool {
	p = strings.ToLower(p)

	sequences := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"0123456789",
		"qwertyuiop",
	}

	for _, seq := range sequences {
		for i := 0; i+3 < len(seq); i++ {
			if strings.Contains(p, seq[i:i+4]) {
				return true
			}
		}
	}
	return false
}
