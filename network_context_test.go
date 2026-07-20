package ferret

import (
	"context"
	"strings"
	"testing"

	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestWithNetworkRejectsNil(t *testing.T) {
	t.Parallel()

	_, err := New(WithNetwork(nil))
	if err == nil {
		t.Fatal("expected nil network to fail")
	}

	if !strings.Contains(err.Error(), "network cannot be nil") {
		t.Fatalf("expected nil network validation error, got: %v", err)
	}
}

func TestSessionRunInjectsConfiguredNetwork(t *testing.T) {
	t.Parallel()

	client := &recordingHTTPClient{body: []byte("session-network")}
	network := mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))
	engine := mustNewEngine(t, WithNetwork(network))
	defer func() { _ = engine.Close() }()

	out, err := engine.Run(context.Background(), source.NewAnonymous(`
RETURN TO_STRING(IO::NET::HTTP::GET({ url: "https://example.test/session" }))
`))
	if err != nil {
		t.Fatal(err)
	}

	if got := string(out.Content); got != `"session-network"` {
		t.Fatalf("expected injected network response, got %s", got)
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	if client.calls != 1 {
		t.Fatalf("expected custom network to be called once, got %d", client.calls)
	}
	if client.lastURL != "https://example.test/session" {
		t.Fatalf("expected custom network URL to be recorded, got %q", client.lastURL)
	}
}

func TestDebugSessionRunInjectsConfiguredNetwork(t *testing.T) {
	t.Parallel()

	client := &recordingHTTPClient{body: []byte("debug-network")}
	network := mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))
	engine := mustNewEngine(t, WithNetwork(network))
	defer func() { _ = engine.Close() }()

	plan, err := engine.CompileDebug(context.Background(), source.New("debug-network.fql", `
LET marker = 1
RETURN TO_STRING(IO::NET::HTTP::GET({ url: "https://example.test/debug" }))
`))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if event.Reason != DebugReasonCompleted {
		t.Fatalf("expected debug session to complete, got %#v", event)
	}
	if event.Output == nil || string(event.Output.Content) != `"debug-network"` {
		t.Fatalf("expected injected network output, got %#v", event.Output)
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	if client.calls != 1 {
		t.Fatalf("expected custom network to be called once, got %d", client.calls)
	}
	if client.lastURL != "https://example.test/debug" {
		t.Fatalf("expected custom network URL to be recorded, got %q", client.lastURL)
	}
}
