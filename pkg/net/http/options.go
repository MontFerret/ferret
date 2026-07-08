package http

import (
	"strings"
	"time"
)

// WithTimeout sets the underlying HTTP client timeout.
func WithTimeout(timeout time.Duration) Policy {
	return func(p *Policies) {
		if timeout <= 0 {
			p.Timeout = 0
			return
		}

		p.Timeout = timeout
	}
}

// WithMaxRequestSize limits request body size in bytes. A non-positive value disables the limit.
func WithMaxRequestSize(size int64) Policy {
	return func(p *Policies) {
		if size < 0 {
			size = 0
		}

		p.MaxRequestSize = size
	}
}

// WithMaxResponseSize limits response body size in bytes. A non-positive value disables the limit.
func WithMaxResponseSize(size int64) Policy {
	return func(p *Policies) {
		if size < 0 {
			size = 0
		}

		p.MaxResponseSize = size
	}
}

// WithFollowRedirects controls whether redirects are followed.
func WithFollowRedirects(follow bool) Policy {
	return func(p *Policies) {
		p.FollowRedirects = follow
	}
}

// WithMaxRedirects limits followed redirects. A value of 0 allows the standard library default.
func WithMaxRedirects(count int) Policy {
	return func(p *Policies) {
		if count < 0 {
			count = 0
		}

		p.MaxRedirects = count
	}
}

// WithAllowedSchemes replaces the set of permitted URL schemes.
func WithAllowedSchemes(schemes ...string) Policy {
	return func(p *Policies) {
		p.AllowedSchemes = append([]string(nil), schemes...)
	}
}

// WithAllowedHosts restricts requests to the provided host names.
func WithAllowedHosts(hosts ...string) Policy {
	return func(p *Policies) {
		p.AllowedHosts = append([]string(nil), hosts...)
	}
}

// WithBlockedHosts blocks requests to the provided host names.
func WithBlockedHosts(hosts ...string) Policy {
	return func(p *Policies) {
		p.BlockedHosts = append([]string(nil), hosts...)
	}
}

// WithAllowLocalhost controls whether localhost and loopback addresses are allowed.
func WithAllowLocalhost(allow bool) Policy {
	return func(p *Policies) {
		p.AllowLocalhost = allow
	}
}

// WithAllowPrivateNetworks controls whether private IP network addresses are allowed.
func WithAllowPrivateNetworks(allow bool) Policy {
	return func(p *Policies) {
		p.AllowPrivateNetworks = allow
	}
}

// WithDefaultHeader sets a default request header when the request does not already provide it.
func WithDefaultHeader(key, value string) Policy {
	return func(p *Policies) {
		key = strings.TrimSpace(key)
		if key == "" {
			return
		}

		if p.DefaultHeaders == nil {
			p.DefaultHeaders = make(map[string]string)
		}

		p.DefaultHeaders[key] = value
	}
}

// WithDefaultHeaders sets default request headers when the request does not already provide them.
func WithDefaultHeaders(headers map[string]string) Policy {
	return func(p *Policies) {
		if len(headers) == 0 {
			return
		}

		if p.DefaultHeaders == nil {
			p.DefaultHeaders = make(map[string]string, len(headers))
		}

		for key, value := range headers {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}

			p.DefaultHeaders[key] = value
		}
	}
}

// WithBlockedRequestHeaders removes the provided header names from outbound requests.
func WithBlockedRequestHeaders(headers ...string) Policy {
	return func(p *Policies) {
		p.BlockedRequestHeaders = append([]string(nil), headers...)
	}
}
