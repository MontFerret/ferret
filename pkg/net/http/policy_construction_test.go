package http

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewPolicyRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		options []PolicyOption
	}{
		{name: "invalid method", options: []PolicyOption{WithAllowedMethods("GET", "BAD METHOD")}},
		{name: "blank method", options: []PolicyOption{WithAllowedMethods("")}},
		{name: "invalid scheme", options: []PolicyOption{WithAllowedSchemes("https", "://")}},
		{name: "blank scheme", options: []PolicyOption{WithAllowedSchemes("")}},
		{name: "blank allowed host", options: []PolicyOption{WithAllowedHosts("")}},
		{name: "blank blocked host", options: []PolicyOption{WithBlockedHosts(" ")}},
		{name: "malformed blocked host", options: []PolicyOption{WithBlockedHosts("example..com")}},
		{name: "invalid blocked header", options: []PolicyOption{WithBlockedRequestHeaders("Bad Header")}},
		{name: "blank blocked header", options: []PolicyOption{WithBlockedRequestHeaders("")}},
		{name: "blank default header", options: []PolicyOption{WithDefaultHeader("", "value")}},
		{name: "invalid default header", options: []PolicyOption{WithDefaultHeader("Bad Header", "value")}},
		{name: "invalid default value", options: []PolicyOption{WithDefaultHeader("X-Test", "safe\r\nInjected: true")}},
		{name: "negative timeout", options: []PolicyOption{WithTimeout(-time.Nanosecond)}},
		{name: "negative redirects", options: []PolicyOption{WithMaxRedirects(-1)}},
		{name: "negative request limit", options: []PolicyOption{WithMaxRequestSize(-1)}},
		{name: "negative response limit", options: []PolicyOption{WithMaxResponseSize(-1)}},
		{name: "negative response header limit", options: []PolicyOption{WithMaxResponseHeaderSize(-1)}},
		{
			name: "default then blocked conflict",
			options: []PolicyOption{
				WithDefaultHeader("Authorization", "Bearer configured-secret"),
				WithBlockedRequestHeaders("authorization"),
			},
		},
		{
			name: "blocked then default conflict",
			options: []PolicyOption{
				WithBlockedRequestHeaders("AUTHORIZATION"),
				WithDefaultHeader("authorization", "Bearer configured-secret"),
			},
		},
		{
			name: "case-equivalent defaults conflict",
			options: []PolicyOption{WithDefaultHeaders(map[string]string{
				"X-Test": "one",
				"x-test": "two",
			})},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := NewPolicy(tt.options...)
			if policy != nil {
				t.Fatalf("expected no policy, got %#v", policy)
			}

			configErr := requirePolicyConfigurationError(t, err)
			if configErr.Option == "" || configErr.Reason == "" {
				t.Fatalf("expected populated configuration error, got %#v", configErr)
			}
			for _, secret := range []string{"configured-secret", "safe\r\nInjected: true"} {
				if strings.Contains(err.Error(), secret) || strings.Contains(configErr.Value, secret) {
					t.Fatalf("configuration error leaked a header value: %v", err)
				}
			}
		})
	}
}

func TestNewPolicyRejectsMalformedHosts(t *testing.T) {
	for _, host := range []string{
		"example..com",
		"-example.com",
		"example-.com",
		"*.example.com",
		"example.com:port",
		"example.com:70000",
		"[::1",
		"[::1]suffix",
		"999.1.1.1",
		"éxample.com",
	} {
		t.Run(host, func(t *testing.T) {
			policy, err := NewPolicy(WithAllowedHosts(host))
			if policy != nil {
				t.Fatalf("expected malformed host %q to return no policy", host)
			}
			requirePolicyConfigurationError(t, err)
		})
	}
}

func TestNewPolicyAcceptsAndDeduplicatesCanonicalHosts(t *testing.T) {
	policy := newTestPolicy(t,
		WithAllowedHosts(
			"Example.COM.",
			"example.com",
			"127.1",
			"127.0.0.1",
			"[::ffff:8.8.8.8]",
			"8.8.8.8",
			"xn--xample-9ua.com",
			"example.net:08443",
			"example.net:8443",
			"[2606:4700:4700::1111]:443",
		),
		WithAllowLocalhost(true),
	)

	if got, want := len(policy.allowedHosts), 6; got != want {
		t.Fatalf("expected %d canonical hosts after deduplication, got %d: %v", want, got, policy.allowedHosts)
	}

	for _, rawURL := range []string{
		"https://example.com",
		"http://127.0.0.1",
		"https://[::ffff:8.8.8.8]",
		"https://xn--xample-9ua.com",
		"https://example.net:8443",
		"https://[2606:4700:4700::1111]:443",
	} {
		if err := policy.Eval(newTestPolicyGETRequest(t, rawURL)); err != nil {
			t.Fatalf("expected canonical host %q to be allowed: %v", rawURL, err)
		}
	}

	requirePolicyError(
		t,
		policy.Eval(newTestPolicyGETRequest(t, "https://sub.example.com")),
		PolicyTargetRequest,
	)
}

func TestNewRejectsInvalidPolicyConfiguration(t *testing.T) {
	client, err := New(WithMaxResponseSize(-1))
	if client != nil {
		t.Fatalf("expected no client, got %T", client)
	}
	requirePolicyConfigurationError(t, err)
}

func TestNewPolicyNumericLimitSemantics(t *testing.T) {
	defaults := newTestPolicy(t,
		WithTimeout(0),
		WithMaxRedirects(0),
		WithMaxRequestSize(0),
		WithMaxResponseSize(0),
		WithMaxResponseHeaderSize(0),
	)
	if defaults.timeout != defaultTimeout ||
		defaults.maxRedirects != defaultMaxRedirects ||
		defaults.maxRequestSize != defaultMaxRequestSize ||
		defaults.maxResponseSize != defaultMaxResponseSize ||
		defaults.maxResponseHeaderSize != defaultMaxResponseHeaderSize {
		t.Fatalf(
			"zero values did not restore defaults: timeout=%s redirects=%d request=%d response=%d headers=%d",
			defaults.timeout,
			defaults.maxRedirects,
			defaults.maxRequestSize,
			defaults.maxResponseSize,
			defaults.maxResponseHeaderSize,
		)
	}

	custom := newTestPolicy(t,
		WithTimeout(time.Second),
		WithMaxRedirects(2),
		WithMaxRequestSize(3),
		WithMaxResponseSize(4),
		WithMaxResponseHeaderSize(5),
	)
	if custom.timeout != time.Second ||
		custom.maxRedirects != 2 ||
		custom.maxRequestSize != 3 ||
		custom.maxResponseSize != 4 ||
		custom.maxResponseHeaderSize != 5 {
		t.Fatalf("positive values were not preserved: %#v", custom)
	}

	unlimited := newTestPolicy(t,
		WithNoTimeout(),
		WithUnlimitedRequestSize(),
		WithUnlimitedResponseSize(),
	)
	if unlimited.timeout != 0 || unlimited.maxRequestSize != 0 || unlimited.maxResponseSize != 0 {
		t.Fatalf(
			"explicit unlimited options did not disable their limits: timeout=%s request=%d response=%d",
			unlimited.timeout,
			unlimited.maxRequestSize,
			unlimited.maxResponseSize,
		)
	}
	if unlimited.maxResponseHeaderSize != defaultMaxResponseHeaderSize {
		t.Fatalf("explicit unlimited options changed response header limit: %d", unlimited.maxResponseHeaderSize)
	}

	redirectsDisabled := newTestPolicy(t, WithFollowRedirects(false), WithMaxRedirects(3))
	if redirectsDisabled.followRedirects || redirectsDisabled.maxRedirects != 3 {
		t.Fatalf("redirect disablement changed redirect count: %#v", redirectsDisabled)
	}
}

func TestExplicitUnlimitedBodyOptionsDisableOnlyTheirLimits(t *testing.T) {
	requestPolicy := newTestPolicy(t,
		WithMaxRequestSize(3),
		WithUnlimitedRequestSize(),
	)
	request := newTestPolicyGETRequest(t, "https://example.com")
	request.Body = io.NopCloser(strings.NewReader("four"))
	request.ContentLength = 4
	if err := requestPolicy.Eval(request); err != nil {
		t.Fatalf("unlimited request policy rejected body: %v", err)
	}
	if requestPolicy.maxResponseSize != defaultMaxResponseSize {
		t.Fatalf("request unlimited option changed response limit: %d", requestPolicy.maxResponseSize)
	}

	responsePolicy := newTestPolicy(t,
		WithMaxResponseSize(3),
		WithUnlimitedResponseSize(),
	)
	body, err := readResponseBody(strings.NewReader("four"), responsePolicy.maxResponseSize)
	if err != nil {
		t.Fatalf("unlimited response policy rejected body: %v", err)
	}
	if got := string(body); got != "four" {
		t.Fatalf("expected response body %q, got %q", "four", got)
	}
	if responsePolicy.maxRequestSize != defaultMaxRequestSize {
		t.Fatalf("response unlimited option changed request limit: %d", responsePolicy.maxRequestSize)
	}
}

func TestNewPolicyNilAndZeroArgumentOptions(t *testing.T) {
	policy := newTestPolicy(t, nil)
	if err := policy.Eval(newTestPolicyGETRequest(t, "https://example.com")); err != nil {
		t.Fatalf("nil option changed defaults: %v", err)
	}

	methodsCleared := newTestPolicy(t, WithAllowedMethods())
	requirePolicyError(
		t,
		methodsCleared.Eval(newTestPolicyGETRequest(t, "https://example.com")),
		PolicyTargetRequest,
	)

	schemesCleared := newTestPolicy(t, WithAllowedSchemes())
	requirePolicyError(
		t,
		schemesCleared.Eval(newTestPolicyGETRequest(t, "https://example.com")),
		PolicyTargetRequest,
	)
}

func TestNewPolicyReturnsLeafForSingleConfigurationError(t *testing.T) {
	policy, err := NewPolicy(WithMaxRequestSize(-1))
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	configErr := requirePolicyConfigurationError(t, err)
	want := PolicyConfigurationError{
		Option: "WithMaxRequestSize",
		Value:  "-1",
		Reason: "must not be negative",
	}
	if got := *configErr; got != want {
		t.Fatalf("unexpected configuration error: got %#v, want %#v", got, want)
	}

	var multiErr *MultiPolicyConfigurationError
	if errors.As(err, &multiErr) {
		t.Fatalf("expected a single leaf error, got aggregate %#v", multiErr)
	}
}

func TestNewPolicyAggregatesMultipleEntriesFromOneOption(t *testing.T) {
	policy, err := NewPolicy(WithAllowedSchemes(
		"://",
		"https",
		"bad scheme",
		"",
		"://",
	))
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	multiErr := requireMultiPolicyConfigurationError(t, err)
	assertPolicyConfigurationErrors(t, multiErr.Errors, []PolicyConfigurationError{
		{
			Option: "WithAllowedSchemes",
			Value:  "://",
			Reason: "must be a valid URL scheme",
		},
		{
			Option: "WithAllowedSchemes",
			Value:  "bad scheme",
			Reason: "must be a valid URL scheme",
		},
		{
			Option: "WithAllowedSchemes",
			Reason: "must be a non-empty URL scheme",
		},
		{
			Option: "WithAllowedSchemes",
			Value:  "://",
			Reason: "must be a valid URL scheme",
		},
	})
}

func TestNewPolicyOrdersMapBackedDefaultHeaderErrorsByKey(t *testing.T) {
	policy, err := NewPolicy(WithDefaultHeaders(map[string]string{
		"Z Bad Header": "last-secret",
		"A Bad Header": "first-secret",
	}))
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	multiErr := requireMultiPolicyConfigurationError(t, err)
	assertPolicyConfigurationErrors(t, multiErr.Errors, []PolicyConfigurationError{
		{
			Option: "WithDefaultHeaders",
			Value:  "A Bad Header",
			Reason: "name is not a valid HTTP field-name token",
		},
		{
			Option: "WithDefaultHeaders",
			Value:  "Z Bad Header",
			Reason: "name is not a valid HTTP field-name token",
		},
	})

	for _, secret := range []string{"first-secret", "last-secret"} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("aggregate configuration error leaked default value %q: %v", secret, err)
		}
	}
}

func TestNewPolicyAggregatesMixedOptionFailures(t *testing.T) {
	policy, err := NewPolicy(
		WithAllowedMethods("BAD METHOD"),
		WithAllowedSchemes("bad scheme"),
		WithAllowedHosts(""),
		WithBlockedHosts("*.example.com"),
		WithBlockedRequestHeaders("Bad Header"),
		WithDefaultHeader("Another Bad Header", "configured-secret"),
		WithMaxResponseSize(-1),
	)
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	multiErr := requireMultiPolicyConfigurationError(t, err)
	assertPolicyConfigurationErrors(t, multiErr.Errors, []PolicyConfigurationError{
		{
			Option: "WithAllowedMethods",
			Value:  "BAD METHOD",
			Reason: "must be a non-empty HTTP method token",
		},
		{
			Option: "WithAllowedSchemes",
			Value:  "bad scheme",
			Reason: "must be a valid URL scheme",
		},
		{
			Option: "WithAllowedHosts",
			Reason: "must not be blank",
		},
		{
			Option: "WithBlockedHosts",
			Value:  "*.example.com",
			Reason: "wildcards are not supported",
		},
		{
			Option: "WithBlockedRequestHeaders",
			Value:  "Bad Header",
			Reason: "name is not a valid HTTP field-name token",
		},
		{
			Option: "WithDefaultHeader",
			Value:  "Another Bad Header",
			Reason: "name is not a valid HTTP field-name token",
		},
		{
			Option: "WithMaxResponseSize",
			Value:  "-1",
			Reason: "must not be negative",
		},
	})

	if strings.Contains(err.Error(), "configured-secret") {
		t.Fatalf("aggregate configuration error leaked a default header value: %v", err)
	}
}

func TestNewPolicyOrdersMixedImmediateAndDeferredConfigurationErrors(t *testing.T) {
	const (
		firstSecret  = "first-configured-secret"
		secondSecret = "second-configured-secret"
	)

	policy, err := NewPolicy(
		WithMaxResponseSize(-2),
		WithDefaultHeader("X-Conflict", firstSecret),
		WithDefaultHeader("x-conflict", secondSecret),
		WithMaxRequestSize(-1),
	)
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	multiErr := requireMultiPolicyConfigurationError(t, err)
	assertPolicyConfigurationErrors(t, multiErr.Errors, []PolicyConfigurationError{
		{
			Option: "WithMaxResponseSize",
			Value:  "-2",
			Reason: "must not be negative",
		},
		{
			Option: "WithMaxRequestSize",
			Value:  "-1",
			Reason: "must not be negative",
		},
		{
			Option: "WithDefaultHeader",
			Value:  "X-Conflict",
			Reason: "conflicts with another default for the same header",
		},
	})

	for _, secret := range []string{firstSecret, secondSecret} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("aggregate configuration error leaked header value %q: %v", secret, err)
		}
		for _, child := range multiErr.Errors {
			if strings.Contains(child.Error(), secret) || strings.Contains(child.Value, secret) {
				t.Fatalf("configuration child leaked header value %q: %#v", secret, child)
			}
		}
	}
}

func TestPolicyConfigurationAggregateSupportsStandardErrorDiscovery(t *testing.T) {
	_, err := NewPolicy(
		WithMaxRequestSize(-1),
		WithMaxResponseSize(-2),
	)
	multiErr := requireMultiPolicyConfigurationError(t, err)
	wrapped := fmt.Errorf("construct HTTP policy: %w", err)

	if !errors.Is(wrapped, ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected wrapped aggregate to match ErrInvalidPolicyConfiguration: %v", wrapped)
	}

	var discoveredMulti *MultiPolicyConfigurationError
	if !errors.As(wrapped, &discoveredMulti) || discoveredMulti != multiErr {
		t.Fatalf("expected wrapped aggregate to expose the exact multi-error, got %#v", discoveredMulti)
	}

	var discoveredLeaf *PolicyConfigurationError
	if !errors.As(wrapped, &discoveredLeaf) || discoveredLeaf != multiErr.Errors[0] {
		t.Fatalf("expected wrapped aggregate to expose its first leaf, got %#v", discoveredLeaf)
	}

	for index, child := range multiErr.Errors {
		if !errors.Is(wrapped, child) {
			t.Fatalf("expected wrapped aggregate to match exact child %d (%#v)", index, child)
		}
	}

	rendered := multiErr.Error()
	if count := strings.Count(rendered, ErrInvalidPolicyConfiguration.Error()); count != 1 {
		t.Fatalf("expected the shared sentinel once, got %d occurrences in %q", count, rendered)
	}
	for _, child := range multiErr.Errors {
		if !strings.Contains(rendered, child.Option) || !strings.Contains(rendered, child.Reason) {
			t.Fatalf("aggregate text omitted child %#v: %q", child, rendered)
		}
	}
}

func TestNewPolicyAggregatesDefaultBlockedConflictsInBothOptionOrders(t *testing.T) {
	const secret = "configured-secret"

	policy, err := NewPolicy(
		WithDefaultHeader("X-Default-First", secret),
		WithBlockedRequestHeaders("x-default-first", "x-blocked-first"),
		WithDefaultHeader("X-Blocked-First", secret),
		WithTimeout(-time.Nanosecond),
	)
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	multiErr := requireMultiPolicyConfigurationError(t, err)
	assertPolicyConfigurationErrors(t, multiErr.Errors, []PolicyConfigurationError{
		{
			Option: "WithTimeout",
			Value:  "-1ns",
			Reason: "must not be negative",
		},
		{
			Option: "WithDefaultHeader",
			Value:  "X-Blocked-First",
			Reason: "default header is also configured as blocked",
		},
		{
			Option: "WithDefaultHeader",
			Value:  "X-Default-First",
			Reason: "default header is also configured as blocked",
		},
	})

	if strings.Contains(err.Error(), secret) {
		t.Fatalf("aggregate configuration error leaked a default header value: %v", err)
	}
}

func TestNewPolicyExcludesInvalidDefaultsFromConflictChecks(t *testing.T) {
	const secret = "configured-secret\r\nInjected: true"

	policy, err := NewPolicy(
		WithDefaultHeader("Authorization", secret),
		WithBlockedRequestHeaders("authorization"),
	)
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	configErr := requirePolicyConfigurationError(t, err)
	if configErr.Option != "WithDefaultHeader" || configErr.Reason != "value contains a newline" {
		t.Fatalf("unexpected configuration error: %#v", configErr)
	}

	var multiErr *MultiPolicyConfigurationError
	if errors.As(err, &multiErr) {
		t.Fatalf("invalid default produced a spurious blocked-header conflict: %#v", multiErr.Errors)
	}
	if strings.Contains(err.Error(), secret) || strings.Contains(configErr.Value, secret) {
		t.Fatalf("configuration error leaked a default header value: %v", err)
	}
}

func TestNewPolicyRetainsErrorsFromOverriddenOptions(t *testing.T) {
	policy, err := NewPolicy(
		WithAllowedSchemes("://"),
		WithAllowedSchemes("https"),
		WithMaxRequestSize(-1),
		WithMaxRequestSize(1024),
	)
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	multiErr := requireMultiPolicyConfigurationError(t, err)
	assertPolicyConfigurationErrors(t, multiErr.Errors, []PolicyConfigurationError{
		{
			Option: "WithAllowedSchemes",
			Value:  "://",
			Reason: "must be a valid URL scheme",
		},
		{
			Option: "WithMaxRequestSize",
			Value:  "-1",
			Reason: "must not be negative",
		},
	})
}

func TestNewReturnsNilClientForAggregatedPolicyConfigurationError(t *testing.T) {
	client, err := New(
		WithMaxRequestSize(-1),
		WithMaxResponseSize(-2),
	)
	if client != nil {
		t.Fatalf("expected no client, got %T", client)
	}
	requireMultiPolicyConfigurationError(t, err)
}

func TestZeroPolicyIsDenyAll(t *testing.T) {
	var policy Policy

	err := policy.Eval(newTestPolicyGETRequest(t, "https://example.com"))
	policyErr := requirePolicyError(t, err, PolicyTargetRequest)
	if policyErr.Subject != `method "GET"` {
		t.Fatalf("expected zero policy to deny the normalized method, got %#v", policyErr)
	}
}

func requirePolicyConfigurationError(t *testing.T, err error) *PolicyConfigurationError {
	t.Helper()

	if err == nil {
		t.Fatal("expected policy configuration error")
	}
	if !errors.Is(err, ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected ErrInvalidPolicyConfiguration, got %v", err)
	}

	var configErr *PolicyConfigurationError
	if !errors.As(err, &configErr) {
		t.Fatalf("expected PolicyConfigurationError, got %T: %v", err, err)
	}

	return configErr
}

func requireMultiPolicyConfigurationError(t *testing.T, err error) *MultiPolicyConfigurationError {
	t.Helper()

	if err == nil {
		t.Fatal("expected multiple policy configuration errors")
	}
	if !errors.Is(err, ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected ErrInvalidPolicyConfiguration, got %v", err)
	}

	var multiErr *MultiPolicyConfigurationError
	if !errors.As(err, &multiErr) {
		t.Fatalf("expected MultiPolicyConfigurationError, got %T: %v", err, err)
	}
	if len(multiErr.Errors) < 2 {
		t.Fatalf("expected aggregate with at least two children, got %#v", multiErr)
	}

	return multiErr
}

func assertPolicyConfigurationErrors(
	t *testing.T,
	got []*PolicyConfigurationError,
	want []PolicyConfigurationError,
) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("unexpected configuration error count: got %d (%#v), want %d", len(got), got, len(want))
	}

	for index := range want {
		if got[index] == nil {
			t.Fatalf("configuration error %d is nil", index)
		}
		if actual := *got[index]; actual != want[index] {
			t.Fatalf(
				"unexpected configuration error %d: got %#v, want %#v",
				index,
				actual,
				want[index],
			)
		}
	}
}
