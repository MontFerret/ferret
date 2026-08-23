package net

import (
	"context"
	"errors"
	"testing"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

func TestNetworkOptionsPreserveDefaultsAndNilInputs(t *testing.T) {
	network, err := New(
		nil,
		WithHTTPClient(nil),
		WithHTTPPolicies(),
		WithHTTPTransport(nil, ferrethttp.WithMaxResponseSize(-1)),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if network == nil {
		t.Fatal("New() network = nil, want default network")
	}

	closer, ok := network.(interface{ CloseIdleConnections() })
	if !ok {
		t.Fatalf("New() network = %T, want idle-connection cleanup", network)
	}
	t.Cleanup(closer.CloseIdleConnections)
}

func TestNetworkScalarOptionsUseLaterValues(t *testing.T) {
	t.Run("HTTP client", func(t *testing.T) {
		first := stubHTTPClient{}
		second := &stubHTTPClient{}
		network, err := New(
			WithHTTPClient(first),
			WithHTTPClient(second),
		)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		if network.HTTP() != second {
			t.Fatalf("HTTP() = %T, want later client", network.HTTP())
		}
	})

	t.Run("HTTP transport", func(t *testing.T) {
		first := &trackingHTTPTransport{body: "first"}
		second := &trackingHTTPTransport{body: "second"}
		network, err := New(
			WithHTTPTransport(first),
			WithHTTPTransport(second),
		)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		closer, ok := network.(interface{ CloseIdleConnections() })
		if !ok {
			t.Fatalf("New() network = %T, want idle-connection cleanup", network)
		}
		t.Cleanup(closer.CloseIdleConnections)

		response, err := network.HTTP().Do(
			context.Background(),
			&ferrethttp.Request{URL: "https://example.com"},
		)
		if err != nil {
			t.Fatalf("HTTP().Do() error = %v", err)
		}
		if got := string(response.Body); got != "second" {
			t.Fatalf("response body = %q, want %q", got, "second")
		}
		if got := first.callCount(); got != 0 {
			t.Fatalf("first transport calls = %d, want 0", got)
		}
		if got := second.callCount(); got != 1 {
			t.Fatalf("second transport calls = %d, want 1", got)
		}
	})
}

func TestNewReturnsJoinedOptionErrors(t *testing.T) {
	first := errors.New("first option failed")
	second := errors.New("second option failed")
	secondApplied := false
	network, err := New(
		func(_ *config) error {
			return first
		},
		nil,
		func(_ *config) error {
			secondApplied = true

			return second
		},
	)
	if network != nil {
		t.Fatalf("New() network = %T, want nil", network)
	}
	if !errors.Is(err, first) || !errors.Is(err, second) {
		t.Fatalf("New() error = %v, want both option errors", err)
	}
	if !secondApplied {
		t.Fatal("later option was not applied after an earlier failure")
	}
}
