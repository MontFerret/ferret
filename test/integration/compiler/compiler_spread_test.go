package compiler_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestLiteralSpreadLowering(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`RETURN [1, ...@first, 2, ...@second]`, expectSpreadOrder(
			bytecode.OpArrayPush,
			bytecode.OpArraySpread,
			bytecode.OpArrayPush,
			bytecode.OpArraySpread,
		), "array entries retain source order"),
		ProgramCheck(`RETURN { a: 1, ...@first, b: 2, ...@second }`, expectSpreadOrder(
			bytecode.OpObjectSetConst,
			bytecode.OpObjectSpread,
			bytecode.OpObjectSetConst,
			bytecode.OpObjectSpread,
		), "object entries retain source order"),
		Opcode(`RETURN [...(true ? [1] : [2]),]`, OpcodeCount{
			Count: map[bytecode.Opcode]int{bytecode.OpArraySpread: 1},
		}, "array spread accepts a full expression and trailing comma"),
		Opcode(`RETURN {...(true ? { a: 1 } : { b: 2 }),}`, OpcodeCount{
			Count: map[bytecode.Opcode]int{bytecode.OpObjectSpread: 1},
		}, "object spread accepts a full expression and trailing comma"),
		Opcode(`RETURN WAITFOR EXISTS [...@source] TIMEOUT 1ms EVERY 0ms`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpArraySpread},
		}, "spread-only arrays remain unknown to wait constant analysis"),
		Opcode(`RETURN WAITFOR EXISTS {...@source} TIMEOUT 1ms EVERY 0ms`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpObjectSpread},
		}, "spread-only objects remain unknown to wait constant analysis"),
		ProgramCheck(`RETURN FOR item IN [...@items] RETURN item.value`, func(program *bytecode.Program) error {
			var hasGenericLoad bool

			for _, instruction := range program.Bytecode {
				if instruction.Opcode == bytecode.OpLoadPropertyConst {
					hasGenericLoad = true
				}

				if instruction.Opcode == bytecode.OpLoadKeyConst {
					return fmt.Errorf("spread-only element inference must not assume Object")
				}
			}

			if !hasGenericLoad {
				return fmt.Errorf("expected generic property load for spread-only element")
			}

			return bytecode.ValidateProgram(program)
		}, "spread-only array element inference remains conservative"),
	}, compiler.O0, compiler.O1)
}

func TestLiteralSpreadSyntaxBoundaries(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`RETURN [..., 1]`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after '...' in array spread",
			Hint:    "Provide an Array expression or none after '...'.",
		}, "array spread requires an operand"),
		Failure(`RETURN {..., a: 1}`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after '...' in object spread",
			Hint:    "Provide an Object expression or none after '...'.",
		}, "object spread requires an operand"),
		Failure(`RETURN T::FN(...[1])`, E{Kind: parserd.SyntaxError}, "call argument spread remains unsupported"),
		Failure(`LET [...rest] = [1] RETURN rest`, E{Kind: parserd.SyntaxError}, "rest destructuring remains unsupported"),
		Failure(`RETURN [1 ...[2]]`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected ',' between array items",
			Hint:    "Separate array items with commas, e.g. [1, 2, 3].",
		}, "array spread entry reports a missing separator"),
		Failure(`RETURN { a: 1 ...{ b: 2 } }`, E{Kind: parserd.SyntaxError}, "object spread entry requires a separator"),
	})
}

func expectSpreadOrder(want ...bytecode.Opcode) func(*bytecode.Program) error {
	return func(program *bytecode.Program) error {
		var got []bytecode.Opcode
		allocation := bytecode.OpLoadArray
		if len(want) > 0 && want[0] == bytecode.OpObjectSetConst {
			allocation = bytecode.OpLoadObject
		}

		allocations := 0

		for _, instruction := range program.Bytecode {
			if instruction.Opcode == allocation {
				allocations++
			}

			switch instruction.Opcode {
			case bytecode.OpArrayPush,
				bytecode.OpArraySpread,
				bytecode.OpObjectSet,
				bytecode.OpObjectSetConst,
				bytecode.OpObjectSpread:
				got = append(got, instruction.Opcode)
			}
		}

		if !slices.Equal(got, want) {
			return fmt.Errorf("unexpected literal operation order: got %v, want %v", got, want)
		}

		if allocations != 1 {
			return fmt.Errorf("expected one destination allocation, got %d", allocations)
		}

		return bytecode.ValidateProgram(program)
	}
}
