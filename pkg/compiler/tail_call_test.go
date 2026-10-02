package compiler

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestTailCallEliminationLevelsAndDebugPolicy(t *testing.T) {
	const query = `func target(value) => value + 1
func forward(value) => target(value)
return forward(2)`
	cases := []struct {
		name    string
		options []Option
		level   OptimizationLevel
	}{
		{"default", nil, Full},
		{"none", []Option{WithOptimizationLevel(None)}, None},
		{"basic", []Option{WithOptimizationLevel(Basic)}, Basic},
		{"full", []Option{WithOptimizationLevel(Full)}, Full},
		{"debug after full", []Option{WithOptimizationLevel(Full), WithDebugInfo()}, None},
		{"full after debug", []Option{WithDebugInfo(), WithOptimizationLevel(Full)}, None},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			program, err := mustNewCompiler(t, tc.options...).Compile(t.Context(), source.NewAnonymous(query))
			if err != nil {
				t.Fatal(err)
			}

			if program.Metadata.OptimizationLevel != int(tc.level) {
				t.Fatalf("effective level = %d, want %v", program.Metadata.OptimizationLevel, tc.level)
			}

			tailCalls, calls := 0, 0
			for _, instruction := range program.Bytecode {
				switch instruction.Opcode {
				case bytecode.OpTailCall:
					tailCalls++
				case bytecode.OpCall:
					calls++
				}
			}

			wantTail, wantCalls := 1, 1
			if tc.level == None {
				wantTail, wantCalls = 0, 2
			}

			if tailCalls != wantTail || calls != wantCalls {
				t.Fatalf("tail/ordinary calls = %d/%d, want %d/%d", tailCalls, calls, wantTail, wantCalls)
			}
		})
	}
}

func TestTailCallEliminationRecoveryAndCellEligibility(t *testing.T) {
	cases := []struct {
		name     string
		query    string
		wantTail bool
	}{
		{"explicit fail", `func target() => 1
func forward() => target() on error fail
return forward()`, true},
		{"returned local", `func target() => 1
func forward() { let result = target() return result }
return forward()`, true},
		{"normalized retry", `func target() => 1
func forward() => target() on error retry 0 or fail
return forward()`, true},
		{"suppressed", `func target() => 1
func forward() => target()?
return forward()`, false},
		{"fallback", `func target() => 1
func forward() => target() on error return 2
return forward()`, false},
		{"retry", `func target() => 1
func forward() => target() on error retry 2 or return 3
return forward()`, false},
		{"distinct", `func target() => [1,1]
func forward() { return distinct target() }
return forward()`, false},
		{"local mutable capture", `func outer() {
var value = 1
func inner() { value += 1 return value }
return inner()
}
return outer()`, false},
		{"borrowed mutable capture", `var value = 1
func inner() { value += 1 return value }
func forward() => inner()
return forward()`, true},
		{"match result", `func recur(n) => match n { 0 => 0, _ => recur(n - 1) }
return recur(3)`, false},
	}

	for _, tc := range cases {
		for _, level := range []OptimizationLevel{None, Basic, Full} {
			t.Run(tc.name+"/"+level.String(), func(t *testing.T) {
				program, err := mustNewCompiler(t, WithOptimizationLevel(level)).Compile(t.Context(), source.NewAnonymous(tc.query))
				if err != nil {
					t.Fatal(err)
				}

				tails := 0
				for _, instruction := range program.Bytecode {
					if instruction.Opcode == bytecode.OpTailCall {
						tails++
					}
				}

				want := 0
				if tc.wantTail && level != None {
					want = 1
				}

				if tails != want {
					t.Fatalf("tail calls = %d, want %d", tails, want)
				}
			})
		}
	}
}
