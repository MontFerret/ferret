package uapi

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/api"
	apidiagnostics "github.com/MontFerret/api/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestNewOwnsConfiguredEngine(t *testing.T) {
	var initialized int
	r, err := New(
		engine.WithEngineInitHook(func() error {
			initialized++

			return nil
		}),
		engine.WithParam("value", 1),
		engine.WithParam("value", 41),
	)
	if err != nil {
		t.Fatal(err)
	}

	var portable api.Runtime = r
	t.Cleanup(func() {
		if err := portable.Close(); err != nil {
			t.Error(err)
		}
	})

	if initialized != 1 {
		t.Fatalf("init calls=%d", initialized)
	}

	out, err := portable.Run(t.Context(), api.NewAnonymousSource("RETURN @value + 1"))
	if err != nil || out == nil || string(out.Content) != "42" {
		t.Fatalf("configured runtime: output=%+v err=%v", out, err)
	}
}

func TestNewIgnoresNilNativeOption(t *testing.T) {
	r, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := r.Close(); err != nil {
			t.Error(err)
		}
	})
	out, err := r.Run(t.Context(), api.NewAnonymousSource("RETURN 42"))
	if err != nil || out == nil || string(out.Content) != "42" {
		t.Fatalf("nil option: output=%+v err=%v", out, err)
	}
}

func TestNewDelegatesOptionFailureRollback(t *testing.T) {
	client := &ownershipHTTPClient{}
	r, err := New(
		engine.WithNetworkOptions(ferretnet.WithHTTPClient(client)),
		engine.WithEngineInitHook(nil),
	)
	if r != nil {
		_ = r.Close()

		t.Fatal("failed construction returned a runtime")
	}

	_, nativeErr := engine.New(engine.WithEngineInitHook(nil))
	if err == nil || nativeErr == nil || err.Error() != nativeErr.Error() {
		t.Fatalf("option error=%v, want Native error=%v", err, nativeErr)
	}

	if got := client.closes.Load(); got != 1 {
		t.Fatalf("rollback network closes=%d, want 1", got)
	}
}

func TestNewProjectsConstructionAndRollbackFailures(t *testing.T) {
	client := &ownershipHTTPClient{}
	initCause, closeCause := errors.New("init cause"), errors.New("close cause")
	src := source.New("init.fql", "RETURN 1")
	initDiagnostic := diagnostics.NewUnexpectedErrorWith(src, "init diagnostic", initCause)
	closeDiagnostic := diagnostics.NewUnexpectedErrorWith(src, "close diagnostic", closeCause)
	var closes int
	r, err := New(
		engine.WithNetworkOptions(ferretnet.WithHTTPClient(client)),
		engine.WithEngineInitHook(func() error { return initDiagnostic }),
		engine.WithEngineCloseHook(func() error {
			closes++

			return closeDiagnostic
		}),
	)
	if r != nil {
		_ = r.Close()

		t.Fatal("failed construction returned a runtime")
	}

	for _, cause := range []error{initCause, initDiagnostic, closeCause, closeDiagnostic} {
		if !errors.Is(err, cause) {
			t.Fatalf("construction error lost %v: %v", cause, err)
		}
	}

	var projected apidiagnostics.Diagnostics
	portableSource := api.NewSource(src.Name(), src.Content())
	if !errors.As(err, &projected) || len(projected) != 2 ||
		projected[0].Message != initDiagnostic.Message || projected[1].Message != closeDiagnostic.Message ||
		projected[0].Source != portableSource || projected[1].Source != portableSource {
		t.Fatalf("construction diagnostics=%+v err=%v", projected, err)
	}

	if closes != 1 || client.closes.Load() != 1 {
		t.Fatalf("rollback: hook calls=%d network closes=%d, want 1 each", closes, client.closes.Load())
	}
}

func TestOwnedRuntimeCloseDelegatesCleanupAndRejection(t *testing.T) {
	client := &ownershipHTTPClient{}
	cause := errors.New("engine cleanup failed")
	diagnostic := diagnostics.NewUnexpectedErrorWith(source.Source{}, "close diagnostic", cause)
	var closes atomic.Int32
	r, err := New(
		engine.WithNetworkOptions(ferretnet.WithHTTPClient(client)),
		engine.WithEngineCloseHook(func() error {
			closes.Add(1)

			return diagnostic
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = r.Close() })
	results := make([]error, 8)
	var wg sync.WaitGroup
	for i := range results {
		wg.Go(func() { results[i] = r.Close() })
	}

	wg.Wait()
	results = append(results, r.Close())
	for _, result := range results {
		var projected apidiagnostics.Diagnostics
		if !errors.Is(result, cause) || !errors.Is(result, diagnostic) ||
			!errors.As(result, &projected) || len(projected) != 1 || projected[0].Message != diagnostic.Message {
			t.Fatalf("close lost Native error or diagnostic: %v", result)
		}
	}

	if closes.Load() != 1 || client.closes.Load() != 1 {
		t.Fatalf("cleanup: hook calls=%d network closes=%d, want 1 each", closes.Load(), client.closes.Load())
	}

	for _, compile := range []func(context.Context, api.Source, ...api.PlanOption) (api.Plan, error){r.Compile, r.CompileDebug} {
		p, err := compile(t.Context(), api.NewAnonymousSource("RETURN 1"))
		if p != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("closed owned runtime: plan=%v err=%v", p, err)
		}
	}

	if out, err := r.Run(t.Context(), api.NewAnonymousSource("RETURN 1")); !errors.Is(err, runtime.ErrInvalidOperation) || out != nil {
		t.Fatalf("closed owned runtime: output=%+v err=%v", out, err)
	}
}
