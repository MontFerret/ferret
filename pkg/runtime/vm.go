package runtime

import (
	"context"
	"errors"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"io"
)

type VM struct {
	ip      int
	stack   *Stack
	globals map[string]core.Value
	env     *Environment
}

func NewVM(opts ...EnvironmentOption) *VM {
	vm := new(VM)
	vm.env = newEnvironment(opts)

	return vm
}

func (vm *VM) Run(ctx context.Context, program *Program) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}

			program = nil
		}
	}()

	// TODO: Add code analysis to calculate the number of operands and variables
	stack := NewStack(len(program.Bytecode), 8)
	vm.stack = stack
	vm.globals = make(map[string]core.Value)
	vm.ip = 0

loop:
	for vm.ip < len(program.Bytecode) {
		op := program.Bytecode[vm.ip]
		arg := program.Arguments[vm.ip]
		vm.ip++

		switch op {
		case OpPush:
			stack.Push(program.Constants[arg])

		case OpPop:
			stack.Pop()

		case OpPopClose:
			closable, ok := stack.Pop().(io.Closer)

			if ok {
				closable.Close()
			}

		case OpStoreGlobal:
			vm.globals[program.Constants[arg].String()] = stack.Pop()

		case OpLoadGlobal:
			stack.Push(vm.globals[program.Constants[arg].String()])

		case OpStoreLocal:
			stack.SetVariable(arg, stack.Pop())

		case OpPopLocal:
			stack.PopVariable()

		case OpLoadLocal:
			stack.Push(stack.GetVariable(arg))

		case OpNone:
			stack.Push(values.None)

		case OpConstant:
			stack.Push(program.Constants[arg])

		case OpCastBool:
			stack.Push(values.ToBoolean(stack.Pop()))

		case OpTrue:
			stack.Push(values.True)

		case OpFalse:
			stack.Push(values.False)

		case OpArray:
			size := arg
			arr := values.NewSizedArray(size)

			// iterate from the end to the beginning
			// because stack is LIFO
			for i := size - 1; i >= 0; i-- {
				arr.MustSet(values.Int(i), stack.Pop())
			}

			stack.Push(arr)

		case OpObject:
			obj := values.NewObject()
			propertyCount := arg

			for i := 0; i < propertyCount; i++ {
				value := stack.Pop()
				key := stack.Pop()
				obj.Set(values.ToString(key), value)
			}

			stack.Push(obj)

		case OpLoadProperty, OpLoadPropertyOptional:
			prop := stack.Pop()
			val := stack.Pop()

			switch getter := prop.(type) {
			case values.String:
				switch src := val.(type) {
				case *values.Object:
					stack.Push(src.MustGetOr(getter, values.None))
				case core.Keyed:
					out, err := src.GetByKey(ctx, getter.String())

					if err == nil {
						stack.Push(out)
					} else if op == OpLoadPropertyOptional {
						stack.Push(values.None)
					} else {
						return nil, err
					}
				default:
					if op != OpLoadPropertyOptional {
						return nil, core.TypeError(src, types.Object, types.Keyed)
					}

					stack.Push(values.None)
				}
			case values.Float, values.Int:
				switch src := val.(type) {
				case *values.Array:
					idx := values.ToInt(getter)

					stack.Push(src.Get(idx))
				case core.Indexed:
					out, err := src.GetByIndex(ctx, int(values.ToInt(getter)))

					if err == nil {
						stack.Push(out)
					} else if op == OpLoadPropertyOptional {
						stack.Push(values.None)
					} else {
						return nil, err
					}
				default:
					if op != OpLoadPropertyOptional {
						return nil, core.TypeError(src, types.Array, types.Indexed)
					}

					stack.Push(values.None)
				}
			}
		case OpNegate:
			stack.Push(values.Negate(stack.Pop()))

		case OpFlipPositive:
			stack.Push(values.Positive(stack.Pop()))

		case OpFlipNegative:
			stack.Push(values.Negative(stack.Pop()))

		case OpNot:
			stack.Push(!values.ToBoolean(stack.Pop()))

		case OpEq:
			left := stack.Pop()
			right := stack.Pop()
			stack.Push(values.NewBoolean(values.Compare(left, right) == 0))

		case OpNeq:
			left := stack.Pop()
			right := stack.Pop()
			stack.Push(values.NewBoolean(values.Compare(left, right) != 0))

		case OpGt:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(values.Compare(left, right) > 0))

		case OpLt:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(values.Compare(left, right) < 0))

		case OpGte:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(values.Compare(left, right) >= 0))

		case OpLte:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(values.Compare(left, right) <= 0))

		case OpIn:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.Contains(right, left))

		case OpNotIn:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(!values.Contains(right, left))

		case OpLike:
			right := stack.Pop()
			left := stack.Pop()
			res, err := operators.Like(left, right)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}

		case OpNotLike:
			right := stack.Pop()
			left := stack.Pop()
			res, err := operators.Like(left, right)

			if err == nil {
				stack.Push(!res)
			} else {
				return nil, err
			}

		case OpAdd:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(operators.Add(left, right))

		case OpSub:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(operators.Subtract(left, right))

		case OpMulti:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(operators.Multiply(left, right))

		case OpDiv:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(operators.Divide(left, right))

		case OpMod:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(operators.Modulus(left, right))

		case OpIncr:
			stack.Push(operators.Increment(stack.Pop()))

		case OpDecr:
			stack.Push(operators.Decrement(stack.Pop()))

		case OpRegexpPositive:
			reg := program.Constants[arg].(*values.Regexp)
			stack.Push(reg.Match(stack.Pop()))

		case OpRegexpNegative:
			reg := program.Constants[arg].(*values.Regexp)
			stack.Push(!reg.Match(stack.Pop()))

		case OpCall:
			fnName := stack.Pop().String()
			res, err := vm.env.GetFunction(fnName)(ctx)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}

		case OpCall1:
			arg := stack.Pop()
			fnName := stack.Pop().String()
			res, err := vm.env.GetFunction(fnName)(ctx, arg)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}

		case OpCall2:
			arg2 := stack.Pop()
			arg1 := stack.Pop()
			fnName := stack.Pop().String()
			res, err := vm.env.GetFunction(fnName)(ctx, arg1, arg2)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}

		case OpCall3:
			arg3 := stack.Pop()
			arg2 := stack.Pop()
			arg1 := stack.Pop()
			fnName := stack.Pop().String()
			res, err := vm.env.GetFunction(fnName)(ctx, arg1, arg2, arg3)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}

		case OpCall4:
			arg4 := stack.Pop()
			arg3 := stack.Pop()
			arg2 := stack.Pop()
			arg1 := stack.Pop()
			fnName := stack.Pop().String()
			res, err := vm.env.GetFunction(fnName)(ctx, arg1, arg2, arg3, arg4)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}
		case OpCallN:
			// pop arguments from the stack
			// and push them to the arguments array
			// in reverse order because stack is LIFO and arguments array is FIFO
			argCount := arg
			args := make([]core.Value, argCount)

			for i := argCount - 1; i >= 0; i-- {
				args[i] = stack.Pop()
			}

			// pop the function name from the stack
			fnName := stack.Pop().String()

			// call the function
			res, err := vm.env.GetFunction(fnName)(ctx, args...)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}

		case OpRange:
			right := stack.Pop()
			left := stack.Pop()
			res, err := operators.Range(left, right)

			if err == nil {
				stack.Push(res)
			} else {
				return nil, err
			}
		case OpLoopDestinationInit:
			ds := NewDataSet(arg == 1)

			stack.Push(values.NewBoxedValue(ds))

		case OpLoopSourceInit:
			// start a new iteration
			// get the data source
			src := stack.Pop()

			switch src := src.(type) {
			case core.Iterable:
				iterator, err := src.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				stack.Push(values.NewBoxedValue(iterator))
			default:
				return nil, core.TypeError(src, types.Iterable)
			}

		case OpLoopHasNext:
			boxed := stack.Peek()
			iterator := boxed.Unwrap().(core.Iterator)
			hasNext, err := iterator.HasNext(ctx)

			if err != nil {
				return nil, err
			}

			stack.Push(values.NewBoolean(hasNext))

		case OpLoopNext, OpLoopNextValue, OpLoopNextCounter:
			boxed := stack.Peek()
			iterator := boxed.Unwrap().(core.Iterator)

			// get the next value from the iterator
			val, key, err := iterator.Next(ctx)

			if err != nil {
				return nil, err
			}

			switch op {
			case OpLoopNextValue:
				stack.Push(val)
			case OpLoopNextCounter:
				stack.Push(key)
			default:
				stack.Push(key)
				stack.Push(val)
			}

		case OpJump:
			vm.ip += arg

		case OpJumpBackward:
			vm.ip -= arg

		case OpJumpIfFalse:
			if !values.ToBoolean(stack.Peek()) {
				vm.ip += arg
			}

		case OpJumpIfTrue:
			if values.ToBoolean(stack.Peek()) {
				vm.ip += arg
			}

		case OpLoopReturn:
			// pop the return value from the stack
			res := stack.Pop()
			pos := stack.Len() - arg
			ds := stack.Get(pos).(*values.Boxed).Unwrap().(*DataSet)
			ds.Push(res)

		case OpReturn:
			break loop
		}
	}

	return stack.Pop().MarshalJSON()
}
