package http

import (
	stdhttp "net/http"
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
