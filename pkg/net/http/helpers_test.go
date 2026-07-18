package http

import (
	stdhttp "net/http"
	"net/url"
	"testing"
)

func newTestPolicy(tb testing.TB, options ...PolicyOption) *Policy {
	tb.Helper()

	policy, err := NewPolicy(options...)
	if err != nil {
		tb.Fatalf("construct HTTP policy: %v", err)
	}

	return policy
}

func newTestClient(tb testing.TB, options ...PolicyOption) Client {
	tb.Helper()

	client, err := New(options...)
	if err != nil {
		tb.Fatalf("construct HTTP client: %v", err)
	}

	return client
}

func newTestPolicyRequest(tb testing.TB, method, rawURL string) *stdhttp.Request {
	tb.Helper()

	u, err := url.Parse(rawURL)
	if err != nil {
		tb.Fatalf("parse policy request URL: %v", err)
	}

	return &stdhttp.Request{
		Method: method,
		URL:    u,
		Header: make(stdhttp.Header),
	}
}

func newTestPolicyGETRequest(tb testing.TB, rawURL string) *stdhttp.Request {
	tb.Helper()

	return newTestPolicyRequest(tb, stdhttp.MethodGet, rawURL)
}

func policyTransportForTest(tb testing.TB, transport stdhttp.RoundTripper) *stdhttp.Transport {
	tb.Helper()

	if policyTransport, ok := transport.(*stdhttp.Transport); ok {
		return policyTransport
	}

	validated, ok := transport.(*responseValidatingTransport)
	if !ok {
		tb.Fatalf("expected policy or response-validating transport, got %T", transport)
	}
	policyTransport, ok := validated.next.(*stdhttp.Transport)
	if !ok {
		tb.Fatalf("expected policy transport, got %T", validated.next)
	}

	return policyTransport
}
