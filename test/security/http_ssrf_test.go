package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/ferret/v2"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestDefaultEngineHTTPPolicyBlocksLoopback(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	query := source.NewAnonymous(`
RETURN TO_STRING(IO::NET::HTTP::GET(@url))
`)

	engine, err := ferret.New(ferret.WithParam("url", server.URL))
	if err != nil {
		t.Fatal(err)
	}

	_, err = engine.Run(context.Background(), query)
	if err == nil || !strings.Contains(err.Error(), "blocked by access policy") {
		t.Fatalf("expected default engine to block loopback, got %v", err)
	}
	if got := requests.Load(); got != 0 {
		t.Fatalf("expected no requests to reach loopback server, got %d", got)
	}
	if err := engine.Close(); err != nil {
		t.Fatal(err)
	}

	localClient, err := ferrethttp.New(ferrethttp.WithAllowLocalhost(true))
	if err != nil {
		t.Fatal(err)
	}

	localNetwork, err := ferretnet.New(ferretnet.WithHTTPClient(localClient))
	if err != nil {
		t.Fatal(err)
	}

	engine, err = ferret.New(
		ferret.WithParam("url", server.URL),
		ferret.WithNetwork(localNetwork),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	result, err := engine.Run(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(result.Content); got != `"ok"` {
		t.Fatalf("expected localhost-enabled response, got %s", got)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("expected one request to reach loopback server, got %d", got)
	}
}
