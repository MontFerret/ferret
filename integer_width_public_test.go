package ferret_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2"
)

func TestIntegerWidthParameterRoundTrip(t *testing.T) {
	engine, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	query := "RETURN [@value, math::floor(@value), math::ceil(@value), math::round(@value), encoding::json_parse(TO_STRING(@value)), TYPENAME(encoding::json_parse(TO_STRING(@value)))]"
	for _, level := range []ferret.OptimizationLevel{ferret.OptimizationNone, ferret.OptimizationBasic, ferret.OptimizationFull} {
		plan, err := engine.Compile(t.Context(), ferret.NewAnonymousSource(query), ferret.WithPlanOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = plan.Close() })

		for _, value := range []int64{math.MinInt64, math.MinInt32 - 1, math.MinInt32, math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1, math.MaxInt64} {
			session, err := plan.NewSession(t.Context(), ferret.WithSessionParam("value", value))
			if err != nil {
				t.Fatal(err)
			}
			output, runErr := session.Run(t.Context())
			closeErr := session.Close()
			want := fmt.Sprintf("[%d,%d,%d,%d,%d,\"Int\"]", value, value, value, value, value)
			if runErr != nil || closeErr != nil || output == nil || string(output.Content) != want {
				t.Fatalf("level %v value %d = %v, run %v, close %v; want %s", level, value, output, runErr, closeErr, want)
			}
		}
	}
}

func TestIntegerWidthRangeErrorRecovery(t *testing.T) {
	engine, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	for _, level := range []ferret.OptimizationLevel{ferret.OptimizationNone, ferret.OptimizationBasic, ferret.OptimizationFull} {
		for _, query := range []string{
			"RETURN math::floor(9223372036854775808.0) ON ERROR RETURN \"range\"",
			"RETURN math::ceil(9223372036854775808.0) ON ERROR RETURN \"range\"",
			"RETURN math::round(9223372036854775808.0) ON ERROR RETURN \"range\"",
			"RETURN datetime::add(datetime::parse(\"2024-01-01T00:00:00Z\"), 9223372036854775807, \"second\") ON ERROR RETURN \"range\"",
		} {
			plan, err := engine.Compile(t.Context(), ferret.NewAnonymousSource(query), ferret.WithPlanOptimizationLevel(level))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = plan.Close() })
			session, err := plan.NewSession(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			output, runErr := session.Run(t.Context())
			closeErr := session.Close()
			if runErr != nil || closeErr != nil || output == nil || string(output.Content) != "\"range\"" {
				t.Fatalf("level %v %s = %v, %v, %v", level, query, output, runErr, closeErr)
			}
		}
	}
}
