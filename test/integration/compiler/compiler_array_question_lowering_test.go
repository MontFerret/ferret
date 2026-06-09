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
		ProgramCheck(`RETURN @arr[?]`, expectBareArrayQuestionLowering, "bare array question uses expansion length"),
		ProgramCheck(`RETURN @arr[? FILTER . > 1]`, expectFilteredArrayQuestionLowering, "filtered array question keeps counting loop"),
		ProgramCheck(`RETURN @arr[? ANY FILTER . > 1]`, expectFilteredArrayQuestionLowering, "quantified array question keeps counting loop"),
	}, compiler.O0, compiler.O1)
}

func expectBareArrayQuestionLowering(prog *bytecode.Program) error {
	// The optimization must eliminate the counting loop (no OpIncr)
	// and use a length/emptiness check instead.
	if got := inspect.CountOpcode(prog, bytecode.OpIncr); got != 0 {
		return fmt.Errorf("expected no OpIncr (counting loop eliminated), got %d", got)
	}

	if !inspect.HasOpcode(prog, bytecode.OpLength) {
		return fmt.Errorf("expected OpLength for emptiness check")
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
