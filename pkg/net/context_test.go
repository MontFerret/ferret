package net

import (
	"context"
	"errors"
	"testing"
)

func TestWithNetworkRoundTrip(t *testing.T) {
	network, err := New()
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	ctx := WithNetwork(context.Background(), network)

	resolved, err := NetworkFrom(ctx)
	if err != nil {
		t.Fatalf("network from context failed: %v", err)
	}

	if resolved != network {
		t.Fatalf("expected same network instance from context")
	}
}

func TestHTTPClientFrom(t *testing.T) {
	client := stubHTTPClient{}
	network, err := New(WithHTTPClient(client))
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	ctx := WithNetwork(context.Background(), network)

	resolved, err := HTTPClientFrom(ctx)
	if err != nil {
		t.Fatalf("http client from context failed: %v", err)
	}

	if resolved != client {
		t.Fatalf("expected custom http client from context")
	}
}

func TestNetworkFromError(t *testing.T) {
	if _, err := NetworkFrom(context.Background()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for background context, got %v", err)
	}

	if _, err := NetworkFrom(nil); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nil context, got %v", err)
	}

	if _, err := HTTPClientFrom(context.Background()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for missing http client, got %v", err)
	}
}

func TestDefaultNetworkForwardsIdleConnectionCleanup(t *testing.T) {
	t.Parallel()

	client := &trackingHTTPClient{}
	network, err := New(WithHTTPClient(client))
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	closer, ok := network.(interface{ CloseIdleConnections() })
	if !ok {
		t.Fatalf("expected default network to expose idle-connection cleanup")
	}

	closer.CloseIdleConnections()
	closer.CloseIdleConnections()

	if got := client.idleCloseCount(); got != 2 {
		t.Fatalf("expected cleanup to be forwarded twice, got %d", got)
	}
}
