package internal

import (
	"regexp"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// Runtime functions
const (
	runtimeTypename = "TYPENAME"
	runtimeLength   = "LENGTH"
	runtimeWait     = "WAIT"
)

// ExprCompiler handles the compilation of expressions in FQL queries.
// It transforms expression operations from the AST into VM instructions.
type ExprCompiler struct {
	ctx *CompilerContext
}

// NewExprCompiler creates a new instance of ExprCompiler with the given compiler context.
func NewExprCompiler(ctx *CompilerContext) *ExprCompiler {
	return &ExprCompiler{ctx: ctx}
}

// Compile processes an expression from the FQL AST and delegates to the appropriate
// compilation method based on the expression type (unary, logical, ternary, or predicate).
// Parameters:
//   - ctx: The expression context from the AST
//
// Returns:
//   - An operand representing the compiled expression
//
// Panics if the expression type is not recognized.
func (ec *ExprCompiler) Compile(ctx fql.IExpressionContext) vm.Operand {
	if c := ctx.UnaryOperator(); c != nil {
		return ec.compileUnary(c, ctx)
	}

	if c := ctx.LogicalAndOperator(); c != nil {
		return ec.compileLogicalAnd(ctx)
	}

	if c := ctx.LogicalOrOperator(); c != nil {
		return ec.compileLogicalOr(ctx)
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

// compileUnary processes a unary operation (NOT, minus, plus) from the FQL AST.
// It compiles the operand expression and applies the appropriate unary operation to it.
// Parameters:
//   - ctx: The unary operator context from the AST
//   - parent: The parent expression context containing the operand
//
// Returns:
//   - An operand representing the result of the unary operation
//
// Panics if the unary operator type is not recognized.
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

// compileLogicalAnd processes a logical AND operation from the FQL AST.
// It implements short-circuit evaluation: if the left operand is falsy, it returns that value
// without evaluating the right operand. Otherwise, it evaluates and returns the right operand.
// Parameters:
//   - ctx: The expression context from the AST containing both operands
//
// Returns:
//   - An operand representing the result of the logical AND operation
func (ec *ExprCompiler) compileLogicalAnd(ctx fql.IExpressionContext) vm.Operand {
	left := ec.Compile(ctx.GetLeft())

	skip := ec.ctx.Emitter.NewLabel("and.false")
	done := ec.ctx.Emitter.NewLabel("and.done")
	dst := ec.ctx.Registers.Allocate(core.Temp)

	// If left is falsy, jump to skip and use left
	ec.ctx.Emitter.EmitJumpIfFalse(left, skip)

	// Otherwise evaluate right and use it
	right := ec.Compile(ctx.GetRight())
	ec.ctx.Emitter.EmitMove(dst, right)
	ec.ctx.Emitter.EmitJump(done)

	// Short-circuit: use left as result
	ec.ctx.Emitter.MarkLabel(skip)
	ec.ctx.Emitter.EmitMove(dst, left)

	ec.ctx.Emitter.MarkLabel(done)

	return dst
}

// TODO: Free temporary registers if needed

// compileLogicalOr processes a logical OR operation from the FQL AST.
// It implements short-circuit evaluation: if the left operand is truthy, it returns that value
// without evaluating the right operand. Otherwise, it evaluates and returns the right operand.
// Parameters:
//   - ctx: The expression context from the AST containing both operands
//
// Returns:
//   - An operand representing the result of the logical OR operation
func (ec *ExprCompiler) compileLogicalOr(ctx fql.IExpressionContext) vm.Operand {
	left := ec.Compile(ctx.GetLeft())

	next := ec.ctx.Emitter.NewLabel("or.false")
	done := ec.ctx.Emitter.NewLabel("or.done")
	dst := ec.ctx.Registers.Allocate(core.Temp)

	// If left is truthy, short-circuit and skip right
	ec.ctx.Emitter.EmitJumpIfTrue(left, next)

	// Otherwise evaluate right
	right := ec.Compile(ctx.GetRight())
	ec.ctx.Emitter.EmitMove(dst, right)
	ec.ctx.Emitter.EmitJump(done)

	// Short-circuit: use left value
	ec.ctx.Emitter.MarkLabel(next)
	ec.ctx.Emitter.EmitMove(dst, left)

	// Common exit
	ec.ctx.Emitter.MarkLabel(done)

	return dst
}

// TODO: Free temporary registers if needed

// compileTernary processes a ternary conditional operation (condition ? trueExpr : falseExpr) from the FQL AST.
// It evaluates the condition and then either the true expression or the false expression based on the result.
// Parameters:
//   - ctx: The expression context from the AST containing the condition, true expression, and false expression
//
// Returns:
//   - An operand representing the result of either the true or false expression
func (ec *ExprCompiler) compileTernary(ctx fql.IExpressionContext) vm.Operand {
	dst := ec.ctx.Registers.Allocate(core.Temp)

	// Compile condition and put result in dst
	condReg := ec.Compile(ctx.GetCondition())
	ec.ctx.Emitter.EmitMove(dst, condReg)

	// Define jump labels
	elseLabel := ec.ctx.Emitter.NewLabel()
	endLabel := ec.ctx.Emitter.NewLabel()

	// EndLabel to 'false' branch if condition is false
	ec.ctx.Emitter.EmitJumpIfFalse(dst, elseLabel)

	// True branch
	if onTrue := ctx.GetOnTrue(); onTrue != nil {
		trueReg := ec.Compile(onTrue)
		// Move result of true branch to dst
		ec.ctx.Emitter.EmitMove(dst, trueReg)
	}

	// EndLabel over false branch
	ec.ctx.Emitter.EmitJump(endLabel)
	// Mark label for 'else' branch
	ec.ctx.Emitter.MarkLabel(elseLabel)

	// False branch
	if onFalse := ctx.GetOnFalse(); onFalse != nil {
		falseReg := ec.Compile(onFalse)
		// Move result of false branch to dst
		ec.ctx.Emitter.EmitMove(dst, falseReg)
	}

	// EndLabel
	ec.ctx.Emitter.MarkLabel(endLabel)

	return dst
}

// TODO: Free temporary registers if needed

// compilePredicate processes a predicate expression from the FQL AST.
// It handles both atomic expressions and comparison operations (equality, array operations, IN, LIKE).
// For atomic expressions, it also handles error catching with the TRY operator.
// Parameters:
//   - ctx: The predicate context from the AST
//
// Returns:
//   - An operand representing the result of the predicate expression
//
// Panics if the operator type is not recognized or not implemented.
func (ec *ExprCompiler) compilePredicate(ctx fql.IPredicateContext) vm.Operand {
	if c := ctx.ExpressionAtom(); c != nil {
		startCatch := ec.ctx.Emitter.Size()
		reg := ec.compileAtom(c)

		if c.ErrorOperator() != nil {
			jump := -1
			endCatch := ec.ctx.Emitter.Size()

			if fe := c.ForExpression(); fe != nil {
				// Since FOR-IN loops depend on custom iterators,
				// We need to handle cleanup before exiting the loop.
				// TODO: Find a better way to handle this. The code assumes the knowledge of the internals of the FOR-IN loop.
				if fe.In() != nil {
					jump = endCatch - 1
				}
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
		switch op.GetText() {
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

// compileAtom processes an atomic expression from the FQL AST.
// It handles various types of expressions including arithmetic operations (*, /, %, +, -),
// regular expression operations (=~, !~), function calls, range operators, literals,
// variables, member expressions, parameters, and nested expressions.
// Parameters:
//   - ctx: The expression atom context from the AST
//
// Returns:
//   - An operand representing the result of the atomic expression
//
// Panics if the expression type is not recognized.
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

// CompileMemberExpression processes a member expression (e.g., obj.prop, arr[idx]) from the FQL AST.
// It handles property access for variables, parameters, object literals, array literals, and function calls.
// It supports both dot notation (obj.prop) and bracket notation (obj["prop"] or arr[idx]),
// as well as optional chaining with the ?. operator.
// Parameters:
//   - ctx: The member expression context from the AST
//
// Returns:
//   - An operand representing the value of the accessed property
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

// CompileVariable processes a variable reference from the FQL AST.
// It resolves the variable name in the symbol table and returns the register or constant
// containing its value. If the variable is not found, it panics with an error.
// Parameters:
//   - ctx: The variable context from the AST
//
// Returns:
//   - An operand representing the variable's value
//
// Panics if the variable is not found in the symbol table.
func (ec *ExprCompiler) CompileVariable(ctx fql.IVariableContext) vm.Operand {
	name := ctx.Identifier().GetText()
	// Just return the register / constant index
	op, _, found := ec.ctx.Symbols.Resolve(name)

	if !found {
		panic(runtime.Error(core.ErrVariableNotFound, name))
	}

	if op.IsRegister() {
		return op
	}

	reg := ec.ctx.Registers.Allocate(core.Temp)
	ec.ctx.Emitter.EmitMove(reg, op)

	return reg
}

// CompileParam processes a parameter reference (e.g., @paramName) from the FQL AST.
// It binds the parameter name in the symbol table and emits instructions to load
// the parameter value at runtime.
// Parameters:
//   - ctx: The parameter context from the AST
//
// Returns:
//   - An operand representing the parameter's value
func (ec *ExprCompiler) CompileParam(ctx fql.IParamContext) vm.Operand {
	name := ctx.Identifier().GetText()
	reg := ec.ctx.Registers.Allocate(core.Temp)
	ec.ctx.Emitter.EmitLoadParam(reg, ec.ctx.Symbols.BindParam(name))

	return reg
}

// CompileFunctionCallExpression processes a function call expression from the FQL AST.
// It handles both regular function calls and protected function calls (with TRY).
// Parameters:
//   - ctx: The function call expression context from the AST
//
// Returns:
//   - An operand representing the function call result
func (ec *ExprCompiler) CompileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) vm.Operand {
	protected := ctx.ErrorOperator() != nil
	call := ctx.FunctionCall()

	return ec.CompileFunctionCall(call, protected)
}

// CompileFunctionCall processes a function call from the FQL AST.
// It compiles the function arguments and delegates to CompileFunctionCallWith.
// Parameters:
//   - ctx: The function call context from the AST
//   - protected: Whether this is a protected call (with TRY)
//
// Returns:
//   - An operand representing the function call result
func (ec *ExprCompiler) CompileFunctionCall(ctx fql.IFunctionCallContext, protected bool) vm.Operand {
	return ec.CompileFunctionCallWith(ctx, protected, ec.CompileArgumentList(ctx.ArgumentList()))
}

// CompileFunctionCallWith processes a function call with pre-compiled arguments from the FQL AST.
// It extracts the function name and delegates to CompileFunctionCallByNameWith.
// Parameters:
//   - ctx: The function call context from the AST
//   - protected: Whether this is a protected call (with TRY)
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
func (ec *ExprCompiler) CompileFunctionCallWith(ctx fql.IFunctionCallContext, protected bool, seq core.RegisterSequence) vm.Operand {
	name := getFunctionName(ctx)

	return ec.CompileFunctionCallByNameWith(name, protected, seq)
}

// CompileFunctionCallByNameWith processes a function call by name with pre-compiled arguments.
// It handles both built-in runtime functions (TYPENAME, LENGTH, WAIT) and user-defined functions.
// For built-in functions, it emits specialized instructions. For user-defined functions,
// it delegates to compileUserFunctionCallWith.
// Parameters:
//   - name: The function name
//   - protected: Whether this is a protected call (with TRY)
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
//
// Panics if a built-in function is called with an incorrect number of arguments.
func (ec *ExprCompiler) CompileFunctionCallByNameWith(name runtime.String, protected bool, seq core.RegisterSequence) vm.Operand {
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
		return ec.compileUserFunctionCallWith(name, protected, seq)
	}
}

// compileUserFunctionCallWith processes a user-defined function call with pre-compiled arguments.
// It loads the function name as a constant, binds the function in the symbol table,
// and emits the appropriate call instruction based on the number of arguments and whether
// the call is protected.
// Parameters:
//   - name: The function name
//   - protected: Whether this is a protected call (with TRY)
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
func (ec *ExprCompiler) compileUserFunctionCallWith(name runtime.String, protected bool, seq core.RegisterSequence) vm.Operand {
	dest := ec.ctx.Registers.Allocate(core.Temp)
	ec.ctx.Emitter.EmitLoadConst(dest, ec.ctx.Symbols.AddConstant(name))
	ec.ctx.Symbols.BindFunction(name.String(), len(seq))

	var opcode vm.Opcode
	var protectedOpcode vm.Opcode

	switch len(seq) {
	case 0:
		opcode = vm.OpCall0
		protectedOpcode = vm.OpProtectedCall0
	case 1:
		opcode = vm.OpCall1
		protectedOpcode = vm.OpProtectedCall1
	case 2:
		opcode = vm.OpCall2
		protectedOpcode = vm.OpProtectedCall2
	case 3:
		opcode = vm.OpCall3
		protectedOpcode = vm.OpProtectedCall3
	case 4:
		opcode = vm.OpCall4
		protectedOpcode = vm.OpProtectedCall4
	default:
		opcode = vm.OpCall
		protectedOpcode = vm.OpProtectedCall
	}

	if !protected {
		ec.ctx.Emitter.EmitAs(opcode, dest, seq)
	} else {
		ec.ctx.Emitter.EmitAs(protectedOpcode, dest, seq)
	}

	return dest
}

// CompileArgumentList processes a list of function arguments from the FQL AST.
// It compiles each argument expression and moves the results into a contiguous sequence
// of registers, which is required for function calls, array literals, and object literals.
// Parameters:
//   - ctx: The argument list context from the AST
//
// Returns:
//   - A sequence of registers containing the compiled argument values,
//     or an empty sequence if there are no arguments
func (ec *ExprCompiler) CompileArgumentList(ctx fql.IArgumentListContext) core.RegisterSequence {
	var seq core.RegisterSequence

	if ctx == nil {
		return seq
	}

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
			// The reason we move is that the argument list must be a contiguous sequence of registers
			// Otherwise, we cannot compileInitialization neither a list nor an object literal with arguments
			ec.ctx.Emitter.EmitMove(seq[i], srcReg)

			// Free source register if temporary
			if srcReg.IsRegister() {
				//ec.ctx.Registers.Free(srcReg)
			}
		}
	}

	return seq
}

// CompileRangeOperator processes a range operator expression (e.g., 1..10) from the FQL AST.
// It compiles the start and end operands and emits instructions to create a range.
// Parameters:
//   - ctx: The range operator context from the AST
//
// Returns:
//   - An operand representing the compiled range
func (ec *ExprCompiler) CompileRangeOperator(ctx fql.IRangeOperatorContext) vm.Operand {
	dst := ec.ctx.Registers.Allocate(core.Temp)
	start := ec.compileRangeOperand(ctx.GetLeft())
	end := ec.compileRangeOperand(ctx.GetRight())

	ec.ctx.Emitter.EmitRange(dst, start, end)

	return dst
}

// compileRangeOperand processes a range operand (start or end value) from the FQL AST.
// It handles variables, parameters, and integer literals as valid range operands.
// Parameters:
//   - ctx: The range operand context from the AST
//
// Returns:
//   - An operand representing the compiled range operand
//
// Panics if the operand type is not recognized.
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
