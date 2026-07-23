package net

import (
	"context"
	"errors"
	"strings"
	"testing"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

func TestWithHTTPTransportConstructsPolicyAwareClient(t *testing.T) {
	transport := &trackingHTTPTransport{body: "four"}
	network, err := New(
		WithHTTPTransport(
			transport,
			ferrethttp.WithMaxResponseSize(3),
		),
	)
	if err != nil {
		t.Fatalf("create network: %v", err)
	}

	response, err := network.HTTP().Do(
		context.Background(),
		&ferrethttp.Request{URL: "https://example.com"},
	)
	if response != nil {
		t.Fatalf("expected no oversized response, got %#v", response)
	}
	var limitErr *ferrethttp.ResponseBodyLimitError
	if !errors.As(err, &limitErr) {
		t.Fatalf("expected ResponseBodyLimitError, got %T: %v", err, err)
	}
	if limitErr.Size != 4 || limitErr.Limit != 3 {
		t.Fatalf("unexpected response limit error: %#v", limitErr)
	}
	if got := transport.callCount(); got != 1 {
		t.Fatalf("expected one transport call, got %d", got)
	}

	closer, ok := network.(interface{ CloseIdleConnections() })
	if !ok {
		t.Fatal("expected network to expose idle-connection cleanup")
	}
	closer.CloseIdleConnections()
	closer.CloseIdleConnections()
	if got := transport.idleCloseCount(); got != 2 {
		t.Fatalf("expected cleanup to reach transport twice, got %d", got)
	}
}

func TestWithHTTPTransportPreservesPolicyOptionOrder(t *testing.T) {
	tests := []struct {
		options func(*trackingHTTPTransport) []Option
		name    string
		wantErr bool
	}{
		{
			name: "transport policy last",
			options: func(transport *trackingHTTPTransport) []Option {
				return []Option{
					WithHTTPPolicies(ferrethttp.WithMaxResponseSize(3)),
					WithHTTPTransport(transport, ferrethttp.WithMaxResponseSize(2)),
				}
			},
			wantErr: true,
		},
		{
			name: "standalone policy last",
			options: func(transport *trackingHTTPTransport) []Option {
				return []Option{
					WithHTTPTransport(transport, ferrethttp.WithMaxResponseSize(2)),
					WithHTTPPolicies(ferrethttp.WithMaxResponseSize(3)),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &trackingHTTPTransport{body: "123"}
			network, err := New(tt.options(transport)...)
			if err != nil {
				t.Fatalf("create network: %v", err)
			}

			response, err := network.HTTP().Do(
				context.Background(),
				&ferrethttp.Request{URL: "https://example.com"},
			)
			if tt.wantErr {
				var limitErr *ferrethttp.ResponseBodyLimitError
				if !errors.As(err, &limitErr) {
					t.Fatalf("expected ResponseBodyLimitError, got %T: %v", err, err)
				}

				return
			}
			if err != nil {
				t.Fatalf("execute request: %v", err)
			}
			if got := string(response.Body); got != "123" {
				t.Fatalf("expected response body %q, got %q", "123", got)
			}
		})
	}
}

func TestWithHTTPClientTakesPrecedenceOverTransport(t *testing.T) {
	client := stubHTTPClient{}

	tests := []struct {
		options func(*trackingHTTPTransport) []Option
		name    string
	}{
		{
			name: "client first",
			options: func(transport *trackingHTTPTransport) []Option {
				return []Option{
					WithHTTPClient(client),
					WithHTTPTransport(
						transport,
						ferrethttp.WithMaxResponseSize(-1),
					),
				}
			},
		},
		{
			name: "transport first",
			options: func(transport *trackingHTTPTransport) []Option {
				return []Option{
					WithHTTPTransport(
						transport,
						ferrethttp.WithMaxResponseSize(-1),
					),
					WithHTTPClient(client),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &trackingHTTPTransport{}
			network, err := New(tt.options(transport)...)
			if err != nil {
				t.Fatalf("create network: %v", err)
			}
			if network.HTTP() != client {
				t.Fatalf("expected supplied client, got %T", network.HTTP())
			}

			if _, err := network.HTTP().Do(
				context.Background(),
				&ferrethttp.Request{URL: "https://example.com"},
			); err != nil {
				t.Fatalf("execute request: %v", err)
			}
			if got := transport.callCount(); got != 0 {
				t.Fatalf("expected transport to remain unused, got %d calls", got)
			}
		})
	}
}

func TestWithHTTPTransportPropagatesPolicyConfigurationError(t *testing.T) {
	network, err := New(
		WithHTTPTransport(
			&trackingHTTPTransport{},
			ferrethttp.WithMaxResponseSize(-1),
		),
	)
	if network != nil {
		t.Fatalf("expected no network, got %T", network)
	}
	if !errors.Is(err, ferrethttp.ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected ErrInvalidPolicyConfiguration, got %v", err)
	}
	if err == nil || !strings.HasPrefix(err.Error(), "http client: ") {
		t.Fatalf("expected wrapped HTTP client error, got %v", err)
	}
}

func TestWithHTTPTransportNilIsNoOp(t *testing.T) {
	network, err := New(
		WithHTTPTransport(nil, ferrethttp.WithMaxResponseSize(-1)),
	)
	if err != nil {
		t.Fatalf("expected nil transport option to be ignored: %v", err)
	}
	if network == nil {
		t.Fatal("expected default network")
	}
}
