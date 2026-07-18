package http

import (
	"errors"
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
		if err := policy.Eval(&Request{URL: rawURL}); err != nil {
			t.Fatalf("expected canonical host %q to be allowed: %v", rawURL, err)
		}
	}

	requirePolicyError(
		t,
		policy.Eval(&Request{URL: "https://sub.example.com"}),
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
	if err := requestPolicy.Eval(&Request{
		URL:  "https://example.com",
		Body: []byte("four"),
	}); err != nil {
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
	if err := policy.Eval(&Request{URL: "https://example.com"}); err != nil {
		t.Fatalf("nil option changed defaults: %v", err)
	}

	methodsCleared := newTestPolicy(t, WithAllowedMethods())
	requirePolicyError(
		t,
		methodsCleared.Eval(&Request{URL: "https://example.com"}),
		PolicyTargetRequest,
	)

	schemesCleared := newTestPolicy(t, WithAllowedSchemes())
	requirePolicyError(
		t,
		schemesCleared.Eval(&Request{URL: "https://example.com"}),
		PolicyTargetRequest,
	)
}

func TestNewPolicyReportsFirstMalformedOption(t *testing.T) {
	policy, err := NewPolicy(
		WithAllowedSchemes("://"),
		WithMaxResponseSize(-1),
	)
	if policy != nil {
		t.Fatalf("expected no policy, got %#v", policy)
	}

	configErr := requirePolicyConfigurationError(t, err)
	if configErr.Option != "WithAllowedSchemes" || configErr.Value != "://" {
		t.Fatalf("expected first malformed option metadata, got %#v", configErr)
	}
}

func TestZeroPolicyIsDenyAll(t *testing.T) {
	var policy Policy

	err := policy.Eval(&Request{URL: "https://example.com"})
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
