package http

import (
	"strings"
	"time"
)

// WithTimeout sets the overall HTTP client timeout. The secure default is 30
// seconds. A zero or negative value explicitly disables it. A request context
// may impose a shorter deadline.
func WithTimeout(timeout time.Duration) PolicyOption {
	return func(p *Policy) {
		if timeout <= 0 {
			p.timeout = 0
			return
		}

		p.timeout = timeout
	}
}

// WithMaxRequestSize limits request body size in bytes. A non-positive value disables the limit.
func WithMaxRequestSize(size int64) PolicyOption {
	return func(p *Policy) {
		if size < 0 {
			size = 0
		}

		p.maxRequestSize = size
	}
}

// WithMaxResponseSize limits response body size in bytes. A non-positive value disables the limit.
func WithMaxResponseSize(size int64) PolicyOption {
	return func(p *Policy) {
		if size < 0 {
			size = 0
		}

		p.maxResponseSize = size
	}
}

// WithMaxResponseHeaderSize limits response headers in bytes. A non-positive
// value restores the secure 1 MiB default.
func WithMaxResponseHeaderSize(size int64) PolicyOption {
	return func(p *Policy) {
		if size <= 0 {
			size = defaultMaxResponseHeaderSize
		}

		p.maxResponseHeaderSize = size
	}
}

// WithFollowRedirects controls whether redirects are followed.
func WithFollowRedirects(follow bool) PolicyOption {
	return func(p *Policy) {
		p.followRedirects = follow
	}
}

// WithMaxRedirects limits followed redirects. A value of 0 allows the standard library default.
func WithMaxRedirects(count int) PolicyOption {
	return func(p *Policy) {
		if count < 0 {
			count = 0
		}

		p.maxRedirects = count
	}
}

// WithAllowedSchemes replaces the set of permitted URL schemes.
func WithAllowedSchemes(schemes ...string) PolicyOption {
	return func(p *Policy) {
		p.allowedSchemes = append([]string(nil), schemes...)
	}
}

// WithAllowedMethods replaces the set of permitted HTTP methods. Values are
// trimmed, normalized to uppercase ASCII, and must be valid HTTP method
// tokens. Invalid and blank entries are ignored; an empty result denies every
// method.
func WithAllowedMethods(methods ...string) PolicyOption {
	return func(p *Policy) {
		p.allowedMethods = append([]string(nil), methods...)
	}
}

// WithAllowedHosts restricts requests to the provided exact host names.
// Entries must use ASCII/punycode. A hostname-only entry matches every port;
// use host:port for a port-specific rule. Subdomains are not matched implicitly.
func WithAllowedHosts(hosts ...string) PolicyOption {
	return func(p *Policy) {
		p.allowedHosts = append([]string(nil), hosts...)
	}
}

// WithBlockedHosts blocks the provided exact host names. Entries must use
// ASCII/punycode. A hostname-only entry matches every port; use host:port for
// a port-specific rule. Subdomains are not matched implicitly.
func WithBlockedHosts(hosts ...string) PolicyOption {
	return func(p *Policy) {
		p.blockedHosts = append([]string(nil), hosts...)
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

// WithDefaultHeader sets a default request header when the request does not already provide it.
func WithDefaultHeader(key, value string) PolicyOption {
	return func(p *Policy) {
		key = strings.TrimSpace(key)
		if key == "" {
			return
		}

		if p.defaultHeaders == nil {
			p.defaultHeaders = make(map[string]string)
		}

		p.defaultHeaders[key] = value
	}
}

// WithDefaultHeaders sets default request headers when the request does not already provide them.
func WithDefaultHeaders(headers map[string]string) PolicyOption {
	return func(p *Policy) {
		if len(headers) == 0 {
			return
		}

		if p.defaultHeaders == nil {
			p.defaultHeaders = make(map[string]string, len(headers))
		}

		for key, value := range headers {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}

			p.defaultHeaders[key] = value
		}
	}
}

// WithBlockedRequestHeaders rejects outbound requests that supply any of the
// provided header names. Transport-reserved headers are always rejected.
func WithBlockedRequestHeaders(headers ...string) PolicyOption {
	return func(p *Policy) {
		p.blockedRequestHeaders = append([]string(nil), headers...)
	}
}
