package uapi

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apidiagnostics "github.com/MontFerret/api/diagnostics"

	ferretdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/engine"
)

func TestWrapRequiresNativeEngine(t *testing.T) {
	defer func() {
		if got := recover(); got != "uapi: nil native engine" {
			t.Fatalf("Wrap(nil) panic = %v, want uapi: nil native engine", got)
		}
	}()

	Wrap(nil)
}

func TestRuntimeTranslatesSourceOptionsAndReusesPlan(t *testing.T) {
	runtime := newTestRuntime(t)

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "value.txt"), []byte("rooted"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	sourcePath := filepath.Join(root, "query.fql")

	compiled, err := runtime.Compile(
		context.Background(),
		api.NewSource(sourcePath, "RETURN [@value, TO_STRING(IO::FS::READ(\"value.txt\"))]"),
	)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	t.Cleanup(func() { _ = compiled.Close() })

	if got := compiled.Params(); len(got) != 1 || got[0] != "value" {
		t.Fatalf("Params = %v, want [value]", got)
	}

	parameters := compiled.Params()
	parameters[0] = "changed"

	if got := compiled.Params(); len(got) != 1 || got[0] != "value" {
		t.Fatalf("Params after caller mutation = %v, want [value]", got)
	}

	for _, value := range []string{"first", "second"} {
		session, err := compiled.NewSession(
			context.Background(),
			api.WithParam("value", value),
			api.WithOutputContentType("application/json"),
			api.WithFSRoot(root),
		)
		if err != nil {
			t.Fatalf("NewSession(%s): %v", value, err)
		}

		output, runErr := session.Run(context.Background())

		closeErr := session.Close()
		if err := errors.Join(runErr, closeErr); err != nil {
			t.Fatalf("session(%s): %v", value, err)
		}

		want := "[\"" + value + "\",\"rooted\"]"
		if output.ContentType != "application/json" || string(output.Content) != want {
			t.Fatalf("output(%s) = %+v, want %s", value, output, want)
		}
	}
}

func TestRuntimeRunUsesUniversalSessionOptions(t *testing.T) {
	runtime := newTestRuntime(t)

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "value.txt"), []byte("rooted"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	output, err := runtime.Run(
		context.Background(),
		api.NewSource(
			filepath.Join(root, "query.fql"),
			"RETURN [@value, TO_STRING(IO::FS::READ(\"value.txt\"))]",
		),
		api.WithParam("value", "direct"),
		api.WithOutputContentType("application/json"),
		api.WithFSRoot(root),
	)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if output.ContentType != "application/json" || string(output.Content) != `["direct","rooted"]` {
		t.Fatalf("output = %+v", output)
	}
}

func TestRuntimeOutputsAreIndependentAndDebugSessionsWork(t *testing.T) {
	runtime := newTestRuntime(t)
	ctx := context.Background()
	sourcePath := filepath.Join(t.TempDir(), "query.fql")

	compiled, err := runtime.Compile(ctx, api.NewSource(sourcePath, "RETURN 1"))
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	t.Cleanup(func() { _ = compiled.Close() })

	first, err := compiled.NewSession(ctx)
	if err != nil {
		t.Fatalf("NewSession first: %v", err)
	}

	firstOutput, err := first.Run(ctx)
	if err != nil {
		t.Fatalf("Run first: %v", err)
	}

	if err := first.Close(); err != nil {
		t.Fatalf("Close first: %v", err)
	}

	firstOutput.Content[0] = '9'

	second, err := compiled.NewSession(ctx)
	if err != nil {
		t.Fatalf("NewSession second: %v", err)
	}

	secondOutput, err := second.Run(ctx)
	if err != nil {
		t.Fatalf("Run second: %v", err)
	}

	if err := second.Close(); err != nil {
		t.Fatalf("Close second: %v", err)
	}

	if string(secondOutput.Content) != "1" {
		t.Fatalf("second output = %q, want independent 1", secondOutput.Content)
	}

	debugPlan, err := runtime.CompileDebug(
		ctx,
		api.NewSource(sourcePath, "RETURN 1"),
		api.WithOptimizationLevel(api.OptimizationNone),
	)
	if err != nil {
		t.Fatalf("CompileDebug: %v", err)
	}

	t.Cleanup(func() { _ = debugPlan.Close() })

	debugSession, err := debugPlan.NewDebugSession(ctx)
	if err != nil {
		t.Fatalf("NewDebugSession: %v", err)
	}

	t.Cleanup(func() { _ = debugSession.Close() })

	event, err := debugSession.Start(ctx)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	if event.Reason != apidebugger.ReasonEntry || event.Location.SourceName != sourcePath {
		t.Fatalf("entry event = %+v", event)
	}

	completed, err := debugSession.Continue(ctx)
	if err != nil {
		t.Fatalf("Continue: %v", err)
	}

	if completed.Reason != apidebugger.ReasonCompleted || completed.Output == nil ||
		string(completed.Output.Content) != "1" {
		t.Fatalf("completed event = %+v", completed)
	}
}

func TestRuntimeConvertsDiagnosticsAndPreservesNativeCause(t *testing.T) {
	runtime := newTestRuntime(t)
	source := api.NewSource("/absolute/query.fql", "RETURN")

	_, err := runtime.Compile(context.Background(), source)
	if err == nil {
		t.Fatal("Compile unexpectedly succeeded")
	}

	var nativeSet *ferretdiagnostics.DiagnosticSet

	var native *ferretdiagnostics.Diagnostic
	if !errors.As(err, &nativeSet) && !errors.As(err, &native) {
		t.Fatalf("error does not preserve native diagnostic cause: %v", err)
	}

	var portable apidiagnostics.Diagnostics
	if !errors.As(err, &portable) || len(portable) == 0 {
		t.Fatalf("portable diagnostics = %+v, want non-empty", portable)
	}

	if portable[0].Source != source || portable[0].Message == "" ||
		len(portable[0].Annotations) == 0 ||
		portable[0].Annotations[0].Range.SourceName != source.Name {
		t.Fatalf("portable diagnostic = %+v", portable[0])
	}
}

func newTestEngine(t testing.TB, options ...engine.Option) *engine.Engine {
	t.Helper()

	native, err := engine.New(options...)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := native.Close(); err != nil {
			t.Errorf("engine Close: %v", err)
		}
	})

	return native
}

func newTestRuntime(t testing.TB, options ...engine.Option) api.Runtime {
	t.Helper()

	runtime := Wrap(newTestEngine(t, options...))
	t.Cleanup(func() {
		if err := runtime.Close(); err != nil {
			t.Errorf("runtime Close: %v", err)
		}
	})

	return runtime
}
