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
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	nativesource "github.com/MontFerret/ferret/v2/pkg/source"
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

	if got, err := compiled.Params(); err != nil || len(got) != 1 || got[0] != "value" {
		t.Fatalf("Params = %v, want [value]", got)
	}

	parameters, err := compiled.Params()
	if err != nil {
		t.Fatal(err)
	}
	parameters[0] = "changed"

	if got, err := compiled.Params(); err != nil || len(got) != 1 || got[0] != "value" {
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
		if output == nil || output.ContentType != "application/json" || string(output.Content) != want {
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

	if output == nil || output.ContentType != "application/json" || string(output.Content) != `["direct","rooted"]` {
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
	if err != nil || firstOutput == nil || len(firstOutput.Content) == 0 {
		t.Fatalf("Run first: output=%+v err=%v", firstOutput, err)
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
	if err != nil || secondOutput == nil {
		t.Fatalf("Run second: output=%+v err=%v", secondOutput, err)
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

	var native *ferretdiagnostics.Diagnostic
	if !errors.As(err, &native) {
		t.Fatalf("error does not preserve native diagnostic cause: %v", err)
	}

	const message = "Expected expression after 'RETURN'"
	const hint = "Did you forget to provide a value to return?"
	const label = "missing return value"
	wantSpan := nativesource.Span{Start: len(source.Content), End: len(source.Content)}
	wantPosition := nativesource.Position{Line: 1, Column: 7}
	if native.Source.Name() != source.Name || native.Source.Content() != source.Content ||
		native.Kind != parserd.SyntaxError || native.Message != message || native.Hint != hint ||
		len(native.Spans) != 1 || native.Spans[0] != ferretdiagnostics.NewMainErrorSpan(wantSpan, label) ||
		native.Source.PositionAt(native.Spans[0].Span) != wantPosition {
		t.Errorf("native diagnostic = %+v, want RETURN insertion at 1:7 / [6,6)", native)
	}

	var portable apidiagnostics.Diagnostics
	if !errors.As(err, &portable) || len(portable) != 1 {
		t.Fatalf("portable diagnostics = %+v, want one", portable)
	}

	if portable[0].Source != source || portable[0].Kind != apidiagnostics.Kind(parserd.SyntaxError) ||
		portable[0].Message != message || portable[0].Hint != hint || len(portable[0].Annotations) != 1 {
		t.Fatalf("portable diagnostic = %+v", portable[0])
	}

	annotation := portable[0].Annotations[0]
	if annotation.Range.SourceName != source.Name || annotation.Range.Position != wantPosition ||
		annotation.Range.Span != wantSpan || annotation.Message != label || !annotation.Primary {
		t.Fatalf("portable annotation = %+v, want primary %q at %s:1:7 / [6,6)", annotation, label, source.Name)
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

func TestPortableParamsRemainDetachedAfterClose(t *testing.T) {
	portable := newTestRuntime(t)
	plan, err := portable.Compile(t.Context(), api.NewAnonymousSource("RETURN @value"))
	if err != nil {
		t.Fatal(err)
	}

	if err := plan.Close(); err != nil {
		t.Fatal(err)
	}

	params, err := plan.Params()
	if err != nil || len(params) != 1 || params[0] != "value" {
		t.Fatalf("Params after close = %v, %v", params, err)
	}

	params[0] = "changed"
	if again, err := plan.Params(); err != nil || len(again) != 1 || again[0] != "value" {
		t.Fatalf("Params snapshot was aliased: %v, %v", again, err)
	}
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
