package http

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ziflex/go-options"
)

func TestNewPolicyRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		options []PolicyOption
	}{
		{name: "invalid method", field: "allowed methods", options: []PolicyOption{WithAllowedMethods("GET", "BAD METHOD")}},
		{name: "blank method", field: "allowed methods", options: []PolicyOption{WithAllowedMethods("")}},
		{name: "invalid scheme", field: "allowed schemes", options: []PolicyOption{WithAllowedSchemes("https", "://")}},
		{name: "blank scheme", field: "allowed schemes", options: []PolicyOption{WithAllowedSchemes("")}},
		{name: "blank allowed host", field: "allowed hosts", options: []PolicyOption{WithAllowedHosts("")}},
		{name: "blank blocked host", field: "blocked hosts", options: []PolicyOption{WithBlockedHosts(" ")}},
		{name: "malformed blocked host", field: "blocked hosts", options: []PolicyOption{WithBlockedHosts("example..com")}},
		{name: "invalid blocked header", field: "blocked request headers", options: []PolicyOption{WithBlockedRequestHeaders("Bad Header")}},
		{name: "blank blocked header", field: "blocked request headers", options: []PolicyOption{WithBlockedRequestHeaders("")}},
		{name: "blank default header", field: "default header", options: []PolicyOption{WithDefaultHeader("", "value")}},
		{name: "invalid default header", field: "default header", options: []PolicyOption{WithDefaultHeader("Bad Header", "value")}},
		{name: "invalid default value", field: "default header", options: []PolicyOption{WithDefaultHeader("X-Test", "safe\r\nInjected: true")}},
		{name: "negative timeout", field: "timeout", options: []PolicyOption{WithTimeout(-time.Nanosecond)}},
		{name: "negative redirects", field: "max redirects", options: []PolicyOption{WithMaxRedirects(-1)}},
		{name: "negative request limit", field: "max request size", options: []PolicyOption{WithMaxRequestSize(-1)}},
		{name: "negative response limit", field: "max response size", options: []PolicyOption{WithMaxResponseSize(-1)}},
		{name: "negative response header limit", field: "max response header size", options: []PolicyOption{WithMaxResponseHeaderSize(-1)}},
		{
			name:  "default then blocked conflict",
			field: "default header",
			options: []PolicyOption{
				WithDefaultHeader("Authorization", "Bearer configured-secret"),
				WithBlockedRequestHeaders("authorization"),
			},
		},
		{
			name:  "blocked then default conflict",
			field: "default header",
			options: []PolicyOption{
				WithBlockedRequestHeaders("AUTHORIZATION"),
				WithDefaultHeader("authorization", "Bearer configured-secret"),
			},
		},
		{
			name:  "case-equivalent defaults conflict",
			field: "default headers",
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

			validationErr := requireOptionValidationError(t, err)
			if validationErr.Field != tt.field || validationErr.Reason == nil {
				t.Fatalf("validation error = %#v, want field %q and a reason", validationErr, tt.field)
			}
			for _, secret := range []string{"configured-secret", "safe\r\nInjected: true"} {
				if strings.Contains(err.Error(), secret) || strings.Contains(validationErr.Value, secret) {
					t.Fatalf("validation error leaked a header value: %v", err)
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
			requireOptionValidationError(t, err)
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
	requireOptionValidationError(t, err)
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
	body, err := responsePolicy.ReadResponseBody(strings.NewReader("four"))
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

func TestNewPolicyReturnsValidationErrorForSingleFailure(t *testing.T) {
	policy, err := NewPolicy(WithMaxRequestSize(-1))
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	validationErr := requireOptionValidationError(t, err)
	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "max request size",
			Value:  "-1",
			Reason: errors.New("must not be negative"),
		},
	})
	if validationErr.Field != "max request size" {
		t.Fatalf("validation field = %q, want %q", validationErr.Field, "max request size")
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

	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "allowed schemes",
			Value:  "://",
			Reason: errors.New("must be a valid URL scheme"),
		},
		{
			Field:  "allowed schemes",
			Value:  "bad scheme",
			Reason: errors.New("must be a valid URL scheme"),
		},
		{
			Field:  "allowed schemes",
			Reason: errors.New("must be a non-empty URL scheme"),
		},
		{
			Field:  "allowed schemes",
			Value:  "://",
			Reason: errors.New("must be a valid URL scheme"),
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

	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "default headers",
			Value:  "A Bad Header",
			Reason: errors.New("name is not a valid HTTP field-name token"),
		},
		{
			Field:  "default headers",
			Value:  "Z Bad Header",
			Reason: errors.New("name is not a valid HTTP field-name token"),
		},
	})

	for _, secret := range []string{"first-secret", "last-secret"} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("joined validation error leaked default value %q: %v", secret, err)
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

	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "allowed methods",
			Value:  "BAD METHOD",
			Reason: errors.New("must be a non-empty HTTP method token"),
		},
		{
			Field:  "allowed schemes",
			Value:  "bad scheme",
			Reason: errors.New("must be a valid URL scheme"),
		},
		{
			Field:  "allowed hosts",
			Reason: errors.New("must not be blank"),
		},
		{
			Field:  "blocked hosts",
			Value:  "*.example.com",
			Reason: errors.New("wildcards are not supported"),
		},
		{
			Field:  "blocked request headers",
			Value:  "Bad Header",
			Reason: errors.New("name is not a valid HTTP field-name token"),
		},
		{
			Field:  "default header",
			Value:  "Another Bad Header",
			Reason: errors.New("name is not a valid HTTP field-name token"),
		},
		{
			Field:  "max response size",
			Value:  "-1",
			Reason: errors.New("must not be negative"),
		},
	})

	if strings.Contains(err.Error(), "configured-secret") {
		t.Fatalf("joined validation error leaked a default header value: %v", err)
	}
}

func TestNewPolicyOrdersMixedImmediateAndFinalizationErrors(t *testing.T) {
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

	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "max response size",
			Value:  "-2",
			Reason: errors.New("must not be negative"),
		},
		{
			Field:  "max request size",
			Value:  "-1",
			Reason: errors.New("must not be negative"),
		},
		{
			Field:  "default header",
			Value:  "X-Conflict",
			Reason: errors.New("conflicts with another default for the same header"),
		},
	})

	for _, secret := range []string{firstSecret, secondSecret} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("joined validation error leaked header value %q: %v", secret, err)
		}
		for _, child := range collectOptionValidationErrors(err) {
			if strings.Contains(child.Error(), secret) || strings.Contains(child.Value, secret) {
				t.Fatalf("validation error leaked header value %q: %#v", secret, child)
			}
		}
	}
}

func TestPolicyValidationErrorsSupportStandardErrorDiscovery(t *testing.T) {
	_, err := NewPolicy(
		WithMaxRequestSize(-1),
		WithMaxResponseSize(-2),
	)
	requireOptionValidationError(t, err)
	wrapped := fmt.Errorf("construct HTTP policy: %w", err)

	if !errors.Is(wrapped, ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected wrapped validation errors to match ErrInvalidPolicyConfiguration: %v", wrapped)
	}

	var discovered options.ValidationError
	if !errors.As(wrapped, &discovered) {
		t.Fatalf("expected wrapped error to expose options.ValidationError, got %T: %v", wrapped, wrapped)
	}
	if discovered.Field != "max request size" || discovered.Value != "-1" {
		t.Fatalf("first discovered validation error = %+v", discovered)
	}

	rendered := wrapped.Error()
	if count := strings.Count(rendered, ErrInvalidPolicyConfiguration.Error()); count != 1 {
		t.Fatalf("expected the shared sentinel once, got %d occurrences in %q", count, rendered)
	}
	for _, child := range collectOptionValidationErrors(wrapped) {
		if !strings.Contains(rendered, child.Field) || !strings.Contains(rendered, child.Reason.Error()) {
			t.Fatalf("joined error text omitted child %#v: %q", child, rendered)
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

	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "timeout",
			Value:  "-1ns",
			Reason: errors.New("must not be negative"),
		},
		{
			Field:  "default header",
			Value:  "X-Blocked-First",
			Reason: errors.New("default header is also configured as blocked"),
		},
		{
			Field:  "default header",
			Value:  "X-Default-First",
			Reason: errors.New("default header is also configured as blocked"),
		},
	})

	if strings.Contains(err.Error(), secret) {
		t.Fatalf("joined validation error leaked a default header value: %v", err)
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

	validationErr := requireOptionValidationError(t, err)
	if validationErr.Field != "default header" || validationErr.Reason.Error() != "value contains a newline" {
		t.Fatalf("unexpected validation error: %#v", validationErr)
	}
	if got := collectOptionValidationErrors(err); len(got) != 1 {
		t.Fatalf("invalid default produced a spurious blocked-header conflict: %#v", got)
	}
	if strings.Contains(err.Error(), secret) || strings.Contains(validationErr.Value, secret) {
		t.Fatalf("validation error leaked a default header value: %v", err)
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

	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "allowed schemes",
			Value:  "://",
			Reason: errors.New("must be a valid URL scheme"),
		},
		{
			Field:  "max request size",
			Value:  "-1",
			Reason: errors.New("must not be negative"),
		},
	})
}

func TestNewReturnsNilClientForJoinedPolicyValidationErrors(t *testing.T) {
	client, err := New(
		WithMaxRequestSize(-1),
		WithMaxResponseSize(-2),
	)
	if client != nil {
		t.Fatalf("expected no client, got %T", client)
	}
	assertOptionValidationErrors(t, err, []options.ValidationError{
		{
			Field:  "max request size",
			Value:  "-1",
			Reason: errors.New("must not be negative"),
		},
		{
			Field:  "max response size",
			Value:  "-2",
			Reason: errors.New("must not be negative"),
		},
	})
}

func TestZeroPolicyIsDenyAll(t *testing.T) {
	var policy Policy

	err := policy.Eval(newTestPolicyGETRequest(t, "https://example.com"))
	policyErr := requirePolicyError(t, err, PolicyTargetRequest)
	if policyErr.Subject != `method "GET"` {
		t.Fatalf("expected zero policy to deny the normalized method, got %#v", policyErr)
	}
}

func requireOptionValidationError(t *testing.T, err error) options.ValidationError {
	t.Helper()

	if err == nil {
		t.Fatal("expected policy validation error")
	}
	if !errors.Is(err, ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected ErrInvalidPolicyConfiguration, got %v", err)
	}

	var validationErr options.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected options.ValidationError, got %T: %v", err, err)
	}

	return validationErr
}

func assertOptionValidationErrors(t *testing.T, err error, want []options.ValidationError) {
	t.Helper()

	requireOptionValidationError(t, err)
	got := collectOptionValidationErrors(err)
	if len(got) != len(want) {
		t.Fatalf("unexpected validation error count: got %d (%#v), want %d", len(got), got, len(want))
	}

	for index := range want {
		if got[index].Field != want[index].Field ||
			got[index].Value != want[index].Value ||
			errorText(got[index].Reason) != errorText(want[index].Reason) {
			t.Fatalf(
				"unexpected validation error %d: got %#v, want %#v",
				index,
				got[index],
				want[index],
			)
		}
	}
}

func collectOptionValidationErrors(err error) []options.ValidationError {
	if err == nil {
		return nil
	}

	switch validationErr := err.(type) {
	case options.ValidationError:
		return []options.ValidationError{validationErr}
	case *options.ValidationError:
		if validationErr == nil {
			return nil
		}

		return []options.ValidationError{*validationErr}
	}

	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		var result []options.ValidationError
		for _, child := range joined.Unwrap() {
			result = append(result, collectOptionValidationErrors(child)...)
		}

		return result
	}

	if wrapped, ok := err.(interface{ Unwrap() error }); ok {
		return collectOptionValidationErrors(wrapped.Unwrap())
	}

	return nil
}

func errorText(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
