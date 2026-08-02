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
		return runtime.Boolean(!runtime.ToBoolean(val)), true
	case bytecode.OpNegate:
		result, err := runtime.NegateChecked(val)
		return result, err == nil
	case bytecode.OpFlipPositive:
		result, err := runtime.PositiveChecked(val)
		return result, err == nil
	case bytecode.OpFlipNegative:
		result, err := runtime.NegativeChecked(val)
		return result, err == nil
	default:
		return nil, false
	}
}

func foldBinary(op bytecode.Opcode, left, right runtime.Value, bg context.Context) (runtime.Value, bool) {
	switch op {
	case bytecode.OpAdd:
		result, err := runtime.AddChecked(bg, left, right)
		return result, err == nil
	case bytecode.OpSub:
		result, err := runtime.SubtractChecked(bg, left, right)
		return result, err == nil
	case bytecode.OpMul:
		result, err := runtime.MultiplyChecked(bg, left, right)
		return result, err == nil
	case bytecode.OpDiv:
		if _, ok := left.(runtime.Duration); ok {
			result, err := runtime.DivideChecked(bg, left, right)

			return result, err == nil
		}

		if _, ok := right.(runtime.Duration); ok {
			return nil, false
		}

		lv := runtime.ToNumberOnly(bg, left)

		if _, ok := lv.(runtime.Int); ok {
			rv := runtime.ToNumberOnly(bg, right)
			if ri, ok := rv.(runtime.Int); ok && ri == 0 {
				return nil, false
			}
		}

		result, err := runtime.DivideChecked(bg, left, right)

		return result, err == nil
	case bytecode.OpMod:
		if _, ok := left.(runtime.Duration); ok {
			return nil, false
		}

		if _, ok := right.(runtime.Duration); ok {
			return nil, false
		}

		if r, _ := runtime.ToInt(bg, right); r == 0 {
			return nil, false
		}

		return runtime.Modulus(bg, left, right), true
	case bytecode.OpCmp:
		result, err := runtime.CompareChecked(bg, right, left)
		return runtime.Int(result), err == nil
	case bytecode.OpEq:
		result, err := runtime.EqualChecked(bg, left, right)
		return result, err == nil
	case bytecode.OpNe:
		result, err := runtime.EqualChecked(bg, left, right)
		if err != nil {
			return nil, false
		}

		return !result, true
	case bytecode.OpGt:
		result, err := runtime.CompareChecked(bg, left, right)
		return runtime.Boolean(result > 0), err == nil
	case bytecode.OpLt:
		result, err := runtime.CompareChecked(bg, left, right)
		return runtime.Boolean(result < 0), err == nil
	case bytecode.OpGte:
		result, err := runtime.CompareChecked(bg, left, right)
		return runtime.Boolean(result >= 0), err == nil
	case bytecode.OpLte:
		result, err := runtime.CompareChecked(bg, left, right)
		return runtime.Boolean(result <= 0), err == nil
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
