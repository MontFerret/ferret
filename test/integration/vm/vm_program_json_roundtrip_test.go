package vm_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	specexec "github.com/MontFerret/ferret/v2/test/spec/exec"
)

type roundTripQueryable struct{}

func (q *roundTripQueryable) Query(_ context.Context, _ runtime.Query) (runtime.List, error) {
	return runtime.NewArrayWith(runtime.NewString("ok")), nil
}

func (q *roundTripQueryable) MarshalJSON() ([]byte, error) {
	return json.Marshal("queryable")
}

func (q *roundTripQueryable) String() string {
	return "queryable"
}

func (q *roundTripQueryable) Hash() uint64 {
	return 0
}

func (q *roundTripQueryable) Copy() runtime.Value {
	return q
}

func TestProgramJSONRoundTrip(t *testing.T) {
	queryable := &roundTripQueryable{}

	specexec.RunRoundTrips(t, []specexec.RoundTrip{
		{
			Input:       spec.NewExpressionInput("RETURN TYPENAME(1.0)"),
			Expected:    "Float",
			Description: "Float constants preserve type",
		},
		{
			Input: spec.NewExpressionInput(`
LET users = [
  { gender: "m", age: 31 },
  { gender: "f", age: 25 },
  { gender: "m", age: 45 }
]
FOR u IN users
  COLLECT gender = u.gender
  AGGREGATE minAge = MIN(u.age)
  RETURN { gender, minAge }
`),
			Expected: []any{
				map[string]any{"gender": "f", "minAge": 25},
				map[string]any{"gender": "m", "minAge": 31},
			},
			Description: "Grouped aggregate round-trip",
		},
		{
			Input: spec.NewExpressionInput(`LET users = [
  { gender: "M", age: 30, salary: 100 },
  { gender: "F", age: 28, salary: 120 },
  { gender: "M", age: 40, salary: 90 }
]

FOR u IN users
  COLLECT g = u.gender
  AGGREGATE
    cnt = COUNT(u.age),
    sum = SUM(u.salary),
    minAge = MIN(u.age)
  RETURN { g, cnt, sum, minAge }`),
			Expected: []any{
				map[string]any{"g": "F", "cnt": 1, "sum": 120, "minAge": 28},
				map[string]any{"g": "M", "cnt": 2, "sum": 190, "minAge": 30},
			},
			Description: "Grouped aggregate round-trip with multiple documents",
		},
		{
			Input:       spec.NewExpressionInput("RETURN @doc[~ text]"),
			Expected:    []any{"ok"},
			Env:         []vm.EnvironmentOption{vm.WithParams(map[string]runtime.Value{"doc": queryable})},
			Description: "Query literal without params",
		},
		{
			Input: spec.NewExpressionInput(`
FUNC add(x, y) => x + y
RETURN add(2, 3)
`),
			Expected:    5,
			Description: "UDF round-trip",
		},
	})
}
