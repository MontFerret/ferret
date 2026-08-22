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
			name: "parenthesized script expression statement",
			query: `(FOR value IN [1, 2] {
  RETURN value
})`,
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
			name: "parenthesized loop body expression statement",
			query: `FOR outer IN [1] {
  (FOR value IN [1, 2] {
    RETURN value
  })
  RETURN outer
}`,
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
			name: "discarded explicit returned loop",
			query: `FOR outer IN [1, 2] {
  RETURN FOR inner IN [outer, outer + 1] {
    RETURN inner
  }
}`,
		},
		{
			name: "consumed outer with nonterminal collecting loop",
			query: `RETURN FOR outer IN [1, 2] {
  FOR inner IN [outer] {
    RETURN inner
  }
  LET after = outer
  RETURN after
}`,
			wantDataSets: 1,
			wantPushes:   1,
		},
		{
			name: "discarded parenthesized explicit returned loop",
			query: `FOR outer IN [1, 2] {
  RETURN (FOR inner IN [outer, outer + 1] {
    RETURN inner
  })
}`,
		},
		{
			name: "consumed explicit returned loop",
			query: `RETURN FOR outer IN [1, 2] {
  RETURN FOR inner IN [outer, outer + 1] {
    RETURN inner
  }
}`,
			wantDataSets: 2,
			wantPushes:   2,
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

func TestReturnlessForLowering(t *testing.T) {
	queries := map[string]string{
		"empty": `FOR value IN [] {}`,
		"for in": `FOR value IN [1, 2] {
  LET copy = value
}`,
		"while": `VAR count = 0
FOR WHILE count < 2 {
  count += 1
}
RETURN count`,
		"do while": `VAR count = 0
FOR DO WHILE false {
  count += 1
}
RETURN count`,
		"nested before statement": `FOR outer IN [1, 2] {
  FOR inner IN [outer] {
    LET copy = inner
  }
  LET after = outer
}`,
		"discarded grouped UDF statement": `FUNC effect() {
  (FOR value IN [1, 2] {
    LET copy = value
  })
}
RETURN effect()`,
		"discarded terminal nested loop": `FOR outer IN [1, 2] {
  FOR inner IN [outer] {
    RETURN inner
  }
}`,
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for name, query := range queries {
			t.Run(fmt.Sprintf("O%d/%s", level, name), func(t *testing.T) {
				program := compileWithLevel(t, level, query)
				if err := bytecode.ValidateProgram(program); err != nil {
					t.Fatalf("validate bytecode: %v", err)
				}

				for _, opcode := range []bytecode.Opcode{bytecode.OpDataSet, bytecode.OpPush} {
					if inspect.HasOpcode(program, opcode) {
						t.Fatalf("returnless FOR unexpectedly contains %s", opcode)
					}
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
			program, err := mustNewCompiler(
				t,
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

	program, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(query))
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

func TestReturnlessForRetainsNestedStatementDebugMetadata(t *testing.T) {
	const query = `FOR outer IN [1] {
  FOR inner IN [outer] {}
  RECORD(outer)
}
RETURN NONE`

	program, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("compile query: %v", err)
	}

	wantLines := map[int]bool{1: false, 2: false, 3: false}
	for _, point := range program.Metadata.DebugPoints {
		line, _ := program.Source.LocationAt(point.Span)
		if point.Kind == bytecode.DebugPointReturn && line < 5 {
			t.Fatalf("returnless FOR unexpectedly has a return debug point on line %d", line)
		}

		if point.Kind == bytecode.DebugPointStatement {
			if _, ok := wantLines[line]; ok {
				wantLines[line] = true
			}
		}
	}

	for line, found := range wantLines {
		if !found {
			t.Fatalf("returnless FOR has no statement debug point on line %d: %#v", line, program.Metadata.DebugPoints)
		}
	}
}

func TestGeneralExpressionStatementsRetainDebugMetadata(t *testing.T) {
	const query = `1 + 2
FUNC effect(value) {
  VAR copy = 0
  copy += 1
  RECORD(value.member) + 1
  RETURN NONE
}
FOR value IN [{ member: 1 }] {
  value.member + 1
}
RETURN NONE`

	program, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("compile query: %v", err)
	}

	wantLines := map[int]bool{1: false, 5: false, 9: false}
	for _, point := range program.Metadata.DebugPoints {
		if point.Kind != bytecode.DebugPointStatement {
			continue
		}

		line, _ := program.Source.LocationAt(point.Span)
		if _, ok := wantLines[line]; ok {
			wantLines[line] = true
		}
	}

	for line, found := range wantLines {
		if !found {
			t.Fatalf("expression statement has no statement debug point on line %d: %#v", line, program.Metadata.DebugPoints)
		}
	}
}

func TestExpressionStatementCallsPreserveReachabilityAndCaptures(t *testing.T) {
	const query = `LET base = 10
FUNC outer() {
  (FOR item IN [base] {
    inner(item) + 1
  })
  FUNC inner(value) => base + value
}
outer() + 1
RETURN NONE`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program := compileWithLevel(t, level, query)

			outer, err := findUserDefined(program, "outer")
			if err != nil {
				t.Fatal(err)
			}

			inner, err := findUserDefined(program, "inner")
			if err != nil {
				t.Fatal(err)
			}

			if got, want := outer.Params, 1; got != want {
				t.Fatalf("outer params/captures = %d, want %d", got, want)
			}

			if got, want := inner.Params, 2; got != want {
				t.Fatalf("inner params/captures = %d, want %d", got, want)
			}
		})
	}
}

func TestParenthesizedLoopExpressionStatementRetainsDebugMetadata(t *testing.T) {
	const query = `(FOR value IN [1] {
  value + 1
  RETURN value
})
RETURN NONE`

	program, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("compile query: %v", err)
	}

	want := map[bytecode.DebugPointKind]map[int]bool{
		bytecode.DebugPointStatement: {1: false, 2: false},
		bytecode.DebugPointReturn:    {3: false, 5: false},
	}

	for _, point := range program.Metadata.DebugPoints {
		lines, ok := want[point.Kind]
		if !ok {
			continue
		}

		line, _ := program.Source.LocationAt(point.Span)
		if _, ok := lines[line]; ok {
			lines[line] = true
		}
	}

	for kind, lines := range want {
		for line, found := range lines {
			if !found {
				t.Fatalf("parenthesized loop has no kind %d debug point on line %d: %#v", kind, line, program.Metadata.DebugPoints)
			}
		}
	}
}

func TestDiscardedExplicitReturnedForRetainsReturnDebugMetadata(t *testing.T) {
	const query = `FOR outer IN [1] {
  RETURN FOR inner IN [outer] {
    RETURN inner
  }
}
RETURN NONE`

	program, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("compile query: %v", err)
	}

	if err := bytecode.ValidateProgram(program); err != nil {
		t.Fatalf("validate bytecode: %v", err)
	}

	returnLines := make(map[int]bool)
	for _, point := range program.Metadata.DebugPoints {
		if point.Kind != bytecode.DebugPointReturn || point.FunctionID != -1 {
			continue
		}

		line, _ := program.Source.LocationAt(point.Span)
		returnLines[line] = true
	}

	for _, line := range []int{2, 3} {
		if !returnLines[line] {
			t.Fatalf("discarded explicit returned loop has no return debug point on line %d: %#v", line, program.Metadata.DebugPoints)
		}
	}
}

func TestEmptySourceRemainsInvalid(t *testing.T) {
	for _, input := range []string{"", " \n\t"} {
		if _, err := mustNewCompiler(t).Compile(source.NewAnonymous(input)); err == nil {
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
