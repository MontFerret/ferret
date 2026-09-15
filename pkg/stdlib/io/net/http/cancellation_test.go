package http_test

import (
	"context"
	"errors"
	"testing"
	"time"

	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	stdlibhttp "github.com/MontFerret/ferret/v2/pkg/stdlib/io/net/http"
)

func TestHTTPRequestCancellationWhileBlocked(t *testing.T) {
	client := &blockingHTTPClient{started: make(chan context.Context, 1)}
	network, err := ferretnet.New(ferretnet.WithHTTPClient(client))
	if err != nil {
		t.Fatal(err)
	}

	defer ferretnet.CloseIdleNetworkConnections(network)

	ctx, cancel := context.WithCancel(ferretnet.WithNetwork(t.Context(), network))
	defer cancel()
	finished := make(chan error, 1)
	go func() {
		_, err := stdlibhttp.GET(ctx, runtime.String("https://example.com"))
		finished <- err
	}()

	select {
	case received := <-client.started:
		if received != ctx {
			t.Fatal("HTTP client did not receive the caller's context")
		}

		cancel()
	case <-time.After(5 * time.Second):
		t.Fatal("request never reached the HTTP client")
	}

	select {
	case err := <-finished:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("request error = %v, want cancellation", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request did not observe cancellation")
	}
}
