package engine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

func TestEngineNetworkMatchesModuleAndExecutionContexts(t *testing.T) {
	t.Parallel()

	var (
		registeredNetwork ferretnet.Network
		registeredClient  ferrethttp.Client
	)
	mod := testModule{registerFn: func(boot module.Bootstrap) error {
		registeredNetwork = boot.Host().Network()
		if registeredNetwork == nil {
			return errors.New("module host has no network")
		}
		registeredClient = registeredNetwork.HTTP()
		boot.Host().Library().Function().A0().Add("VERIFY_NETWORK_IDENTITY", func(ctx context.Context) (runtime.Value, error) {
			network, err := ferretnet.NetworkFrom(ctx)
			if err != nil {
				return runtime.None, err
			}
			if network != registeredNetwork {
				return runtime.None, errors.New("execution network differs from the module host network")
			}
			client, err := ferretnet.HTTPClientFrom(ctx)
			if err != nil {
				return runtime.None, err
			}
			if client != registeredClient {
				return runtime.None, errors.New("execution HTTP client differs from the module host client")
			}

			return runtime.NewBoolean(true), nil
		})

		return nil
	}}
	eng := mustNewEngine(t, WithModules(mod))
	defer eng.Close()

	if got := eng.Network(); got != registeredNetwork || got.HTTP() != registeredClient {
		t.Fatal("expected the same network and HTTP client as the module host")
	}

	ctx := t.Context()
	src := source.New("network-identity.fql", "RETURN VERIFY_NETWORK_IDENTITY()")
	out, err := collectEngine(eng, ctx, src)
	if err != nil || out == nil || string(out.Data) != "true" {
		t.Fatalf("expected ordinary execution to verify network identity, got output %#v, error %v", out, err)
	}

	plan, err := eng.CompileDebug(ctx, src)
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(ctx); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if event == nil || event.Reason != debugger.ReasonCompleted || event.Output == nil || string(event.Output.Data) != "true" {
		t.Fatalf("expected debugger execution to verify network identity, got %#v", event)
	}
}

func TestSessionRunInjectsConfiguredNetwork(t *testing.T) {
	t.Parallel()

	client := &recordingHTTPClient{body: []byte("session-network")}
	network := mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))
	engine := mustNewEngine(t, WithNetwork(network))
	defer func() { _ = engine.Close() }()

	out, err := collectEngine(engine, context.Background(), source.NewAnonymous(`
RETURN TO_STRING(IO::NET::HTTP::GET({ url: "https://example.test/session" }))
`))
	if err != nil {
		t.Fatal(err)
	}

	if got := string(out.Data); got != `"session-network"` {
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

	if event.Reason != debugger.ReasonCompleted {
		t.Fatalf("expected debug session to complete, got %#v", event)
	}
	if event.Output == nil || string(event.Output.Data) != `"debug-network"` {
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
