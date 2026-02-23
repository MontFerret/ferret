package vm_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

type roundTripQueryable struct{}

func (q *roundTripQueryable) Query(_ context.Context, _ runtime.Query) (runtime.Value, error) {
	return runtime.NewString("ok"), nil
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
	expression  string
	expected    any
	options     []vm.EnvironmentOption
	description string
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
				convey.Convey(tc.expression, t, func() {
					prog, err := compilerInstance.Compile(file.NewSource(name, tc.expression))
					convey.So(err, convey.ShouldBeNil)

					opts := []vm.EnvironmentOption{
						vm.WithFunctions(base.Stdlib()),
					}
					opts = append(opts, tc.options...)

					expectedJSON, err := json.Marshal(tc.expected)
					convey.So(err, convey.ShouldBeNil)

					originalOut, err := base.Run(prog, opts...)
					convey.So(err, convey.ShouldBeNil)

					encoded, err := json.Marshal(prog)
					convey.So(err, convey.ShouldBeNil)

					var decoded bytecode.Program
					err = json.Unmarshal(encoded, &decoded)
					convey.So(err, convey.ShouldBeNil)

					roundTripOut, err := base.Run(&decoded, opts...)
					convey.So(err, convey.ShouldBeNil)

					convey.So(string(originalOut), convey.ShouldEqualJSON, string(expectedJSON))
					convey.So(string(roundTripOut), convey.ShouldEqualJSON, string(expectedJSON))
				})
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
			expression:  "RETURN @doc[~ text]",
			expected:    "ok",
			options:     []vm.EnvironmentOption{vm.WithParams(map[string]runtime.Value{"doc": queryable})},
			description: "Query literal without params",
		},
	})
}
