package security

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/ferret/v2"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type localHTTPServiceEndpoint struct {
	name string
	path string
	body string
}

func TestDefaultEngineHTTPPolicyBlocksLoopback(t *testing.T) {
	endpoints := localHTTPServiceEndpoints()
	baseURL, requests := newLoopbackHTTPService(t, endpoints)

	engine, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := engine.Close(); err != nil {
			t.Errorf("close default engine: %v", err)
		}
	})

	query := source.NewAnonymous(`
RETURN TO_STRING(IO::NET::HTTP::GET(@url))
`)

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			_, err := engine.Run(
				context.Background(),
				query,
				ferret.WithSessionParam("url", baseURL+endpoint.path),
			)

			if err == nil {
				t.Fatal("expected the default HTTP policy to reject loopback")
			}

			if !errors.Is(err, ferrethttp.ErrPolicyDenied) {
				t.Fatalf("expected ErrPolicyDenied, got %T: %v", err, err)
			}

			var policyErr *ferrethttp.PolicyError
			if !errors.As(err, &policyErr) {
				t.Fatalf("expected PolicyError, got %T: %v", err, err)
			}

			if policyErr.Target != ferrethttp.PolicyTargetRequest {
				t.Fatalf("expected request-target denial, got %q", policyErr.Target)
			}

			if policyErr.Reason != "localhost is not allowed" {
				t.Fatalf("expected localhost denial, got %q", policyErr.Reason)
			}
		})
	}

	if got := requests.Load(); got != 0 {
		t.Fatalf("expected no requests to reach the loopback service, got %d", got)
	}
}

func TestExplicitLocalhostOptInAllowsLoopbackHTTP(t *testing.T) {
	endpoints := localHTTPServiceEndpoints()
	baseURL, requests := newLoopbackHTTPService(t, endpoints)

	localClient, err := ferrethttp.New(ferrethttp.WithAllowLocalhost(true))
	if err != nil {
		t.Fatal(err)
	}

	localNetwork, err := ferretnet.New(ferretnet.WithHTTPClient(localClient))
	if err != nil {
		t.Fatal(err)
	}

	engine, err := ferret.New(ferret.WithNetwork(localNetwork))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := engine.Close(); err != nil {
			t.Errorf("close localhost-enabled engine: %v", err)
		}
	})

	query := source.NewAnonymous(`
RETURN TO_STRING(IO::NET::HTTP::GET(@url))
`)

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			result, err := engine.Run(
				context.Background(),
				query,
				ferret.WithSessionParam("url", baseURL+endpoint.path),
			)

			if err != nil {
				t.Fatal(err)
			}

			if result == nil {
				t.Fatal("expected an HTTP response")
			}

			var body string
			if err := json.Unmarshal(result.Content, &body); err != nil {
				t.Fatalf("decode FQL result %q: %v", result.Content, err)
			}

			if body != endpoint.body {
				t.Fatalf("expected response %q, got %q", endpoint.body, body)
			}
		})
	}

	if got, want := requests.Load(), int64(len(endpoints)); got != want {
		t.Fatalf("expected %d requests to reach the opted-in loopback service, got %d", want, got)
	}
}

func localHTTPServiceEndpoints() []localHTTPServiceEndpoint {
	return []localHTTPServiceEndpoint{
		{
			name: "test endpoint",
			path: "/test",
			body: "localhost-ssrf-ok",
		},
		{
			name: "mock metadata role",
			path: "/latest/meta-data/iam/security-credentials/",
			body: "ferret-internal-role",
		},
		{
			name: "mock metadata credentials",
			path: "/latest/meta-data/iam/security-credentials/ferret-internal-role",
			body: `{"AccessKeyId":"AKIA_TEST_ONLY","SecretAccessKey":"test-secret","Token":"test-token"}`,
		},
		{
			name: "mock admin panel",
			path: "/admin",
			body: "INTERNAL ADMIN PANEL - localhost only",
		},
	}
}

func newLoopbackHTTPService(
	t *testing.T,
	endpoints []localHTTPServiceEndpoint,
) (string, *atomic.Int64) {
	t.Helper()

	requests := &atomic.Int64{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		for _, endpoint := range endpoints {
			if r.URL.Path == endpoint.path {
				_, _ = w.Write([]byte(endpoint.body))

				return
			}
		}

		http.NotFound(w, r)
	})

	server := httptest.NewUnstartedServer(handler)
	if err := server.Listener.Close(); err != nil {
		t.Fatalf("close default test listener: %v", err)
	}

	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on IPv4 loopback: %v", err)
	}

	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)

	return server.URL, requests
}
