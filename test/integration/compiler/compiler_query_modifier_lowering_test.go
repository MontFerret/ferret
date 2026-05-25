package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
	"github.com/MontFerret/ferret/v2/test/spec/compile/inspect"
)

func threeSlotQueryDescriptorFor(code []bytecode.Instruction, opcode bytecode.Opcode) error {
	applyIdx, ok := inspect.FindFirstOpcodeIndex(code, opcode)
	if !ok {
		return fmt.Errorf("expected %s in bytecode", opcode)
	}

	size, ok := inspect.FindApplyQueryDescriptorSize(code, applyIdx)
	if !ok {
		return fmt.Errorf("expected OpLoadArray for query descriptor before OpQuery")
	}

	if size != 3 {
		return fmt.Errorf("expected 3-slot query descriptor, got %d", size)
	}

	return nil
}

func threeSlotQueryDescriptor(code []bytecode.Instruction) error {
	return threeSlotQueryDescriptorFor(code, bytecode.OpQuery)
}

func TestQueryExpressionSourceImplicitCurrentCompiles(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET sections = @sections
LET linksBySection = sections[* RETURN (QUERY "a" IN . USING css)]
RETURN linksBySection[**]`, func(prog *bytecode.Program) error {
			return threeSlotQueryDescriptor(prog.Bytecode)
		}, "Should compile query expression with implicit current source"),
	})
}

func TestQueryModifierLowering_OneUsesDirectOpcode(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`RETURN QUERY ONE ".items" IN @doc USING css`, func(prog *bytecode.Program) error {
			if err := threeSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQueryOne); err != nil {
				return err
			}
			if inspect.HasOpcode(prog, bytecode.OpLength) {
				return fmt.Errorf("did not expect OpLength for direct QUERY ONE lowering")
			}
			if inspect.HasOpcode(prog, bytecode.OpJumpIfEqConst) {
				return fmt.Errorf("did not expect OpJumpIfEqConst for direct QUERY ONE lowering")
			}
			if inspect.HasOpcode(prog, bytecode.OpLoadIndexConst) {
				return fmt.Errorf("did not expect OpLoadIndexConst for direct QUERY ONE lowering")
			}
			if inspect.HasOpcode(prog, bytecode.OpFail) {
				return fmt.Errorf("did not expect OpFail for direct QUERY ONE lowering")
			}

			return nil
		}, "query one lowering"),
	})
}

func TestQueryModifierLowering_ExistsCount(t *testing.T) {
	cases := []struct {
		name   string
		expr   string
		opcode bytecode.Opcode
		absent bytecode.Opcode
	}{
		{name: "exists", expr: `RETURN QUERY EXISTS ".items" IN @doc USING css`, opcode: bytecode.OpQueryExists, absent: bytecode.OpExists},
		{name: "count", expr: `RETURN QUERY COUNT ".items" IN @doc USING css`, opcode: bytecode.OpQueryCount, absent: bytecode.OpLength},
	}

	specs := make([]spec.Spec, 0, len(cases))
	for _, tc := range cases {
		specs = append(specs, ProgramCheck(tc.expr, func(prog *bytecode.Program) error {
			if err := threeSlotQueryDescriptorFor(prog.Bytecode, tc.opcode); err != nil {
				return err
			}

			if !inspect.HasOpcode(prog, tc.opcode) {
				return fmt.Errorf("expected opcode %s for QUERY %s lowering", tc.opcode, tc.name)
			}

			if tc.absent != 0 && inspect.HasOpcode(prog, tc.absent) {
				return fmt.Errorf("did not expect opcode %s for direct QUERY %s lowering", tc.absent, tc.name)
			}

			if inspect.HasOpcode(prog, bytecode.OpFail) {
				return fmt.Errorf("did not expect OpFail in QUERY %s lowering", tc.name)
			}

			return nil
		}, tc.name))
	}

	RunSpecs(t, specs)
}
