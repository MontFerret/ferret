package runtime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/operators"
)

type VM struct {
	env          *Environment
	globals      map[string]core.Value
	frames       []*Frame
	currentFrame *Frame
	pc           int
}

func NewVM(opts ...EnvironmentOption) *VM {
	vm := new(VM)
	vm.env = newEnvironment(opts)

	return vm
}

// TODO: Move program to the constructor. No need to pass it as an argument since the VM is stateful.
// But the environment can be passed as an argument.
func (vm *VM) Run(ctx context.Context, program *Program) ([]byte, error) {
	//tryCatch := func(pos int) bool {
	//	for _, pair := range program.CatchTable {
	//		if pos >= pair[0] && pos <= pair[1] {
	//			return true
	//		}
	//	}
	//
	//	return false
	//}

	loadData := func(op Operand) core.Value {
		if op.IsRegister() {
			return vm.currentFrame.registers[op.Register()]
		}

		return program.Constants[op.Constant()]
	}

	vm.currentFrame = newFrame(64, 0, nil)
	vm.frames = make([]*Frame, 4)
	vm.globals = make(map[string]core.Value)
	vm.pc = 0

loop:
	for vm.pc < len(program.Bytecode) {
		inst := program.Bytecode[vm.pc]
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		reg := vm.currentFrame.registers
		vm.pc++

		switch inst.Opcode {
		case OpMove:
			reg[dst] = loadData(src1)
		case OpLoadConst:
			reg[dst] = program.Constants[src1.Constant()]
		case OpStoreGlobal:
			vm.globals[program.Constants[dst.Constant()].String()] = loadData(src1)
		case OpLoadGlobal:
			reg[dst] = vm.globals[program.Constants[src1].String()]
		case OpAdd:
			reg[dst] = operators.Add(loadData(src1), loadData(src2))
		case OpSub:
			reg[dst] = operators.Subtract(loadData(src1), loadData(src2))
		case OpMulti:
			reg[dst] = operators.Multiply(loadData(src1), loadData(src2))
		case OpDiv:
			reg[dst] = operators.Divide(loadData(src1), loadData(src2))
		case OpMod:
			reg[dst] = operators.Modulus(loadData(src1), loadData(src2))
		case OpIncr:
			reg[dst] = operators.Increment(reg[dst])
		case OpDecr:
			reg[dst] = operators.Decrement(reg[dst])

		case OpCastBool:
			//stack.Push(values.ToBoolean(stack.Pop()))

		case OpArray:
			//size := arg
			//arr := values.NewSizedArray(size)
			//
			//// iterate from the end to the beginning
			//// because stack is LIFO
			//for i := size - 1; i >= 0; i-- {
			//	arr.MustSet(values.Int(i), stack.Pop())
			//}
			//
			//stack.Push(arr)

		case OpObject:
			//obj := values.NewObject()
			//propertyCount := arg
			//
			//for i := 0; i < propertyCount; i++ {
			//	value := stack.Pop()
			//	key := stack.Pop()
			//	obj.Set(values.ToString(key), value)
			//}
			//
			//stack.Push(obj)

		case OpLoadProperty, OpLoadPropertyOptional:
			//prop := stack.Pop()
			//val := stack.Pop()
			//
			//switch getter := prop.(type) {
			//case values.String:
			//	switch src := val.(type) {
			//	case *values.Object:
			//		stack.Push(src.MustGetOr(getter, values.None))
			//	case core.Keyed:
			//		out, err := src.GetByKey(ctx, getter.String())
			//
			//		if err == nil {
			//			stack.Push(out)
			//		} else if op == OpLoadPropertyOptional {
			//			stack.Push(values.None)
			//		} else {
			//			return nil, err
			//		}
			//	default:
			//		if op != OpLoadPropertyOptional {
			//			return nil, core.TypeError(src, types.Object, types.Keyed)
			//		}
			//
			//		stack.Push(values.None)
			//	}
			//case values.Float, values.Int:
			//	switch src := val.(type) {
			//	case *values.Array:
			//		idx := values.ToInt(getter)
			//
			//		stack.Push(src.Get(idx))
			//	case core.Indexed:
			//		out, err := src.GetByIndex(ctx, int(values.ToInt(getter)))
			//
			//		if err == nil {
			//			stack.Push(out)
			//		} else if op == OpLoadPropertyOptional {
			//			stack.Push(values.None)
			//		} else {
			//			return nil, err
			//		}
			//	default:
			//		if op != OpLoadPropertyOptional {
			//			return nil, core.TypeError(src, types.Array, types.Indexed)
			//		}
			//
			//		stack.Push(values.None)
			//	}
			//}
		case OpNegate:
			//stack.Push(values.Negate(stack.Pop()))

		case OpFlipPositive:
			//stack.Push(values.Positive(stack.Pop()))

		case OpFlipNegative:
			//stack.Push(values.Negative(stack.Pop()))

		case OpNot:
			//stack.Push(!values.ToBoolean(stack.Pop()))

		case OpEq:
			//left := stack.Pop()
			//right := stack.Pop()
			//stack.Push(values.NewBoolean(values.Compare(left, right) == 0))

		case OpNeq:
			//left := stack.Pop()
			//right := stack.Pop()
			//stack.Push(values.NewBoolean(values.Compare(left, right) != 0))

		case OpGt:
			//right := stack.Pop()
			//left := stack.Pop()
			//stack.Push(values.NewBoolean(values.Compare(left, right) > 0))

		case OpLt:
			//right := stack.Pop()
			//left := stack.Pop()
			//stack.Push(values.NewBoolean(values.Compare(left, right) < 0))

		case OpGte:
			//right := stack.Pop()
			//left := stack.Pop()
			//stack.Push(values.NewBoolean(values.Compare(left, right) >= 0))

		case OpLte:
			//right := stack.Pop()
			//left := stack.Pop()
			//stack.Push(values.NewBoolean(values.Compare(left, right) <= 0))

		case OpIn:
			//right := stack.Pop()
			//left := stack.Pop()
			//stack.Push(values.Contains(right, left))

		case OpNotIn:
			//right := stack.Pop()
			//left := stack.Pop()
			//stack.Push(!values.Contains(right, left))

		case OpLike:
			//right := stack.Pop()
			//left := stack.Pop()
			//res, err := operators.Like(left, right)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else {
			//	return nil, err
			//}

		case OpNotLike:
			//right := stack.Pop()
			//left := stack.Pop()
			//res, err := operators.Like(left, right)
			//
			//if err == nil {
			//	stack.Push(!res)
			//} else {
			//	return nil, err
			//}

		case OpRegexpPositive:
			//reg, err := values.ToRegexp(stack.Pop())
			//
			//if err == nil {
			//	stack.Push(reg.Match(stack.Pop()))
			//} else if tryCatch(vm.pc) {
			//	stack.Push(values.False)
			//} else {
			//	return nil, err
			//}

		case OpRegexpNegative:
			//reg, err := values.ToRegexp(stack.Pop())
			//
			//if err == nil {
			//	stack.Push(!reg.Match(stack.Pop()))
			//} else if tryCatch(vm.pc) {
			//	stack.Push(values.True)
			//} else {
			//	return nil, err
			//}

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
			//right := stack.Pop()
			//left := stack.Pop()
			//res, err := operators.Range(left, right)
			//
			//if err == nil {
			//	stack.Push(res)
			//} else {
			//	return nil, err
			//}
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

		case OpJump:
			//vm.pc += arg

		case OpJumpBackward:
			//vm.pc -= arg

		case OpJumpIfFalse:
			//if !values.ToBoolean(stack.Peek()) {
			//	vm.pc += arg
			//}

		case OpJumpIfTrue:
			//if values.ToBoolean(stack.Peek()) {
			//	vm.pc += arg
			//}

		case OpLoopReturn:
			// pop the return value from the stack
			//res := stack.Pop()
			//ds := stack.Get(arg).(*DataSet)
			//ds.Push(res)

		case OpReturn:
			break loop
		}
	}

	return vm.currentFrame.registers[ResultOperand].MarshalJSON()
}
