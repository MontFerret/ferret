package http

import (
	"net"
	stdhttp "net/http"
	"net/url"
	"strings"
)

func normalizeValues(values []string) []string {
	if values == nil {
		return nil
	}

	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}

		normalized = append(normalized, value)
	}

	return normalized
}

func normalizeHeaders(headers []string) []string {
	if headers == nil {
		return nil
	}

	normalized := make([]string, 0, len(headers))
	for _, header := range headers {
		header = strings.TrimSpace(header)
		if header == "" {
			continue
		}

		normalized = append(normalized, stdhttp.CanonicalHeaderKey(header))
	}

	return normalized
}

func containsValue(values []string, needle string) bool {
	needle = strings.ToLower(strings.TrimSpace(needle))
	for _, value := range values {
		if strings.ToLower(strings.TrimSpace(value)) == needle {
			return true
		}
	}

	return false
}

func containsHost(hosts []string, u *url.URL) bool {
	rawHost := strings.ToLower(strings.TrimSpace(u.Host))
	hostname := strings.ToLower(strings.TrimSpace(u.Hostname()))

	for _, host := range hosts {
		if host == rawHost || host == hostname {
			return true
		}
	}

	return false
}

func isLocalhost(hostname string, ip net.IP) bool {
	return hostname == "localhost" || (ip != nil && ip.IsLoopback())
}
