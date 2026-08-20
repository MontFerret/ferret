package optimization

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
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
	case runtime.Int, runtime.Float, runtime.Duration, runtime.String, runtime.Boolean:
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
	case runtime.Duration:
		if bv, ok := b.(runtime.Duration); ok {
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
		result, err := runtime.Not(val)
		return result, err == nil
	case bytecode.OpNegate:
		result, err := runtime.Negative(val)
		return result, err == nil
	case bytecode.OpFlipPositive:
		result, err := runtime.Positive(val)
		return result, err == nil
	case bytecode.OpFlipNegative:
		result, err := runtime.Negative(val)
		return result, err == nil
	default:
		return nil, false
	}
}

func foldBinary(op bytecode.Opcode, left, right runtime.Value, bg context.Context) (runtime.Value, bool) {
	switch op {
	case bytecode.OpAdd:
		result, err := runtime.Add(bg, left, right)
		return result, err == nil
	case bytecode.OpSub:
		result, err := runtime.Subtract(bg, left, right)
		return result, err == nil
	case bytecode.OpMul:
		result, err := runtime.Multiply(bg, left, right)
		return result, err == nil
	case bytecode.OpDiv:
		result, err := runtime.Divide(bg, left, right)

		return result, err == nil
	case bytecode.OpMod:
		result, err := runtime.Mod(bg, left, right)
		return result, err == nil
	case bytecode.OpCmp:
		result, err := runtime.CompareValues(bg, right, left)
		return runtime.Int(result), err == nil
	case bytecode.OpEq, bytecode.OpNe, bytecode.OpGt, bytecode.OpLt, bytecode.OpGte, bytecode.OpLte:
		comparison, ok := comparisonOperatorForOpcode(op)
		if !ok {
			return nil, false
		}

		result, err := runtime.EvaluateComparison(bg, comparison, left, right)
		return result, err == nil
	default:
		return nil, false
	}
}

func comparisonOperatorForOpcode(op bytecode.Opcode) (operator.Binary, bool) {
	switch op {
	case bytecode.OpEq:
		return operator.Equal, true
	case bytecode.OpNe:
		return operator.NotEqual, true
	case bytecode.OpGt:
		return operator.Greater, true
	case bytecode.OpLt:
		return operator.Less, true
	case bytecode.OpGte:
		return operator.GreaterOrEqual, true
	case bytecode.OpLte:
		return operator.LessOrEqual, true
	default:
		return operator.Unknown, false
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
	case runtime.Duration:
		return fmt.Sprintf("d:%d", int64(v)), true
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
