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
func (c *ExprCompiler) Compile(ctx fql.IExpressionContext) vm.Operand {
	if uo := ctx.UnaryOperator(); uo != nil {
		return c.compileUnary(uo, ctx)
	}

	if and := ctx.LogicalAndOperator(); and != nil {
		return c.compileLogicalAnd(ctx)
	}

	if or := ctx.LogicalOrOperator(); or != nil {
		return c.compileLogicalOr(ctx)
	}

	if t := ctx.GetTernaryOperator(); t != nil {
		return c.compileTernary(ctx)
	}

	if p := ctx.Predicate(); p != nil {
		return c.compilePredicate(p)
	}

	return vm.NoopOperand
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
func (c *ExprCompiler) compileUnary(ctx fql.IUnaryOperatorContext, parent fql.IExpressionContext) vm.Operand {
	src := c.Compile(parent.GetRight())
	dst := c.ctx.Registers.Allocate(core.Temp)

	var op vm.Opcode

	if ctx.Not() != nil {
		op = vm.OpNot
	} else if ctx.Minus() != nil {
		op = vm.OpFlipNegative
	} else if ctx.Plus() != nil {
		op = vm.OpFlipPositive
	} else {
		return vm.NoopOperand
	}

	// We do not overwrite the source register
	c.ctx.Emitter.EmitAB(op, dst, src)

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
func (c *ExprCompiler) compileLogicalAnd(ctx fql.IExpressionContext) vm.Operand {
	left := c.Compile(ctx.GetLeft())

	skip := c.ctx.Emitter.NewLabel("and.false")
	done := c.ctx.Emitter.NewLabel("and.done")
	dst := c.ctx.Registers.Allocate(core.Temp)

	// If left is falsy, jump to skip and use left
	c.ctx.Emitter.EmitJumpIfFalse(left, skip)

	// Otherwise evaluate right and use it
	right := c.Compile(ctx.GetRight())
	c.ctx.Emitter.EmitMove(dst, right)
	c.ctx.Emitter.EmitJump(done)

	// Short-circuit: use left as result
	c.ctx.Emitter.MarkLabel(skip)
	c.ctx.Emitter.EmitMove(dst, left)

	c.ctx.Emitter.MarkLabel(done)

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
func (c *ExprCompiler) compileLogicalOr(ctx fql.IExpressionContext) vm.Operand {
	left := c.Compile(ctx.GetLeft())

	next := c.ctx.Emitter.NewLabel("or.false")
	done := c.ctx.Emitter.NewLabel("or.done")
	dst := c.ctx.Registers.Allocate(core.Temp)

	// If left is truthy, short-circuit and skip right
	c.ctx.Emitter.EmitJumpIfTrue(left, next)

	// Otherwise evaluate right
	right := c.Compile(ctx.GetRight())
	c.ctx.Emitter.EmitMove(dst, right)
	c.ctx.Emitter.EmitJump(done)

	// Short-circuit: use left value
	c.ctx.Emitter.MarkLabel(next)
	c.ctx.Emitter.EmitMove(dst, left)

	// Common exit
	c.ctx.Emitter.MarkLabel(done)

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
func (c *ExprCompiler) compileTernary(ctx fql.IExpressionContext) vm.Operand {
	dst := c.ctx.Registers.Allocate(core.Temp)

	// Compile condition and put result in dst
	condReg := c.Compile(ctx.GetCondition())
	c.ctx.Emitter.EmitMove(dst, condReg)

	// Define jump labels
	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	// endLabel to 'false' branch if condition is false
	c.ctx.Emitter.EmitJumpIfFalse(dst, elseLabel)

	// True branch
	if onTrue := ctx.GetOnTrue(); onTrue != nil {
		trueReg := c.Compile(onTrue)
		// Move result of true branch to dst
		c.ctx.Emitter.EmitMove(dst, trueReg)
	}

	// endLabel over false branch
	c.ctx.Emitter.EmitJump(endLabel)
	// Mark label for 'else' branch
	c.ctx.Emitter.MarkLabel(elseLabel)

	// False branch
	if onFalse := ctx.GetOnFalse(); onFalse != nil {
		falseReg := c.Compile(onFalse)
		// Move result of false branch to dst
		c.ctx.Emitter.EmitMove(dst, falseReg)
	}

	// endLabel
	c.ctx.Emitter.MarkLabel(endLabel)

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
func (c *ExprCompiler) compilePredicate(ctx fql.IPredicateContext) vm.Operand {
	if atom := ctx.ExpressionAtom(); atom != nil {
		startCatch := c.ctx.Emitter.Size()
		reg := c.compileAtom(atom)

		if atom.ErrorOperator() != nil {
			jump := -1
			endCatch := c.ctx.Emitter.Size()

			if fe := atom.ForExpression(); fe != nil {
				// Since FOR-IN loops depend on custom iterators,
				// We need to handle cleanup before exiting the loop.
				// TODO: Find a better way to handle this. The code assumes the knowledge of the internals of the FOR-IN loop.
				if fe.In() != nil {
					jump = endCatch - 1
				}
			}

			c.ctx.CatchTable.Push(startCatch, endCatch, jump)
		}

		return reg
	}

	var opcode vm.Opcode
	var isNegated bool
	dest := c.ctx.Registers.Allocate(core.Temp)
	left := c.compilePredicate(ctx.Predicate(0))
	right := c.compilePredicate(ctx.Predicate(1))

	if op := ctx.EqualityOperator(); op != nil {
		switch op.GetText() {
		case "==":
			opcode = vm.OpEq
		case "!=":
			opcode = vm.OpNe
		case ">":
			opcode = vm.OpGt
		case ">=":
			opcode = vm.OpGte
		case "<":
			opcode = vm.OpLt
		case "<=":
			opcode = vm.OpLte
		default:
			return vm.NoopOperand
		}
	} else if op := ctx.InOperator(); op != nil {
		opcode = vm.OpIn
		isNegated = op.Not() != nil
	} else if op := ctx.LikeOperator(); op != nil {
		opcode = vm.OpLike
		isNegated = op.Not() != nil
	} else if op := ctx.ArrayOperator(); op != nil {
		var pos int

		if op.All() != nil {
			pos = int(vm.OpAllEq)
		} else if op.Any() != nil {
			pos = int(vm.OpAnyEq)
		} else if op.None() != nil {
			pos = int(vm.OpNoneEq)
		}

		if eo := op.EqualityOperator(); eo != nil {
			switch eo.GetText() {
			case "!=":
				pos += int(vm.OpAllNe) - int(vm.OpAllEq)
			case ">":
				pos += int(vm.OpAllGt) - int(vm.OpAllEq)
			case ">=":
				pos += int(vm.OpAllGte) - int(vm.OpAllEq)
			case "<":
				pos += int(vm.OpAllLt) - int(vm.OpAllEq)
			case "<=":
				pos += int(vm.OpAllLte) - int(vm.OpAllEq)
			default:
				break
			}
		} else if inOp := op.InOperator(); inOp != nil {
			pos += int(vm.OpAllIn) - int(vm.OpAllEq)
		} else {
			return vm.NoopOperand
		}

		opcode = vm.Opcode(pos)
	}

	c.ctx.Emitter.EmitABC(opcode, dest, left, right)

	if isNegated {
		// If the operator is negated, we need to invert the result
		c.ctx.Emitter.EmitAB(vm.OpNot, dest, dest)
	}

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
func (c *ExprCompiler) compileAtom(ctx fql.IExpressionAtomContext) vm.Operand {
	var opcode vm.Opcode
	var isSet bool
	var isNegated bool
	var isRegexp bool

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
			return vm.NoopOperand
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "+":
			opcode = vm.OpAdd
		case "-":
			opcode = vm.OpSub
		default:
			return vm.NoopOperand
		}

	} else if op := ctx.RegexpOperator(); op != nil {
		isSet = true
		opcode = vm.OpRegexp
		isNegated = op.GetText() == "!~"
		isRegexp = true
	}

	if isSet {
		regLeft := c.compileAtom(ctx.ExpressionAtom(0))
		regRight := c.compileAtom(ctx.ExpressionAtom(1))
		dst := c.ctx.Registers.Allocate(core.Temp)

		c.ctx.Emitter.EmitABC(opcode, dst, regLeft, regRight)

		if isNegated {
			// If the operator is negated, we need to invert the result
			c.ctx.Emitter.EmitAB(vm.OpNot, dst, dst)
		}

		if isRegexp {
			if regRight.IsConstant() {
				val := c.ctx.Symbols.Constant(regRight)
				exp := val.String()

				// Verify that the expression is a valid regular expression
				if _, err := regexp.Compile(exp); err != nil {
					c.ctx.Errors.InvalidRegexExpression(ctx, exp)
				}
			}
		}

		return dst
	}

	if fex := ctx.FunctionCallExpression(); fex != nil {
		return c.CompileFunctionCallExpression(fex)
	} else if r := ctx.RangeOperator(); r != nil {
		return c.CompileRangeOperator(r)
	} else if l := ctx.Literal(); l != nil {
		return c.ctx.LiteralCompiler.Compile(l)
	} else if v := ctx.Variable(); v != nil {
		return c.CompileVariable(v)
	} else if me := ctx.MemberExpression(); me != nil {
		return c.CompileMemberExpression(me)
	} else if p := ctx.Param(); p != nil {
		return c.CompileParam(p)
	} else if fe := ctx.ForExpression(); fe != nil {
		return c.ctx.LoopCompiler.Compile(fe)
	} else if wfe := ctx.WaitForExpression(); wfe != nil {
		return c.ctx.WaitCompiler.Compile(wfe)
	} else if e := ctx.Expression(); e != nil {
		return c.Compile(e)
	}

	//c.ctx.Errors.UnexpectedToken(ctx)

	return vm.NoopOperand
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
func (c *ExprCompiler) CompileMemberExpression(ctx fql.IMemberExpressionContext) vm.Operand {
	mes := ctx.MemberExpressionSource()
	segments := ctx.AllMemberExpressionPath()

	var src1 vm.Operand

	if v := mes.Variable(); v != nil {
		src1 = c.CompileVariable(v)
	} else if p := mes.Param(); p != nil {
		src1 = c.CompileParam(p)
	} else if ol := mes.ObjectLiteral(); ol != nil {
		src1 = c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	} else if al := mes.ArrayLiteral(); al != nil {
		src1 = c.ctx.LiteralCompiler.CompileArrayLiteral(al)
	} else if fc := mes.FunctionCall(); fc != nil {
		// FOO()?.bar
		segment := segments[0]
		src1 = c.CompileFunctionCall(fc, segment.ErrorOperator() != nil)
	}

	var dst vm.Operand

	for _, segment := range segments {
		var src2 vm.Operand
		p := segment.(*fql.MemberExpressionPathContext)

		if pn := p.PropertyName(); pn != nil {
			src2 = c.ctx.LiteralCompiler.CompilePropertyName(pn)
		} else if cpn := p.ComputedPropertyName(); cpn != nil {
			src2 = c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn)
		}

		dst = c.ctx.Registers.Allocate(core.Temp)

		// TODO: Replace with EmitLoadKey
		if p.ErrorOperator() != nil {
			c.ctx.Emitter.EmitLoadPropertyOptional(dst, src1, src2)
		} else {
			c.ctx.Emitter.EmitLoadProperty(dst, src1, src2)
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
func (c *ExprCompiler) CompileVariable(ctx fql.IVariableContext) vm.Operand {
	// Check if the context is valid (in case of parser errors)
	if ctx.Identifier() == nil {
		return vm.NoopOperand
	}

	name := ctx.Identifier().GetText()
	// Just return the register / constant index
	op, _, found := c.ctx.Symbols.Resolve(name)

	if !found {
		c.ctx.Errors.VariableNotFound(ctx.Identifier().GetSymbol(), name)

		return vm.NoopOperand
	}

	if op.IsRegister() {
		return op
	}

	reg := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitMove(reg, op)

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
func (c *ExprCompiler) CompileParam(ctx fql.IParamContext) vm.Operand {
	name := ctx.Identifier().GetText()
	reg := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitLoadParam(reg, c.ctx.Symbols.BindParam(name))

	return reg
}

// CompileFunctionCallExpression processes a function call expression from the FQL AST.
// It handles both regular function calls and protected function calls (with TRY).
// Parameters:
//   - ctx: The function call expression context from the AST
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) vm.Operand {
	protected := ctx.ErrorOperator() != nil
	call := ctx.FunctionCall()

	return c.CompileFunctionCall(call, protected)
}

// CompileFunctionCall processes a function call from the FQL AST.
// It compiles the function arguments and delegates to CompileFunctionCallWith.
// Parameters:
//   - ctx: The function call context from the AST
//   - protected: Whether this is a protected call (with TRY)
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCall(ctx fql.IFunctionCallContext, protected bool) vm.Operand {
	return c.CompileFunctionCallWith(ctx, protected, c.CompileArgumentList(ctx.ArgumentList()))
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
func (c *ExprCompiler) CompileFunctionCallWith(ctx fql.IFunctionCallContext, protected bool, seq core.RegisterSequence) vm.Operand {
	name := getFunctionName(ctx)

	return c.CompileFunctionCallByNameWith(name, protected, seq)
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
func (c *ExprCompiler) CompileFunctionCallByNameWith(name runtime.String, protected bool, seq core.RegisterSequence) vm.Operand {
	switch name {
	case runtimeLength:
		dst := c.ctx.Registers.Allocate(core.Temp)

		if seq == nil || len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeLength+": expected 1 argument"))
		}

		c.ctx.Emitter.EmitAB(vm.OpLength, dst, seq[0])

		return dst
	case runtimeTypename:
		dst := c.ctx.Registers.Allocate(core.Temp)

		if seq == nil || len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeTypename+": expected 1 argument"))
		}

		c.ctx.Emitter.EmitAB(vm.OpType, dst, seq[0])

		return dst
	case runtimeWait:
		if len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeWait+": expected 1 argument"))
		}

		c.ctx.Emitter.EmitA(vm.OpSleep, seq[0])

		return seq[0]
	default:
		return c.compileUserFunctionCallWith(name, protected, seq)
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
func (c *ExprCompiler) compileUserFunctionCallWith(name runtime.String, protected bool, seq core.RegisterSequence) vm.Operand {
	dest := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(name))
	c.ctx.Symbols.BindFunction(name.String(), len(seq))

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
		c.ctx.Emitter.EmitAs(opcode, dest, seq)
	} else {
		c.ctx.Emitter.EmitAs(protectedOpcode, dest, seq)
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
func (c *ExprCompiler) CompileArgumentList(ctx fql.IArgumentListContext) core.RegisterSequence {
	var seq core.RegisterSequence

	if ctx == nil {
		return seq
	}

	// Get all array element expressions
	exps := ctx.AllExpression()
	size := len(exps)

	if size > 0 {
		// Allocate seq for function arguments
		seq = c.ctx.Registers.AllocateSequence(size)

		// Evaluate each element into seq Registers
		for i, exp := range exps {
			// Compile expression and move to seq register
			srcReg := c.Compile(exp)

			// TODO: Figure out how to remove OpMove and use Registers returned from each expression
			// The reason we move is that the argument list must be a contiguous sequence of registers
			// Otherwise, we cannot compileInitialization neither a list nor an object literal with arguments
			c.ctx.Emitter.EmitMove(seq[i], srcReg)

			// Free source register if temporary
			//if srcReg.IsRegister() {
			//	c.ctx.Registers.Free(srcReg)
			//}
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
func (c *ExprCompiler) CompileRangeOperator(ctx fql.IRangeOperatorContext) vm.Operand {
	dst := c.ctx.Registers.Allocate(core.Temp)
	start := c.compileRangeOperand(ctx.GetLeft())
	end := c.compileRangeOperand(ctx.GetRight())

	c.ctx.Emitter.EmitRange(dst, start, end)

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
func (c *ExprCompiler) compileRangeOperand(ctx fql.IRangeOperandContext) vm.Operand {
	if v := ctx.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.CompileParam(p)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	return vm.NoopOperand
}
