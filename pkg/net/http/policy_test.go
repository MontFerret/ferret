package http

import (
	"errors"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"
)

func TestNewPoliciesDefaultsAndOptions(t *testing.T) {
	policies := NewPolicies(
		WithTimeout(time.Second),
		WithAllowedSchemes("HTTPS", " "),
		WithAllowedHosts("Example.COM"),
		WithBlockedHosts("Blocked.EXAMPLE"),
		WithBlockedRequestHeaders("x-token"),
		WithDefaultHeader("x-default", "value"),
		WithAllowLocalhost(false),
		WithAllowPrivateNetworks(false),
		WithAllowLinkLocal(false),
		WithFollowRedirects(false),
		WithMaxRedirects(2),
		WithMaxRequestSize(3),
		WithMaxResponseSize(4),
	)

	if policies.timeout != time.Second {
		t.Fatalf("expected timeout to be configured, got %s", policies.timeout)
	}
	if policies.followRedirects {
		t.Fatal("expected redirects to be disabled")
	}
	if policies.maxRedirects != 2 || policies.maxRequestSize != 3 || policies.maxResponseSize != 4 {
		t.Fatalf("expected configured limits, got redirects=%d request=%d response=%d",
			policies.maxRedirects,
			policies.maxRequestSize,
			policies.maxResponseSize,
		)
	}
	if policies.allowLocalhost || policies.allowPrivateNetworks || policies.allowLinkLocal {
		t.Fatal("expected localhost, private networks, and link-local addresses to be disabled")
	}
	if got := policies.allowedSchemes; len(got) != 1 || got[0] != "https" {
		t.Fatalf("expected normalized allowed scheme, got %v", got)
	}
	if got := policies.allowedHosts; len(got) != 1 || got[0] != "example.com" {
		t.Fatalf("expected normalized allowed host, got %v", got)
	}
	if got := policies.blockedHosts; len(got) != 1 || got[0] != "blocked.example" {
		t.Fatalf("expected normalized blocked host, got %v", got)
	}
	if got := policies.blockedRequestHeaders; len(got) != 1 || got[0] != "X-Token" {
		t.Fatalf("expected normalized blocked request header, got %v", got)
	}
	if got := policies.defaultHeaders["X-Default"]; got != "value" {
		t.Fatalf("expected normalized default header, got %q", got)
	}
}

func TestNewPoliciesDefaultValues(t *testing.T) {
	policies := NewPolicies()

	if !policies.followRedirects {
		t.Fatal("expected redirects to be enabled by default")
	}
	if policies.allowLocalhost || policies.allowPrivateNetworks || policies.allowLinkLocal {
		t.Fatal("expected localhost, private networks, and link-local addresses to be denied by default")
	}
	if got := policies.allowedSchemes; len(got) != 2 || got[0] != "http" || got[1] != "https" {
		t.Fatalf("expected default http/https schemes, got %v", got)
	}
}

func TestPoliciesEvalAddressClasses(t *testing.T) {
	allAddressOptIns := NewPolicies(
		WithAllowLocalhost(true),
		WithAllowPrivateNetworks(true),
		WithAllowLinkLocal(true),
	)

	tests := []struct {
		name     string
		url      string
		want     string
		policies *Policies
	}{
		{name: "localhost", url: "http://localhost", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "localhost port", url: "http://localhost:8080", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "localhost subdomain", url: "http://service.localhost", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "localhost trailing dot", url: "http://localhost.", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "ipv4 loopback", url: "http://127.0.0.1", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "abbreviated ipv4 loopback", url: "http://127.1", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "decimal ipv4 loopback", url: "http://2130706433", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "ipv6 loopback", url: "http://[::1]", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "private 10", url: "http://10.0.0.1", want: "private networks are not allowed", policies: NewPolicies()},
		{name: "private 172", url: "http://172.16.0.1", want: "private networks are not allowed", policies: NewPolicies()},
		{name: "private 192", url: "http://192.168.0.1", want: "private networks are not allowed", policies: NewPolicies()},
		{name: "private ipv6 ula", url: "http://[fc00::1]", want: "private networks are not allowed", policies: NewPolicies()},
		{name: "carrier grade nat", url: "http://100.64.0.1", want: "private networks are not allowed", policies: NewPolicies()},
		{name: "ipv4 link local", url: "http://169.254.169.254", want: "link-local addresses are not allowed", policies: NewPolicies()},
		{name: "ipv6 link local", url: "http://[fe80::1]", want: "link-local addresses are not allowed", policies: NewPolicies()},
		{name: "ipv4 unspecified", url: "http://0.0.0.0", want: "unspecified addresses are not allowed", policies: NewPolicies()},
		{name: "ipv4 current network", url: "http://0.1.2.3", want: "unspecified addresses are not allowed", policies: NewPolicies()},
		{name: "ipv6 unspecified", url: "http://[::]", want: "unspecified addresses are not allowed", policies: NewPolicies()},
		{name: "ipv4 multicast", url: "http://224.0.0.1", want: "multicast addresses are not allowed", policies: NewPolicies()},
		{name: "ipv6 multicast", url: "http://[ff02::1]", want: "multicast addresses are not allowed", policies: NewPolicies()},
		{name: "ipv4 reserved", url: "http://240.0.0.1", want: "reserved addresses are not allowed", policies: NewPolicies()},
		{name: "ipv4 limited broadcast", url: "http://255.255.255.255", want: "reserved addresses are not allowed", policies: NewPolicies()},
		{name: "ipv6 site local", url: "http://[fec0::1]", want: "reserved addresses are not allowed", policies: NewPolicies()},
		{name: "mapped loopback", url: "http://[::ffff:127.0.0.1]", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "mapped private", url: "http://[::ffff:10.0.0.1]", want: "private networks are not allowed", policies: NewPolicies()},
		{name: "mapped link local", url: "http://[::ffff:169.254.169.254]", want: "link-local addresses are not allowed", policies: NewPolicies()},
		{name: "mapped public", url: "http://[::ffff:8.8.8.8]", policies: NewPolicies()},
		{name: "ietf protocol assignments", url: "http://192.0.0.8", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "pcp anycast", url: "http://192.0.0.9", policies: NewPolicies()},
		{name: "turn anycast", url: "http://192.0.0.10", policies: NewPolicies()},
		{name: "ietf protocol assignment remainder", url: "http://192.0.0.11", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "documentation ipv4 one", url: "http://192.0.2.1", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "deprecated relay anycast", url: "http://192.88.99.1", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "benchmarking ipv4 first", url: "http://198.18.0.1", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "benchmarking ipv4 last", url: "http://198.19.255.254", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "documentation ipv4 two", url: "http://198.51.100.1", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "documentation ipv4 three", url: "http://203.0.113.1", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "nat64 loopback", url: "http://[64:ff9b::127.0.0.1]", want: "localhost is not allowed", policies: NewPolicies()},
		{name: "nat64 private", url: "http://[64:ff9b::10.0.0.1]", want: "private networks are not allowed", policies: NewPolicies()},
		{name: "nat64 link local", url: "http://[64:ff9b::169.254.169.254]", want: "link-local addresses are not allowed", policies: NewPolicies()},
		{name: "nat64 public", url: "http://[64:ff9b::8.8.8.8]", policies: NewPolicies()},
		{name: "allow nat64 loopback", url: "http://[64:ff9b::127.0.0.1]", policies: NewPolicies(WithAllowLocalhost(true))},
		{name: "allow nat64 private", url: "http://[64:ff9b::10.0.0.1]", policies: NewPolicies(WithAllowPrivateNetworks(true))},
		{name: "allow nat64 link local", url: "http://[64:ff9b::169.254.169.254]", policies: NewPolicies(WithAllowLinkLocal(true))},
		{name: "local use nat64", url: "http://[64:ff9b:1::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "discard only ipv6", url: "http://[100::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "dummy ipv6", url: "http://[100:0:0:1::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "teredo ipv6", url: "http://[2001::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "benchmarking ipv6", url: "http://[2001:2::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "deprecated orchid", url: "http://[2001:10::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "documentation ipv6", url: "http://[2001:db8::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "six to four", url: "http://[2002::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "documentation ipv6 two", url: "http://[3fff::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "segment routing ipv6", url: "http://[5f00::1]", want: "non-public addresses are not allowed", policies: NewPolicies()},
		{name: "public ipv4", url: "http://8.8.8.8", policies: NewPolicies()},
		{name: "public ipv6 with port", url: "http://[2606:4700:4700::1111]:8080", policies: NewPolicies()},
		{name: "allow localhost", url: "http://127.1", policies: NewPolicies(WithAllowLocalhost(true))},
		{name: "allow private", url: "http://10.0.0.1", policies: NewPolicies(WithAllowPrivateNetworks(true))},
		{name: "allow carrier grade nat", url: "http://100.64.0.1", policies: NewPolicies(WithAllowPrivateNetworks(true))},
		{name: "allow link local", url: "http://169.254.169.254", policies: NewPolicies(WithAllowLinkLocal(true))},
		{name: "allow ipv6 link local", url: "http://[fe80::1]", policies: NewPolicies(WithAllowLinkLocal(true))},
		{name: "localhost does not allow private", url: "http://10.0.0.1", want: "private networks are not allowed", policies: NewPolicies(WithAllowLocalhost(true))},
		{name: "private does not allow link local", url: "http://169.254.169.254", want: "link-local addresses are not allowed", policies: NewPolicies(WithAllowPrivateNetworks(true))},
		{name: "link local does not allow private", url: "http://10.0.0.1", want: "private networks are not allowed", policies: NewPolicies(WithAllowLinkLocal(true))},
		{name: "allowlist does not allow loopback", url: "http://127.0.0.1", want: "localhost is not allowed", policies: NewPolicies(WithAllowedHosts("127.0.0.1"))},
		{name: "canonical trailing dot allowlist", url: "http://example.com.", policies: NewPolicies(WithAllowedHosts("example.com"))},
		{name: "opt ins do not allow unspecified", url: "http://0.1.2.3", want: "unspecified addresses are not allowed", policies: allAddressOptIns},
		{name: "opt ins do not allow multicast", url: "http://224.0.0.1", want: "multicast addresses are not allowed", policies: allAddressOptIns},
		{name: "opt ins do not allow reserved", url: "http://240.0.0.1", want: "reserved addresses are not allowed", policies: allAddressOptIns},
		{name: "opt ins do not allow benchmarking", url: "http://198.18.0.1", want: "non-public addresses are not allowed", policies: allAddressOptIns},
		{name: "opt ins do not allow documentation ipv6", url: "http://[2001:db8::1]", want: "non-public addresses are not allowed", policies: allAddressOptIns},
		{name: "opt ins do not allow nat64 documentation ipv4", url: "http://[64:ff9b::192.0.2.1]", want: "non-public addresses are not allowed", policies: allAddressOptIns},
		{name: "opt ins do not allow local nat64", url: "http://[64:ff9b:1::1]", want: "non-public addresses are not allowed", policies: allAddressOptIns},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policies.Eval(&Request{URL: tt.url})
			if tt.want == "" {
				if err != nil {
					t.Fatalf("expected address to be allowed, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
		})
	}
}

func TestPoliciesCanonicalHostPolicies(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		want     string
		policies *Policies
	}{
		{
			name:     "blocked integer ipv4",
			url:      "http://134744072",
			want:     "host is blocked",
			policies: NewPolicies(WithBlockedHosts("8.8.8.8")),
		},
		{
			name: "blocked abbreviated ipv4 with port",
			url:  "http://127.1:8080",
			want: "host is blocked",
			policies: NewPolicies(
				WithAllowLocalhost(true),
				WithBlockedHosts("127.0.0.1:8080"),
			),
		},
		{
			name:     "blocked mapped ipv4",
			url:      "http://[::ffff:8.8.8.8]",
			want:     "host is blocked",
			policies: NewPolicies(WithBlockedHosts("8.8.8.8")),
		},
		{
			name:     "blocked port remains specific",
			url:      "http://134744072:80",
			policies: NewPolicies(WithBlockedHosts("8.8.8.8:443")),
		},
		{
			name:     "allowlist port remains specific",
			url:      "http://134744072:80",
			want:     "host is not allowed",
			policies: NewPolicies(WithAllowedHosts("8.8.8.8:443")),
		},
		{
			name:     "allowed integer ipv4",
			url:      "http://134744072",
			policies: NewPolicies(WithAllowedHosts("8.8.8.8")),
		},
		{
			name: "allowed abbreviated ipv4 with port",
			url:  "http://127.1:8080",
			policies: NewPolicies(
				WithAllowLocalhost(true),
				WithAllowedHosts("127.0.0.1:8080"),
			),
		},
		{
			name:     "allowed mapped ipv4",
			url:      "http://[::ffff:8.8.8.8]",
			policies: NewPolicies(WithAllowedHosts("8.8.8.8")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policies.Eval(&Request{URL: tt.url})
			if tt.want == "" {
				if err != nil {
					t.Fatalf("expected host policy to allow request, got %v", err)
				}
				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
		})
	}
}

func TestPoliciesEval(t *testing.T) {
	tests := []struct {
		req      *Request
		name     string
		want     string
		policies *Policies
	}{
		{
			name:     "valid default get",
			policies: NewPolicies(),
			req:      &Request{URL: "http://example.com"},
		},
		{
			name:     "nil request",
			policies: NewPolicies(),
			req:      nil,
			want:     ErrNilRequest.Error(),
		},
		{
			name:     "invalid method",
			policies: NewPolicies(),
			req:      &Request{Method: "BAD METHOD", URL: "http://example.com"},
			want:     "invalid method",
		},
		{
			name:     "missing url",
			policies: NewPolicies(),
			req:      &Request{},
			want:     "url is required",
		},
		{
			name:     "missing scheme",
			policies: NewPolicies(),
			req:      &Request{URL: "example.com"},
			want:     "url scheme is required",
		},
		{
			name:     "missing host",
			policies: NewPolicies(),
			req:      &Request{URL: "http:///path"},
			want:     "url host is required",
		},
		{
			name:     "disallowed scheme",
			policies: NewPolicies(WithAllowedSchemes("https")),
			req:      &Request{URL: "http://example.com"},
			want:     "scheme",
		},
		{
			name:     "blocked host",
			policies: NewPolicies(WithBlockedHosts("example.com")),
			req:      &Request{URL: "http://example.com"},
			want:     "blocked",
		},
		{
			name:     "not allowed host",
			policies: NewPolicies(WithAllowedHosts("allowed.example")),
			req:      &Request{URL: "http://other.example"},
			want:     "not allowed",
		},
		{
			name:     "localhost",
			policies: NewPolicies(WithAllowLocalhost(false)),
			req:      &Request{URL: "http://127.0.0.1"},
			want:     "localhost is not allowed",
		},
		{
			name:     "private network",
			policies: NewPolicies(WithAllowPrivateNetworks(false)),
			req:      &Request{URL: "http://10.0.0.1"},
			want:     "private network",
		},
		{
			name:     "request body limit",
			policies: NewPolicies(WithMaxRequestSize(3)),
			req:      &Request{Method: stdhttp.MethodPost, URL: "http://example.com", Body: []byte("four")},
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
				if !errors.Is(err, ErrNilRequest) {
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
	policies := NewPolicies(
		WithDefaultHeader("X-Default", "default"),
		WithBlockedRequestHeaders("X-Blocked"),
	)
	req := &Request{
		URL: "HTTP://EXAMPLE.COM/path",
		Headers: Headers{
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
