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

		case OpTrue:
			stack.Push(values.True)

		case OpFalse:
			stack.Push(values.False)

		case OpPop:
			stack.Pop()

		case OpDefineGlobal:
			vm.globals[program.Constants[arg].String()] = stack.Pop()

		case OpGetGlobal:
			stack.Push(vm.globals[program.Constants[arg].String()])

		case OpNegate:
			stack.Push(values.Negate(stack.Pop()))

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

		case OpIncrement:
			stack.Push(operators.Increment(stack.Pop()))

		case OpDecrement:
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

		case OpReturn:
			res := stack.Pop()

			return res.MarshalJSON()
		}
	}

	// TODO: return error
	// program should exit with return statement
	return nil, nil
}
