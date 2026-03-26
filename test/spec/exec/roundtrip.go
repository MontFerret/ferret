package exec

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	specassert "github.com/MontFerret/ferret/v2/test/spec/assert"
)

type RoundTrip struct {
	Input       spec.Input
	Expected    any
	Description string
	Env         []vm.EnvironmentOption
}

func (c RoundTrip) name() string {
	if c.Description != "" {
		return c.Description
	}

	return c.Input.String()
}

func RunRoundTripsWith(t *testing.T, name string, c *compiler.Compiler, cases []RoundTrip, opts ...vm.EnvironmentOption) {
	t.Helper()

	std := spec.Stdlib()

	for _, tc := range cases {
		caseName := fmt.Sprintf("%s/%s", name, tc.name())

		t.Run(caseName, func(t *testing.T) {
			prog, err := tc.Input.ResolveProgram(caseName, c)
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			envOpts := []vm.EnvironmentOption{
				vm.WithNamespace(std),
			}
			envOpts = append(envOpts, opts...)
			envOpts = append(envOpts, tc.Env...)

			originalOut, err := spec.Run(prog, envOpts...)
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

			roundTripOut, err := spec.Run(&decoded, envOpts...)
			if err != nil {
				t.Fatalf("round-trip run failed: %v", err)
			}

			specassert.Expect(t, specassert.EqualJSON, json.RawMessage(originalOut), tc.Expected, "original program output mismatch")
			specassert.Expect(t, specassert.EqualJSON, json.RawMessage(roundTripOut), tc.Expected, "round-trip program output mismatch")
			specassert.Expect(t, specassert.EqualJSON, json.RawMessage(originalOut), json.RawMessage(roundTripOut), "original and round-trip outputs differ")
		})
	}
}

func RunRoundTrips(t *testing.T, cases []RoundTrip, opts ...vm.EnvironmentOption) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		RunRoundTripsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			cases,
			opts...,
		)
	}
}
