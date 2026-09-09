package universal

import (
	"context"
	"errors"
	"reflect"
	"testing"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSessionOptionsMergeSnapshotAndValidateImmediately(t *testing.T) {
	r := newTestRuntime(t, engine.WithParam("inherited", 7))
	p, err := r.Compile(t.Context(), api.NewAnonymousSource("RETURN [@inherited, @first, @second, @nil]"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = p.Close() })
	params := map[string]any{"first": 1, "second": 2, "nil": nil}
	var order []int
	s, err := p.NewSession(t.Context(), nil, func(opts api.SessionOptions) error {
		order = append(order, 1)

		return opts.SetParams(params)
	}, func(opts api.SessionOptions) error {
		order = append(order, 2)
		params["second"] = 99
		if err := opts.SetParam("first", 3); err != nil {
			return err
		}

		if err := opts.SetOutputContentType(" application/json "); err != nil {
			return err
		}

		for _, setter := range []func() error{
			func() error { return opts.SetParam("", 1) },
			func() error { return opts.SetParam("first", nil) },
			func() error { return opts.SetParams(map[string]any{"first": make(chan int)}) },
			func() error { return opts.SetOutputContentType(" \t") },
			func() error { return opts.SetFSRoot(" \t") },
		} {
			if err := setter(); err == nil {
				t.Fatal("setter did not report invalid input immediately")
			}
		}

		return nil // Rejected setters must leave the valid earlier options intact.
	}, api.WithParams(nil), api.WithParams(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = s.Close() })
	params["first"] = 100
	out, err := s.Run(t.Context())
	if err != nil || string(out.Content) != "[7,3,2,null]" || out.ContentType != "application/json" || !reflect.DeepEqual(order, []int{1, 2}) {
		t.Fatalf("output=%+v err=%v order=%v", out, err, order)
	}
}

func TestInvalidOptionsAggregateBeforeNativeCalls(t *testing.T) {
	compileCalls := 0
	r := newTestRuntime(t, engine.WithBeforeCompileHook(func(context.Context) error {
		compileCalls++

		return nil
	}))
	first, second := errors.New("first"), errors.New("unsupported runtime extension")
	var order []int
	_, err := r.Run(t.Context(), api.NewAnonymousSource("RETURN 1"),
		func(api.SessionOptions) error {
			order = append(order, 1)

			return first
		},
		nil,
		api.WithOutputContentType(" "),
		func(opts api.SessionOptions) error {
			order = append(order, 2)
			if _, ok := opts.(interface{ SetRemoteWorkspace(string) error }); ok {
				t.Fatal("adapter exposes an unrelated extension")
			}

			return second
		},
	)
	var invalid gooptions.ValidationError
	if !errors.Is(err, first) || !errors.Is(err, second) || !errors.As(err, &invalid) || invalid.Field != "output content type" || compileCalls != 0 || !reflect.DeepEqual(order, []int{1, 2}) {
		t.Fatalf("err=%v validation=%+v compile=%d order=%v", err, invalid, compileCalls, order)
	}

	_, err = r.Compile(t.Context(), api.NewAnonymousSource("RETURN 1"),
		func(api.PlanOptions) error { return first }, nil,
		func(api.PlanOptions) error { return second },
		api.WithOptimizationLevel(api.OptimizationAggressive))
	if !errors.Is(err, first) || !errors.Is(err, second) || compileCalls != 0 {
		t.Fatalf("plan options reached Native or lost errors: %v calls=%d", err, compileCalls)
	}
}

func TestContextAndClosedAdmissionSkipOptions(t *testing.T) {
	r := newTestRuntime(t)
	p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = p.Close() })
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	calls := 0
	po := func(api.PlanOptions) error {
		calls++

		return nil
	}

	so := func(api.SessionOptions) error {
		calls++

		return nil
	}

	for _, mode := range []string{"nil", "canceled", "closed"} {
		ctx := t.Context()
		want := runtime.ErrInvalidOperation
		switch mode {
		case "nil":
			ctx = nil
			want = runtime.ErrInvalidArgument
		case "canceled":
			ctx = canceled
			want = context.Canceled
		case "closed":
			_ = p.Close()
			_ = r.Close()
		}

		for _, call := range []func() error{
			func() error {
				_, err := r.Compile(ctx, api.NewAnonymousSource("RETURN 1"), po)

				return err
			},
			func() error {
				_, err := r.CompileDebug(ctx, api.NewAnonymousSource("RETURN 1"), po)

				return err
			},
			func() error {
				_, err := r.Run(ctx, api.NewAnonymousSource("RETURN 1"), so)

				return err
			},
			func() error {
				_, err := p.NewSession(ctx, so)

				return err
			},
			func() error {
				_, err := p.NewDebugSession(ctx, so)

				return err
			},
		} {
			if err := call(); !errors.Is(err, want) {
				t.Fatalf("%s admission: %v, want %v", mode, err, want)
			}
		}
	}

	if calls != 0 {
		t.Fatalf("rejected admission invoked %d options", calls)
	}
}

func TestPlanOptionsValidateEverySetterWithoutMutatingDefaults(t *testing.T) {
	r := newTestRuntime(t, engine.WithOptimizationLevel(engine.OptimizationBasic))
	calls := 0
	p, err := r.Compile(t.Context(), api.NewAnonymousSource("RETURN 1"), nil, func(opts api.PlanOptions) error {
		calls++
		if err := opts.SetOptimizationLevel(api.OptimizationNone); err != nil {
			return err
		}

		if err := opts.SetOptimizationLevel(api.OptimizationAggressive); err == nil {
			t.Fatal("unsupported optimization accepted")
		}

		return opts.SetOptimizationLevel(api.OptimizationFull)
	})
	if err != nil || calls != 1 {
		t.Fatalf("compile err=%v callbacks=%d", err, calls)
	}

	t.Cleanup(func() { _ = p.Close() })
	if debug, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 1")); err != nil {
		t.Fatal(err)
	} else {
		_ = debug.Close()
	}

	if s, err := p.NewDebugSession(t.Context()); err == nil || s != nil {
		t.Fatalf("ordinary plan accepted debugging: session=%v err=%v", s, err)
	}
}
