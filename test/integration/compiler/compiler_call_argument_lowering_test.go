package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
	"github.com/MontFerret/ferret/v2/test/spec/compile/inspect"
)

func TestUdfCallConstantArgsDirectLoadO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
FUNC f2(x, y) => x + y
RETURN f2(1, 2)
`, func(prog *bytecode.Program) error {
			callIndex, ok := inspect.FindFirstOpcodeIndex(prog.Bytecode, bytecode.OpCall)
			if !ok {
				return fmt.Errorf("expected OpCall in bytecode")
			}

			if err := CallArgsLoadedFromConsts(prog, CallArgsLoadedExpectation{
				Index:     callIndex,
				ArgsCount: 2,
			}); err != nil {
				return err
			}

			if inspect.HasOpcode(prog, bytecode.OpMove) || inspect.HasOpcode(prog, bytecode.OpMoveTracked) {
				return fmt.Errorf("expected no MOVE/MOVET instructions for constant-only UDF call setup")
			}

			return nil
		}, "udf constant args direct load"),
	})
}

func TestHostCallConstantArgsDirectLoadO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`RETURN TEST(1, 2)`, func(prog *bytecode.Program) error {
			callIndex, ok := inspect.FindFirstOpcodeIndex(prog.Bytecode, bytecode.OpHCall)
			if !ok {
				return fmt.Errorf("expected OpHCall in bytecode")
			}

			if err := CallArgsLoadedFromConsts(prog, CallArgsLoadedExpectation{
				Index:     callIndex,
				ArgsCount: 2,
			}); err != nil {
				return err
			}

			if inspect.HasOpcode(prog, bytecode.OpMove) || inspect.HasOpcode(prog, bytecode.OpMoveTracked) {
				return fmt.Errorf("expected no MOVE/MOVET instructions for constant-only host call setup")
			}

			return nil
		}, "host constant args direct load"),
	})
}

func TestCallArgumentLoweringKeepsMoveForNonLiteralArgO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET x = 1
RETURN TEST(x, 2)
`, func(prog *bytecode.Program) error {
			if !inspect.HasOpcode(prog, bytecode.OpMoveTracked) {
				return fmt.Errorf("expected MOVET instruction for non-literal argument setup")
			}

			return nil
		}, "non-literal arg keeps tracked move"),
	})
}

func TestCallArgumentSpansRecordedForCallInstructions(t *testing.T) {
	const query = "RETURN TEST(1 + 2, [3, 4])"

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			prog := compileWithLevel(t, level, query)

			callIndex, ok := inspect.FindFirstOpcodeIndex(prog.Bytecode, bytecode.OpHCall)
			if !ok {
				t.Fatal("expected OpHCall in bytecode")
			}

			if got, want := len(prog.Metadata.CallArgumentSpans), len(prog.Bytecode); got != want {
				t.Fatalf("unexpected call argument span metadata length: got %d, want %d", got, want)
			}

			for i, spans := range prog.Metadata.CallArgumentSpans {
				if i == callIndex {
					continue
				}

				if len(spans) != 0 {
					t.Fatalf("expected non-call instruction %d to have no call argument spans, got %#v", i, spans)
				}
			}

			spans := prog.Metadata.CallArgumentSpans[callIndex]
			if got, want := len(spans), 2; got != want {
				t.Fatalf("unexpected call argument span count: got %d, want %d", got, want)
			}

			wantFragments := []string{"1 + 2", "[3, 4]"}
			for i, span := range spans {
				if got := query[span.Start:span.End]; got != wantFragments[i] {
					t.Fatalf("unexpected call argument %d span: got %q, want %q", i, got, wantFragments[i])
				}
			}
		})
	}
}
