package ferret_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const randomSequenceQuery = `RETURN [random::float(), random::int(-5, @max), random::choice([10, 20, 30]), random::bool(), random::shuffle([1, 2, 3, 4]), rand(), rand(8), rand(6, 1), random::float(-3, 9)]`

func TestRandomSeededSessions(t *testing.T) {
	poison := rnd.NewSeed(999)
	eng, err := ferret.New(ferret.WithBeforeRunHook(func(context.Context) (context.Context, error) {
		return rnd.WithContext(context.Background(), poison), nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = eng.Close() })

	for _, level := range []ferret.OptimizationLevel{ferret.OptimizationNone, ferret.OptimizationBasic, ferret.OptimizationFull} {
		plan, err := eng.Compile(t.Context(), ferret.NewAnonymousSource(randomSequenceQuery), ferret.WithPlanOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = plan.Close() })

		for _, seed := range []int64{0, -1, 42} {
			option := ferret.WithSessionRandomSeed(seed)
			var sessions []*ferret.Session
			for range 2 {
				session, err := plan.NewSession(t.Context(), ferret.WithSessionParam("max", 7), ferret.WithSessionRandomSeed(123), option)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = session.Close() })
				sessions = append(sessions, session)
			}

			control := rnd.NewSeed(seed)
			for range 3 {
				want := expectedRandomSequence(t, control)
				for _, session := range sessions {
					output, err := session.Run(rnd.WithContext(t.Context(), poison))
					if err != nil || output == nil || string(output.Content) != want {
						t.Fatalf("level %v seed %d: output %v, error %v; want %s", level, seed, output, err, want)
					}
				}
			}

			for _, session := range sessions {
				if err := session.Close(); err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	if poison.Float64() != rnd.NewSeed(999).Float64() {
		t.Fatal("execution consumed caller or hook RNG instead of Session RNG")
	}
}

func TestRandomSeedSurvivesDebugResumes(t *testing.T) {
	eng, err := ferret.New(ferret.WithBeforeRunHook(func(context.Context) (context.Context, error) {
		return context.Background(), nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = eng.Close() })
	plan, err := eng.CompileDebug(t.Context(), ferret.NewAnonymousSource(`
LET first = random::float()
LET second = random::int(-5, @max)
LET chosen = random::choice([10, 20, 30])
LET third = random::bool()
LET shuffled = random::shuffle([1, 2, 3, 4])
LET fourth = rand()
LET fifth = rand(8)
LET sixth = rand(6, 1)
LET seventh = random::float(-3, 9)
RETURN [first, second, chosen, third, shuffled, fourth, fifth, sixth, seventh]`))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = plan.Close() })

	for _, step := range []bool{false, true} {
		session, err := plan.NewDebugSession(t.Context(), ferret.WithSessionRandomSeed(42), ferret.WithSessionParam("max", 7))
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = session.Close() })
		event, err := session.Start(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		for resumes := 0; event.Reason != debugger.ReasonCompleted; resumes++ {
			if resumes > 100 {
				t.Fatal("debug execution did not finish")
			}

			ctx := rnd.WithContext(t.Context(), rnd.NewSeed(123))
			if step {
				event, err = session.StepIn(ctx)
			} else {
				event, err = session.Continue(ctx)
			}

			if err != nil {
				t.Fatal(err)
			}
		}

		if want := expectedRandomSequence(t, rnd.NewSeed(42)); event.Output == nil || string(event.Output.Content) != want {
			t.Fatalf("debug output %v, want %s", event.Output, want)
		}
	}
}

func TestDefaultSessionsOwnDistinctSources(t *testing.T) {
	var captured []*rnd.Source
	eng, err := ferret.New(ferret.WithMaxIdleVMsPerPlan(1), ferret.WithMaxVMsPerPlan(1),
		ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
			ns.Function().A0().Add("capture_rng", func(ctx context.Context) (runtime.Value, error) {
				src, ok := rnd.FromContext(ctx)
				if !ok {
					return runtime.None, runtime.Error(runtime.ErrUnexpected, "missing Session source")
				}

				captured = append(captured, src)

				return runtime.True, nil
			})
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = eng.Close() })
	plan, err := eng.CompileDebug(t.Context(), ferret.NewAnonymousSource("RETURN CAPTURE_RNG()"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = plan.Close() })

	for range 2 {
		session, err := plan.NewSession(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = session.Close() })
		for range 2 {
			if _, err := session.Run(t.Context()); err != nil {
				t.Fatal(err)
			}
		}

		if err := session.Close(); err != nil {
			t.Fatal(err)
		}
	}

	for range 2 {
		session, err := plan.NewDebugSession(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = session.Close() })
		if _, err := session.Start(t.Context()); err != nil {
			t.Fatal(err)
		}

		if event, err := session.Continue(t.Context()); err != nil || event.Reason != debugger.ReasonCompleted {
			t.Fatalf("debug completion: %v, %v", event, err)
		}
	}

	if len(captured) != 6 || captured[0] != captured[1] || captured[2] != captured[3] {
		t.Fatal("a Session replaced its source between runs")
	}

	unique := map[*rnd.Source]bool{}
	for _, index := range []int{0, 2, 4, 5} {
		if unique[captured[index]] {
			t.Fatal("ordinary, debug, or pooled-VM Sessions shared mutable state")
		}

		unique[captured[index]] = true
	}
}

func TestRandomWaitForJitterUsesSessionSequence(t *testing.T) {
	attempts := 0
	eng, err := ferret.New(ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("ready", func(context.Context) (runtime.Value, error) {
			attempts++

			return runtime.Boolean(attempts == 3), nil
		})
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = eng.Close() })

	for _, level := range []ferret.OptimizationLevel{ferret.OptimizationNone, ferret.OptimizationBasic, ferret.OptimizationFull} {
		plan, err := eng.Compile(t.Context(), ferret.NewAnonymousSource(`
LET first = random::float()
LET ok = WAITFOR READY() EVERY 1ms JITTER 0.5
RETURN [first, ok, random::float()]`), ferret.WithPlanOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = plan.Close() })
		for range 2 {
			attempts = 0
			session, err := plan.NewSession(t.Context(), ferret.WithSessionRandomSeed(42))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = session.Close() })
			control := rnd.NewSeed(42)
			first := control.Float64()
			control.Float64()
			control.Float64()
			want, _ := json.Marshal([]any{first, true, control.Float64()})
			out, err := session.Run(t.Context())
			if err != nil || out == nil || string(out.Content) != string(want) || attempts != 3 {
				t.Fatalf("jitter output %v, error %v, attempts %d; want %s", out, err, attempts, want)
			}
		}
	}
}

func TestRandomFQLArity(t *testing.T) {
	eng, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = eng.Close() })
	for _, call := range []string{"random::float(1)", "random::float(1, 2, 3)", "random::int()", "random::int(1)", "random::int(1, 2, 3)", "random::bool(1)", "random::choice()", "random::choice([], [])", "random::shuffle()", "random::shuffle([], [])", "rand(1, 2, 3)"} {
		if _, err := eng.Run(t.Context(), ferret.NewAnonymousSource("RETURN "+call), ferret.WithSessionRandomSeed(42)); err == nil {
			t.Fatalf("accepted unsupported arity %s", call)
		}
	}
}

func TestRandomFQLNumericContracts(t *testing.T) {
	eng, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = eng.Close() })
	for _, level := range []ferret.OptimizationLevel{ferret.OptimizationNone, ferret.OptimizationBasic, ferret.OptimizationFull} {
		for _, test := range []struct {
			min     runtime.Value
			max     runtime.Value
			wantErr error
			name    string
		}{
			{name: "int", min: runtime.Int(1), max: runtime.Int(6), wantErr: nil},
			{name: "int", min: runtime.Int(-10), max: runtime.Int(-1), wantErr: nil},
			{name: "int", min: runtime.Int(-10), max: runtime.Int(10), wantErr: nil},
			{name: "int", min: runtime.Int(math.MinInt64), max: runtime.Int(math.MaxInt64), wantErr: nil},
			{name: "int", min: runtime.Int(math.MaxInt64 - 1), max: runtime.Int(math.MaxInt64), wantErr: nil},
			{name: "int", min: runtime.Int(math.MinInt64), max: runtime.Int(math.MinInt64), wantErr: nil},
			{name: "int", min: runtime.Int(math.MaxInt64), max: runtime.Int(math.MaxInt64), wantErr: nil},
			{name: "int", min: runtime.Float(1), max: runtime.Int(6), wantErr: runtime.ErrInvalidType},
			{name: "int", min: runtime.Int(1), max: runtime.Float(6), wantErr: runtime.ErrInvalidType},
			{name: "int", min: runtime.Int(6), max: runtime.Int(1), wantErr: runtime.ErrInvalidArgument},
			{name: "float", min: runtime.Float(-7.5), max: runtime.Int(-1), wantErr: nil},
			{name: "float", min: runtime.Int(-10), max: runtime.Float(1.5), wantErr: nil},
			{name: "float", min: runtime.Int(5), max: runtime.Float(5), wantErr: nil},
			{name: "float", min: runtime.Float(-math.MaxFloat64), max: runtime.Float(math.MaxFloat64), wantErr: nil},
			{name: "float", min: runtime.Int(1<<53 + 1), max: runtime.Int(1<<53 + 4), wantErr: nil},
			{name: "float", min: runtime.Int(1<<53 + 1), max: runtime.Int(1<<53 + 2), wantErr: runtime.ErrInvalidArgument},
			{name: "float", min: runtime.Int(6), max: runtime.Int(1), wantErr: runtime.ErrInvalidArgument},
			{name: "float", min: runtime.Float(math.NaN()), max: runtime.Int(1), wantErr: runtime.ErrInvalidArgument},
			{name: "float", min: runtime.Int(0), max: runtime.Float(math.NaN()), wantErr: runtime.ErrInvalidArgument},
			{name: "float", min: runtime.Float(math.Inf(-1)), max: runtime.Int(1), wantErr: runtime.ErrInvalidArgument},
			{name: "float", min: runtime.Int(0), max: runtime.Float(math.Inf(1)), wantErr: runtime.ErrInvalidArgument},
			{name: "float", min: runtime.String("0"), max: runtime.Int(1), wantErr: runtime.ErrInvalidType},
			{name: "float", min: runtime.Int(0), max: runtime.True, wantErr: runtime.ErrInvalidType},
		} {
			comparison := "value <= @max"
			if test.name == "float" {
				comparison = "(@min == @max ? value == @max : value < @max)"
			}
			query := "LET value = random::" + test.name + "(@min, @max) RETURN value >= @min AND " + comparison
			plan, err := eng.Compile(t.Context(), ferret.NewAnonymousSource(query), ferret.WithPlanOptimizationLevel(level))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = plan.Close() })
			session, err := plan.NewSession(t.Context(), ferret.WithSessionRandomSeed(42), ferret.WithSessionRuntimeParams(runtime.Params{"min": test.min, "max": test.max}))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = session.Close() })
			out, err := session.Run(t.Context())
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("level %v %s(%v,%v): %v, want %v", level, test.name, test.min, test.max, err, test.wantErr)
				}
			} else if err != nil || out == nil || string(out.Content) != "true" {
				t.Fatalf("level %v %s(%v,%v): output %v, error %v", level, test.name, test.min, test.max, out, err)
			}
			if err := session.Close(); err != nil {
				t.Fatal(err)
			}
			if err := plan.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func expectedRandomSequence(t *testing.T, src *rnd.Source) string {
	t.Helper()

	values := []any{src.Float64(), src.Int64(-5, 7)}
	choice := 10
	for ordinal := int64(2); ordinal <= 3; ordinal++ {
		if src.Int64(0, ordinal-1) == 0 {
			choice = int(ordinal) * 10
		}
	}

	values = append(values, choice, src.Bool())
	shuffled := []int{1, 2, 3, 4}
	for i := len(shuffled) - 1; i > 0; i-- {
		j := int(src.Int64(0, int64(i)))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	values = append(values, shuffled, src.Float64(), src.LegacyFloat64(16, 4), src.LegacyFloat64(6, 1), -3+src.Float64()*12)
	content, err := json.Marshal(values)
	if err != nil {
		t.Fatal(fmt.Errorf("encode expected random sequence: %w", err))
	}

	return string(content)
}
