package http

import (
	"io"
	stdhttp "net/http"
	"strings"
	"testing"
)

func BenchmarkClientDoPolicy(b *testing.B) {
	policies := NewPolicies()
	client := &defaultHTTPClient{
		policy: &policies,
		transport: stdhttp.Client{
			Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
				return &stdhttp.Response{
					StatusCode: stdhttp.StatusOK,
					Status:     "200 OK",
					Header:     make(stdhttp.Header),
					Body:       io.NopCloser(strings.NewReader("ok")),
				}, nil
			}),
		},
	}
	req := &Request{URL: "https://api.example.test/resource"}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := client.Do(b.Context(), req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPolicyDialerControl(b *testing.B) {
	policies := NewPolicies()
	dialer := newPolicyDialer(&policies)
	ctx := b.Context()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if err := dialer.controlContext(ctx, "tcp4", "93.184.216.34:443", nil); err != nil {
			b.Fatal(err)
		}
	}
}
