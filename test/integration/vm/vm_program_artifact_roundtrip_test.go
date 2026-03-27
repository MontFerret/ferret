package vm_test

import (
	"encoding/json"
	"maps"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	specassert "github.com/MontFerret/ferret/v2/test/spec/assert"
)

type artifactRoundTripCase struct {
	Expected    any
	Check       func(t *testing.T, original *bytecode.Program, decoded *bytecode.Program)
	Input       spec.Input
	Description string
	Env         []vm.EnvironmentOption
}

func TestProgramArtifactRoundTrip(t *testing.T) {
	queryable := &roundTripQueryable{}
	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}
	cases := []artifactRoundTripCase{
		{
			Input:       spec.NewExpressionInput("RETURN TYPENAME(1.0)"),
			Expected:    "Float",
			Description: "Float constants preserve type",
			Check: func(t *testing.T, original *bytecode.Program, decoded *bytecode.Program) {
				t.Helper()

				if !maps.Equal(original.Functions.Host, decoded.Functions.Host) {
					t.Fatalf("host metadata mismatch: got %#v, want %#v", decoded.Functions.Host, original.Functions.Host)
				}
			},
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
			Check: func(t *testing.T, original *bytecode.Program, decoded *bytecode.Program) {
				t.Helper()

				if !reflect.DeepEqual(original.Metadata.AggregatePlans, decoded.Metadata.AggregatePlans) {
					t.Fatalf("aggregate plans mismatch: got %#v, want %#v", decoded.Metadata.AggregatePlans, original.Metadata.AggregatePlans)
				}

				if !reflect.DeepEqual(original.Metadata.AggregateSelectorSlots, decoded.Metadata.AggregateSelectorSlots) {
					t.Fatalf("aggregate selector slots mismatch: got %#v, want %#v", decoded.Metadata.AggregateSelectorSlots, original.Metadata.AggregateSelectorSlots)
				}
			},
		},
		{
			Input:       spec.NewExpressionInput("RETURN @doc[~ text]"),
			Expected:    []any{"ok"},
			Env:         []vm.EnvironmentOption{vm.WithParams(runtime.Params{"doc": queryable})},
			Description: "Query literal without params",
		},
		{
			Input: spec.NewExpressionInput(`
FUNC add(x, y) => x + y
RETURN add(2, 3)
`),
			Expected:    5,
			Description: "UDF round-trip",
			Check: func(t *testing.T, original *bytecode.Program, decoded *bytecode.Program) {
				t.Helper()

				if !reflect.DeepEqual(original.Functions.UserDefined, decoded.Functions.UserDefined) {
					t.Fatalf("udf metadata mismatch: got %#v, want %#v", decoded.Functions.UserDefined, original.Functions.UserDefined)
				}
			},
		},
		{
			Input:       spec.NewExpressionInput(`RETURN 1 + 2`),
			Expected:    3,
			Description: "Source and debug metadata round-trip",
			Check: func(t *testing.T, original *bytecode.Program, decoded *bytecode.Program) {
				t.Helper()

				if original.Source == nil || decoded.Source == nil {
					t.Fatalf("expected source to be preserved")
				}

				if got, want := decoded.Source.Name(), original.Source.Name(); got != want {
					t.Fatalf("source name mismatch: got %q, want %q", got, want)
				}

				if got, want := decoded.Source.Content(), original.Source.Content(); got != want {
					t.Fatalf("source content mismatch: got %q, want %q", got, want)
				}

				if !reflect.DeepEqual(original.Metadata.DebugSpans, decoded.Metadata.DebugSpans) {
					t.Fatalf("debug spans mismatch: got %#v, want %#v", decoded.Metadata.DebugSpans, original.Metadata.DebugSpans)
				}

				if !reflect.DeepEqual(original.Metadata.Labels, decoded.Metadata.Labels) {
					t.Fatalf("labels mismatch: got %#v, want %#v", decoded.Metadata.Labels, original.Metadata.Labels)
				}
			},
		},
	}

	std := spec.Stdlib()

	for _, level := range levels {
		compilerInstance := compiler.New(compiler.WithOptimizationLevel(level))

		for _, tc := range cases {
			caseName := tc.Description
			if caseName == "" {
				caseName = tc.Input.String()
			}

			t.Run("artifact/O"+string(rune('0'+level))+"/"+caseName, func(t *testing.T) {
				program, err := tc.Input.ResolveProgram(caseName, compilerInstance)
				if err != nil {
					t.Fatalf("compile failed: %v", err)
				}

				envOpts := []vm.EnvironmentOption{vm.WithNamespace(std)}
				envOpts = append(envOpts, tc.Env...)

				originalOut, err := spec.Run(program, envOpts...)
				if err != nil {
					t.Fatalf("original run failed: %v", err)
				}

				data, err := artifact.Marshal(program, artifact.Options{})
				if err != nil {
					t.Fatalf("artifact marshal failed: %v", err)
				}

				decoded, err := artifact.Unmarshal(data)
				if err != nil {
					t.Fatalf("artifact unmarshal failed: %v", err)
				}

				decodedOut, err := spec.Run(decoded, envOpts...)
				if err != nil {
					t.Fatalf("artifact round-trip run failed: %v", err)
				}

				specassert.Expect(t, specassert.EqualJSON, jsonRaw(originalOut), tc.Expected, "original program output mismatch")
				specassert.Expect(t, specassert.EqualJSON, jsonRaw(decodedOut), tc.Expected, "artifact round-trip output mismatch")
				specassert.Expect(t, specassert.EqualJSON, jsonRaw(originalOut), jsonRaw(decodedOut), "original and artifact outputs differ")

				if tc.Check != nil {
					tc.Check(t, program, decoded)
				}
			})
		}
	}
}

func jsonRaw(value []byte) json.RawMessage {
	return json.RawMessage(value)
}
