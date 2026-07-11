package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
	"github.com/MontFerret/ferret/v2/test/spec/compile/inspect"
)

func fourSlotQueryDescriptorFor(code []bytecode.Instruction, opcode bytecode.Opcode) error {
	applyIdx, ok := inspect.FindFirstOpcodeIndex(code, opcode)
	if !ok {
		return fmt.Errorf("expected %s in bytecode", opcode)
	}

	size, ok := inspect.FindApplyQueryDescriptorSize(code, applyIdx)
	if !ok {
		return fmt.Errorf("expected OpLoadArray for query descriptor before OpQuery")
	}

	if size != 4 {
		return fmt.Errorf("expected 4-slot query descriptor, got %d", size)
	}

	return nil
}

func fourSlotQueryDescriptor(code []bytecode.Instruction) error {
	return fourSlotQueryDescriptorFor(code, bytecode.OpQuery)
}

func TestQueryExpressionSourceImplicitCurrentCompiles(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET sections = @sections
LET linksBySection = sections[* RETURN (QUERY "a" IN . USING css)]
RETURN linksBySection[**]`, func(prog *bytecode.Program) error {
			return fourSlotQueryDescriptor(prog.Bytecode)
		}, "Should compile query expression with implicit current source"),
	})
}

func TestQueryExpressionMemberPayloadCompiles(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET email = { body: ".dynamic-member" }
LET model = @doc
RETURN QUERY ONE email.body IN model USING summarize`, func(prog *bytecode.Program) error {
			return fourSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQueryOne)
		}, "Should compile query expression with member payload"),
	})
}

func TestQueryExpressionAtomicPayloadsCompile(t *testing.T) {
	cases := []string{
		`RETURN QUERY "div" IN @doc`,
		`LET selector = ".item" RETURN QUERY selector IN @doc`,
		`RETURN QUERY @selector IN @doc`,
		`LET config = { selector: ".item" } RETURN QUERY config.selector IN @doc`,
		`LET selectors = [".item"] LET index = 0 RETURN QUERY selectors[index] IN @doc`,
		`FUNC GET_SELECTOR() => ".item" RETURN QUERY GET_SELECTOR() IN @doc`,
		`FUNC factory() => { selector: ".item" } RETURN QUERY factory().selector IN @doc`,
	}

	specs := make([]spec.Spec, 0, len(cases))
	for _, query := range cases {
		specs = append(specs, ProgramCheck(query, func(prog *bytecode.Program) error {
			return fourSlotQueryDescriptor(prog.Bytecode)
		}, query))
	}

	RunSpecs(t, specs)
}

func TestQueryExpressionComputedPayloadsCompile(t *testing.T) {
	cases := []string{
		`LET prefix = ".item-" LET selector = "card" RETURN QUERY (prefix + selector) IN @doc`,
		`LET enabled = TRUE LET primary = ".primary" LET fallback = ".fallback" RETURN QUERY (enabled ? primary : fallback) IN @doc`,
		`LET selector = ".item" LET selectors = [".item"] RETURN QUERY (selector IN selectors) IN @doc`,
		`FUNC BUILD_SELECTOR(options) => ".built" RETURN QUERY (BUILD_SELECTOR({})) IN @doc`,
	}

	specs := make([]spec.Spec, 0, len(cases))
	for _, query := range cases {
		specs = append(specs, ProgramCheck(query, func(prog *bytecode.Program) error {
			return fourSlotQueryDescriptor(prog.Bytecode)
		}, query))
	}

	specs = append(specs, ProgramCheck(`
LET prefix = ".item-"
LET selector = "card"
RETURN QUERY ONE (prefix + selector) IN @doc USING css`, func(prog *bytecode.Program) error {
		return fourSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQueryOne)
	}, "Should compile QUERY ONE with computed payload using direct opcode"))

	RunSpecs(t, specs)
}

func TestQueryExpressionOptionalUsingCompiles(t *testing.T) {
	cases := []string{
		`RETURN QUERY "x" IN @doc`,
		`RETURN QUERY "x" IN @doc WITH {}`,
		`RETURN QUERY "x" IN @doc OPTIONS {}`,
		`RETURN QUERY "x" IN @doc WITH {} OPTIONS {}`,
		`RETURN QUERY "x" IN @doc USING css`,
		`RETURN QUERY "x" IN @doc USING css WITH {}`,
		`RETURN QUERY "x" IN @doc USING css OPTIONS {}`,
		`RETURN QUERY "x" IN @doc USING css WITH {} OPTIONS {}`,
	}

	specs := make([]spec.Spec, 0, len(cases))
	for _, query := range cases {
		specs = append(specs, ProgramCheck(query, func(prog *bytecode.Program) error {
			return fourSlotQueryDescriptor(prog.Bytecode)
		}, query))
	}

	RunSpecs(t, specs)
}

func TestQueryModifierLowering_OneUsesDirectOpcode(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`RETURN QUERY ONE ".items" IN @doc USING css`, func(prog *bytecode.Program) error {
			if err := fourSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQueryOne); err != nil {
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

func TestQueryShorthandLowering(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck("RETURN @doc[~ \"x\"]", func(prog *bytecode.Program) error {
			if err := fourSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQuery); err != nil {
				return err
			}
			if inspect.HasOpcode(prog, bytecode.OpQueryOne) {
				return fmt.Errorf("did not expect OpQueryOne for regular raw-string query shorthand")
			}

			return nil
		}, "regular raw-string query shorthand lowering"),
		ProgramCheck("RETURN @doc[~ css`.items`]", func(prog *bytecode.Program) error {
			if err := fourSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQuery); err != nil {
				return err
			}
			if inspect.HasOpcode(prog, bytecode.OpQueryOne) {
				return fmt.Errorf("did not expect OpQueryOne for regular query shorthand")
			}

			return nil
		}, "regular query shorthand lowering"),
		ProgramCheck("RETURN @doc[~? \"x\"]", func(prog *bytecode.Program) error {
			if err := fourSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQueryOne); err != nil {
				return err
			}
			if inspect.HasOpcode(prog, bytecode.OpQuery) {
				return fmt.Errorf("did not expect OpQuery for raw-string query-one shorthand")
			}
			if inspect.HasOpcode(prog, bytecode.OpLength) {
				return fmt.Errorf("did not expect OpLength for raw-string query-one shorthand")
			}
			if inspect.HasOpcode(prog, bytecode.OpLoadIndexConst) {
				return fmt.Errorf("did not expect OpLoadIndexConst for raw-string query-one shorthand")
			}
			if inspect.HasOpcode(prog, bytecode.OpFail) {
				return fmt.Errorf("did not expect OpFail for raw-string query-one shorthand")
			}

			return nil
		}, "raw-string query-one shorthand lowering"),
		ProgramCheck("RETURN @doc[~? css`.items`]", func(prog *bytecode.Program) error {
			if err := fourSlotQueryDescriptorFor(prog.Bytecode, bytecode.OpQueryOne); err != nil {
				return err
			}
			if inspect.HasOpcode(prog, bytecode.OpQuery) {
				return fmt.Errorf("did not expect OpQuery for query-one shorthand")
			}
			if inspect.HasOpcode(prog, bytecode.OpLength) {
				return fmt.Errorf("did not expect OpLength for query-one shorthand")
			}
			if inspect.HasOpcode(prog, bytecode.OpLoadIndexConst) {
				return fmt.Errorf("did not expect OpLoadIndexConst for query-one shorthand")
			}
			if inspect.HasOpcode(prog, bytecode.OpFail) {
				return fmt.Errorf("did not expect OpFail for query-one shorthand")
			}

			return nil
		}, "query-one shorthand lowering"),
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
			if err := fourSlotQueryDescriptorFor(prog.Bytecode, tc.opcode); err != nil {
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
