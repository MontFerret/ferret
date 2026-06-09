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

func TestArrayQuestionLowering(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`RETURN [1][?]`, expectBareArrayQuestionLengthLowering, "bare array question uses OpLength on measurable source"),
		ProgramCheck(`RETURN @arr[?]`, expectBareArrayQuestionLowering, "bare array question avoids counting loop"),
		ProgramCheck(`RETURN @arr[? FILTER . > 1]`, expectFilteredArrayQuestionLowering, "filtered array question keeps counting loop"),
		ProgramCheck(`RETURN @arr[? ANY FILTER . > 1]`, expectFilteredArrayQuestionLowering, "quantified array question keeps counting loop"),
	}, compiler.O0, compiler.O1)
}

func expectBareArrayQuestionLengthLowering(prog *bytecode.Program) error {
	if err := expectBareArrayQuestionLowering(prog); err != nil {
		return err
	}

	if !inspect.HasOpcode(prog, bytecode.OpLength) {
		return fmt.Errorf("expected OpLength for measurable bare array question")
	}

	if inspect.HasOpcode(prog, bytecode.OpIter) {
		return fmt.Errorf("did not expect OpIter when OpLength fast path is applicable")
	}

	return nil
}

func expectBareArrayQuestionLowering(prog *bytecode.Program) error {
	// The bare array question should NOT materialize a list (no OpDataSet, no OpPush).
	// It should take exactly one of two strategies: OpLength fast path or OpIter early-exit.
	if inspect.HasOpcode(prog, bytecode.OpDataSet) {
		return fmt.Errorf("did not expect OpDataSet for bare array question (should not materialize a list)")
	}

	if inspect.HasOpcode(prog, bytecode.OpPush) {
		return fmt.Errorf("did not expect OpPush for bare array question (should not materialize a list)")
	}

	if inspect.HasOpcode(prog, bytecode.OpIncr) {
		return fmt.Errorf("did not expect OpIncr for bare array question (should not use counting loop)")
	}

	// Must have exactly one of OpLength (measurable fast path) or OpIter (general iterator path).
	hasLength := inspect.HasOpcode(prog, bytecode.OpLength)
	hasIter := inspect.HasOpcode(prog, bytecode.OpIter)

	if !hasLength && !hasIter {
		return fmt.Errorf("expected OpLength or OpIter for bare array question, found neither")
	}

	if hasLength && hasIter {
		return fmt.Errorf("expected exactly one of OpLength or OpIter, found both")
	}

	// When using the iterator strategy, OpIterNext must be present.
	if hasIter && !inspect.HasOpcode(prog, bytecode.OpIterNext) {
		return fmt.Errorf("OpIter present but OpIterNext missing for iterator early-exit strategy")
	}

	return nil
}

func expectFilteredArrayQuestionLowering(prog *bytecode.Program) error {
	if inspect.HasOpcode(prog, bytecode.OpLength) {
		return fmt.Errorf("did not expect OpLength for filtered array question")
	}

	if got := inspect.CountOpcode(prog, bytecode.OpIncr); got != 2 {
		return fmt.Errorf("expected 2 INCR ops for counting loop, got %d", got)
	}

	return nil
}
