package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	stdhttp "net/http"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPolicyDialerValidatesConcreteAddresses(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    string
	}{
		{name: "loopback", address: "127.0.0.1", want: "localhost is not allowed"},
		{name: "private", address: "10.0.0.10", want: "private networks are not allowed"},
		{name: "link local", address: "169.254.169.254", want: "link-local addresses are not allowed"},
		{name: "benchmarking", address: "198.18.0.1", want: "non-public addresses are not allowed"},
		{name: "documentation ipv6", address: "2001:db8::1", want: "non-public addresses are not allowed"},
		{name: "nat64 loopback", address: "64:ff9b::127.0.0.1", want: "localhost is not allowed"},
		{name: "nat64 private", address: "64:ff9b::10.0.0.1", want: "private networks are not allowed"},
		{name: "invalid", address: "not-an-address", want: "invalid address is not allowed"},
		{name: "public ipv4", address: "93.184.216.34"},
		{name: "public ipv6", address: "2606:4700:4700::1111"},
		{name: "public protocol anycast", address: "192.0.0.9"},
		{name: "nat64 public", address: "64:ff9b::8.8.8.8"},
	}

	policies := newTestPolicy(t)
	dialer := newPolicyDialer(policies)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dialer.controlContext(
				context.Background(),
				"tcp",
				net.JoinHostPort(tt.address, "443"),
				nil,
			)
			if tt.want == "" {
				if err != nil {
					t.Fatalf("expected address to be allowed, got %v", err)
				}
				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
			requirePolicyError(t, err, PolicyTargetConnection)
		})
	}
}

func TestPolicyDialerRejectsReboundConcreteAddress(t *testing.T) {
	policies := newTestPolicy(t)
	dialer := newPolicyDialer(policies)

	if err := dialer.controlContext(
		context.Background(),
		"tcp4",
		"93.184.216.34:443",
		nil,
	); err != nil {
		t.Fatalf("expected initial public address to pass, got %v", err)
	}

	err := dialer.controlContext(context.Background(), "tcp4", "127.0.0.1:443", nil)
	if err == nil || !strings.Contains(err.Error(), "localhost is not allowed") {
		t.Fatalf("expected rebound loopback address to be rejected, got %v", err)
	}
	requirePolicyError(t, err, PolicyTargetConnection)
}

func TestPolicyDialerRejectsMalformedDialAddressAsConnectionTarget(t *testing.T) {
	dialer := newPolicyDialer(newTestPolicy(t))

	err := dialer.controlContext(context.Background(), "tcp", "missing-port", nil)
	policyErr := requirePolicyError(t, err, PolicyTargetConnection)
	if policyErr.Reason != "invalid address is not allowed" {
		t.Fatalf("unexpected malformed-address error: %#v", policyErr)
	}
}

func TestDefaultHTTPClientValidatesRedirectDestinations(t *testing.T) {
	policies := newTestPolicy(t)
	requested := make(map[string]int)
	client := newDefaultHTTPClient(policies, stdhttp.Client{
		Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
			requested[req.URL.String()]++
			switch req.URL.Path {
			case "/loopback":
				return responseWithBody(stdhttp.StatusFound, "", stdhttp.Header{"Location": {"http://127.0.0.1/metadata"}}), nil
			case "/private":
				return responseWithBody(stdhttp.StatusFound, "", stdhttp.Header{"Location": {"http://10.0.0.10/private"}}), nil
			case "/allowed":
				return responseWithBody(stdhttp.StatusFound, "", stdhttp.Header{"Location": {"https://1.1.1.1/done"}}), nil
			default:
				return responseWithBody(stdhttp.StatusOK, "done", nil), nil
			}
		})},
	)

	for _, path := range []string{"/loopback", "/private"} {
		_, err := client.Do(context.Background(), &Request{URL: "https://93.184.216.34" + path})
		if err == nil || !strings.Contains(err.Error(), "redirect destination blocked by access policy") {
			t.Fatalf("expected redirect policy error for %s, got %v", path, err)
		}
	}

	res, err := client.Do(context.Background(), &Request{URL: "https://93.184.216.34/allowed"})
	if err != nil {
		t.Fatalf("expected allowed redirect to pass, got %v", err)
	}
	if got := string(res.Body); got != "done" {
		t.Fatalf("expected allowed redirect response %q, got %q", "done", got)
	}
	if requested["http://127.0.0.1/metadata"] != 0 || requested["http://10.0.0.10/private"] != 0 {
		t.Fatalf("expected forbidden redirect targets not to be requested, got %v", requested)
	}
}

func TestDefaultHTTPClientNoFollowSkipsRedirectValidation(t *testing.T) {
	policies := newTestPolicy(t, WithFollowRedirects(false))
	client := newDefaultHTTPClient(policies, stdhttp.Client{
		Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			return responseWithBody(stdhttp.StatusFound, "", stdhttp.Header{"Location": {"http://10.0.0.10/private"}}), nil
		})},
	)

	res, err := client.Do(context.Background(), &Request{URL: "https://93.184.216.34/start"})
	if err != nil {
		t.Fatalf("expected no-follow response, got %v", err)
	}
	if res.StatusCode != stdhttp.StatusFound {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusFound, res.StatusCode)
	}
}

func TestDefaultHTTPClientUsesStandardRedirectLimit(t *testing.T) {
	policies := newTestPolicy(t)
	roundTrips := 0
	client := newDefaultHTTPClient(policies, stdhttp.Client{
		Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			roundTrips++
			return responseWithBody(stdhttp.StatusFound, "", stdhttp.Header{"Location": {"/next"}}), nil
		})},
	)

	_, err := client.Do(context.Background(), &Request{URL: "https://93.184.216.34/start"})
	if err == nil || !strings.Contains(err.Error(), "stopped after 10 redirect(s)") {
		t.Fatalf("expected standard redirect limit error, got %v", err)
	}
	if roundTrips != defaultMaxRedirects+1 {
		t.Fatalf("expected %d requests before stopping, got %d", defaultMaxRedirects+1, roundTrips)
	}
}

func TestDefaultHTTPClientFollowsExactlyConfiguredRedirects(t *testing.T) {
	for _, limit := range []int{1, 2, 3} {
		t.Run(fmt.Sprintf("limit_%d", limit), func(t *testing.T) {
			roundTrips := 0
			client := newDefaultHTTPClient(
				newTestPolicy(t, WithMaxRedirects(limit)),
				stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
					roundTrips++
					return responseWithBody(
						stdhttp.StatusFound,
						"",
						stdhttp.Header{"Location": {"/next"}},
					), nil
				})},
			)

			_, err := client.Do(context.Background(), &Request{URL: "https://example.com/start"})
			var limitErr *RedirectLimitError
			if !errors.As(err, &limitErr) {
				t.Fatalf("expected RedirectLimitError, got %T: %v", err, err)
			}
			if limitErr.Limit != limit {
				t.Fatalf("expected redirect limit %d, got %d", limit, limitErr.Limit)
			}
			if roundTrips != limit+1 {
				t.Fatalf("expected initial request plus %d redirects, got %d request(s)", limit, roundTrips)
			}
		})
	}
}

func TestRedirectPrivateResolutionUsesConnectionTarget(t *testing.T) {
	policy := newTestPolicy(t, WithTimeout(2*time.Second))
	dialer := newPolicyDialer(policy)
	resolver, dnsQueries := newLoopbackResolver(t)
	dialer.dialer.Resolver = resolver

	policyTransport := newPolicyTransport(dialer, policy.maxResponseHeaderSize)
	t.Cleanup(policyTransport.CloseIdleConnections)

	client := newDefaultHTTPClient(policy, stdhttp.Client{
		Transport: testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
			if req.URL.Hostname() == "public.example" {
				return responseWithBody(
					stdhttp.StatusFound,
					"",
					stdhttp.Header{"Location": {"https://rebound.example/private"}},
				), nil
			}

			return policyTransport.RoundTrip(req)
		})},
	)

	_, err := client.Do(context.Background(), &Request{URL: "https://public.example/start"})
	policyErr := requirePolicyError(t, err, PolicyTargetConnection)
	if policyErr.Reason != "localhost is not allowed" {
		t.Fatalf("unexpected private-resolution denial: %#v", policyErr)
	}
	if dnsQueries.Load() == 0 {
		t.Fatal("expected redirected hostname to be resolved through the policy transport")
	}
}

func TestDefaultHTTPClientConcurrentDo(t *testing.T) {
	client := newDefaultHTTPClient(newTestPolicy(t), stdhttp.Client{
		Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
		})},
	)

	const (
		goroutines = 32
		iterations = 25
	)

	errs := make(chan error, goroutines)
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for range iterations {
				response, err := client.Do(context.Background(), &Request{URL: "https://example.com"})
				if err != nil {
					errs <- err
					return
				}
				if string(response.Body) != "ok" {
					errs <- fmt.Errorf("unexpected response body %q", response.Body)
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatal(err)
	}
}

func TestDefaultHTTPClientTimeoutCoversDNSResolution(t *testing.T) {
	policies := newTestPolicy(t, WithTimeout(50*time.Millisecond))
	dialer := newPolicyDialer(policies)
	lookupStarted := make(chan struct{})

	var startOnce sync.Once
	dialer.dialer.Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			startOnce.Do(func() { close(lookupStarted) })
			<-ctx.Done()

			return nil, ctx.Err()
		},
	}
	client := newDefaultHTTPClient(policies, stdhttp.Client{
		Transport: newPolicyTransport(dialer, policies.maxResponseHeaderSize),
	})

	started := time.Now()
	_, err := client.Do(context.Background(), &Request{URL: "http://timeout.example"})
	elapsed := time.Since(started)

	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DNS resolution deadline error, got %v", err)
	}

	if elapsed >= time.Second {
		t.Fatalf("expected configured timeout to cancel DNS promptly, took %s", elapsed)
	}

	select {
	case <-lookupStarted:
	default:
		t.Fatal("expected custom DNS resolver to be called")
	}
}

func TestPolicyRejectsInvalidAddress(t *testing.T) {
	policies := newTestPolicy(t,
		WithAllowLocalhost(true),
		WithAllowPrivateNetworks(true),
		WithAllowLinkLocal(true),
	)

	err := policies.validateAddress(PolicyTargetRequest, "destination address", netip.Addr{})
	if err == nil || !strings.Contains(err.Error(), "invalid address is not allowed") {
		t.Fatalf("expected invalid address rejection, got %v", err)
	}
}

func TestNewUsesDedicatedPolicyTransportWithoutProxy(t *testing.T) {
	client, ok := newTestClient(t).(*defaultHTTPClient)
	if !ok {
		t.Fatalf("expected built-in client, got %T", client)
	}

	if _, ok := any(client).(IdleConnectionCloser); !ok {
		t.Fatalf("expected built-in client to implement IdleConnectionCloser")
	}

	transport := policyTransportForTest(t, client.client.Transport)

	if transport.Proxy != nil {
		t.Fatal("expected ambient proxy lookup to be disabled")
	}

	if transport.DialContext == nil {
		t.Fatal("expected policy-aware dialer")
	}
}

func TestDefaultHTTPClientClosesIdleConnections(t *testing.T) {
	transport := &trackingIdleTransport{}
	client := newDefaultHTTPClient(newTestPolicy(t), stdhttp.Client{
		Transport: transport,
	})

	client.CloseIdleConnections()
	if !transport.closed.Load() {
		t.Fatal("expected idle connection cleanup to reach transport")
	}
}

func responseWithBody(status int, body string, headers stdhttp.Header) *stdhttp.Response {
	if headers == nil {
		headers = make(stdhttp.Header)
	}

	return &stdhttp.Response{
		StatusCode: status,
		Status:     stdhttp.StatusText(status),
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
