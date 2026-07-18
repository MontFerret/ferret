package http

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"
	"reflect"
	"strings"
	"testing"
)

func TestPolicyPrepareAddsOnlyMissingDefaultsAndIsIdempotent(t *testing.T) {
	tests := []struct {
		headers stdhttp.Header
		name    string
		wantKey string
		want    []string
	}{
		{name: "nil headers", wantKey: "X-Value", want: []string{"default"}},
		{name: "missing", headers: stdhttp.Header{"X-Other": {"other"}}, wantKey: "X-Value", want: []string{"default"}},
		{name: "explicit nil values", headers: stdhttp.Header{"X-Value": nil}, wantKey: "X-Value"},
		{name: "explicit empty values", headers: stdhttp.Header{"X-Value": {}}, wantKey: "X-Value", want: []string{}},
		{name: "explicit empty string", headers: stdhttp.Header{"X-Value": {""}}, wantKey: "X-Value", want: []string{""}},
		{name: "explicit value", headers: stdhttp.Header{"X-Value": {"request"}}, wantKey: "X-Value", want: []string{"request"}},
		{name: "case variant is present", headers: stdhttp.Header{"x-value": {"request"}}, wantKey: "x-value", want: []string{"request"}},
	}

	policy := newTestPolicy(t, WithDefaultHeader("X-Value", "default"))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestPolicyGETRequest(t, "https://example.com")
			if tt.headers == nil {
				req.Header = nil
			} else {
				req.Header = make(stdhttp.Header, len(tt.headers))
				for key, values := range tt.headers {
					if values == nil {
						req.Header[key] = nil
						continue
					}

					req.Header[key] = append(make([]string, 0, len(values)), values...)
				}
			}

			if err := policy.Prepare(req); err != nil {
				t.Fatalf("prepare request: %v", err)
			}
			if got, exists := req.Header[tt.wantKey]; !exists || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %q=%v, got %v (exists=%t)", tt.wantKey, tt.want, got, exists)
			}
			if tt.wantKey == "x-value" {
				if _, exists := req.Header["X-Value"]; exists {
					t.Fatalf("case-equivalent default overwrote existing header: %v", req.Header)
				}
			}

			want := req.Header.Clone()
			if err := policy.Prepare(req); err != nil {
				t.Fatalf("prepare request again: %v", err)
			}
			if !reflect.DeepEqual(req.Header, want) {
				t.Fatalf("Prepare is not idempotent: before=%v after=%v", want, req.Header)
			}
		})
	}
}

func TestPolicyPrepareOnErrorOnlyAddsDefaults(t *testing.T) {
	type contextKey struct{}

	ctx := context.WithValue(context.Background(), contextKey{}, "caller")
	req := newTestPolicyRequest(t, "BAD METHOD", "https://example.com/path")
	req = req.WithContext(ctx)
	req.Header = stdhttp.Header{"X-Existing": {"one", "two"}}
	req.Body = io.NopCloser(strings.NewReader("body"))
	req.ContentLength = 4

	wantMethod := req.Method
	wantURL := *req.URL
	wantExisting := append([]string(nil), req.Header["X-Existing"]...)
	wantBody := req.Body
	wantContentLength := req.ContentLength
	wantContext := req.Context()

	err := newTestPolicy(t, WithDefaultHeader("X-Default", "default")).Prepare(req)
	var methodErr *InvalidMethodError
	if !errors.As(err, &methodErr) {
		t.Fatalf("expected InvalidMethodError, got %T: %v", err, err)
	}
	if got := req.Header.Values("X-Default"); !reflect.DeepEqual(got, []string{"default"}) {
		t.Fatalf("expected default added before validation failure, got %v", got)
	}
	if len(req.Header) != 2 {
		t.Fatalf("expected Prepare to add only one default, got %v", req.Header)
	}
	if !reflect.DeepEqual(req.Header["X-Existing"], wantExisting) {
		t.Fatalf("Prepare changed existing headers: %v", req.Header)
	}
	if req.Method != wantMethod || !reflect.DeepEqual(*req.URL, wantURL) ||
		req.Body != wantBody || req.ContentLength != wantContentLength || req.Context() != wantContext {
		t.Fatalf("Prepare mutated non-header request state: %#v", req)
	}
}

func TestPolicyEvalRejectsNonClientStandardRequestState(t *testing.T) {
	tests := []struct {
		mutate func(*stdhttp.Request)
		name   string
	}{
		{name: "request URI", mutate: func(req *stdhttp.Request) { req.RequestURI = "/server-form" }},
		{name: "mismatched host", mutate: func(req *stdhttp.Request) { req.Host = "other.example" }},
		{name: "host with whitespace", mutate: func(req *stdhttp.Request) { req.Host = " example.com " }},
		{name: "bracketed DNS host", mutate: func(req *stdhttp.Request) { req.Host = "[example.com]" }},
		{name: "host with empty port", mutate: func(req *stdhttp.Request) { req.Host = "example.com:" }},
		{name: "close", mutate: func(req *stdhttp.Request) { req.Close = true }},
		{name: "transfer encoding", mutate: func(req *stdhttp.Request) { req.TransferEncoding = []string{"chunked"} }},
		{name: "trailer", mutate: func(req *stdhttp.Request) { req.Trailer = stdhttp.Header{"X-Trailer": {"value"}} }},
	}

	policy := newTestPolicy(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestPolicyGETRequest(t, "https://example.com/path")
			tt.mutate(req)
			requirePolicyError(t, policy.Eval(req), PolicyTargetRequest)
		})
	}

	for _, tt := range []struct {
		host   string
		rawURL string
	}{
		{host: "example.com", rawURL: "https://example.com/path"},
		{host: "EXAMPLE.COM", rawURL: "https://example.com/path"},
		{host: "example.com:443", rawURL: "https://example.com:443/path"},
		{host: "EXAMPLE.COM.:443", rawURL: "https://example.com:443/path"},
	} {
		req := newTestPolicyGETRequest(t, tt.rawURL)
		req.Host = tt.host
		if err := policy.Eval(req); err != nil {
			t.Fatalf("expected authority-equivalent Host %q to pass: %v", tt.host, err)
		}
	}

	req := newTestPolicyGETRequest(t, "https://example.com/path")
	req.URL.Host = "[example.com]"
	req.Host = req.URL.Host
	requirePolicyError(t, policy.Eval(req), PolicyTargetRequest)
}

func TestPolicyEvalValidationOrder(t *testing.T) {
	policy := newTestPolicy(t, WithMaxRequestSize(3))

	t.Run("method before URL", func(t *testing.T) {
		err := policy.Eval(&stdhttp.Request{Method: "BAD METHOD"})
		var methodErr *InvalidMethodError
		if !errors.As(err, &methodErr) {
			t.Fatalf("expected InvalidMethodError, got %T: %v", err, err)
		}
	})

	t.Run("URL before strict state", func(t *testing.T) {
		req := &stdhttp.Request{Method: stdhttp.MethodGet, RequestURI: "/server-form"}
		err := policy.Eval(req)
		var validationErr *URLValidationError
		if !errors.As(err, &validationErr) || validationErr.Field != "url" {
			t.Fatalf("expected URLValidationError for URL, got %T: %v", err, err)
		}
	})

	t.Run("strict state before headers", func(t *testing.T) {
		req := newTestPolicyGETRequest(t, "https://example.com")
		req.RequestURI = "/server-form"
		req.Header = stdhttp.Header{"Bad Header": {"value"}}
		requirePolicyError(t, policy.Eval(req), PolicyTargetRequest)
	})

	t.Run("headers before body", func(t *testing.T) {
		req := newTestPolicyGETRequest(t, "https://example.com")
		req.Header = stdhttp.Header{"Bad Header": {"value"}}
		req.Body = io.NopCloser(strings.NewReader("unknown"))
		req.ContentLength = -1

		err := policy.Eval(req)
		var headerErr *HeaderValidationError
		if !errors.As(err, &headerErr) {
			t.Fatalf("expected HeaderValidationError, got %T: %v", err, err)
		}
	})

	t.Run("body after valid headers", func(t *testing.T) {
		req := newTestPolicyGETRequest(t, "https://example.com")
		req.Header = stdhttp.Header{"X-Test": {"value"}}
		req.Body = io.NopCloser(strings.NewReader("unknown"))
		req.ContentLength = -1

		err := policy.Eval(req)
		var lengthErr *RequestBodyLengthError
		if !errors.As(err, &lengthErr) {
			t.Fatalf("expected RequestBodyLengthError, got %T: %v", err, err)
		}
	})
}
