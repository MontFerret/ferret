package http_test

import (
	"errors"
	stdhttp "net/http"
	"net/netip"
	"strings"
	"testing"
	"time"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

func TestPolicySupportsCustomBackendIntegration(t *testing.T) {
	policy, err := ferrethttp.NewPolicy(
		ferrethttp.WithAllowedHosts("api.example.com"),
		ferrethttp.WithTimeout(2*time.Second),
		ferrethttp.WithMaxResponseSize(3),
		ferrethttp.WithMaxResponseHeaderSize(2048),
	)
	if err != nil {
		t.Fatalf("construct policy: %v", err)
	}

	request, err := stdhttp.NewRequest(stdhttp.MethodGet, "https://api.example.com/start", nil)
	if err != nil {
		t.Fatalf("construct request: %v", err)
	}
	if err := policy.Prepare(request); err != nil {
		t.Fatalf("prepare request: %v", err)
	}

	redirect, err := stdhttp.NewRequest(stdhttp.MethodGet, "https://api.example.com/next", nil)
	if err != nil {
		t.Fatalf("construct redirect request: %v", err)
	}
	if err := policy.CheckRedirect(redirect, []*stdhttp.Request{request}); err != nil {
		t.Fatalf("check redirect: %v", err)
	}
	if err := policy.EvalConnection(netip.MustParseAddr("93.184.216.34")); err != nil {
		t.Fatalf("evaluate concrete address: %v", err)
	}

	if policy.Timeout() != 2*time.Second {
		t.Fatalf("expected timeout %s, got %s", 2*time.Second, policy.Timeout())
	}
	if policy.MaxResponseSize() != 3 {
		t.Fatalf("expected response limit 3, got %d", policy.MaxResponseSize())
	}
	if policy.MaxResponseHeaderSize() != 2048 {
		t.Fatalf("expected response header limit 2048, got %d", policy.MaxResponseHeaderSize())
	}

	body, err := policy.ReadResponseBody(strings.NewReader("ok"))
	if err != nil {
		t.Fatalf("materialize response body: %v", err)
	}
	if string(body) != "ok" {
		t.Fatalf("expected body %q, got %q", "ok", body)
	}

	err = policy.EvalResponseSize(4)
	var limitErr *ferrethttp.ResponseBodyLimitError
	if !errors.As(err, &limitErr) {
		t.Fatalf("expected response body limit error, got %T: %v", err, err)
	}
}
