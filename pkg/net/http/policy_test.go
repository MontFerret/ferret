package http

import (
	"errors"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"
)

func TestNewPolicyDefaultsAndOptions(t *testing.T) {
	policies := newTestPolicy(t,
		WithTimeout(time.Second),
		WithAllowedSchemes("HTTPS", "https"),
		WithAllowedMethods("post", "CUSTOM-METHOD", "POST"),
		WithAllowedHosts("Example.COM", "example.com"),
		WithBlockedHosts("Blocked.EXAMPLE", "blocked.example"),
		WithBlockedRequestHeaders("x-token", "X-Token"),
		WithDefaultHeader("x-default", "value"),
		WithAllowLocalhost(false),
		WithAllowPrivateNetworks(false),
		WithAllowLinkLocal(false),
		WithFollowRedirects(false),
		WithMaxRedirects(2),
		WithMaxRequestSize(3),
		WithMaxResponseSize(4),
		WithMaxResponseHeaderSize(5),
	)

	if policies.timeout != time.Second {
		t.Fatalf("expected timeout to be configured, got %s", policies.timeout)
	}
	if policies.followRedirects {
		t.Fatal("expected redirects to be disabled")
	}
	if policies.maxRedirects != 2 || policies.maxRequestSize != 3 || policies.maxResponseSize != 4 || policies.maxResponseHeaderSize != 5 {
		t.Fatalf("expected configured limits, got redirects=%d request=%d response=%d response-headers=%d",
			policies.maxRedirects,
			policies.maxRequestSize,
			policies.maxResponseSize,
			policies.maxResponseHeaderSize,
		)
	}
	if policies.allowLocalhost || policies.allowPrivateNetworks || policies.allowLinkLocal {
		t.Fatal("expected localhost, private networks, and link-local addresses to be disabled")
	}
	if got := policies.allowedSchemes; len(got) != 1 || got[0] != "https" {
		t.Fatalf("expected normalized allowed scheme, got %v", got)
	}
	if got := policies.allowedMethods; len(got) != 2 || got[0] != "POST" || got[1] != "CUSTOM-METHOD" {
		t.Fatalf("expected normalized allowed methods, got %v", got)
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

func TestNewPolicyDefaultValues(t *testing.T) {
	policies := newTestPolicy(t)

	if !policies.followRedirects {
		t.Fatal("expected redirects to be enabled by default")
	}
	if policies.allowLocalhost || policies.allowPrivateNetworks || policies.allowLinkLocal {
		t.Fatal("expected localhost, private networks, and link-local addresses to be denied by default")
	}
	if got := policies.allowedSchemes; len(got) != 2 || got[0] != "http" || got[1] != "https" {
		t.Fatalf("expected default http/https schemes, got %v", got)
	}
	if policies.timeout != defaultTimeout {
		t.Fatalf("expected default timeout %s, got %s", defaultTimeout, policies.timeout)
	}
	if policies.maxResponseHeaderSize != defaultMaxResponseHeaderSize {
		t.Fatalf("expected default response header limit %d, got %d", defaultMaxResponseHeaderSize, policies.maxResponseHeaderSize)
	}
	if policies.maxRedirects != defaultMaxRedirects {
		t.Fatalf("expected default redirect limit %d, got %d", defaultMaxRedirects, policies.maxRedirects)
	}
	if policies.maxRequestSize != defaultMaxRequestSize {
		t.Fatalf("expected default request size %d, got %d", defaultMaxRequestSize, policies.maxRequestSize)
	}
	if policies.maxResponseSize != defaultMaxResponseSize {
		t.Fatalf("expected default response size %d, got %d", defaultMaxResponseSize, policies.maxResponseSize)
	}
	if len(policies.allowedMethods) != 7 {
		t.Fatalf("expected seven default methods, got %v", policies.allowedMethods)
	}
}

func TestPolicyEvalAddressClasses(t *testing.T) {
	allAddressOptIns := newTestPolicy(t,
		WithAllowLocalhost(true),
		WithAllowPrivateNetworks(true),
		WithAllowLinkLocal(true),
	)

	tests := []struct {
		policies *Policy
		name     string
		url      string
		want     string
	}{
		{name: "localhost", url: "http://localhost", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "localhost port", url: "http://localhost:8080", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "localhost subdomain", url: "http://service.localhost", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "localhost trailing dot", url: "http://localhost.", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "ipv4 loopback", url: "http://127.0.0.1", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "abbreviated ipv4 loopback", url: "http://127.1", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "decimal ipv4 loopback", url: "http://2130706433", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "ipv6 loopback", url: "http://[::1]", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "private 10", url: "http://10.0.0.1", want: "private networks are not allowed", policies: newTestPolicy(t)},
		{name: "private 172", url: "http://172.16.0.1", want: "private networks are not allowed", policies: newTestPolicy(t)},
		{name: "private 192", url: "http://192.168.0.1", want: "private networks are not allowed", policies: newTestPolicy(t)},
		{name: "private ipv6 ula", url: "http://[fc00::1]", want: "private networks are not allowed", policies: newTestPolicy(t)},
		{name: "carrier grade nat", url: "http://100.64.0.1", want: "private networks are not allowed", policies: newTestPolicy(t)},
		{name: "ipv4 link local", url: "http://169.254.169.254", want: "link-local addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv6 link local", url: "http://[fe80::1]", want: "link-local addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv4 unspecified", url: "http://0.0.0.0", want: "unspecified addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv4 current network", url: "http://0.1.2.3", want: "unspecified addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv6 unspecified", url: "http://[::]", want: "unspecified addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv4 multicast", url: "http://224.0.0.1", want: "multicast addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv6 multicast", url: "http://[ff02::1]", want: "multicast addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv4 reserved", url: "http://240.0.0.1", want: "reserved addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv4 limited broadcast", url: "http://255.255.255.255", want: "reserved addresses are not allowed", policies: newTestPolicy(t)},
		{name: "ipv6 site local", url: "http://[fec0::1]", want: "reserved addresses are not allowed", policies: newTestPolicy(t)},
		{name: "mapped loopback", url: "http://[::ffff:127.0.0.1]", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "mapped private", url: "http://[::ffff:10.0.0.1]", want: "private networks are not allowed", policies: newTestPolicy(t)},
		{name: "mapped link local", url: "http://[::ffff:169.254.169.254]", want: "link-local addresses are not allowed", policies: newTestPolicy(t)},
		{name: "mapped public", url: "http://[::ffff:8.8.8.8]", policies: newTestPolicy(t)},
		{name: "ietf protocol assignments", url: "http://192.0.0.8", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "pcp anycast", url: "http://192.0.0.9", policies: newTestPolicy(t)},
		{name: "turn anycast", url: "http://192.0.0.10", policies: newTestPolicy(t)},
		{name: "ietf protocol assignment remainder", url: "http://192.0.0.11", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "documentation ipv4 one", url: "http://192.0.2.1", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "deprecated relay anycast", url: "http://192.88.99.1", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "benchmarking ipv4 first", url: "http://198.18.0.1", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "benchmarking ipv4 last", url: "http://198.19.255.254", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "documentation ipv4 two", url: "http://198.51.100.1", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "documentation ipv4 three", url: "http://203.0.113.1", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "nat64 loopback", url: "http://[64:ff9b::127.0.0.1]", want: "localhost is not allowed", policies: newTestPolicy(t)},
		{name: "nat64 private", url: "http://[64:ff9b::10.0.0.1]", want: "private networks are not allowed", policies: newTestPolicy(t)},
		{name: "nat64 link local", url: "http://[64:ff9b::169.254.169.254]", want: "link-local addresses are not allowed", policies: newTestPolicy(t)},
		{name: "nat64 public", url: "http://[64:ff9b::8.8.8.8]", policies: newTestPolicy(t)},
		{name: "allow nat64 loopback", url: "http://[64:ff9b::127.0.0.1]", policies: newTestPolicy(t, WithAllowLocalhost(true))},
		{name: "allow nat64 private", url: "http://[64:ff9b::10.0.0.1]", policies: newTestPolicy(t, WithAllowPrivateNetworks(true))},
		{name: "allow nat64 link local", url: "http://[64:ff9b::169.254.169.254]", policies: newTestPolicy(t, WithAllowLinkLocal(true))},
		{name: "local use nat64", url: "http://[64:ff9b:1::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "discard only ipv6", url: "http://[100::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "dummy ipv6", url: "http://[100:0:0:1::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "teredo ipv6", url: "http://[2001::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "benchmarking ipv6", url: "http://[2001:2::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "deprecated orchid", url: "http://[2001:10::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "documentation ipv6", url: "http://[2001:db8::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "six to four", url: "http://[2002::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "documentation ipv6 two", url: "http://[3fff::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "segment routing ipv6", url: "http://[5f00::1]", want: "non-public addresses are not allowed", policies: newTestPolicy(t)},
		{name: "public ipv4", url: "http://8.8.8.8", policies: newTestPolicy(t)},
		{name: "public ipv6 with port", url: "http://[2606:4700:4700::1111]:8080", policies: newTestPolicy(t)},
		{name: "allow localhost", url: "http://127.1", policies: newTestPolicy(t, WithAllowLocalhost(true))},
		{name: "allow private", url: "http://10.0.0.1", policies: newTestPolicy(t, WithAllowPrivateNetworks(true))},
		{name: "allow carrier grade nat", url: "http://100.64.0.1", policies: newTestPolicy(t, WithAllowPrivateNetworks(true))},
		{name: "allow link local", url: "http://169.254.169.254", policies: newTestPolicy(t, WithAllowLinkLocal(true))},
		{name: "allow ipv6 link local", url: "http://[fe80::1]", policies: newTestPolicy(t, WithAllowLinkLocal(true))},
		{name: "localhost does not allow private", url: "http://10.0.0.1", want: "private networks are not allowed", policies: newTestPolicy(t, WithAllowLocalhost(true))},
		{name: "private does not allow link local", url: "http://169.254.169.254", want: "link-local addresses are not allowed", policies: newTestPolicy(t, WithAllowPrivateNetworks(true))},
		{name: "link local does not allow private", url: "http://10.0.0.1", want: "private networks are not allowed", policies: newTestPolicy(t, WithAllowLinkLocal(true))},
		{name: "allowlist does not allow loopback", url: "http://127.0.0.1", want: "localhost is not allowed", policies: newTestPolicy(t, WithAllowedHosts("127.0.0.1"))},
		{name: "canonical trailing dot allowlist", url: "http://example.com.", policies: newTestPolicy(t, WithAllowedHosts("example.com"))},
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

func TestPolicyCanonicalHostPolicies(t *testing.T) {
	tests := []struct {
		policies *Policy
		name     string
		url      string
		want     string
	}{
		{
			name:     "blocked integer ipv4",
			url:      "http://134744072",
			want:     "host is blocked",
			policies: newTestPolicy(t, WithBlockedHosts("8.8.8.8")),
		},
		{
			name: "blocked abbreviated ipv4 with port",
			url:  "http://127.1:8080",
			want: "host is blocked",
			policies: newTestPolicy(t,
				WithAllowLocalhost(true),
				WithBlockedHosts("127.0.0.1:8080"),
			),
		},
		{
			name:     "blocked mapped ipv4",
			url:      "http://[::ffff:8.8.8.8]",
			want:     "host is blocked",
			policies: newTestPolicy(t, WithBlockedHosts("8.8.8.8")),
		},
		{
			name:     "blocked port remains specific",
			url:      "http://134744072:80",
			policies: newTestPolicy(t, WithBlockedHosts("8.8.8.8:443")),
		},
		{
			name:     "allowlist port remains specific",
			url:      "http://134744072:80",
			want:     "host is not allowed",
			policies: newTestPolicy(t, WithAllowedHosts("8.8.8.8:443")),
		},
		{
			name:     "allowed integer ipv4",
			url:      "http://134744072",
			policies: newTestPolicy(t, WithAllowedHosts("8.8.8.8")),
		},
		{
			name: "allowed abbreviated ipv4 with port",
			url:  "http://127.1:8080",
			policies: newTestPolicy(t,
				WithAllowLocalhost(true),
				WithAllowedHosts("127.0.0.1:8080"),
			),
		},
		{
			name:     "allowed mapped ipv4",
			url:      "http://[::ffff:8.8.8.8]",
			policies: newTestPolicy(t, WithAllowedHosts("8.8.8.8")),
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

func TestPolicyEval(t *testing.T) {
	tests := []struct {
		req      *Request
		policies *Policy
		name     string
		want     string
	}{
		{
			name:     "valid default get",
			policies: newTestPolicy(t),
			req:      &Request{URL: "http://example.com"},
		},
		{
			name:     "nil request",
			policies: newTestPolicy(t),
			req:      nil,
			want:     ErrNilRequest.Error(),
		},
		{
			name:     "invalid method",
			policies: newTestPolicy(t),
			req:      &Request{Method: "BAD METHOD", URL: "http://example.com"},
			want:     "invalid method",
		},
		{
			name:     "missing url",
			policies: newTestPolicy(t),
			req:      &Request{},
			want:     "url is required",
		},
		{
			name:     "missing scheme",
			policies: newTestPolicy(t),
			req:      &Request{URL: "example.com"},
			want:     "url scheme is required",
		},
		{
			name:     "missing host",
			policies: newTestPolicy(t),
			req:      &Request{URL: "http:///path"},
			want:     "url host is required",
		},
		{
			name:     "disallowed scheme",
			policies: newTestPolicy(t, WithAllowedSchemes("https")),
			req:      &Request{URL: "http://example.com"},
			want:     "scheme",
		},
		{
			name:     "blocked host",
			policies: newTestPolicy(t, WithBlockedHosts("example.com")),
			req:      &Request{URL: "http://example.com"},
			want:     "blocked",
		},
		{
			name:     "not allowed host",
			policies: newTestPolicy(t, WithAllowedHosts("allowed.example")),
			req:      &Request{URL: "http://other.example"},
			want:     "not allowed",
		},
		{
			name:     "localhost",
			policies: newTestPolicy(t, WithAllowLocalhost(false)),
			req:      &Request{URL: "http://127.0.0.1"},
			want:     "localhost is not allowed",
		},
		{
			name:     "private network",
			policies: newTestPolicy(t, WithAllowPrivateNetworks(false)),
			req:      &Request{URL: "http://10.0.0.1"},
			want:     "private network",
		},
		{
			name:     "request body limit",
			policies: newTestPolicy(t, WithMaxRequestSize(3)),
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

func TestPolicyEvalDoesNotMutateRequest(t *testing.T) {
	policies := newTestPolicy(t,
		WithDefaultHeader("X-Default", "default"),
		WithBlockedRequestHeaders("X-Blocked"),
	)
	req := &Request{
		URL: "HTTP://EXAMPLE.COM/path",
		Headers: Headers{
			"X-Blocked": {"secret"},
		},
	}

	err := policies.Eval(req)
	requirePolicyError(t, err, PolicyTargetRequest)

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
