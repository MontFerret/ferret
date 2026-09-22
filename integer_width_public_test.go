package ferret_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2"
)

func TestCalendarBoundaryParameter(t *testing.T) {
	const maxEpochSeconds int64 = math.MaxInt64 - 62_135_596_800
	engine, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	location := time.FixedZone("UTC+01", 3600)
	for _, level := range []ferret.OptimizationLevel{ferret.OptimizationNone, ferret.OptimizationBasic, ferret.OptimizationFull} {
		for _, function := range []string{"add", "subtract"} {
			query := fmt.Sprintf("RETURN TO_INT(datetime::%s(@date, @amount, \"day\")) ON ERROR RETURN \"range\"", function)
			plan, err := engine.Compile(t.Context(), ferret.NewAnonymousSource(query), ferret.WithPlanOptimizationLevel(level))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = plan.Close() })

			for _, tc := range []struct {
				want    string
				seconds int64
				amount  int64
			}{
				{seconds: maxEpochSeconds, amount: 0, want: fmt.Sprint(maxEpochSeconds)},
				{seconds: maxEpochSeconds - 86400, amount: 1, want: fmt.Sprint(maxEpochSeconds)},
				{seconds: maxEpochSeconds, amount: -1, want: fmt.Sprint(maxEpochSeconds - 86400)},
				{seconds: maxEpochSeconds, amount: 1, want: `"range"`},
			} {
				amount := tc.amount
				if function == "subtract" {
					amount = -amount
				}

				session, err := plan.NewSession(t.Context(),
					ferret.WithSessionParam("date", time.Unix(tc.seconds, 123).In(location)),
					ferret.WithSessionParam("amount", amount))
				if err != nil {
					t.Fatal(err)
				}

				output, runErr := session.Run(t.Context())
				closeErr := session.Close()
				want := tc.want
				if math.MaxInt == math.MaxInt32 {
					want = `"range"`
				}

				if runErr != nil || closeErr != nil || output == nil || string(output.Content) != want {
					t.Fatalf("level %v %s(%d, %d) = %v, %v, %v; want %s", level, function, tc.seconds, amount, output, runErr, closeErr, want)
				}
			}
		}
	}
}

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
