package universal

import (
	"context"
	"testing"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apisource "github.com/MontFerret/api/source"
	"github.com/MontFerret/ferret/v2/pkg/engine"
)

func TestRuntimeOptimizationOptions(t *testing.T) {
	for _, nativeLevel := range []engine.OptimizationLevel{
		engine.OptimizationNone, engine.OptimizationBasic, engine.OptimizationFull,
	} {
		engine, err := engine.New(engine.WithOptimizationLevel(nativeLevel))
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() { _ = engine.Close() })
		runtime := New(engine)
		t.Cleanup(func() { _ = runtime.Close() })
		for _, debug := range []bool{false, true} {
			compile := runtime.Compile

			if debug {
				compile = runtime.CompileDebug
			}

			for _, level := range []api.OptimizationLevel{
				api.OptimizationNone, api.OptimizationBasic, api.OptimizationFull,
				api.OptimizationAggressive, api.OptimizationLevel(-1),
			} {
				plan, err := compile(context.Background(), api.NewAnonymousSource("RETURN 1"), api.WithOptimizationLevel(level))

				wantSuccess := (debug && level == api.OptimizationNone) || (!debug && level >= api.OptimizationNone && level <= api.OptimizationFull)
				if (err == nil) != wantSuccess || (plan != nil) != wantSuccess {
					t.Fatalf("native=%v debug=%t option=%v: plan=%v error=%v", nativeLevel, debug, level, plan, err)
				}

				if plan != nil {
					if err := plan.Close(); err != nil {
						t.Fatal(err)
					}
				}
			}

			plan, err := compile(context.Background(), api.NewAnonymousSource("RETURN 1"))
			if err != nil {
				t.Fatalf("omitted optimization: %v", err)
			}

			if err := plan.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestDebuggerRejectsUnknownBreakpointModes(t *testing.T) {
	runtime := newTestRuntime(t)

	plan, err := runtime.CompileDebug(context.Background(), api.NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = session.Close() })

	for _, mode := range []apidebugger.BreakpointBindingMode{-1, 99} {
		_, err := session.SetBreakpointAt(apisource.Location{Position: apisource.Position{Line: 1, Column: 1}}, apidebugger.BreakpointOptions{BindingMode: mode})
		if err == nil {
			t.Fatalf("mode %d was silently accepted", mode)
		}
	}

	if values := session.Breakpoints(); len(values) != 0 {
		t.Fatalf("invalid options installed breakpoints: %+v", values)
	}
}

func TestRuntimeRejectsInvalidParameterValues(t *testing.T) {
	runtime := newTestRuntime(t)

	plan, err := runtime.CompileDebug(context.Background(), api.NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	option := api.WithParams(map[string]any{"invalid": make(chan int)})
	if session, err := plan.NewSession(context.Background(), option); err == nil || session != nil {
		t.Fatalf("ordinary invalid parameters: session=%v error=%v", session, err)
	}

	if session, err := plan.NewDebugSession(context.Background(), option); err == nil || session != nil {
		t.Fatalf("debug invalid parameters: session=%v error=%v", session, err)
	}
}
