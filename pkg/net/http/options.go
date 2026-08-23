package http

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/ziflex/go-options"
)

const (
	fieldTimeout               = "timeout"
	fieldMaxRequestSize        = "max request size"
	fieldMaxResponseSize       = "max response size"
	fieldMaxResponseHeaderSize = "max response header size"
	fieldMaxRedirects          = "max redirects"
	fieldAllowedSchemes        = "allowed schemes"
	fieldAllowedMethods        = "allowed methods"
	fieldAllowedHosts          = "allowed hosts"
	fieldBlockedHosts          = "blocked hosts"
	fieldDefaultHeader         = "default header"
	fieldDefaultHeaders        = "default headers"
	fieldBlockedRequestHeaders = "blocked request headers"
)

// WithTimeout sets the overall HTTP client timeout. Zero restores the secure
// 30-second default; negative values make policy construction fail.
func WithTimeout(timeout time.Duration) PolicyOption {
	return options.New(func(p *Policy, timeout time.Duration) {
		if timeout == 0 {
			timeout = defaultTimeout
		}

		p.timeout = timeout
	}).
		Named(fieldTimeout).
		Validators(options.NonNegative[time.Duration]()).
		Value(timeout).
		Build()
}

// WithNoTimeout explicitly disables the overall HTTP client timeout. Request
// contexts may still impose deadlines.
func WithNoTimeout() PolicyOption {
	return func(p *Policy) error {
		p.timeout = 0

		return nil
	}
}

// WithMaxRequestSize limits request bodies in bytes. Zero restores the secure
// 16 MiB default; negative values make policy construction fail.
func WithMaxRequestSize(size int64) PolicyOption {
	return options.New(func(p *Policy, size int64) {
		if size == 0 {
			size = defaultMaxRequestSize
		}

		p.maxRequestSize = size
	}).
		Named(fieldMaxRequestSize).
		Validators(options.NonNegative[int64]()).
		Value(size).
		Build()
}

// WithUnlimitedRequestSize explicitly disables the request body size limit.
func WithUnlimitedRequestSize() PolicyOption {
	return func(p *Policy) error {
		p.maxRequestSize = 0

		return nil
	}
}

// WithMaxResponseSize limits materialized response bodies in bytes. Zero
// restores the secure 16 MiB default; negative values make construction fail.
func WithMaxResponseSize(size int64) PolicyOption {
	return options.New(func(p *Policy, size int64) {
		if size == 0 {
			size = defaultMaxResponseSize
		}

		p.maxResponseSize = size
	}).
		Named(fieldMaxResponseSize).
		Validators(options.NonNegative[int64]()).
		Value(size).
		Build()
}

// WithUnlimitedResponseSize explicitly disables the response body size limit.
func WithUnlimitedResponseSize() PolicyOption {
	return func(p *Policy) error {
		p.maxResponseSize = 0

		return nil
	}
}

// WithMaxResponseHeaderSize limits response headers in bytes. Zero restores
// the secure 1 MiB default; negative values make policy construction fail.
// Response headers cannot be configured as unlimited.
func WithMaxResponseHeaderSize(size int64) PolicyOption {
	return options.New(func(p *Policy, size int64) {
		if size == 0 {
			size = defaultMaxResponseHeaderSize
		}

		p.maxResponseHeaderSize = size
	}).
		Named(fieldMaxResponseHeaderSize).
		Validators(options.NonNegative[int64]()).
		Value(size).
		Build()
}

// WithFollowRedirects controls whether redirects are followed. Disabling
// redirects is distinct from configuring the maximum redirect count.
func WithFollowRedirects(follow bool) PolicyOption {
	return func(p *Policy) error {
		p.followRedirects = follow

		return nil
	}
}

// WithMaxRedirects limits how many redirects may be followed. Zero restores
// the secure default of 10; negative values make policy construction fail.
func WithMaxRedirects(count int) PolicyOption {
	return options.New(func(p *Policy, count int) {
		if count == 0 {
			count = defaultMaxRedirects
		}

		p.maxRedirects = count
	}).
		Named(fieldMaxRedirects).
		Validators(options.NonNegative[int]()).
		Value(count).
		Build()
}

// WithAllowedSchemes replaces the set of permitted URL schemes. Entries are
// trimmed, normalized to lowercase ASCII, validated, and deduplicated.
func WithAllowedSchemes(schemes ...string) PolicyOption {
	return func(p *Policy) error {
		var errs []error

		for _, scheme := range schemes {
			if err := validateConfiguredScheme(scheme); err != nil {
				errs = append(errs, newOptionValidationError(fieldAllowedSchemes, scheme, err))
			}
		}

		if err := errors.Join(errs...); err != nil {
			return err
		}

		p.allowedSchemes = normalizeValues(schemes)

		return nil
	}
}

// WithAllowedMethods replaces the set of permitted HTTP methods. Entries are
// trimmed, normalized to uppercase ASCII, validated, and deduplicated.
func WithAllowedMethods(methods ...string) PolicyOption {
	return func(p *Policy) error {
		var errs []error

		for _, method := range methods {
			if !isValidMethod(normalizeMethod(method)) {
				errs = append(
					errs,
					newOptionValidationError(
						fieldAllowedMethods,
						method,
						errors.New("must be a non-empty HTTP method token"),
					),
				)
			}
		}

		if err := errors.Join(errs...); err != nil {
			return err
		}

		p.allowedMethods = normalizeMethods(methods)

		return nil
	}
}

// WithAllowedHosts restricts requests to the provided exact hosts. Entries
// must be ASCII DNS names or IP literals, optionally with a numeric port. A
// hostname-only entry matches every port; subdomains are not implicit.
func WithAllowedHosts(hosts ...string) PolicyOption {
	return configuredHostsOption(fieldAllowedHosts, func(p *Policy, hosts []string) {
		p.allowedHosts = hosts
	}, hosts)
}

// WithBlockedHosts blocks the provided exact hosts. Entries must be ASCII DNS
// names or IP literals, optionally with a numeric port. A hostname-only entry
// matches every port; subdomains are not implicit.
func WithBlockedHosts(hosts ...string) PolicyOption {
	return configuredHostsOption(fieldBlockedHosts, func(p *Policy, hosts []string) {
		p.blockedHosts = hosts
	}, hosts)
}

// WithAllowLocalhost controls whether localhost and loopback addresses are allowed.
func WithAllowLocalhost(allow bool) PolicyOption {
	return func(p *Policy) error {
		p.allowLocalhost = allow

		return nil
	}
}

// WithAllowPrivateNetworks controls whether private IP network addresses are allowed.
func WithAllowPrivateNetworks(allow bool) PolicyOption {
	return func(p *Policy) error {
		p.allowPrivateNetworks = allow

		return nil
	}
}

// WithAllowLinkLocal controls whether IPv4 and IPv6 link-local addresses are allowed.
func WithAllowLinkLocal(allow bool) PolicyOption {
	return func(p *Policy) error {
		p.allowLinkLocal = allow

		return nil
	}
}

// WithDefaultHeader sets a validated default request header when a request
// does not already supply that header.
func WithDefaultHeader(key, value string) PolicyOption {
	return func(p *Policy) error {
		return p.addDefaultHeader(fieldDefaultHeader, key, value)
	}
}

// WithDefaultHeaders adds validated default request headers. Empty maps are a
// no-op. Keys are processed deterministically.
func WithDefaultHeaders(headers map[string]string) PolicyOption {
	return func(p *Policy) error {
		if len(headers) == 0 {
			return nil
		}

		keys := make([]string, 0, len(headers))

		for key := range headers {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		var errs []error

		for _, key := range keys {
			if err := p.addDefaultHeader(fieldDefaultHeaders, key, headers[key]); err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Join(errs...)
	}
}

// WithBlockedRequestHeaders rejects outbound requests that supply any of the
// provided header names. Transport-controlled headers are always rejected.
func WithBlockedRequestHeaders(headers ...string) PolicyOption {
	return func(p *Policy) error {
		valid := make([]string, 0, len(headers))
		var errs []error

		for _, header := range headers {
			rawHeader := header
			header = strings.TrimSpace(rawHeader)

			if err := validateHeaderName(header); err != nil {
				errs = append(
					errs,
					newOptionValidationError(
						fieldBlockedRequestHeaders,
						rawHeader,
						errors.New(err.Reason),
					),
				)

				continue
			}

			valid = append(valid, header)
		}

		p.blockedRequestHeaders = normalizeHeaders(valid)

		return errors.Join(errs...)
	}
}

func configuredHostsOption(
	field string,
	setter func(*Policy, []string),
	hosts []string,
) PolicyOption {
	return func(p *Policy) error {
		var errs []error

		for _, host := range hosts {
			if _, err := normalizeConfiguredHost(host); err != nil {
				errs = append(errs, newOptionValidationError(field, host, err))
			}
		}

		if err := errors.Join(errs...); err != nil {
			return err
		}

		setter(p, normalizeHosts(hosts))

		return nil
	}
}

func newOptionValidationError(field, value string, reason error) options.ValidationError {
	return options.ValidationError{
		Field:  field,
		Value:  value,
		Reason: reason,
	}
}
