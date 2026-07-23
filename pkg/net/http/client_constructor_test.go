package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewWithClientRejectsNilClientBeforePolicyValidation(t *testing.T) {
	client, err := NewWithClient(nil, WithMaxRequestSize(-1))
	if client != nil {
		t.Fatalf("expected no client, got %T", client)
	}
	if !errors.Is(err, ErrNilClient) {
		t.Fatalf("expected ErrNilClient, got %v", err)
	}
	if errors.Is(err, ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected nil client validation to take precedence, got %v", err)
	}
}

func TestNewWithClientRejectsInvalidPolicyConfiguration(t *testing.T) {
	client, err := NewWithClient(
		&stdhttp.Client{},
		WithMaxRequestSize(-1),
		WithMaxResponseSize(-2),
	)
	if client != nil {
		t.Fatalf("expected no client, got %T", client)
	}
	if !errors.Is(err, ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected ErrInvalidPolicyConfiguration, got %v", err)
	}
}

func TestNewWithClientInstallsPolicyTransportWhenMissing(t *testing.T) {
	const maxResponseHeaderSize = 2048

	stdClient := &stdhttp.Client{}
	client, err := NewWithClient(
		stdClient,
		WithMaxResponseHeaderSize(maxResponseHeaderSize),
	)
	if err != nil {
		t.Fatalf("construct HTTP client: %v", err)
	}

	if stdClient.Transport != nil {
		t.Fatalf("expected supplied client to remain unchanged, got %T", stdClient.Transport)
	}

	adapted, ok := client.(*defaultHTTPClient)
	if !ok {
		t.Fatalf("expected built-in client, got %T", client)
	}

	transport := policyTransportForTest(t, adapted.client.Transport)
	if transport.Proxy != nil {
		t.Fatal("expected ambient proxy lookup to be disabled")
	}
	if transport.DialContext == nil {
		t.Fatal("expected policy-aware dialer")
	}
	if transport.MaxResponseHeaderBytes != maxResponseHeaderSize {
		t.Fatalf(
			"expected max response header size %d, got %d",
			maxResponseHeaderSize,
			transport.MaxResponseHeaderBytes,
		)
	}
	if adapted.client.Timeout != defaultTimeout {
		t.Fatalf("expected policy timeout %s, got %s", defaultTimeout, adapted.client.Timeout)
	}
	if adapted.client.CheckRedirect == nil {
		t.Fatal("expected policy redirect callback")
	}
}

func TestNewWithClientSnapshotsClientAndSharesResources(t *testing.T) {
	originalTransport := &trackingIdleTransport{}
	replacementTransport := &trackingIdleTransport{}
	originalJar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("construct original cookie jar: %v", err)
	}
	replacementJar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("construct replacement cookie jar: %v", err)
	}

	stdClient := &stdhttp.Client{
		Transport: originalTransport,
		Jar:       originalJar,
		Timeout:   time.Hour,
	}
	client, err := NewWithClient(stdClient)
	if err != nil {
		t.Fatalf("construct HTTP client: %v", err)
	}
	if stdClient.Transport != originalTransport {
		t.Fatal("expected supplied transport to remain unchanged")
	}
	if stdClient.Jar != originalJar {
		t.Fatal("expected supplied cookie jar to remain unchanged")
	}
	if stdClient.Timeout != time.Hour {
		t.Fatalf("expected supplied timeout %s, got %s", time.Hour, stdClient.Timeout)
	}

	stdClient.Transport = replacementTransport
	stdClient.Jar = replacementJar
	stdClient.Timeout = 2 * time.Hour

	adapted, ok := client.(*defaultHTTPClient)
	if !ok {
		t.Fatalf("expected built-in client, got %T", client)
	}
	validatingTransport, ok := adapted.client.Transport.(*responseValidatingTransport)
	if !ok {
		t.Fatalf("expected response-validating transport, got %T", adapted.client.Transport)
	}
	if validatingTransport.next != originalTransport {
		t.Fatalf("expected original underlying transport, got %T", validatingTransport.next)
	}
	if adapted.client.Jar != originalJar {
		t.Fatal("expected original cookie jar")
	}
	if adapted.client.Timeout != defaultTimeout {
		t.Fatalf("expected policy timeout %s, got %s", defaultTimeout, adapted.client.Timeout)
	}

	closer, ok := client.(IdleConnectionCloser)
	if !ok {
		t.Fatalf("expected client to implement IdleConnectionCloser, got %T", client)
	}
	closer.CloseIdleConnections()

	if !originalTransport.closed.Load() {
		t.Fatal("expected idle connection cleanup to reach original transport")
	}
	if replacementTransport.closed.Load() {
		t.Fatal("expected replacement transport to remain untouched")
	}
}

func TestNewWithClientPolicyTimeoutOverridesClientTimeout(t *testing.T) {
	var remaining time.Duration
	transport := testRoundTripper(func(req *stdhttp.Request) (*stdhttp.Response, error) {
		deadline, ok := req.Context().Deadline()
		if !ok {
			t.Fatal("expected request deadline")
		}

		remaining = time.Until(deadline)

		return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
	})

	client, err := NewWithClient(
		&stdhttp.Client{
			Transport: transport,
			Timeout:   time.Hour,
		},
		WithTimeout(10*time.Second),
	)
	if err != nil {
		t.Fatalf("construct HTTP client: %v", err)
	}

	if _, err := client.Do(
		context.Background(),
		&Request{URL: "https://example.com"},
	); err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if remaining <= 0 || remaining >= time.Minute {
		t.Fatalf("expected policy timeout to override client timeout, got %s", remaining)
	}
}

func TestNewWithClientPolicyRedirectOverridesClientCallback(t *testing.T) {
	mux := stdhttp.NewServeMux()
	mux.HandleFunc("/start", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		stdhttp.Redirect(w, r, "/done", stdhttp.StatusFound)
	})
	mux.HandleFunc("/done", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		_, _ = w.Write([]byte("done"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	suppliedCallbackCalled := false
	client, err := NewWithClient(
		&stdhttp.Client{
			CheckRedirect: func(*stdhttp.Request, []*stdhttp.Request) error {
				suppliedCallbackCalled = true

				return errors.New("supplied redirect callback")
			},
		},
		WithAllowLocalhost(true),
	)
	if err != nil {
		t.Fatalf("construct HTTP client: %v", err)
	}

	res, err := client.Do(
		context.Background(),
		&Request{URL: server.URL + "/start"},
	)
	if err != nil {
		t.Fatalf("redirect request failed: %v", err)
	}
	if got := string(res.Body); got != "done" {
		t.Fatalf("expected followed redirect body %q, got %q", "done", got)
	}
	if suppliedCallbackCalled {
		t.Fatal("expected policy redirect callback to replace supplied callback")
	}
}
