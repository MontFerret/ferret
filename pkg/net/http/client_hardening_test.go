package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientNormalizesConfiguredMethod(t *testing.T) {
	var seenMethod string
	policy := NewPolicy(WithAllowedMethods("custom-method"))
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
			seenMethod = req.Method
			return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
		})},
	}

	_, err := client.Do(context.Background(), &Request{
		Method: " custom-method ",
		URL:    "https://example.com",
	})
	if err != nil {
		t.Fatalf("expected configured custom method to succeed, got %v", err)
	}
	if seenMethod != "CUSTOM-METHOD" {
		t.Fatalf("expected normalized method %q, got %q", "CUSTOM-METHOD", seenMethod)
	}
}

func TestClientRejectsCredentialedRedirect(t *testing.T) {
	requested := make(map[string]int)
	policy := NewPolicy()
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
			requested[req.URL.String()]++
			return responseWithBody(
				stdhttp.StatusFound,
				"",
				stdhttp.Header{"Location": {"https://user:password@1.1.1.1/secret"}},
			), nil
		})},
	}

	_, err := client.Do(context.Background(), &Request{URL: "https://1.1.1.1/start"})
	policyErr := requirePolicyError(t, err, PolicyTargetRedirect)
	if policyErr.Subject != "URL credentials" {
		t.Fatalf("unexpected redirect policy subject: %q", policyErr.Subject)
	}
	if strings.Contains(err.Error(), "user:password@") || strings.Contains(err.Error(), "password@") {
		t.Fatalf("redirect policy error leaked URL credentials: %v", err)
	}
	if requested["https://user:password@1.1.1.1/secret"] != 0 {
		t.Fatalf("credentialed redirect was requested: %v", requested)
	}
}

func TestClientAppliesMethodPolicyToRedirects(t *testing.T) {
	t.Run("rewritten method denied", func(t *testing.T) {
		roundTrips := 0
		policy := NewPolicy(WithAllowedMethods(stdhttp.MethodPost))
		client := &defaultHTTPClient{
			policy: policy,
			client: stdhttp.Client{Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
				roundTrips++
				if roundTrips == 1 {
					if req.Method != stdhttp.MethodPost {
						t.Fatalf("expected initial POST, got %q", req.Method)
					}
					return responseWithBody(
						stdhttp.StatusFound,
						"",
						stdhttp.Header{"Location": {"/next"}},
					), nil
				}

				return responseWithBody(stdhttp.StatusOK, "unexpected", nil), nil
			})},
		}

		_, err := client.Do(context.Background(), &Request{
			Method: stdhttp.MethodPost,
			URL:    "https://example.com/start",
		})
		policyErr := requirePolicyError(t, err, PolicyTargetRedirect)
		if policyErr.Subject != `method "GET"` {
			t.Fatalf("unexpected redirect method subject: %q", policyErr.Subject)
		}
		if roundTrips != 1 {
			t.Fatalf("expected redirect to be rejected before round trip, got %d requests", roundTrips)
		}
	})

	t.Run("preserved method allowed", func(t *testing.T) {
		roundTrips := 0
		policy := NewPolicy(WithAllowedMethods(stdhttp.MethodPost))
		client := &defaultHTTPClient{
			policy: policy,
			client: stdhttp.Client{Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
				roundTrips++
				if req.Method != stdhttp.MethodPost {
					t.Fatalf("expected POST on round trip %d, got %q", roundTrips, req.Method)
				}
				if roundTrips == 1 {
					return responseWithBody(
						stdhttp.StatusTemporaryRedirect,
						"",
						stdhttp.Header{"Location": {"/next"}},
					), nil
				}

				return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
			})},
		}

		res, err := client.Do(context.Background(), &Request{
			Method: stdhttp.MethodPost,
			URL:    "https://example.com/start",
			Body:   []byte("payload"),
		})
		if err != nil {
			t.Fatalf("expected preserved redirect method to be allowed, got %v", err)
		}
		if string(res.Body) != "ok" || roundTrips != 2 {
			t.Fatalf("unexpected redirect result body=%q requests=%d", res.Body, roundTrips)
		}
	})
}

func TestClientAppliesHeaderPolicyToRedirects(t *testing.T) {
	roundTrips := 0
	policy := NewPolicy(WithBlockedRequestHeaders("Referer"))
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			roundTrips++
			if roundTrips == 1 {
				return responseWithBody(
					stdhttp.StatusFound,
					"",
					stdhttp.Header{"Location": {"/next"}},
				), nil
			}

			return responseWithBody(stdhttp.StatusOK, "unexpected", nil), nil
		})},
	}

	_, err := client.Do(context.Background(), &Request{URL: "https://example.com/start"})
	policyErr := requirePolicyError(t, err, PolicyTargetRedirect)
	if policyErr.Subject != `header "Referer"` {
		t.Fatalf("unexpected redirect header subject: %q", policyErr.Subject)
	}
	if roundTrips != 1 {
		t.Fatalf("expected redirected header to be rejected before round trip, got %d requests", roundTrips)
	}
}

func TestClientRejectsNonASCIIRedirect(t *testing.T) {
	policy := NewPolicy()
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			return responseWithBody(
				stdhttp.StatusFound,
				"",
				stdhttp.Header{"Location": {"https://K.example/secret"}},
			), nil
		})},
	}

	_, err := client.Do(context.Background(), &Request{URL: "https://1.1.1.1/start"})
	policyErr := requirePolicyError(t, err, PolicyTargetRedirect)
	if policyErr.Reason != "internationalized hostnames must use ASCII/punycode" {
		t.Fatalf("unexpected redirect host reason: %q", policyErr.Reason)
	}
}

func TestClientRejectsBlockedHeaderBeforeRoundTrip(t *testing.T) {
	called := false
	policy := NewPolicy(
		WithDefaultHeader("Authorization", "Bearer trusted"),
		WithBlockedRequestHeaders("Authorization"),
	)
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			called = true
			return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
		})},
	}

	_, err := client.Do(context.Background(), &Request{
		URL:     "https://example.com",
		Headers: Headers{"authorization": {"Bearer secret"}},
	})
	policyErr := requirePolicyError(t, err, PolicyTargetRequest)
	if policyErr.Subject != `header "Authorization"` {
		t.Fatalf("unexpected blocked-header subject: %q", policyErr.Subject)
	}
	if called {
		t.Fatal("blocked request reached the round tripper")
	}
}

func TestClientAppliesDefaultHeadersByPresence(t *testing.T) {
	tests := []struct {
		headers    Headers
		name       string
		wantValues []string
		wantExists bool
	}{
		{name: "missing", wantExists: true, wantValues: []string{"default"}},
		{name: "explicit empty value", headers: Headers{"X-Value": {""}}, wantExists: true, wantValues: []string{""}},
		{name: "explicit nil values", headers: Headers{"X-Value": nil}, wantExists: true, wantValues: nil},
		{name: "request override", headers: Headers{"X-Value": {"request"}}, wantExists: true, wantValues: []string{"request"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				gotValues []string
				gotExists bool
			)
			policy := NewPolicy(WithDefaultHeader("x-value", "default"))
			client := &defaultHTTPClient{
				policy: policy,
				client: stdhttp.Client{Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
					gotValues, gotExists = req.Header["X-Value"]
					gotValues = append([]string(nil), gotValues...)
					return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
				})},
			}

			_, err := client.Do(context.Background(), &Request{
				URL:     "https://example.com",
				Headers: tt.headers,
			})
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if gotExists != tt.wantExists {
				t.Fatalf("expected header presence %t, got %t", tt.wantExists, gotExists)
			}
			if !equalStrings(gotValues, tt.wantValues) {
				t.Fatalf("expected header values %v, got %v", tt.wantValues, gotValues)
			}
		})
	}
}

func TestClientBlockedAndReservedHeadersTakePrecedenceOverDefaults(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{name: "blocked", header: "Authorization"},
		{name: "reserved", header: "Connection"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := []PolicyOption{WithDefaultHeader(tt.header, "trusted")}
			if tt.name == "blocked" {
				options = append(options, WithBlockedRequestHeaders(tt.header))
			}

			policy := NewPolicy(options...)
			client := &defaultHTTPClient{
				policy: policy,
				client: stdhttp.Client{Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
					if _, exists := req.Header[stdhttp.CanonicalHeaderKey(tt.header)]; exists {
						t.Fatalf("expected contradictory default %q to be omitted", tt.header)
					}
					return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
				})},
			}

			if _, err := client.Do(context.Background(), &Request{URL: "https://example.com"}); err != nil {
				t.Fatalf("request failed: %v", err)
			}
		})
	}
}

func TestClientResponseHeaderLimit(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.URL.Path == "/large" {
			w.Header().Set("X-Large", strings.Repeat("x", 2<<10))
		} else {
			w.Header().Set("X-Small", "ok")
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := New(
		WithAllowLocalhost(true),
		WithMaxResponseHeaderSize(1<<10),
	)
	if _, err := client.Do(context.Background(), &Request{URL: server.URL + "/small"}); err != nil {
		t.Fatalf("small response headers failed: %v", err)
	}

	_, err := client.Do(context.Background(), &Request{URL: server.URL + "/large"})
	if err == nil || !strings.Contains(err.Error(), "response headers exceeded 1024 bytes") {
		t.Fatalf("expected response header limit error, got %v", err)
	}
	if errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected transport limit error, got policy denial: %v", err)
	}
}

func TestClientSecureTransportDefaultsAndOverrides(t *testing.T) {
	tests := []struct {
		name           string
		options        []PolicyOption
		wantHeaderSize int64
		wantTimeout    time.Duration
	}{
		{name: "defaults", wantHeaderSize: defaultMaxResponseHeaderSize, wantTimeout: defaultTimeout},
		{name: "custom", options: []PolicyOption{WithMaxResponseHeaderSize(2048), WithTimeout(time.Second)}, wantHeaderSize: 2048, wantTimeout: time.Second},
		{name: "zero restores header default and disables timeout", options: []PolicyOption{WithMaxResponseHeaderSize(0), WithTimeout(0)}, wantHeaderSize: defaultMaxResponseHeaderSize},
		{name: "negative restores header default and disables timeout", options: []PolicyOption{WithMaxResponseHeaderSize(-1), WithTimeout(-1)}, wantHeaderSize: defaultMaxResponseHeaderSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ok := New(tt.options...).(*defaultHTTPClient)
			if !ok {
				t.Fatalf("expected built-in client, got %T", client)
			}
			if client.client.Timeout != tt.wantTimeout {
				t.Fatalf("expected timeout %s, got %s", tt.wantTimeout, client.client.Timeout)
			}

			transport, ok := client.client.Transport.(*stdhttp.Transport)
			if !ok {
				t.Fatalf("expected built-in transport, got %T", client.client.Transport)
			}
			if transport.MaxResponseHeaderBytes != tt.wantHeaderSize {
				t.Fatalf("expected response header limit %d, got %d", tt.wantHeaderSize, transport.MaxResponseHeaderBytes)
			}
			if transport.MaxIdleConns != defaultMaxIdleConnections ||
				transport.MaxIdleConnsPerHost != defaultMaxIdleConnectionsPerHost ||
				transport.MaxConnsPerHost != defaultMaxConnectionsPerHost {
				t.Fatalf(
					"unexpected connection limits: idle=%d idle-per-host=%d per-host=%d",
					transport.MaxIdleConns,
					transport.MaxIdleConnsPerHost,
					transport.MaxConnsPerHost,
				)
			}
		})
	}
}

func TestClientRequestContextCanShortenPolicyTimeout(t *testing.T) {
	policy := NewPolicy(WithTimeout(time.Hour))
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
			<-req.Context().Done()
			return nil, req.Context().Err()
		})},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Do(ctx, &Request{URL: "https://example.com"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected request context deadline, got %v", err)
	}
	if errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected deadline error, got policy denial: %v", err)
	}
}

func TestClientOrdinaryTransportErrorIsNotPolicyDenial(t *testing.T) {
	want := errors.New("transport failed")
	policy := NewPolicy()
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			return nil, want
		})},
	}

	_, err := client.Do(context.Background(), &Request{URL: "https://example.com"})
	if !errors.Is(err, want) {
		t.Fatalf("expected transport error, got %v", err)
	}
	if errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected ordinary transport error, got policy denial: %v", err)
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for idx := range left {
		if left[idx] != right[idx] {
			return false
		}
	}

	return true
}
