package runtime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type VM struct {
	env          *Environment
	program      *Program
	globals      map[string]core.Value
	frames       []*Frame
	currentFrame *Frame
	pc           int
}

func NewVM(program *Program) *VM {
	vm := new(VM)
	vm.program = program

	return vm
}

func (vm *VM) Run(ctx context.Context, opts ...EnvironmentOption) (core.Value, error) {
	tryCatch := func(pos int) bool {
		for _, pair := range vm.program.CatchTable {
			if pos >= pair[0] && pos <= pair[1] {
				return true
			}
		}

		return false
	}

	vm.env = newEnvironment(opts)
	vm.currentFrame = newFrame(vm.program.Registers, 0, nil)
	vm.frames = make([]*Frame, 4)
	vm.globals = make(map[string]core.Value)
	vm.pc = 0
	program := vm.program

loop:
	for vm.pc < len(program.Bytecode) {
		inst := program.Bytecode[vm.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		reg := vm.currentFrame.registers
		vm.pc++

		switch op {
		case OpMove:
			reg[dst] = reg[src1]
		case OpLoadConst:
			reg[dst] = program.Constants[src1.Constant()]
		case OpStoreGlobal:
			vm.globals[program.Constants[dst.Constant()].String()] = reg[src1]
		case OpLoadGlobal:
			reg[dst] = vm.globals[program.Constants[src1.Constant()].String()]
		case OpJump:
			vm.pc += int(dst)
		case OpJumpBackward:
			vm.pc -= int(dst)
		case OpJumpIfFalse:
			if !values.ToBoolean(reg[dst]) {
				vm.pc += int(src1)
			}
		case OpJumpIfTrue:
			if values.ToBoolean(reg[dst]) {
				vm.pc += int(src1)
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

					reg[dst] = src.Get(idx)
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
		case OpCall, OpCallSafe:
			//fnName := stack.Pop().String()
			//res, err := vm.env.GetFunction(fnName)(ctx)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else if op == OpCallSafe || tryCatch(vm.pc) {
			//	stack.Push(values.None)
			//} else {
			//	return nil, err
			//}
		case OpCall1, OpCall1Safe:
			//arg := stack.Pop()
			//fnName := stack.Pop().String()
			//res, err := vm.env.GetFunction(fnName)(ctx, arg)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else if op == OpCall1Safe || tryCatch(vm.pc) {
			//	stack.Push(values.None)
			//} else {
			//	return nil, err
			//}
		case OpCall2, OpCall2Safe:
			//arg2 := stack.Pop()
			//arg1 := stack.Pop()
			//fnName := stack.Pop().String()
			//res, err := vm.env.GetFunction(fnName)(ctx, arg1, arg2)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else if op == OpCall2Safe || tryCatch(vm.pc) {
			//	stack.Push(values.None)
			//} else {
			//	return nil, err
			//}
		case OpCall3, OpCall3Safe:
			//arg3 := stack.Pop()
			//arg2 := stack.Pop()
			//arg1 := stack.Pop()
			//fnName := stack.Pop().String()
			//res, err := vm.env.GetFunction(fnName)(ctx, arg1, arg2, arg3)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else if op == OpCall3Safe || tryCatch(vm.pc) {
			//	stack.Push(values.None)
			//} else {
			//	return nil, err
			//}
		case OpCall4, OpCall4Safe:
			//arg4 := stack.Pop()
			//arg3 := stack.Pop()
			//arg2 := stack.Pop()
			//arg1 := stack.Pop()
			//fnName := stack.Pop().String()
			//res, err := vm.env.GetFunction(fnName)(ctx, arg1, arg2, arg3, arg4)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else if op == OpCall4Safe || tryCatch(vm.pc) {
			//	stack.Push(values.None)
			//} else {
			//	return nil, err
			//}
		case OpCallN, OpCallNSafe:
			//// pop arguments from the stack
			//// and push them to the arguments array
			//// in reverse order because stack is LIFO and arguments array is FIFO
			//argCount := arg
			//args := make([]core.Value, argCount)
			//
			//for i := argCount - 1; i >= 0; i-- {
			//	args[i] = stack.Pop()
			//}
			//
			//// pop the function name from the stack
			//fnName := stack.Pop().String()
			//
			//// call the function
			//res, err := vm.env.GetFunction(fnName)(ctx, args...)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else if op == OpCallNSafe || tryCatch(vm.pc) {
			//	stack.Push(values.None)
			//} else {
			//	return nil, err
			//}
		case OpRange:
			res, err := operators.Range(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpLoopInitOutput:
			//stack.Push(NewDataSet(arg == 1))

		case OpLoopUnwrapOutput:
			//ds := stack.Pop().(*DataSet)
			//stack.Push(ds.ToArray())

		case OpForLoopInitInput:
			//// start a new iteration
			//// get the data source
			//input := stack.Pop()
			//
			//switch src := input.(type) {
			//case core.Iterable:
			//	iterator, err := src.Iterate(ctx)
			//
			//	if err != nil {
			//		return nil, err
			//	}
			//
			//	stack.Push(values.NewBoxedValue(iterator))
			//default:
			//	return nil, core.TypeError(src, types.Iterable)
			//}

		case OpForLoopHasNext:
			//boxed := stack.Peek()
			//iterator := boxed.Unwrap().(core.Iterator)
			//hasNext, err := iterator.HasNext(ctx)
			//
			//if err != nil {
			//	return nil, err
			//}
			//
			//stack.Push(values.NewBoolean(hasNext))

		case OpForLoopNext, OpForLoopNextValue, OpForLoopNextCounter:
			//boxed := stack.Peek()
			//iterator := boxed.Unwrap().(core.Iterator)
			//
			//// get the next value from the iterator
			//val, key, err := iterator.Next(ctx)
			//
			//if err != nil {
			//	return nil, err
			//}
			//
			//switch op {
			//case OpForLoopNextValue:
			//	stack.Push(val)
			//case OpForLoopNextCounter:
			//	stack.Push(key)
			//default:
			//	stack.Push(key)
			//	stack.Push(val)
			//}

		case OpWhileLoopInitCounter:
			//stack.Push(values.ZeroInt)

		case OpWhileLoopNext:
			//counter := stack.Pop().(values.Int)
			//// increment the counter for the next iteration
			//stack.Push(counter + 1)
			//// put the current counter value
			//stack.Push(counter)
		case OpLoopReturn:
			// pop the return value from the stack
			//res := stack.Pop()
			//ds := stack.Get(arg).(*DataSet)
			//ds.Push(res)

		case OpReturn:
			break loop
		}
	}

	return vm.currentFrame.registers[ResultOperand], nil
}
