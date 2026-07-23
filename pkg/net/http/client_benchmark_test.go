package http

import (
	"context"
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

func BenchmarkToStdRequest(b *testing.B) {
	req := &Request{
		URL: "https://api.example.test/resource",
		Headers: Headers{
			"Accept":  {"application/json"},
			"X-Trace": {"one", "two"},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := toStdRequest(context.Background(), req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPolicyEval(b *testing.B) {
	policy := newTestPolicy(b)
	req := newTestPolicyGETRequest(b, "https://api.example.test/resource")
	req.Header = stdhttp.Header{
		"Accept":  {"application/json"},
		"X-Trace": {"one", "two"},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if err := policy.Eval(req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPolicyPrepareDefaults(b *testing.B) {
	policy := newTestPolicy(b,
		WithDefaultHeader("Accept", "application/json"),
		WithDefaultHeader("X-Trace", "benchmark"),
	)
	base := newTestPolicyGETRequest(b, "https://api.example.test/resource")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		req := *base
		req.Header = make(stdhttp.Header)
		if err := policy.Prepare(&req); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkClientDoPolicy(b *testing.B, policy *Policy) {
	benchmarkClientDoPolicyRequest(b, policy, &Request{URL: "https://api.example.test/resource"})
}

func benchmarkClientDoPolicyRequest(b *testing.B, policy *Policy, req *Request) {
	client := newDefaultHTTPClient(policy, stdhttp.Client{
		Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			return &stdhttp.Response{
				StatusCode: stdhttp.StatusOK,
				Status:     "200 OK",
				Header:     make(stdhttp.Header),
				Body:       io.NopCloser(strings.NewReader("ok")),
			}, nil
		}),
	})

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
	policy := newTestPolicy(b, WithMaxResponseSize(8<<10))

	b.ReportAllocs()
	b.SetBytes(int64(len(body)))
	b.ResetTimer()

	for b.Loop() {
		if _, err := policy.ReadResponseBody(strings.NewReader(body)); err != nil {
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
