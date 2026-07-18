package http

import (
	"io"
	stdhttp "net/http"
	"strings"
	"testing"
)

func BenchmarkClientDoPolicy(b *testing.B) {
	benchmarkClientDoPolicy(b, newTestPolicy(b))
}

func BenchmarkClientDoPolicyWithoutTimeout(b *testing.B) {
	benchmarkClientDoPolicy(b, newTestPolicy(b, WithNoTimeout()))
}

func BenchmarkClientDoPolicyWithHeaders(b *testing.B) {
	benchmarkClientDoPolicyRequest(b, newTestPolicy(b), &Request{
		URL: "https://api.example.test/resource",
		Headers: Headers{
			"Accept":        {"application/json"},
			"Authorization": {"Bearer benchmark-token"},
			"X-Trace":       {"one", "two"},
		},
	})
}

func benchmarkClientDoPolicy(b *testing.B, policy *Policy) {
	benchmarkClientDoPolicyRequest(b, policy, &Request{URL: "https://api.example.test/resource"})
}

func benchmarkClientDoPolicyRequest(b *testing.B, policy *Policy, req *Request) {
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{
			Transport: newResponseValidatingTransport(testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
				return &stdhttp.Response{
					StatusCode: stdhttp.StatusOK,
					Status:     "200 OK",
					Header:     make(stdhttp.Header),
					Body:       io.NopCloser(strings.NewReader("ok")),
				}, nil
			})),
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := client.Do(b.Context(), req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadResponseBodyBounded(b *testing.B) {
	body := strings.Repeat("x", 4<<10)

	b.ReportAllocs()
	b.SetBytes(int64(len(body)))
	b.ResetTimer()

	for b.Loop() {
		if _, err := readResponseBody(strings.NewReader(body), 8<<10); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPolicyDialerControl(b *testing.B) {
	policies := newTestPolicy(b)
	dialer := newPolicyDialer(policies)
	ctx := b.Context()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if err := dialer.controlContext(ctx, "tcp4", "93.184.216.34:443", nil); err != nil {
			b.Fatal(err)
		}
	}
}
