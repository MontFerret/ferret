package http

import (
	"fmt"
	"net"
	stdhttp "net/http"
	"net/netip"
	"net/url"
	"strings"
)

type policyTarget string

const (
	policyTargetRequest  policyTarget = "request"
	policyTargetRedirect policyTarget = "redirect destination"
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

func normalizeHosts(hosts []string) []string {
	if hosts == nil {
		return nil
	}

	normalized := make([]string, 0, len(hosts))

	for _, host := range hosts {
		host = normalizeHostValue(host)

		if host == "" {
			continue
		}

		normalized = append(normalized, host)
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
	value = strings.ToLower(strings.TrimSpace(value))
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
	hostname = strings.ToLower(strings.TrimSpace(hostname))

	return strings.TrimSuffix(hostname, ".")
}

func canonicalHostKey(hostname string) string {
	hostname = canonicalHostname(hostname)
	if addr, ok := parseIPAddress(hostname); ok {
		return addr.String()
	}

	return hostname
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

func newPolicyError(target policyTarget, subject, reason string) error {
	return fmt.Errorf("http: %s blocked by access policy: %s: %s", target, subject, reason)
}

func addressSubject(addr netip.Addr) string {
	return fmt.Sprintf("address %q", addr.Unmap().String())
}
