package http

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// WithTimeout sets the overall HTTP client timeout. Zero restores the secure
// 30-second default; negative values make policy construction fail.
func WithTimeout(timeout time.Duration) PolicyOption {
	return func(p *Policy) {
		if timeout < 0 {
			p.setConfigurationError("WithTimeout", timeout.String(), "must not be negative")

			return
		}

		if timeout == 0 {
			timeout = defaultTimeout
		}

		p.timeout = timeout
	}
}

// WithNoTimeout explicitly disables the overall HTTP client timeout. Request
// contexts may still impose deadlines.
func WithNoTimeout() PolicyOption {
	return func(p *Policy) {
		p.timeout = 0
	}
}

// WithMaxRequestSize limits request bodies in bytes. Zero restores the secure
// 16 MiB default; negative values make policy construction fail.
func WithMaxRequestSize(size int64) PolicyOption {
	return func(p *Policy) {
		if size < 0 {
			p.setConfigurationError(
				"WithMaxRequestSize",
				strconv.FormatInt(size, 10),
				"must not be negative",
			)

			return
		}

		if size == 0 {
			size = defaultMaxRequestSize
		}

		p.maxRequestSize = size
	}
}

// WithUnlimitedRequestSize explicitly disables the request body size limit.
func WithUnlimitedRequestSize() PolicyOption {
	return func(p *Policy) {
		p.maxRequestSize = 0
	}
}

// WithMaxResponseSize limits materialized response bodies in bytes. Zero
// restores the secure 16 MiB default; negative values make construction fail.
func WithMaxResponseSize(size int64) PolicyOption {
	return func(p *Policy) {
		if size < 0 {
			p.setConfigurationError(
				"WithMaxResponseSize",
				strconv.FormatInt(size, 10),
				"must not be negative",
			)

			return
		}

		if size == 0 {
			size = defaultMaxResponseSize
		}

		p.maxResponseSize = size
	}
}

// WithUnlimitedResponseSize explicitly disables the response body size limit.
func WithUnlimitedResponseSize() PolicyOption {
	return func(p *Policy) {
		p.maxResponseSize = 0
	}
}

// WithMaxResponseHeaderSize limits response headers in bytes. Zero restores
// the secure 1 MiB default; negative values make policy construction fail.
// Response headers cannot be configured as unlimited.
func WithMaxResponseHeaderSize(size int64) PolicyOption {
	return func(p *Policy) {
		if size < 0 {
			p.setConfigurationError(
				"WithMaxResponseHeaderSize",
				strconv.FormatInt(size, 10),
				"must not be negative",
			)

			return
		}

		if size == 0 {
			size = defaultMaxResponseHeaderSize
		}

		p.maxResponseHeaderSize = size
	}
}

// WithFollowRedirects controls whether redirects are followed. Disabling
// redirects is distinct from configuring the maximum redirect count.
func WithFollowRedirects(follow bool) PolicyOption {
	return func(p *Policy) {
		p.followRedirects = follow
	}
}

// WithMaxRedirects limits how many redirects may be followed. Zero restores
// the secure default of 10; negative values make policy construction fail.
func WithMaxRedirects(count int) PolicyOption {
	return func(p *Policy) {
		if count < 0 {
			p.setConfigurationError(
				"WithMaxRedirects",
				strconv.Itoa(count),
				"must not be negative",
			)

			return
		}

		if count == 0 {
			count = defaultMaxRedirects
		}

		p.maxRedirects = count
	}
}

// WithAllowedSchemes replaces the set of permitted URL schemes. Entries are
// trimmed, normalized to lowercase ASCII, validated, and deduplicated.
func WithAllowedSchemes(schemes ...string) PolicyOption {
	return func(p *Policy) {
		p.allowedSchemes = append([]string(nil), schemes...)

		for _, scheme := range schemes {
			if err := validateConfiguredScheme(scheme); err != nil {
				p.setConfigurationError("WithAllowedSchemes", scheme, err.Error())

				return
			}
		}
	}
}

// WithAllowedMethods replaces the set of permitted HTTP methods. Entries are
// trimmed, normalized to uppercase ASCII, validated, and deduplicated.
func WithAllowedMethods(methods ...string) PolicyOption {
	return func(p *Policy) {
		p.allowedMethods = append([]string(nil), methods...)

		for _, method := range methods {
			if !isValidMethod(normalizeMethod(method)) {
				p.setConfigurationError(
					"WithAllowedMethods",
					method,
					"must be a non-empty HTTP method token",
				)

				return
			}
		}
	}
}

// WithAllowedHosts restricts requests to the provided exact hosts. Entries
// must be ASCII DNS names or IP literals, optionally with a numeric port. A
// hostname-only entry matches every port; subdomains are not implicit.
func WithAllowedHosts(hosts ...string) PolicyOption {
	return func(p *Policy) {
		p.allowedHosts = append([]string(nil), hosts...)

		validateConfiguredHosts(p, "WithAllowedHosts", hosts)
	}
}

// WithBlockedHosts blocks the provided exact hosts. Entries must be ASCII DNS
// names or IP literals, optionally with a numeric port. A hostname-only entry
// matches every port; subdomains are not implicit.
func WithBlockedHosts(hosts ...string) PolicyOption {
	return func(p *Policy) {
		p.blockedHosts = append([]string(nil), hosts...)

		validateConfiguredHosts(p, "WithBlockedHosts", hosts)
	}
}

// WithAllowLocalhost controls whether localhost and loopback addresses are allowed.
func WithAllowLocalhost(allow bool) PolicyOption {
	return func(p *Policy) {
		p.allowLocalhost = allow
	}
}

// WithAllowPrivateNetworks controls whether private IP network addresses are allowed.
func WithAllowPrivateNetworks(allow bool) PolicyOption {
	return func(p *Policy) {
		p.allowPrivateNetworks = allow
	}
}

// WithAllowLinkLocal controls whether IPv4 and IPv6 link-local addresses are allowed.
func WithAllowLinkLocal(allow bool) PolicyOption {
	return func(p *Policy) {
		p.allowLinkLocal = allow
	}
}

// WithDefaultHeader sets a validated default request header when a request
// does not already supply that header.
func WithDefaultHeader(key, value string) PolicyOption {
	return func(p *Policy) {
		p.addDefaultHeader("WithDefaultHeader", key, value)
	}
}

// WithDefaultHeaders adds validated default request headers. Empty maps are a
// no-op. Keys are processed deterministically.
func WithDefaultHeaders(headers map[string]string) PolicyOption {
	return func(p *Policy) {
		if len(headers) == 0 {
			return
		}

		keys := make([]string, 0, len(headers))

		for key := range headers {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			p.addDefaultHeader("WithDefaultHeaders", key, headers[key])
		}
	}
}

// WithBlockedRequestHeaders rejects outbound requests that supply any of the
// provided header names. Transport-controlled headers are always rejected.
func WithBlockedRequestHeaders(headers ...string) PolicyOption {
	return func(p *Policy) {
		p.blockedRequestHeaders = append([]string(nil), headers...)

		for _, header := range headers {
			if err := validateHeaderName(strings.TrimSpace(header)); err != nil {
				p.setConfigurationError("WithBlockedRequestHeaders", header, err.Reason)

				return
			}
		}
	}
}
