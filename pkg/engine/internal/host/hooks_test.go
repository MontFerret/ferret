package host

import (
	"context"
	"errors"
	"slices"
	"testing"
)

type hookStage struct {
	register func(*Hooks, func() error)
	run      func(*Hooks) error
	name     string
	unwind   bool
}

func TestHookOrderingAndFailures(t *testing.T) {
	for _, stage := range hookStages(t.Context()) {
		t.Run(stage.name, func(t *testing.T) {
			hooks := NewHooks()
			second, third := errors.New("second"), errors.New("third")
			var calls []int

			for i, failure := range []error{nil, second, third} {
				stage.register(hooks, func() error {
					calls = append(calls, i+1)

					return failure
				})
			}

			err := stage.run(hooks)
			want := []int{1, 2}
			if stage.unwind {
				want = []int{3, 2, 1}
			}

			if !slices.Equal(calls, want) {
				t.Fatalf("calls = %v, want %v", calls, want)
			}

			if !errors.Is(err, second) || errors.Is(err, third) != stage.unwind {
				t.Fatalf("lost hook error identity: %v", err)
			}

			if !stage.unwind && err != second {
				t.Fatalf("before hook failure was wrapped: %v", err)
			}
		})
	}
}

func TestHookSnapshotsAreIndependent(t *testing.T) {
	for _, stage := range hookStages(t.Context()) {
		t.Run(stage.name, func(t *testing.T) {
			hooks := NewHooks()
			var calls []int
			register := func(target *Hooks, value int) {
				stage.register(target, func() error {
					calls = append(calls, value)

					return nil
				})
			}

			// Three entries leave spare capacity in the mutable registrar. Appending
			// on both sides catches a snapshot that accidentally shares its backing array.
			for i := 1; i <= 3; i++ {
				register(hooks, i)
			}

			snapshot := hooks.Clone()
			register(hooks, 4)
			register(snapshot, 5)

			for i, target := range []*Hooks{hooks, snapshot} {
				calls = nil
				if err := stage.run(target); err != nil {
					t.Fatal(err)
				}

				want := []int{1, 2, 3, 4 + i}
				if stage.unwind {
					slices.Reverse(want)
				}

				if !slices.Equal(calls, want) {
					t.Fatalf("snapshot %d calls = %v, want %v", i, calls, want)
				}
			}
		})
	}
}

func TestHookContextsAndPrimaryErrors(t *testing.T) {
	type contextKey struct{}

	hooks := NewHooks()
	ctx := t.Context()
	derived := context.WithValue(ctx, contextKey{}, "derived")
	hooks.Session().BeforeRun(func(got context.Context) (context.Context, error) {
		if got != ctx {
			t.Error("first hook lost caller context")
		}

		return derived, nil
	})
	hooks.Session().BeforeRun(func(got context.Context) (context.Context, error) {
		if got != derived {
			t.Error("second hook lost replacement context")
		}

		return got, nil
	})

	got, err := hooks.SessionHooks.RunBeforeRun(ctx)
	if err != nil || got != derived {
		t.Fatalf("hook context = %v, error = %v", got, err)
	}

	primary, failure := errors.New("primary"), errors.New("hook")
	for _, stage := range []string{"compile", "run"} {
		t.Run(stage, func(t *testing.T) {
			calls := 0
			callback := func(got context.Context, err error) error {
				calls++
				if got != derived || err != primary {
					t.Errorf("after hook context/error = %v/%v", got, err)
				}

				return failure
			}

			var err error
			if stage == "compile" {
				hooks.Plan().AfterCompile(callback)
				hooks.Plan().AfterCompile(callback)
				err = hooks.PlanHooks.RunAfterCompile(derived, primary)
			} else {
				hooks.Session().AfterRun(callback)
				hooks.Session().AfterRun(callback)
				err = hooks.SessionHooks.RunAfterRun(derived, primary)
			}

			if calls != 2 || !errors.Is(err, failure) || errors.Is(err, primary) {
				t.Fatalf("after hooks calls/error = %d/%v", calls, err)
			}
		})
	}
}

func TestNilHooksAreIgnored(t *testing.T) {
	hooks := NewHooks()
	hooks.Engine().OnInit(nil)
	hooks.Engine().OnClose(nil)
	hooks.Plan().BeforeCompile(nil)
	hooks.Plan().AfterCompile(nil)
	hooks.Plan().OnClose(nil)
	hooks.Session().BeforeRun(nil)
	hooks.Session().AfterRun(nil)
	hooks.Session().OnClose(nil)

	for _, stage := range hookStages(t.Context()) {
		if err := stage.run(hooks.Clone()); err != nil {
			t.Fatalf("%s: %v", stage.name, err)
		}
	}
}

func hookStages(ctx context.Context) []hookStage {
	return []hookStage{
		{name: "init", register: func(h *Hooks, f func() error) { h.Engine().OnInit(f) }, run: func(h *Hooks) error { return h.EngineHooks.RunInit() }},
		{name: "engine close", unwind: true, register: func(h *Hooks, f func() error) { h.Engine().OnClose(f) }, run: func(h *Hooks) error { return h.EngineHooks.RunClose() }},
		{name: "plan close", unwind: true, register: func(h *Hooks, f func() error) { h.Plan().OnClose(f) }, run: func(h *Hooks) error { return h.PlanHooks.RunClose() }},
		{name: "session close", unwind: true, register: func(h *Hooks, f func() error) { h.Session().OnClose(f) }, run: func(h *Hooks) error { return h.SessionHooks.RunClose() }},
		{name: "before compile", register: func(h *Hooks, f func() error) {
			h.Plan().BeforeCompile(func(context.Context) error { return f() })
		}, run: func(h *Hooks) error { return h.PlanHooks.RunBeforeCompile(ctx) }},
		{name: "after compile", unwind: true, register: func(h *Hooks, f func() error) {
			h.Plan().AfterCompile(func(context.Context, error) error { return f() })
		}, run: func(h *Hooks) error { return h.PlanHooks.RunAfterCompile(ctx, nil) }},
		{name: "before run", register: func(h *Hooks, f func() error) {
			h.Session().BeforeRun(func(ctx context.Context) (context.Context, error) { return ctx, f() })
		}, run: func(h *Hooks) error {
			_, err := h.SessionHooks.RunBeforeRun(ctx)

			return err
		}},
		{name: "after run", unwind: true, register: func(h *Hooks, f func() error) {
			h.Session().AfterRun(func(context.Context, error) error { return f() })
		}, run: func(h *Hooks) error { return h.SessionHooks.RunAfterRun(ctx, nil) }},
	}
}
