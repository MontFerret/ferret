package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"net/url"
	"strings"
	"testing"
)

func TestNewPolicyRejectsTransportControlledDefaultHeaders(t *testing.T) {
	for _, header := range transportControlledRequestHeadersForTest() {
		t.Run(header, func(t *testing.T) {
			policy, err := NewPolicy(WithDefaultHeader(strings.ToLower(header), "configured-secret"))
			if policy != nil {
				t.Fatalf("expected reserved default %q to return no policy", header)
			}

			configErr := requirePolicyConfigurationError(t, err)
			if strings.Contains(err.Error(), "configured-secret") || strings.Contains(configErr.Value, "configured-secret") {
				t.Fatalf("configuration error leaked default value: %v", err)
			}
		})
	}
}

func TestClientRejectsMalformedRequestHeadersBeforeTransport(t *testing.T) {
	tests := []struct {
		name   string
		header string
		value  string
	}{
		{name: "blank name", header: "", value: "value"},
		{name: "space in name", header: "Bad Header", value: "value"},
		{name: "separator in name", header: "Bad@Header", value: "value"},
		{name: "carriage return", header: "Authorization", value: "Bearer secret-value\rInjected: true"},
		{name: "line feed", header: "Authorization", value: "Bearer secret-value\nInjected: true"},
		{name: "nul", header: "X-Test", value: "safe\x00unsafe"},
		{name: "control", header: "X-Test", value: "safe\x1funsafe"},
		{name: "delete", header: "X-Test", value: "safe\x7funsafe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			client := newDefaultHTTPClient(newTestPolicy(t), stdhttp.Client{
				Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
					called = true
					return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
				})},
			)

			_, err := client.Do(context.Background(), &Request{
				URL:     "https://example.com",
				Headers: Headers{tt.header: {tt.value}},
			})
			headerErr := requireHeaderValidationError(t, err)
			if headerErr.Header != tt.header {
				t.Fatalf("expected rejected header %q, got %#v", tt.header, headerErr)
			}
			if headerErr.Reason == "" {
				t.Fatalf("expected validation reason, got %#v", headerErr)
			}
			if called {
				t.Fatal("malformed header reached the transport")
			}
			for _, secret := range []string{"secret-value", "Injected: true"} {
				if strings.Contains(err.Error(), secret) {
					t.Fatalf("header error leaked request value: %v", err)
				}
			}
		})
	}
}

func TestClientAllowsHorizontalTabInRequestAndDefaultHeaders(t *testing.T) {
	policy := newTestPolicy(t, WithDefaultHeader("X-Default", "left\tright"))
	client := newDefaultHTTPClient(policy, stdhttp.Client{
		Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
			if got := req.Header.Get("X-Default"); got != "left\tright" {
				t.Fatalf("expected default HTAB value, got %q", got)
			}
			if got := req.Header.Get("X-Request"); got != "one\ttwo" {
				t.Fatalf("expected request HTAB value, got %q", got)
			}
			return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
		})},
	)

	if _, err := client.Do(context.Background(), &Request{
		URL:     "https://example.com",
		Headers: Headers{"x-request": {"one\ttwo"}},
	}); err != nil {
		t.Fatalf("expected HTAB values to be accepted: %v", err)
	}
}

func TestClientRejectsTransportControlledHeadersBeforeTransport(t *testing.T) {
	for _, header := range transportControlledRequestHeadersForTest() {
		t.Run(header, func(t *testing.T) {
			called := false
			client := newDefaultHTTPClient(newTestPolicy(t), stdhttp.Client{
				Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
					called = true
					return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
				})},
			)

			_, err := client.Do(context.Background(), &Request{
				URL:     "https://example.com",
				Headers: Headers{strings.ToLower(header): {"value"}},
			})
			policyErr := requirePolicyError(t, err, PolicyTargetRequest)
			if policyErr.Reason != "request header is reserved for the transport" {
				t.Fatalf("unexpected reserved-header error: %#v", policyErr)
			}
			if called {
				t.Fatalf("transport-controlled header %q reached transport", header)
			}
		})
	}
}

func TestClientValidatesRedirectHeaders(t *testing.T) {
	redirectURL, err := url.Parse("https://example.com/next")
	if err != nil {
		t.Fatalf("parse redirect URL: %v", err)
	}

	client := newDefaultHTTPClient(newTestPolicy(t), stdhttp.Client{})

	for _, header := range transportControlledRequestHeadersForTest() {
		t.Run("reserved_"+header, func(t *testing.T) {
			err := client.checkRedirect(&stdhttp.Request{
				Method: stdhttp.MethodGet,
				URL:    redirectURL,
				Header: stdhttp.Header{strings.ToLower(header): {"value"}},
			}, []*stdhttp.Request{{}})
			policyErr := requirePolicyError(t, err, PolicyTargetRedirect)
			if policyErr.Reason != "request header is reserved for the transport" {
				t.Fatalf("unexpected redirect header error: %#v", policyErr)
			}
		})
	}

	t.Run("invalid value", func(t *testing.T) {
		err := client.checkRedirect(&stdhttp.Request{
			Method: stdhttp.MethodGet,
			URL:    redirectURL,
			Header: stdhttp.Header{"Authorization": {"Bearer redirect-secret\r\nInjected: true"}},
		}, []*stdhttp.Request{{}})
		requireHeaderValidationError(t, err)
		if strings.Contains(err.Error(), "redirect-secret") || strings.Contains(err.Error(), "Injected: true") {
			t.Fatalf("redirect header error leaked value: %v", err)
		}
	})
}

func TestHeaderValidationErrorSurvivesClientWrapping(t *testing.T) {
	client := newTestClient(t)
	_, err := client.Do(context.Background(), &Request{
		URL:     "https://example.com",
		Headers: Headers{"Authorization": {"Bearer wrapped-secret\nInjected: true"}},
	})
	requireHeaderValidationError(t, err)
	if strings.Contains(err.Error(), "wrapped-secret") {
		t.Fatalf("wrapped header validation error leaked value: %v", err)
	}
}

func requireHeaderValidationError(t *testing.T, err error) *HeaderValidationError {
	t.Helper()

	if err == nil {
		t.Fatal("expected header validation error")
	}
	if errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected structural header error, got policy denial: %v", err)
	}

	var headerErr *HeaderValidationError
	if !errors.As(err, &headerErr) {
		t.Fatalf("expected HeaderValidationError, got %T: %v", err, err)
	}

	return headerErr
}

func transportControlledRequestHeadersForTest() []string {
	return []string{
		"Connection",
		"Content-Length",
		"Host",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Proxy-Connection",
		"TE",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
	}
}
