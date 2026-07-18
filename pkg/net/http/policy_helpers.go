package http

import (
	"fmt"
	"net"
	stdhttp "net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
)

var (
	ipv4CurrentNetwork  = netip.MustParsePrefix("0.0.0.0/8")
	carrierGradeNAT     = netip.MustParsePrefix("100.64.0.0/10")
	ipv4Reserved        = netip.MustParsePrefix("240.0.0.0/4")
	wellKnownNAT64      = netip.MustParsePrefix("64:ff9b::/96")
	ipv6AllocatedGlobal = netip.MustParsePrefix("2000::/3")
	ipv6SiteLocal       = netip.MustParsePrefix("fec0::/10")
	// Keep the public exceptions and non-public ranges aligned with the IANA
	// IPv4 and IPv6 Special-Purpose Address Registries.
	ipv4PublicExceptions = []netip.Prefix{
		netip.MustParsePrefix("192.0.0.9/32"),
		netip.MustParsePrefix("192.0.0.10/32"),
	}
	ipv4NonPublic = []netip.Prefix{
		netip.MustParsePrefix("192.0.0.0/24"),
		netip.MustParsePrefix("192.0.2.0/24"),
		netip.MustParsePrefix("192.88.99.0/24"),
		netip.MustParsePrefix("198.18.0.0/15"),
		netip.MustParsePrefix("198.51.100.0/24"),
		netip.MustParsePrefix("203.0.113.0/24"),
	}
	ipv6NonPublic = []netip.Prefix{
		netip.MustParsePrefix("2001::/32"),
		netip.MustParsePrefix("2001:2::/48"),
		netip.MustParsePrefix("2001:10::/28"),
		netip.MustParsePrefix("2001:db8::/32"),
		netip.MustParsePrefix("2002::/16"),
		netip.MustParsePrefix("3fff::/20"),
	}
)

func normalizeValues(values []string) []string {
	if values == nil {
		return nil
	}

	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, value := range values {
		value = asciiLower(strings.TrimSpace(value))

		if _, exists := seen[value]; exists {
			continue
		}

		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}

	return normalized
}

func normalizeMethods(methods []string) []string {
	if methods == nil {
		return nil
	}

	normalized := make([]string, 0, len(methods))
	seen := make(map[string]struct{}, len(methods))

	for _, method := range methods {
		method = normalizeMethod(method)

		if _, exists := seen[method]; exists {
			continue
		}

		seen[method] = struct{}{}
		normalized = append(normalized, method)
	}

	return normalized
}

func normalizeRequestMethod(method string) string {
	method = normalizeMethod(method)

	if method == "" {
		return stdhttp.MethodGet
	}

	return method
}

func normalizeMethod(method string) string {
	return asciiUpper(strings.TrimSpace(method))
}

func normalizeHeaders(headers []string) []string {
	if headers == nil {
		return nil
	}

	normalized := make([]string, 0, len(headers))
	seen := make(map[string]struct{}, len(headers))

	for _, header := range headers {
		header = strings.TrimSpace(header)
		header = stdhttp.CanonicalHeaderKey(header)

		if _, exists := seen[header]; exists {
			continue
		}

		seen[header] = struct{}{}
		normalized = append(normalized, header)
	}

	return normalized
}

func isReservedRequestHeader(key string) bool {
	switch stdhttp.CanonicalHeaderKey(key) {
	case "Connection",
		"Content-Length",
		"Host",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Proxy-Connection",
		"Te",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade":
		return true
	default:
		return false
	}
}

func normalizeHosts(hosts []string) []string {
	if hosts == nil {
		return nil
	}

	normalized := make([]string, 0, len(hosts))
	seen := make(map[string]struct{}, len(hosts))

	for _, host := range hosts {
		host, _ = normalizeConfiguredHost(host)

		if _, exists := seen[host]; exists {
			continue
		}

		seen[host] = struct{}{}
		normalized = append(normalized, host)
	}

	return normalized
}

func validateConfiguredScheme(value string) error {
	scheme := strings.TrimSpace(value)
	if scheme == "" {
		return fmt.Errorf("must be a non-empty URL scheme")
	}

	if !isASCII(scheme) || !isASCIIAlpha(scheme[0]) {
		return fmt.Errorf("must be a valid URL scheme")
	}

	for idx := 1; idx < len(scheme); idx++ {
		char := scheme[idx]

		if !isASCIIAlpha(char) && (char < '0' || char > '9') && char != '+' && char != '-' && char != '.' {
			return fmt.Errorf("must be a valid URL scheme")
		}
	}

	return nil
}

func validateConfiguredHosts(p *Policy, option string, hosts []string) {
	for _, host := range hosts {
		if _, err := normalizeConfiguredHost(host); err != nil {
			p.setConfigurationError(option, host, err.Error())

			return
		}
	}
}

func normalizeConfiguredHost(value string) (string, error) {
	raw := strings.TrimSpace(value)

	if raw == "" {
		return "", fmt.Errorf("must not be blank")
	}

	if !isASCII(raw) {
		return "", fmt.Errorf("must use ASCII/punycode")
	}

	if strings.Contains(raw, "*") {
		return "", fmt.Errorf("wildcards are not supported")
	}

	host := raw
	port := ""
	hasPort := false

	if strings.HasPrefix(raw, "[") {
		if strings.HasSuffix(raw, "]") {
			host = raw[1 : len(raw)-1]

			if !strings.Contains(host, ":") {
				return "", fmt.Errorf("brackets require an IPv6 literal")
			}

			if _, ok := parseIPAddress(host); !ok {
				return "", fmt.Errorf("must contain a valid bracketed IP literal")
			}
		} else {
			var err error
			host, port, err = net.SplitHostPort(raw)

			if err != nil {
				return "", fmt.Errorf("must be a valid host or host:port")
			}

			if !strings.Contains(host, ":") {
				return "", fmt.Errorf("brackets require an IPv6 literal")
			}

			hasPort = true
		}
	} else if addr, ok := parseIPAddress(raw); ok {
		if addr.Zone() != "" {
			return "", fmt.Errorf("IPv6 zones are not supported")
		}

		return addr.Unmap().String(), nil
	} else if strings.Count(raw, ":") == 1 {
		var err error

		host, port, err = net.SplitHostPort(raw)

		if err != nil {
			return "", fmt.Errorf("must be a valid host or host:port")
		}

		hasPort = true
	} else if strings.Contains(raw, ":") {
		return "", fmt.Errorf("must be a valid IP literal")
	}

	host = strings.TrimSpace(host)

	if host == "" {
		return "", fmt.Errorf("host must not be blank")
	}

	if strings.Contains(host, "%") {
		return "", fmt.Errorf("IPv6 zones are not supported")
	}

	canonicalHost := ""
	canonicalInput := canonicalHostname(host)

	if addr, ok := parseIPAddress(canonicalInput); ok {
		if addr.Zone() != "" {
			return "", fmt.Errorf("IPv6 zones are not supported")
		}

		canonicalHost = addr.Unmap().String()
	} else {
		if err := validateDNSHostname(host); err != nil {
			return "", err
		}

		canonicalHost = canonicalHostname(host)
	}

	if !hasPort {
		return canonicalHost, nil
	}

	parsedPort, err := strconv.ParseUint(port, 10, 16)

	if err != nil || port == "" {
		return "", fmt.Errorf("port must be a number from 0 through 65535")
	}

	return net.JoinHostPort(canonicalHost, strconv.FormatUint(parsedPort, 10)), nil
}

func validateDNSHostname(value string) error {
	hostname := canonicalHostname(value)
	if hostname == "" || len(hostname) > 253 {
		return fmt.Errorf("must be a valid DNS name")
	}

	numeric := true
	for idx := range len(hostname) {
		if hostname[idx] < '0' || hostname[idx] > '9' {
			if hostname[idx] != '.' {
				numeric = false
			}
		}
	}

	if numeric {
		return fmt.Errorf("must be a valid IP literal")
	}

	for _, label := range strings.Split(hostname, ".") {
		if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return fmt.Errorf("must be a valid DNS name")
		}

		for idx := range len(label) {
			char := label[idx]

			if !isASCIIAlpha(char) && (char < '0' || char > '9') && char != '-' {
				return fmt.Errorf("must be a valid DNS name")
			}
		}
	}

	return nil
}

func validateHeaderName(header string) *HeaderValidationError {
	if header == "" {
		return &HeaderValidationError{
			Header: header,
			Reason: "name is not a valid HTTP field-name token",
		}
	}

	for idx := range len(header) {
		if !isTokenByte(header[idx]) {
			return &HeaderValidationError{
				Header: header,
				Reason: "name is not a valid HTTP field-name token",
			}
		}
	}

	return nil
}

func validateHeaderValue(header, value string) *HeaderValidationError {
	for idx := range len(value) {
		char := value[idx]

		switch {
		case char == '\r' || char == '\n':
			return &HeaderValidationError{Header: header, Reason: "value contains a newline"}
		case (char < 0x20 && char != '\t') || char == 0x7f:
			return &HeaderValidationError{
				Header: header,
				Reason: "value contains a prohibited control character",
			}
		}
	}

	return nil
}

func isASCIIAlpha(value byte) bool {
	return value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z'
}

func isTokenByte(value byte) bool {
	if isASCIIAlpha(value) || value >= '0' && value <= '9' {
		return true
	}

	return strings.ContainsRune("!#$%&'*+-.^_`|~", rune(value))
}

func containsValue(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}

	return false
}

func containsHost(hosts []string, u *url.URL) bool {
	if len(hosts) == 0 {
		return false
	}

	rawHost := normalizeHostValue(u.Host)
	hostname := canonicalHostKey(u.Hostname())

	for _, host := range hosts {
		if host == rawHost || host == hostname {
			return true
		}
	}

	return false
}

func normalizeHostValue(value string) string {
	if normalized, err := normalizeConfiguredHost(value); err == nil {
		return normalized
	}

	value = asciiLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}

	if host, port, err := net.SplitHostPort(value); err == nil {
		return net.JoinHostPort(canonicalHostKey(host), port)
	}

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")
	}

	return canonicalHostKey(value)
}

func canonicalHostname(hostname string) string {
	hostname = asciiLower(strings.TrimSpace(hostname))

	return strings.TrimSuffix(hostname, ".")
}

func canonicalHostKey(hostname string) string {
	hostname = canonicalHostname(hostname)
	if addr, ok := parseIPAddress(hostname); ok {
		return addr.String()
	}

	return hostname
}

func isASCII(value string) bool {
	for idx := range len(value) {
		if value[idx] >= 0x80 {
			return false
		}
	}

	return true
}

func asciiLower(value string) string {
	for idx := range len(value) {
		if value[idx] >= 'A' && value[idx] <= 'Z' {
			data := []byte(value)

			for dataIdx := idx; dataIdx < len(data); dataIdx++ {
				if data[dataIdx] >= 'A' && data[dataIdx] <= 'Z' {
					data[dataIdx] += 'a' - 'A'
				}
			}

			return string(data)
		}
	}

	return value
}

func asciiUpper(value string) string {
	for idx := range len(value) {
		if value[idx] >= 'a' && value[idx] <= 'z' {
			data := []byte(value)

			for dataIdx := idx; dataIdx < len(data); dataIdx++ {
				if data[dataIdx] >= 'a' && data[dataIdx] <= 'z' {
					data[dataIdx] -= 'a' - 'A'
				}
			}

			return string(data)
		}
	}

	return value
}

func isLocalhostName(hostname string) bool {
	return hostname == "localhost" || strings.HasSuffix(hostname, ".localhost")
}

func parseIPAddress(hostname string) (netip.Addr, bool) {
	addr, err := netip.ParseAddr(hostname)

	if err == nil {
		return addr.Unmap(), true
	}

	return parseAbbreviatedIPv4(hostname)
}

func parseAbbreviatedIPv4(hostname string) (netip.Addr, bool) {
	if hostname == "" {
		return netip.Addr{}, false
	}

	var (
		values   [4]uint64
		part     int
		hasDigit bool
	)

	for idx := range len(hostname) {
		char := hostname[idx]
		if char == '.' {
			if !hasDigit || part == len(values)-1 {
				return netip.Addr{}, false
			}

			part++
			hasDigit = false

			continue
		}

		if char < '0' || char > '9' {
			return netip.Addr{}, false
		}

		value := values[part]*10 + uint64(char-'0')
		if value > 0xffffffff {
			return netip.Addr{}, false
		}

		values[part] = value
		hasDigit = true
	}

	if !hasDigit {
		return netip.Addr{}, false
	}

	var value uint64

	switch part + 1 {
	case 1:
		value = values[0]
	case 2:
		if values[0] > 0xff || values[1] > 0xffffff {
			return netip.Addr{}, false
		}

		value = values[0]<<24 | values[1]
	case 3:
		if values[0] > 0xff || values[1] > 0xff || values[2] > 0xffff {
			return netip.Addr{}, false
		}

		value = values[0]<<24 | values[1]<<16 | values[2]
	case 4:
		for _, part := range values {
			if part > 0xff {
				return netip.Addr{}, false
			}
		}

		value = values[0]<<24 | values[1]<<16 | values[2]<<8 | values[3]
	}

	return netip.AddrFrom4([4]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	}), true
}

func wellKnownNAT64IPv4(addr netip.Addr) (netip.Addr, bool) {
	if !wellKnownNAT64.Contains(addr) {
		return netip.Addr{}, false
	}

	value := addr.As16()

	return netip.AddrFrom4([4]byte{value[12], value[13], value[14], value[15]}), true
}

func isAlwaysBlockedAddress(addr netip.Addr) bool {
	if addr.Is4() {
		if prefixContains(ipv4PublicExceptions, addr) {
			return false
		}

		return prefixContains(ipv4NonPublic, addr)
	}

	if wellKnownNAT64.Contains(addr) {
		return false
	}

	return !ipv6AllocatedGlobal.Contains(addr) || prefixContains(ipv6NonPublic, addr)
}

func prefixContains(prefixes []netip.Prefix, addr netip.Addr) bool {
	for _, prefix := range prefixes {
		if prefix.Contains(addr) {
			return true
		}
	}

	return false
}

func addressSubject(addr netip.Addr) string {
	return fmt.Sprintf("address %q", addr.Unmap().String())
}
