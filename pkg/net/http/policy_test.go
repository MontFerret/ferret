package http_test

import (
	"errors"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

func TestNewPoliciesDefaultsAndOptions(t *testing.T) {
	policies := ferrethttp.NewPolicies(
		ferrethttp.WithTimeout(time.Second),
		ferrethttp.WithAllowedSchemes("HTTPS", " "),
		ferrethttp.WithAllowedHosts("Example.COM"),
		ferrethttp.WithBlockedHosts("Blocked.EXAMPLE"),
		ferrethttp.WithBlockedRequestHeaders("x-token"),
		ferrethttp.WithDefaultHeader("x-default", "value"),
		ferrethttp.WithAllowLocalhost(false),
		ferrethttp.WithAllowPrivateNetworks(false),
		ferrethttp.WithFollowRedirects(false),
		ferrethttp.WithMaxRedirects(2),
		ferrethttp.WithMaxRequestSize(3),
		ferrethttp.WithMaxResponseSize(4),
	)

	if policies.Timeout != time.Second {
		t.Fatalf("expected timeout to be configured, got %s", policies.Timeout)
	}
	if policies.FollowRedirects {
		t.Fatal("expected redirects to be disabled")
	}
	if policies.MaxRedirects != 2 || policies.MaxRequestSize != 3 || policies.MaxResponseSize != 4 {
		t.Fatalf("expected configured limits, got redirects=%d request=%d response=%d",
			policies.MaxRedirects,
			policies.MaxRequestSize,
			policies.MaxResponseSize,
		)
	}
	if policies.AllowLocalhost || policies.AllowPrivateNetworks {
		t.Fatal("expected localhost and private networks to be disabled")
	}
	if got := policies.AllowedSchemes; len(got) != 1 || got[0] != "https" {
		t.Fatalf("expected normalized allowed scheme, got %v", got)
	}
	if got := policies.AllowedHosts; len(got) != 1 || got[0] != "example.com" {
		t.Fatalf("expected normalized allowed host, got %v", got)
	}
	if got := policies.BlockedHosts; len(got) != 1 || got[0] != "blocked.example" {
		t.Fatalf("expected normalized blocked host, got %v", got)
	}
	if got := policies.BlockedRequestHeaders; len(got) != 1 || got[0] != "X-Token" {
		t.Fatalf("expected normalized blocked request header, got %v", got)
	}
	if got := policies.DefaultHeaders["X-Default"]; got != "value" {
		t.Fatalf("expected normalized default header, got %q", got)
	}
}

func TestNewPoliciesDefaultValues(t *testing.T) {
	policies := ferrethttp.NewPolicies()

	if !policies.FollowRedirects {
		t.Fatal("expected redirects to be enabled by default")
	}
	if !policies.AllowLocalhost || !policies.AllowPrivateNetworks {
		t.Fatal("expected localhost and private networks to be allowed by default")
	}
	if got := policies.AllowedSchemes; len(got) != 2 || got[0] != "http" || got[1] != "https" {
		t.Fatalf("expected default http/https schemes, got %v", got)
	}
}

func TestPoliciesEval(t *testing.T) {
	tests := []struct {
		req      *ferrethttp.Request
		name     string
		want     string
		policies ferrethttp.Policies
	}{
		{
			name:     "valid default get",
			policies: ferrethttp.NewPolicies(),
			req:      &ferrethttp.Request{URL: "http://example.com"},
		},
		{
			name:     "nil request",
			policies: ferrethttp.NewPolicies(),
			req:      nil,
			want:     ferrethttp.ErrNilRequest.Error(),
		},
		{
			name:     "invalid method",
			policies: ferrethttp.NewPolicies(),
			req:      &ferrethttp.Request{Method: "BAD METHOD", URL: "http://example.com"},
			want:     "invalid method",
		},
		{
			name:     "missing url",
			policies: ferrethttp.NewPolicies(),
			req:      &ferrethttp.Request{},
			want:     "url is required",
		},
		{
			name:     "missing scheme",
			policies: ferrethttp.NewPolicies(),
			req:      &ferrethttp.Request{URL: "example.com"},
			want:     "url scheme is required",
		},
		{
			name:     "missing host",
			policies: ferrethttp.NewPolicies(),
			req:      &ferrethttp.Request{URL: "http:///path"},
			want:     "url host is required",
		},
		{
			name:     "disallowed scheme",
			policies: ferrethttp.NewPolicies(ferrethttp.WithAllowedSchemes("https")),
			req:      &ferrethttp.Request{URL: "http://example.com"},
			want:     "scheme",
		},
		{
			name:     "blocked host",
			policies: ferrethttp.NewPolicies(ferrethttp.WithBlockedHosts("example.com")),
			req:      &ferrethttp.Request{URL: "http://example.com"},
			want:     "blocked",
		},
		{
			name:     "not allowed host",
			policies: ferrethttp.NewPolicies(ferrethttp.WithAllowedHosts("allowed.example")),
			req:      &ferrethttp.Request{URL: "http://other.example"},
			want:     "not allowed",
		},
		{
			name:     "localhost",
			policies: ferrethttp.NewPolicies(ferrethttp.WithAllowLocalhost(false)),
			req:      &ferrethttp.Request{URL: "http://127.0.0.1"},
			want:     "localhost is not allowed",
		},
		{
			name:     "private network",
			policies: ferrethttp.NewPolicies(ferrethttp.WithAllowPrivateNetworks(false)),
			req:      &ferrethttp.Request{URL: "http://10.0.0.1"},
			want:     "private network",
		},
		{
			name:     "request body limit",
			policies: ferrethttp.NewPolicies(ferrethttp.WithMaxRequestSize(3)),
			req:      &ferrethttp.Request{Method: stdhttp.MethodPost, URL: "http://example.com", Body: []byte("four")},
			want:     "request body exceeds limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policies.Eval(tt.req)
			if tt.want == "" {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			if tt.req == nil {
				if !errors.Is(err, ferrethttp.ErrNilRequest) {
					t.Fatalf("expected ErrNilRequest, got %v", err)
				}
				return
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
		})
	}
}

func TestPoliciesEvalDoesNotMutateRequest(t *testing.T) {
	policies := ferrethttp.NewPolicies(
		ferrethttp.WithDefaultHeader("X-Default", "default"),
		ferrethttp.WithBlockedRequestHeaders("X-Blocked"),
	)
	req := &ferrethttp.Request{
		URL: "HTTP://EXAMPLE.COM/path",
		Headers: ferrethttp.Headers{
			"X-Blocked": {"secret"},
		},
	}

	if err := policies.Eval(req); err != nil {
		t.Fatalf("expected policy evaluation to pass, got %v", err)
	}

	if req.Method != "" {
		t.Fatalf("expected method to remain empty, got %q", req.Method)
	}
	if req.URL != "HTTP://EXAMPLE.COM/path" {
		t.Fatalf("expected URL to remain unchanged, got %q", req.URL)
	}
	if got := req.Headers["X-Blocked"]; len(got) != 1 || got[0] != "secret" {
		t.Fatalf("expected blocked header to remain unchanged, got %v", got)
	}
	if _, ok := req.Headers["X-Default"]; ok {
		t.Fatalf("expected default header not to be added during Eval")
	}
}
