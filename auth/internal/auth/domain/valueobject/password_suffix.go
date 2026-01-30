package valueobject

import "unicode"

func hasWeakSuffix(p string) bool {
	r := []rune(p)
	n := len(r)

	if n < 4 {
		return false
	}

	count := 0
	for i := n - 1; i >= 0; i-- {
		if unicode.IsDigit(r[i]) || unicode.IsSymbol(r[i]) || unicode.IsPunct(r[i]) {
			count++
		} else {
			break
		}
	}

	return count >= 3
}
