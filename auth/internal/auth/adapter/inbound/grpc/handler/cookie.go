package grpc

import "strings"

func extractRefreshToken(cookieHeader string) string {

	if cookieHeader == "" {
		return ""
	}

	cookies := strings.Split(cookieHeader, ";")

	for _, cookie := range cookies {

		parts := strings.SplitN(cookie, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "refresh_token" {
			return value
		}
	}

	return ""
}