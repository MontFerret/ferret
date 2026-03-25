package vm_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	specassert "github.com/MontFerret/ferret/v2/test/spec/assert"
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

type roundTripCase struct {
	expected    any
	expression  string
	description string
	options     []vm.EnvironmentOption
}

func (c roundTripCase) name() string {
	if c.description != "" {
		return c.description
	}

	return c.expression
}

func runProgramRoundTrip(t *testing.T, cases []roundTripCase) {
	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		compilerInstance := compiler.New(compiler.WithOptimizationLevel(level))

		for _, tc := range cases {
			name := fmt.Sprintf("Program JSON RoundTrip: %s (O%d)", tc.name(), level)

			t.Run(name, func(t *testing.T) {
				prog, err := compilerInstance.Compile(file.NewSource(name, tc.expression))
				if err != nil {
					t.Fatalf("compile failed: %v", err)
				}

				opts := []vm.EnvironmentOption{
					vm.WithNamespace(spec.Stdlib()),
				}
				opts = append(opts, tc.options...)

				originalOut, err := spec.Run(prog, opts...)
				if err != nil {
					t.Fatalf("original run failed: %v", err)
				}

				encoded, err := json.Marshal(prog)
				if err != nil {
					t.Fatalf("program JSON marshal failed: %v", err)
				}

				var decoded bytecode.Program
				if err := json.Unmarshal(encoded, &decoded); err != nil {
					t.Fatalf("program JSON unmarshal failed: %v", err)
				}

				roundTripOut, err := spec.Run(&decoded, opts...)
				if err != nil {
					t.Fatalf("round-trip run failed: %v", err)
				}

				specassert.Expect(t, specassert.EqualJSON, json.RawMessage(originalOut), tc.expected, "original program output mismatch")
				specassert.Expect(t, specassert.EqualJSON, json.RawMessage(roundTripOut), tc.expected, "round-trip program output mismatch")
			})
		}
	}
}

func TestProgramJSONRoundTrip(t *testing.T) {
	queryable := &roundTripQueryable{}

	runProgramRoundTrip(t, []roundTripCase{
		{
			expression:  "RETURN TYPENAME(1.0)",
			expected:    "Float",
			description: "Float constants preserve type",
		},
		{
			expression: `
LET users = [
  { gender: "m", age: 31 },
  { gender: "f", age: 25 },
  { gender: "m", age: 45 }
]
FOR u IN users
  COLLECT gender = u.gender
  AGGREGATE minAge = MIN(u.age)
  RETURN { gender, minAge }
`,
			expected: []any{
				map[string]any{"gender": "f", "minAge": 25},
				map[string]any{"gender": "m", "minAge": 31},
			},
			description: "Grouped aggregate round-trip",
		},
		{
			expression: `LET users = [
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
  RETURN { g, cnt, sum, minAge }`,
			expected: []any{
				map[string]any{"g": "F", "cnt": 1, "sum": 120, "minAge": 28},
				map[string]any{"g": "M", "cnt": 2, "sum": 190, "minAge": 30},
			},
			description: "Grouped aggregate round-trip with multiple documents",
		},
		{
			expression:  "RETURN @doc[~ text]",
			expected:    []any{"ok"},
			options:     []vm.EnvironmentOption{vm.WithParams(map[string]runtime.Value{"doc": queryable})},
			description: "Query literal without params",
		},
		{
			expression: `
FUNC add(x, y) => x + y
RETURN add(2, 3)
`,
			expected:    5,
			description: "UDF round-trip",
		},
	})
}
