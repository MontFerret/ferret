package uapi

import (
	"context"
	"errors"
	"reflect"
	"testing"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestSessionOptionsMergeAndConvertWhenNativeAppliesThem(t *testing.T) {
	r := newTestRuntime(t, engine.WithParam("inherited", 7))
	p, err := r.Compile(t.Context(), api.NewAnonymousSource("RETURN [@inherited, @first, @second, @nil]"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = p.Close() })
	params := map[string]any{"first": 1, "second": 2, "nil": nil}
	empty := map[string]any{}
	var order []int
	s, err := p.NewSession(t.Context(), nil, api.WithParams(empty), func(opts api.SessionOptions) error {
		order = append(order, 1)
		empty["inherited"] = 8

		return opts.SetParams(params)
	}, func(opts api.SessionOptions) error {
		order = append(order, 2)
		params["second"] = 99
		if err := opts.SetParam("first", 3); err != nil {
			return err
		}

		return opts.SetOutputContentType(" application/json ")
	}, api.WithParams(nil), api.WithParams(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = s.Close() })
	params["first"] = 100
	out, err := s.Run(t.Context())
	if err != nil || string(out.Content) != "[8,3,99,null]" || out.ContentType != "application/json" || !reflect.DeepEqual(order, []int{1, 2}) {
		t.Fatalf("output=%+v err=%v order=%v", out, err, order)
	}
}

func TestCallbackFailuresAggregateBeforeNativeCalls(t *testing.T) {
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
	if !errors.Is(err, first) || !errors.Is(err, second) || errors.As(err, &invalid) || compileCalls != 0 || !reflect.DeepEqual(order, []int{1, 2}) {
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

func TestSessionValidationIsDelegatedToNative(t *testing.T) {
	native := newTestEngine(t)
	compiled, err := native.CompileDebug(t.Context(), source.NewAnonymous("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = compiled.Close() })
	portable := newPlan(compiled)
	for _, tc := range []struct {
		name     string
		portable []api.SessionOption
		native   []engine.SessionOption
	}{
		{"empty name", []api.SessionOption{api.WithParam("", 1)}, []engine.SessionOption{engine.WithSessionParam("", 1)}},
		{"nil value", []api.SessionOption{api.WithParam("value", nil)}, []engine.SessionOption{engine.WithSessionParam("value", nil)}},
		{"unsupported value", []api.SessionOption{api.WithParam("value", make(chan int))}, []engine.SessionOption{engine.WithSessionParam("value", make(chan int))}},
		{"unsupported map", []api.SessionOption{api.WithParams(map[string]any{"value": make(chan int)})}, []engine.SessionOption{engine.WithSessionParams(map[string]any{"value": make(chan int)})}},
		{"blank content type", []api.SessionOption{api.WithOutputContentType(" \t")}, []engine.SessionOption{engine.WithOutputContentType(" \t")}},
		{"blank root", []api.SessionOption{api.WithFSRoot(" \t")}, []engine.SessionOption{engine.WithSessionFSRoot(" \t")}},
		{"joined failures", []api.SessionOption{api.WithParam("", 1), api.WithOutputContentType(" ")}, []engine.SessionOption{engine.WithSessionParam("", 1), engine.WithOutputContentType(" ")}},
		{"invalid then valid", []api.SessionOption{api.WithOutputContentType(" "), api.WithOutputContentType("application/json")}, []engine.SessionOption{engine.WithOutputContentType(" "), engine.WithOutputContentType("application/json")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, debug := range []bool{false, true} {
				calls := 0
				option := func(opts api.SessionOptions) error {
					calls++
					for _, setter := range tc.portable {
						if err := setter(opts); err != nil {
							t.Fatalf("setter validated before Native application: %v", err)
						}
					}

					return nil
				}

				var got, want error
				if debug {
					s, err := portable.NewDebugSession(t.Context(), option)
					if s != nil {
						_ = s.Close()
						t.Fatal("invalid options published a debug session")
					}

					got = err
					sn, err := compiled.NewDebugSession(t.Context(), tc.native...)
					if sn != nil {
						_ = sn.Close()
						t.Fatal("invalid Native options published a debug session")
					}

					want = err
				} else {
					s, err := portable.NewSession(t.Context(), option)
					if s != nil {
						_ = s.Close()
						t.Fatal("invalid options published a session")
					}

					got = err
					sn, err := compiled.NewSession(t.Context(), tc.native...)
					if sn != nil {
						_ = sn.Close()
						t.Fatal("invalid Native options published a session")
					}

					want = err
				}

				if got == nil || want == nil || got.Error() != want.Error() || calls != 1 {
					t.Fatalf("debug=%t calls=%d got=%v want=%v", debug, calls, got, want)
				}
			}
		})
	}
}

func TestRuntimeRunValidatesNativeOptionsAfterCompileAndClosesPlan(t *testing.T) {
	var order []string
	closeErr := errors.New("plan cleanup")
	r := newTestRuntime(t,
		engine.WithBeforeCompileHook(func(context.Context) error {
			order = append(order, "compile")

			return nil
		}),
		engine.WithPlanCloseHook(func() error {
			order = append(order, "close plan")

			return closeErr
		}),
		engine.WithSessionCloseHook(func() error {
			t.Error("invalid options acquired a Native session")

			return nil
		}),
	)
	out, err := r.Run(t.Context(), api.NewAnonymousSource("RETURN 1"), func(opts api.SessionOptions) error {
		order = append(order, "options")

		return opts.SetOutputContentType(" ")
	})
	var invalid gooptions.ValidationError
	if !errors.As(err, &invalid) || invalid.Field != "output content type" || !errors.Is(err, closeErr) || out.Content != nil || out.ContentType != "" {
		t.Fatalf("output=%+v err=%v", out, err)
	}

	if !reflect.DeepEqual(order, []string{"options", "compile", "close plan"}) {
		t.Fatalf("Native execution order=%v", order)
	}
}

func TestDebugOptimizationValidationIsDelegatedToNative(t *testing.T) {
	for _, tc := range []struct {
		name     string
		portable []api.OptimizationLevel
		native   []engine.OptimizationLevel
		canceled bool
	}{
		{name: "basic", portable: []api.OptimizationLevel{api.OptimizationBasic}, native: []engine.OptimizationLevel{engine.OptimizationBasic}},
		{name: "full", portable: []api.OptimizationLevel{api.OptimizationFull}, native: []engine.OptimizationLevel{engine.OptimizationFull}},
		{name: "basic then none", portable: []api.OptimizationLevel{api.OptimizationBasic, api.OptimizationNone}, native: []engine.OptimizationLevel{engine.OptimizationBasic, engine.OptimizationNone}},
		{name: "canceled", portable: []api.OptimizationLevel{api.OptimizationBasic}, native: []engine.OptimizationLevel{engine.OptimizationBasic}, canceled: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var compileCalls int
			native := newTestEngine(t, engine.WithBeforeCompileHook(func(context.Context) error {
				compileCalls++

				return nil
			}))
			portable := Wrap(native)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tc.canceled {
				cancel()
			}

			calls := 0
			p, got := portable.CompileDebug(ctx, api.NewAnonymousSource("RETURN 1"), func(opts api.PlanOptions) error {
				calls++
				for _, level := range tc.portable {
					if err := opts.SetOptimizationLevel(level); err != nil {
						t.Fatalf("portable setter duplicated Native debug validation: %v", err)
					}
				}

				return nil
			})
			if p != nil {
				_ = p.Close()

				t.Fatal("invalid debug options published a plan")
			}

			options := make([]engine.PlanOption, 0, len(tc.native))
			for _, level := range tc.native {
				options = append(options, engine.WithPlanOptimizationLevel(level))
			}

			pn, want := native.CompileDebug(ctx, source.NewAnonymous("RETURN 1"), options...)
			if pn != nil {
				_ = pn.Close()

				t.Fatal("invalid Native debug options published a plan")
			}

			if got == nil || want == nil || got.Error() != want.Error() || calls != 1 || compileCalls != 0 {
				t.Fatalf("Native debug validation: got=%v want=%v callbacks=%d compile hooks=%d", got, want, calls, compileCalls)
			}
		})
	}
}
