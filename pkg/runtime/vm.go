package runtime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type VM struct {
	env       *Environment
	program   *Program
	globals   map[string]core.Value
	registers []core.Value
	pc        int
}

func NewVM(program *Program) *VM {
	vm := new(VM)
	vm.program = program

	return vm
}

func (vm *VM) Run(ctx context.Context, opts []EnvironmentOption) (core.Value, error) {
	// TODO: Return jump position if an error occurred within a wrapped loop
	tryCatch := func(pos int) bool {
		for _, pair := range vm.program.CatchTable {
			if pos >= pair[0] && pos <= pair[1] {
				return true
			}
		}

		return false
	}

	vm.env = newEnvironment(opts)
	vm.registers = make([]core.Value, vm.program.Registers)
	vm.globals = make(map[string]core.Value)
	vm.pc = 0
	program := vm.program

	// TODO: Add panic handling and snapshot the last instruction and frame that caused it
loop:
	for vm.pc < len(program.Bytecode) {
		inst := program.Bytecode[vm.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		reg := vm.registers
		vm.pc++

		switch op {
		case OpLoadNone:
			reg[dst] = values.None
		case OpLoadBool:
			reg[dst] = values.NewBoolean(src1 == 1)
		case OpMove:
			reg[dst] = reg[src1]
		case OpLoadConst:
			reg[dst] = program.Constants[src1.Constant()]
		case OpStoreGlobal:
			vm.globals[program.Constants[dst.Constant()].String()] = reg[src1]
		case OpLoadGlobal:
			reg[dst] = vm.globals[program.Constants[src1.Constant()].String()]
		case OpJump:
			vm.pc = int(dst)
		case OpJumpIfFalse:
			if !values.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case OpJumpIfTrue:
			if values.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case OpAdd:
			reg[dst] = operators.Add(reg[src1], reg[src2])
		case OpSub:
			reg[dst] = operators.Subtract(reg[src1], reg[src2])
		case OpMulti:
			reg[dst] = operators.Multiply(reg[src1], reg[src2])
		case OpDiv:
			reg[dst] = operators.Divide(reg[src1], reg[src2])
		case OpMod:
			reg[dst] = operators.Modulus(reg[src1], reg[src2])
		case OpIncr:
			reg[dst] = operators.Increment(reg[dst])
		case OpDecr:
			reg[dst] = operators.Decrement(reg[dst])
		case OpCastBool:
			reg[dst] = values.ToBoolean(reg[src1])
		case OpNegate:
			reg[dst] = values.Negate(reg[src1])
		case OpFlipPositive:
			reg[dst] = values.Positive(reg[src1])
		case OpFlipNegative:
			reg[dst] = values.Negative(reg[src1])
		case OpNot:
			reg[dst] = !values.ToBoolean(reg[src1])
		case OpEq:
			reg[dst] = values.NewBoolean(values.Compare(reg[src1], reg[src2]) == 0)
		case OpNeq:
			reg[dst] = values.NewBoolean(values.Compare(reg[src1], reg[src2]) != 0)
		case OpGt:
			reg[dst] = values.NewBoolean(values.Compare(reg[src1], reg[src2]) > 0)
		case OpLt:
			reg[dst] = values.NewBoolean(values.Compare(reg[src1], reg[src2]) < 0)
		case OpGte:
			reg[dst] = values.NewBoolean(values.Compare(reg[src1], reg[src2]) >= 0)
		case OpLte:
			reg[dst] = values.NewBoolean(values.Compare(reg[src1], reg[src2]) <= 0)
		case OpIn:
			reg[dst] = values.Contains(reg[src2], reg[src1])
		case OpNotIn:
			reg[dst] = !values.Contains(reg[src2], reg[src1])
		case OpLike:
			res, err := operators.Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpNotLike:
			res, err := operators.Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = !res
			} else {
				return nil, err
			}
		case OpRegexpPositive:
			r, err := values.ToRegexp(reg[src1])

			if err == nil {
				reg[dst] = r.Match(reg[src2])
			} else if tryCatch(vm.pc) {
				reg[dst] = values.False
			} else {
				return nil, err
			}
		case OpRegexpNegative:
			r, err := values.ToRegexp(reg[src1])

			if err == nil {
				reg[dst] = !r.Match(reg[src2])
			} else if tryCatch(vm.pc) {
				reg[dst] = values.False
			} else {
				return nil, err
			}
		case OpArray:
			var size int

			if src1 > 0 {
				size = src2.Register() - src1.Register() + 1
			}

			arr := values.NewArray(size)
			start := int(src1)
			end := int(src1) + size

			// Iterate over registers starting from src1 and up to the src2
			for i := start; i < end; i++ {
				arr.Push(reg[i])
			}

			reg[dst] = arr
		case OpObject:
			obj := values.NewObject()
			var args int

			if src1 > 0 {
				args = src2.Register() - src1.Register() + 1
			}

			start := int(src1)
			end := int(src1) + args

			for i := start; i < end; i += 2 {
				key := reg[i]
				value := reg[i+1]

				obj.Set(values.ToString(key), value)
			}

			reg[dst] = obj
		case OpLoadProperty, OpLoadPropertyOptional:
			val := reg[src1]
			prop := reg[src2]

			switch getter := prop.(type) {
			case values.String:
				switch src := val.(type) {
				case *values.Object:
					reg[dst] = src.MustGetOr(getter, values.None)
				case core.Keyed:
					out, err := src.GetByKey(ctx, getter.String())

					if err == nil {
						reg[dst] = out
					} else if op == OpLoadPropertyOptional {
						reg[dst] = values.None
					} else {
						return nil, err
					}
				default:
					if op != OpLoadPropertyOptional {
						return nil, core.TypeError(src, types.Object, types.Keyed)
					}

					reg[dst] = values.None
				}
			case values.Float, values.Int:
				switch src := val.(type) {
				case *values.Array:
					idx := values.ToInt(getter)

					reg[dst] = src.Get(int(idx))
				case core.Indexed:
					out, err := src.GetByIndex(ctx, int(values.ToInt(getter)))

					if err == nil {
						reg[dst] = out
					} else if op == OpLoadPropertyOptional {
						reg[dst] = values.None
					} else {
						return nil, err
					}
				default:
					if op != OpLoadPropertyOptional {
						return nil, core.TypeError(src, types.Array, types.Indexed)
					}

					reg[dst] = values.None
				}
			}
		case OpCall, OpCallSafe:
			var size int

			if src1 > 0 {
				size = src2.Register() - src1.Register() + 1
			}

			start := int(src1)
			end := int(src1) + size
			args := make([]core.Value, size)

			// Iterate over registers starting from src1 and up to the src2
			for i := start; i < end; i++ {
				args[i-start] = reg[i]
			}

			fnName := reg[dst].String()
			fn := vm.env.GetFunction(fnName)

			out, err := fn(ctx, args...)

			if err == nil {
				reg[dst] = out
			} else if op == OpCallSafe || tryCatch(vm.pc) {
				reg[dst] = values.None
			} else {
				return nil, err
			}
		case OpLength:
			val, ok := reg[src1].(core.Measurable)

			if ok {
				reg[dst] = values.NewInt(val.Length())
			} else if tryCatch(vm.pc) {
				reg[dst] = values.ZeroInt
			} else {
				return values.None, core.TypeError(reg[src1],
					types.String,
					types.Array,
					types.Object,
					types.Binary,
					types.Measurable,
				)
			}
		case OpRange:
			res, err := operators.Range(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpLoopBegin:
			reg[dst] = NewDataSet(src1 == 1)
		case OpLoopEnd:
			ds := reg[src1].(*DataSet)
			reg[dst] = ds.ToArray()
		case OpForLoopInit:
			input := reg[src1]

			switch src := input.(type) {
			case core.Iterable:
				iterator, err := src.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				reg[dst] = values.NewBoxedValue(iterator)
			default:
				if tryCatch(vm.pc) {
					// Fall back to an empty iterator
					reg[dst] = values.NewBoxedValue(values.NoopIter)
				} else {
					return nil, core.TypeError(src, types.Iterable)
				}
			}
		case OpForLoopNext:
			boxed := reg[src1]
			// TODO: Remove boxed value
			iterator := boxed.Unwrap().(core.Iterator)
			hasNext, err := iterator.HasNext(ctx)

			if err != nil {
				return nil, err
			}

			if hasNext {
				if err := iterator.Next(ctx); err != nil {
					return nil, err
				}
			} else {
				vm.pc = int(dst)
			}
		case OpForLoopValue:
			// TODO: Remove boxed value!!!
			iter := reg[src1].(*values.Boxed).Unwrap().(core.Iterator)
			reg[dst] = iter.Value()
		case OpForLoopKey:
			// TODO: Remove boxed value!!!
			iter := reg[src1].(*values.Boxed).Unwrap().(core.Iterator)
			reg[dst] = iter.Key()
		case OpWhileLoopInit:
			reg[dst] = values.Int(-1)
		case OpWhileLoopNext:
			cond := values.ToBoolean(reg[src1])

			if cond {
				reg[dst] = operators.Increment(reg[dst])
			} else {
				vm.pc = int(src2)
			}
		case OpWhileLoopValue:
			reg[dst] = reg[src1]
		case OpLoopPush:
			ds := reg[dst].(*DataSet)
			ds.Push(reg[src1])
		case OpReturn:
			break loop
		}
	}

	return vm.registers[NoopOperand], nil
}
