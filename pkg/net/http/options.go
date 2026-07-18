package http

import (
	"strings"
	"time"
)

// WithTimeout sets the underlying HTTP client timeout.
func WithTimeout(timeout time.Duration) Policy {
	return func(p *Policies) {
		if timeout <= 0 {
			p.timeout = 0
			return
		}

		p.timeout = timeout
	}
}

// WithMaxRequestSize limits request body size in bytes. A non-positive value disables the limit.
func WithMaxRequestSize(size int64) Policy {
	return func(p *Policies) {
		if size < 0 {
			size = 0
		}

		p.maxRequestSize = size
	}
}

// WithMaxResponseSize limits response body size in bytes. A non-positive value disables the limit.
func WithMaxResponseSize(size int64) Policy {
	return func(p *Policies) {
		if size < 0 {
			size = 0
		}

		p.maxResponseSize = size
	}
}

// WithFollowRedirects controls whether redirects are followed.
func WithFollowRedirects(follow bool) Policy {
	return func(p *Policies) {
		p.followRedirects = follow
	}
}

// WithMaxRedirects limits followed redirects. A value of 0 allows the standard library default.
func WithMaxRedirects(count int) Policy {
	return func(p *Policies) {
		if count < 0 {
			count = 0
		}

		p.maxRedirects = count
	}
}

// WithAllowedSchemes replaces the set of permitted URL schemes.
func WithAllowedSchemes(schemes ...string) Policy {
	return func(p *Policies) {
		p.allowedSchemes = append([]string(nil), schemes...)
	}
}

// WithAllowedHosts restricts requests to the provided host names.
func WithAllowedHosts(hosts ...string) Policy {
	return func(p *Policies) {
		p.allowedHosts = append([]string(nil), hosts...)
	}
}

// WithBlockedHosts blocks requests to the provided host names.
func WithBlockedHosts(hosts ...string) Policy {
	return func(p *Policies) {
		p.blockedHosts = append([]string(nil), hosts...)
	}
}

// WithAllowLocalhost controls whether localhost and loopback addresses are allowed.
func WithAllowLocalhost(allow bool) Policy {
	return func(p *Policies) {
		p.allowLocalhost = allow
	}
}

// WithAllowPrivateNetworks controls whether private IP network addresses are allowed.
func WithAllowPrivateNetworks(allow bool) Policy {
	return func(p *Policies) {
		p.allowPrivateNetworks = allow
	}
}

// WithAllowLinkLocal controls whether IPv4 and IPv6 link-local addresses are allowed.
func WithAllowLinkLocal(allow bool) Policy {
	return func(p *Policies) {
		p.allowLinkLocal = allow
	}
}

// WithDefaultHeader sets a default request header when the request does not already provide it.
func WithDefaultHeader(key, value string) Policy {
	return func(p *Policies) {
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
func WithDefaultHeaders(headers map[string]string) Policy {
	return func(p *Policies) {
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

// WithBlockedRequestHeaders removes the provided header names from outbound requests.
func WithBlockedRequestHeaders(headers ...string) Policy {
	return func(p *Policies) {
		p.blockedRequestHeaders = append([]string(nil), headers...)
	}
}
