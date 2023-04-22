package runtime_v2

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime_v2/operators"
)

const DefaultStackSize = 128

type VM struct {
	stack   *Stack
	globals map[string]core.Value
	ip      int
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) Run(ctx context.Context, program *Program) ([]byte, error) {
	stack := NewStack(DefaultStackSize)
	vm.stack = stack
	vm.globals = make(map[string]core.Value)
	vm.ip = 0

	for vm.ip < len(program.Bytecode) {
		op := program.Bytecode[vm.ip]
		arg := program.Arguments[vm.ip]
		vm.ip++

		switch op {
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

		case OpPop:
			stack.Pop()

		case OpDefineGlobal:
			vm.globals[program.Constants[arg].String()] = stack.Pop()

		case OpGetGlobal:
			stack.Push(vm.globals[program.Constants[arg].String()])

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
			stack.Push(values.NewBoolean(left.Compare(right) == 0))

		case OpNeq:
			left := stack.Pop()
			right := stack.Pop()
			stack.Push(values.NewBoolean(left.Compare(right) != 0))

		case OpGt:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(left.Compare(right) > 0))

		case OpLt:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(left.Compare(right) < 0))

		case OpGte:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(left.Compare(right) >= 0))

		case OpLte:
			right := stack.Pop()
			left := stack.Pop()
			stack.Push(values.NewBoolean(left.Compare(right) <= 0))

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

			if err != nil {
				panic(err)
			}

			stack.Push(res)

		case OpNotLike:
			right := stack.Pop()
			left := stack.Pop()
			res, err := operators.Like(left, right)

			if err != nil {
				panic(err)
			}

			stack.Push(!res)

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

		case OpRange:
			right := stack.Pop()
			left := stack.Pop()
			res, err := operators.Range(left, right)

			if err != nil {
				panic(err)
			}

			stack.Push(res)

		case OpJumpIfFalse:
			if !values.ToBoolean(stack.Peek()) {
				vm.ip += arg
			}
		case OpJumpIfTrue:
			if values.ToBoolean(stack.Peek()) {
				vm.ip += arg
			}

		case OpJump:
			vm.ip += arg

		case OpReturn:
			res := stack.Pop()

			return res.MarshalJSON()
		}
	}

	// TODO: return error
	// program should exit with return statement
	return nil, nil
}
