package exec

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestRunRoundTripsSupportsSourceInput(t *testing.T) {
	RunRoundTrips(t, []RoundTrip{
		{
			Input:    spec.NewExpressionInput("RETURN 1 + 2"),
			Expected: 3,
		},
	})
}

func TestRunRoundTripsAppliesEnvironmentOptionsToBothExecutions(t *testing.T) {
	callCount := 0

	RunRoundTrips(t, []RoundTrip{
		{
			Input:       spec.NewExpressionInput("RETURN PING()"),
			Expected:    7,
			Description: "env options apply to original and decoded programs",
			Env: []vm.EnvironmentOption{
				vm.WithFunction("PING", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					callCount++
					return runtime.NewInt(7), nil
				}),
			},
		},
	})

	if got, want := callCount, 4; got != want {
		t.Fatalf("unexpected host function call count: got %d, want %d", got, want)
	}
}

func TestRunRoundTripsSupportsProgramInput(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Constants:  []runtime.Value{runtime.NewString("ok")},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Registers: 1,
	}

	RunRoundTrips(t, []RoundTrip{
		{
			Input:       spec.NewProgramInput(prog, "prebuilt program"),
			Expected:    "ok",
			Description: "prebuilt programs round trip",
		},
	})
}
