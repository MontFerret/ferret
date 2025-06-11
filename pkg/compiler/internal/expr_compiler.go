package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	"regexp"
	"strings"
)

// Runtime functions
const (
	runtimeTypename = "TYPENAME"
	runtimeLength   = "LENGTH"
	runtimeWait     = "WAIT"
)

type ExprCompiler struct {
	ctx *FuncContext
}

func NewExprCompiler(ctx *FuncContext) *ExprCompiler {
	return &ExprCompiler{ctx: ctx}
}

func (ec *ExprCompiler) Compile(ctx fql.IExpressionContext) vm.Operand {
	if c := ctx.UnaryOperator(); c != nil {
		return ec.compileUnary(c, ctx)
	}

	if c := ctx.LogicalAndOperator(); c != nil {
		return ec.compileLogicalAnd(ctx.Predicate())
	}

	if c := ctx.LogicalOrOperator(); c != nil {
		return ec.compileLogicalOr(ctx.Predicate())
	}

	if c := ctx.GetTernaryOperator(); c != nil {
		return ec.compileTernary(ctx)
	}

	if c := ctx.Predicate(); c != nil {
		return ec.compilePredicate(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

// TODO: Free temporary registers if needed
func (ec *ExprCompiler) compileUnary(ctx fql.IUnaryOperatorContext, parent fql.IExpressionContext) vm.Operand {
	src := ec.Compile(parent.GetRight())
	dst := ec.ctx.Registers.Allocate(core.Temp)

	var op vm.Opcode

	if ctx.Not() != nil {
		op = vm.OpNot
	} else if ctx.Minus() != nil {
		op = vm.OpFlipNegative
	} else if ctx.Plus() != nil {
		op = vm.OpFlipPositive
	} else {
		panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
	}

	// We do not overwrite the source register
	ec.ctx.Emitter.EmitAB(op, dst, src)

	return dst
}

// TODO: Free temporary registers if needed
func (ec *ExprCompiler) compileLogicalAnd(ctx fql.IPredicateContext) vm.Operand {
	dst := ec.ctx.Registers.Allocate(core.Temp)

	// Execute left expression
	left := ec.compilePredicate(ctx.GetLeft())

	// Execute left expression
	ec.ctx.Emitter.EmitMove(dst, left)

	// Test if left is false and jump to the end
	end := ec.ctx.Emitter.EmitJumpIfFalse(dst, core.JumpPlaceholder)

	// If left is true, execute right expression
	right := ec.compilePredicate(ctx.GetRight())

	// And move the result to the destination register
	ec.ctx.Emitter.EmitMove(dst, right)

	ec.ctx.Emitter.PatchJumpNext(end)

	return dst
}

// TODO: Free temporary registers if needed
func (ec *ExprCompiler) compileLogicalOr(ctx fql.IPredicateContext) vm.Operand {
	dst := ec.ctx.Registers.Allocate(core.Temp)

	// Execute left expression
	left := ec.compilePredicate(ctx.GetLeft())

	// Execute left expression
	ec.ctx.Emitter.EmitMove(dst, left)

	// Test if left is true and jump to the end
	end := ec.ctx.Emitter.EmitJumpIfTrue(dst, core.JumpPlaceholder)

	// If left is false, execute right expression
	right := ec.compilePredicate(ctx.GetRight())

	// And move the result to the destination register
	ec.ctx.Emitter.EmitMove(dst, right)

	ec.ctx.Emitter.PatchJumpNext(end)

	return dst
}

// TODO: Free temporary registers if needed
func (ec *ExprCompiler) compileTernary(ctx fql.IExpressionContext) vm.Operand {
	dst := ec.ctx.Registers.Allocate(core.Temp)

	// Compile condition and put result in dst
	condReg := ec.Compile(ctx.GetCondition())
	ec.ctx.Emitter.EmitMove(dst, condReg)

	// Jump to 'false' branch if condition is false
	otherwise := ec.ctx.Emitter.EmitJumpIfFalse(dst, core.JumpPlaceholder)

	// True branch
	if onTrue := ctx.GetOnTrue(); onTrue != nil {
		trueReg := ec.Compile(onTrue)
		// Move result of true branch to dst
		ec.ctx.Emitter.EmitMove(dst, trueReg)
	}

	// Jump over false branch
	end := ec.ctx.Emitter.EmitJump(core.JumpPlaceholder)
	ec.ctx.Emitter.PatchJumpNext(otherwise)

	// False branch
	if onFalse := ctx.GetOnFalse(); onFalse != nil {
		falseReg := ec.Compile(onFalse)
		// Move result of false branch to dst
		ec.ctx.Emitter.EmitMove(dst, falseReg)
	}

	ec.ctx.Emitter.PatchJumpNext(end)

	return dst
}

// TODO: Free temporary registers if needed
func (ec *ExprCompiler) compilePredicate(ctx fql.IPredicateContext) vm.Operand {
	if c := ctx.ExpressionAtom(); c != nil {
		startCatch := ec.ctx.Emitter.Size()
		reg := ec.compileAtom(c)

		if c.ErrorOperator() != nil {
			jump := -1
			endCatch := ec.ctx.Emitter.Size()

			if c.ForExpression() != nil {
				// We jump back to finalize the loop before exiting
				jump = endCatch - 1
			}

			ec.ctx.CatchTable.Push(startCatch, endCatch, jump)
		}

		return reg
	}

	var opcode vm.Opcode
	dest := ec.ctx.Registers.Allocate(core.Temp)
	left := ec.compilePredicate(ctx.Predicate(0))
	right := ec.compilePredicate(ctx.Predicate(1))

	if op := ctx.EqualityOperator(); op != nil {
		switch ctx.GetText() {
		case "==":
			opcode = vm.OpEq
		case "!=":
			opcode = vm.OpNeq
		case ">":
			opcode = vm.OpGt
		case ">=":
			opcode = vm.OpGte
		case "<":
			opcode = vm.OpLt
		case "<=":
			opcode = vm.OpLte
		default:
			panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
		}
	} else if op := ctx.ArrayOperator(); op != nil {
		// TODO: Implement me
		panic(runtime.Error(runtime.ErrNotImplemented, "array operator"))
	} else if op := ctx.InOperator(); op != nil {
		if op.Not() == nil {
			opcode = vm.OpIn
		} else {
			opcode = vm.OpNotIn
		}
	} else if op := ctx.LikeOperator(); op != nil {
		if op.(*fql.LikeOperatorContext).Not() == nil {
			opcode = vm.OpLike
		} else {
			opcode = vm.OpNotLike
		}
	}

	ec.ctx.Emitter.EmitABC(opcode, dest, left, right)

	return dest
}

// TODO: Free temporary registers if needed
func (ec *ExprCompiler) compileAtom(ctx fql.IExpressionAtomContext) vm.Operand {
	var opcode vm.Opcode
	var isSet bool

	if op := ctx.MultiplicativeOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "*":
			opcode = vm.OpMulti
		case "/":
			opcode = vm.OpDiv
		case "%":
			opcode = vm.OpMod
		default:
			panic(runtime.Error(core.ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "+":
			opcode = vm.OpAdd
		case "-":
			opcode = vm.OpSub
		default:
			panic(runtime.Error(core.ErrUnexpectedToken, op.GetText()))
		}

	} else if op := ctx.RegexpOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "=~":
			opcode = vm.OpRegexpPositive
		case "!~":
			opcode = vm.OpRegexpNegative
		default:
			panic(runtime.Error(core.ErrUnexpectedToken, op.GetText()))
		}
	}

	if isSet {
		regLeft := ec.compileAtom(ctx.ExpressionAtom(0))
		regRight := ec.compileAtom(ctx.ExpressionAtom(1))
		dst := ec.ctx.Registers.Allocate(core.Temp)

		if opcode == vm.OpRegexpPositive || opcode == vm.OpRegexpNegative {
			if regRight.IsConstant() {
				val := ec.ctx.Symbols.Constant(regRight)

				// Verify that the expression is a valid regular expression
				regexp.MustCompile(val.String())
			}
		}

		ec.ctx.Emitter.EmitABC(opcode, dst, regLeft, regRight)

		return dst
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return ec.CompileFunctionCallExpression(c)
	} else if c := ctx.RangeOperator(); c != nil {
		return ec.CompileRangeOperator(c)
	} else if c := ctx.Literal(); c != nil {
		return ec.ctx.LiteralCompiler.Compile(c)
	} else if c := ctx.Variable(); c != nil {
		return ec.CompileVariable(c)
	} else if c := ctx.MemberExpression(); c != nil {
		return ec.CompileMemberExpression(c)
	} else if c := ctx.Param(); c != nil {
		return ec.CompileParam(c)
	} else if c := ctx.ForExpression(); c != nil {
		return ec.ctx.LoopCompiler.Compile(c)
	} else if c := ctx.WaitForExpression(); c != nil {
		return ec.ctx.WaitCompiler.Compile(c)
	} else if c := ctx.Expression(); c != nil {
		return ec.Compile(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (ec *ExprCompiler) CompileMemberExpression(ctx fql.IMemberExpressionContext) vm.Operand {
	mes := ctx.MemberExpressionSource()
	segments := ctx.AllMemberExpressionPath()

	var src1 vm.Operand

	if c := mes.Variable(); c != nil {
		src1 = ec.CompileVariable(c)
	} else if c := mes.Param(); c != nil {
		src1 = ec.CompileParam(c)
	} else if c := mes.ObjectLiteral(); c != nil {
		src1 = ec.ctx.LiteralCompiler.CompileObjectLiteral(c)
	} else if c := mes.ArrayLiteral(); c != nil {
		src1 = ec.ctx.LiteralCompiler.CompileArrayLiteral(c)
	} else if c := mes.FunctionCall(); c != nil {
		// FOO()?.bar
		segment := segments[0]
		src1 = ec.CompileFunctionCall(c, segment.ErrorOperator() != nil)
	}

	var dst vm.Operand

	for _, segment := range segments {
		var src2 vm.Operand
		p := segment.(*fql.MemberExpressionPathContext)

		if c := p.PropertyName(); c != nil {
			src2 = ec.ctx.LiteralCompiler.CompilePropertyName(c)
		} else if c := p.ComputedPropertyName(); c != nil {
			src2 = ec.ctx.LiteralCompiler.CompileComputedPropertyName(c)
		}

		dst = ec.ctx.Registers.Allocate(core.Temp)

		// TODO: Replace with EmitLoadKey
		if p.ErrorOperator() != nil {
			ec.ctx.Emitter.EmitLoadPropertyOptional(dst, src1, src2)
		} else {
			ec.ctx.Emitter.EmitLoadProperty(dst, src1, src2)
		}

		src1 = dst
	}

	return dst
}

func (ec *ExprCompiler) CompileVariable(ctx fql.IVariableContext) vm.Operand {
	// Just return the register / constant index
	op, _, found := ec.ctx.Symbols.Resolve(ctx.GetText())

	if !found {
		panic(runtime.Error(core.ErrVariableNotFound, ctx.GetText()))
	}

	if op.IsRegister() {
		return op
	}

	reg := ec.ctx.Registers.Allocate(core.Temp)
	ec.ctx.Emitter.EmitLoadGlobal(reg, op)

	return reg
}

func (ec *ExprCompiler) CompileParam(ctx fql.IParamContext) vm.Operand {
	name := ctx.Identifier().GetText()
	reg := ec.ctx.Symbols.BindParam(name)
	ec.ctx.Emitter.EmitLoadParam(reg)

	return reg
}

func (ec *ExprCompiler) CompileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) vm.Operand {
	protected := ctx.ErrorOperator() != nil
	call := ctx.FunctionCall()

	return ec.CompileFunctionCall(call, protected)
}

func (ec *ExprCompiler) CompileFunctionCall(ctx fql.IFunctionCallContext, protected bool) vm.Operand {
	name := ec.functionName(ctx)
	seq := ec.CompileArgumentList(ctx.ArgumentList())

	switch name {
	case runtimeLength:
		dst := ec.ctx.Registers.Allocate(core.Temp)

		if seq == nil || len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeLength+": expected 1 argument"))
		}

		ec.ctx.Emitter.EmitAB(vm.OpLength, dst, seq[0])

		return dst
	case runtimeTypename:
		dst := ec.ctx.Registers.Allocate(core.Temp)

		if seq == nil || len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeTypename+": expected 1 argument"))
		}

		ec.ctx.Emitter.EmitAB(vm.OpType, dst, seq[0])

		return dst
	case runtimeWait:
		if len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeWait+": expected 1 argument"))
		}

		ec.ctx.Emitter.EmitA(vm.OpSleep, seq[0])

		return seq[0]
	default:
		dest := ec.ctx.Registers.Allocate(core.Temp)
		ec.ctx.Emitter.EmitLoadConst(dest, ec.ctx.Symbols.AddConstant(name))

		if !protected {
			ec.ctx.Emitter.EmitAs(vm.OpCall, dest, seq)
		} else {
			ec.ctx.Emitter.EmitAs(vm.OpProtectedCall, dest, seq)
		}

		return dest
	}
}

func (ec *ExprCompiler) CompileArgumentList(ctx fql.IArgumentListContext) core.RegisterSequence {
	var seq core.RegisterSequence
	// Get all array element expressions
	exps := ctx.AllExpression()
	size := len(exps)

	if size > 0 {
		// Allocate seq for function arguments
		seq = ec.ctx.Registers.AllocateSequence(size)

		// Evaluate each element into seq Registers
		for i, exp := range exps {
			// Compile expression and move to seq register
			srcReg := ec.Compile(exp)

			// TODO: Figure out how to remove OpMove and use Registers returned from each expression
			ec.ctx.Emitter.EmitMove(seq[i], srcReg)

			// Free source register if temporary
			if srcReg.IsRegister() {
				//ec.ctx.Registers.Free(srcReg)
			}
		}
	}

	return seq
}

func (ec *ExprCompiler) CompileRangeOperator(ctx fql.IRangeOperatorContext) vm.Operand {
	dst := ec.ctx.Registers.Allocate(core.Temp)
	start := ec.compileRangeOperand(ctx.GetLeft())
	end := ec.compileRangeOperand(ctx.GetRight())

	ec.ctx.Emitter.EmitRange(dst, start, end)

	return dst
}

func (ec *ExprCompiler) compileRangeOperand(ctx fql.IRangeOperandContext) vm.Operand {
	if c := ctx.Variable(); c != nil {
		return ec.CompileVariable(c)
	}

	if c := ctx.Param(); c != nil {
		return ec.CompileParam(c)
	}

	if c := ctx.IntegerLiteral(); c != nil {
		return ec.ctx.LiteralCompiler.CompileIntegerLiteral(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (ec *ExprCompiler) functionName(ctx fql.IFunctionCallContext) runtime.String {
	var name string
	funcNS := ctx.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += ctx.FunctionName().GetText()

	return runtime.NewString(strings.ToUpper(name))
}
