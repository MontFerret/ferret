package optimization

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func isSimpleConst(val runtime.Value) bool {
	if val == nil {
		return false
	}
	if val == runtime.None {
		return true
	}
	switch val.(type) {
	case runtime.Int, runtime.Float, runtime.String, runtime.Boolean:
		return true
	default:
		return false
	}
}

func constEqual(a, b runtime.Value) bool {
	if a == b {
		return true
	}
	switch av := a.(type) {
	case runtime.Int:
		if bv, ok := b.(runtime.Int); ok {
			return av == bv
		}
	case runtime.Float:
		if bv, ok := b.(runtime.Float); ok {
			return av == bv
		}
	case runtime.String:
		if bv, ok := b.(runtime.String); ok {
			return av == bv
		}
	case runtime.Boolean:
		if bv, ok := b.(runtime.Boolean); ok {
			return av == bv
		}
	}

	return false
}

func foldUnary(op bytecode.Opcode, val runtime.Value, bg context.Context) (runtime.Value, bool) {
	switch op {
	case bytecode.OpCastBool:
		return runtime.ToBoolean(val), true
	case bytecode.OpNot:
		return runtime.Boolean(!runtime.ToBoolean(val)), true
	case bytecode.OpNegate:
		return negate(val), true
	case bytecode.OpFlipPositive:
		return positive(val), true
	case bytecode.OpFlipNegative:
		return negative(val), true
	default:
		return nil, false
	}
}

func foldBinary(op bytecode.Opcode, left, right runtime.Value, bg context.Context) (runtime.Value, bool) {
	switch op {
	case bytecode.OpAdd:
		return runtime.Add(bg, left, right), true
	case bytecode.OpSub:
		return runtime.Subtract(bg, left, right), true
	case bytecode.OpMul:
		return runtime.Multiply(bg, left, right), true
	case bytecode.OpDiv:
		lv := runtime.ToNumberOnly(bg, left)

		if _, ok := lv.(runtime.Int); ok {
			rv := runtime.ToNumberOnly(bg, right)
			if ri, ok := rv.(runtime.Int); ok && ri == 0 {
				return nil, false
			}
		}

		return runtime.Divide(bg, left, right), true
	case bytecode.OpMod:
		if r, _ := runtime.ToInt(bg, right); r == 0 {
			return nil, false
		}

		return runtime.Modulus(bg, left, right), true
	case bytecode.OpCmp:
		return compare(bg, left, right), true
	case bytecode.OpEq:
		return equals(bg, left, right), true
	case bytecode.OpNe:
		return notEquals(bg, left, right), true
	case bytecode.OpGt:
		return greaterThan(bg, left, right), true
	case bytecode.OpLt:
		return lessThan(bg, left, right), true
	case bytecode.OpGte:
		return greaterThanOrEqual(bg, left, right), true
	case bytecode.OpLte:
		return lessThanOrEqual(bg, left, right), true
	default:
		return nil, false
	}
}

func buildConstIndex(constants []runtime.Value) map[string]int {
	index := make(map[string]int, len(constants))

	for i, val := range constants {
		if key, ok := constKey(val); ok {
			index[key] = i
		}
	}

	return index
}

func constKey(val runtime.Value) (string, bool) {
	if val == runtime.None {
		return "none", true
	}
	switch v := val.(type) {
	case runtime.Int:
		return fmt.Sprintf("i:%s", v.String()), true
	case runtime.Float:
		return fmt.Sprintf("f:%s", v.String()), true
	case runtime.String:
		return fmt.Sprintf("s:%s", v.String()), true
	case runtime.Boolean:
		if v {
			return "b:true", true
		}
		return "b:false", true
	default:
		return "", false
	}
}

func replaceWithConstLoad(inst *bytecode.Instruction, dst int, val runtime.Value, program *bytecode.Program, constIndex map[string]int) bool {
	newInst := buildConstLoad(dst, val, program, constIndex)

	if inst.Opcode == newInst.Opcode && inst.Operands == newInst.Operands {
		return false
	}

	*inst = newInst

	return true
}

func buildConstLoad(dst int, val runtime.Value, program *bytecode.Program, constIndex map[string]int) bytecode.Instruction {
	if val == runtime.None {
		return bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(dst))
	}

	switch v := val.(type) {
	case runtime.Boolean:
		if v {
			return bytecode.NewInstruction(bytecode.OpLoadBool, bytecode.NewRegister(dst), bytecode.Operand(1))
		}

		return bytecode.NewInstruction(bytecode.OpLoadBool, bytecode.NewRegister(dst), bytecode.Operand(0))
	case runtime.Int:
		if v == 0 {
			return bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(dst))
		}
	}

	key, ok := constKey(val)
	if !ok {
		return bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(dst), bytecode.NewConstant(appendConst(program, constIndex, val)))
	}

	if idx, ok := constIndex[key]; ok {
		return bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(dst), bytecode.NewConstant(idx))
	}

	return bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(dst), bytecode.NewConstant(appendConst(program, constIndex, val)))
}

func appendConst(program *bytecode.Program, constIndex map[string]int, val runtime.Value) int {
	idx := len(program.Constants)
	program.Constants = append(program.Constants, val)

	if key, ok := constKey(val); ok {
		constIndex[key] = idx
	}

	return idx
}

func compare(_ context.Context, left, right runtime.Value) runtime.Int {
	return runtime.Int(runtime.CompareValues(right, left))
}

func equals(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) == 0
}

func notEquals(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) != 0
}

func greaterThan(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) > 0
}

func greaterThanOrEqual(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) >= 0
}

func lessThan(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) < 0
}

func lessThanOrEqual(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) <= 0
}

func negate(input runtime.Value) runtime.Value {
	switch val := input.(type) {
	case runtime.Int:
		return -val
	case runtime.Float:
		return -val
	case runtime.Boolean:
		return !val
	default:
		return runtime.None
	}
}

func negative(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return -value
	case runtime.Float:
		return -value
	default:
		return runtime.None
	}
}

func positive(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return +value
	case runtime.Float:
		return +value
	default:
		return runtime.None
	}
}

func increment(ctx context.Context, input runtime.Value) runtime.Value {
	left := runtime.ToNumberOnly(ctx, input)
	switch value := left.(type) {
	case runtime.Int:
		return value + 1
	case runtime.Float:
		return value + 1
	default:
		return runtime.None
	}
}

func decrement(ctx context.Context, input runtime.Value) runtime.Value {
	left := runtime.ToNumberOnly(ctx, input)
	switch value := left.(type) {
	case runtime.Int:
		return value - 1
	case runtime.Float:
		return value - 1
	default:
		return runtime.None
	}
}
