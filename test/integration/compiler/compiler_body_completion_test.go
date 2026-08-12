package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec/compile/inspect"
)

func TestBodyCompletionEmitsNoneReturns(t *testing.T) {
	const query = `
FUNC effect() {
  LET value = 1
}
effect()
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program := compileWithLevel(t, level, query)
			if got, want := len(program.Functions.UserDefined), 1; got != want {
				t.Fatalf("UDF count = %d, want %d", got, want)
			}

			entry := program.Functions.UserDefined[0].Entry
			if entry <= 0 {
				t.Fatalf("UDF entry = %d, want a main-body prefix", entry)
			}

			assertNoneReturn(t, program.Bytecode[entry-1], "main body")
			assertNoneReturn(t, program.Bytecode[len(program.Bytecode)-1], "UDF body")
		})
	}
}

func TestReturnedAndStandaloneForBytecodeResults(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			standalone := compileWithLevel(t, level, `FOR value IN [1] { RETURN value }`)
			assertNoneReturn(t, standalone.Bytecode[len(standalone.Bytecode)-1], "standalone FOR body")
			assertLoopResultOpcodes(t, standalone, false, "standalone FOR body")

			returned := compileWithLevel(t, level, `RETURN FOR value IN [1] { RETURN value }`)
			assertLoopResultOpcodes(t, returned, true, "returned FOR body")

			last := returned.Bytecode[len(returned.Bytecode)-1]
			if last.Opcode != bytecode.OpReturn {
				t.Fatalf("returned FOR terminator = %s, want %s", last.Opcode, bytecode.OpReturn)
			}

			if last.Operands[0] == bytecode.NoopOperand {
				t.Fatal("returned FOR unexpectedly returns the NONE register")
			}
		})
	}
}

func TestDiscardedForLowering(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		wantDataSets  int
		wantPushes    int
		wantSorter    bool
		wantCollector bool
	}{
		{
			name: "standalone script",
			query: `FOR value IN [1, 2] {
  RETURN value
}`,
		},
		{
			name: "standalone UDF",
			query: `FUNC effect() {
  FOR value IN [1, 2] {
    RETURN value
  }
}
RETURN effect()`,
		},
		{
			name: "parenthesized UDF expression statement",
			query: `FUNC effect() {
  (FOR value IN [1, 2] {
    RETURN value
  })
}
RETURN effect()`,
		},
		{
			name: "parenthesized UDF recovered expression statement",
			query: `FUNC effect() {
  (FOR value IN [1, 2] {
    RETURN value
  }) ON ERROR RETRY 1 OR RETURN NONE
}
RETURN effect()`,
		},
		{
			name: "discarded pass-through chain",
			query: `FOR outer IN [1, 2] {
  FOR inner IN [outer, outer + 1] {
    RETURN inner
  }
}`,
		},
		{
			name: "assignment requires result",
			query: `LET values = (FOR value IN [1, 2] {
  RETURN value
})
RETURN values`,
			wantDataSets: 1,
			wantPushes:   1,
		},
		{
			name: "argument requires result",
			query: `RETURN LENGTH((FOR value IN [1, 2] {
  RETURN value
}))`,
			wantDataSets: 1,
			wantPushes:   1,
		},
		{
			name: "returned pass-through chain",
			query: `RETURN FOR outer IN [1, 2] {
  FOR inner IN [outer, outer + 1] {
    RETURN inner
  }
}`,
			wantDataSets: 1,
			wantPushes:   1,
		},
		{
			name: "explicit returned loop is value boundary",
			query: `FOR outer IN [1, 2] {
  RETURN FOR inner IN [outer, outer + 1] {
    RETURN inner
  }
}`,
			wantDataSets: 1,
			wantPushes:   1,
		},
		{
			name: "discarded sorted loop",
			query: `FOR value IN [2, 1] {
  SORT value
  RETURN value
}`,
			wantSorter: true,
		},
		{
			name: "discarded collected loop",
			query: `FOR value IN [1, 2] {
  COLLECT key = value
  RETURN key
}`,
			wantCollector: true,
		},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("O%d/%s", level, test.name), func(t *testing.T) {
				program := compileWithLevel(t, level, test.query)
				if err := bytecode.ValidateProgram(program); err != nil {
					t.Fatalf("validate bytecode: %v", err)
				}

				if got := inspect.CountOpcode(program, bytecode.OpDataSet); got != test.wantDataSets {
					t.Fatalf("OpDataSet count = %d, want %d", got, test.wantDataSets)
				}

				if got := inspect.CountOpcode(program, bytecode.OpPush); got != test.wantPushes {
					t.Fatalf("OpPush count = %d, want %d", got, test.wantPushes)
				}

				if got := inspect.HasOpcode(program, bytecode.OpDataSetSorter); got != test.wantSorter {
					t.Fatalf("OpDataSetSorter present = %t, want %t", got, test.wantSorter)
				}

				if got := inspect.HasOpcode(program, bytecode.OpDataSetCollector); got != test.wantCollector {
					t.Fatalf("OpDataSetCollector present = %t, want %t", got, test.wantCollector)
				}
			})
		}
	}
}

func TestImplicitNoneReturnsAreSourceInvisible(t *testing.T) {
	const query = `
FUNC effect() {
  LET value = 1
}
effect()
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(
				compiler.WithOptimizationLevel(level),
				compiler.WithDebugInfo(),
			).Compile(source.NewAnonymous(query))
			if err != nil {
				t.Fatalf("compile query: %v", err)
			}

			entry := program.Functions.UserDefined[0].Entry
			for _, pc := range []int{entry - 1, len(program.Bytecode) - 1} {
				if got, want := program.Metadata.DebugSpans[pc], (source.Span{Start: -1, End: -1}); got != want {
					t.Fatalf("implicit return span at pc %d = %+v, want %+v", pc, got, want)
				}
			}
		})
	}
}

func TestStandaloneForRetainsStatementDebugMetadata(t *testing.T) {
	const query = `FOR value IN [1] {
  RETURN value
}
RETURN NONE`

	program, err := compiler.New(compiler.WithDebugInfo()).Compile(source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("compile query: %v", err)
	}

	for _, point := range program.Metadata.DebugPoints {
		line, _ := program.Source.LocationAt(point.Span)
		if line == 1 && point.FunctionID == -1 && point.Kind == bytecode.DebugPointStatement {
			return
		}
	}

	t.Fatal("standalone FOR has no top-level statement debug point")
}

func TestEmptySourceRemainsInvalid(t *testing.T) {
	for _, input := range []string{"", " \n\t"} {
		if _, err := compiler.New().Compile(source.NewAnonymous(input)); err == nil {
			t.Fatalf("expected empty source %q to remain invalid", input)
		}
	}
}

func TestNonBodySourceCanFallThrough(t *testing.T) {
	for _, input := range []string{"// comment only", "USE FOO AS F"} {
		for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
			program := compileWithLevel(t, level, input)
			assertNoneReturn(t, program.Bytecode[len(program.Bytecode)-1], "non-body source")
		}
	}
}

func assertNoneReturn(t *testing.T, instruction bytecode.Instruction, body string) {
	t.Helper()

	if instruction.Opcode != bytecode.OpReturn {
		t.Fatalf("%s terminator = %s, want %s", body, instruction.Opcode, bytecode.OpReturn)
	}

	if instruction.Operands[0] != bytecode.NoopOperand {
		t.Fatalf("%s return operand = %d, want NONE register %d", body, instruction.Operands[0], bytecode.NoopOperand)
	}
}

func assertLoopResultOpcodes(t *testing.T, program *bytecode.Program, want bool, body string) {
	t.Helper()

	if got := inspect.HasOpcode(program, bytecode.OpDataSet); got != want {
		t.Fatalf("%s OpDataSet present = %t, want %t", body, got, want)
	}

	if got := inspect.HasOpcode(program, bytecode.OpPush); got != want {
		t.Fatalf("%s OpPush present = %t, want %t", body, got, want)
	}
}
