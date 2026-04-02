package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Runtime functions
const (
	runtimeTypename = "TYPENAME"
	runtimeLength   = "LENGTH"
	runtimeWait     = "WAIT"
)

type (
	// ExprCompiler handles the compilation of expressions in FQL queries.
	// It transforms expression operations from the AST into VM instructions.
	ExprCompiler struct {
		ctx                  *CompilerContext
		implicitCurrentDepth int
	}

	// queryModifier represents modifiers that can be applied to queries for altering their behavior or result interpretation.
	queryModifier string

	// atomBinaryOperator represents binary operators used in FQL expressions.
	atomBinaryOperator struct {
		opcode  bytecode.Opcode
		negated bool
		regexp  bool
	}

	matchResultGroup struct {
		result fql.IExpressionContext
		arms   []int
		label  core.Label
	}

	optionalMemberChainState struct {
		endLabel  core.Label
		stickyDst bool
		hasJump   bool
	}
)

const (
	queryModifierUnknown queryModifier = ""
	queryModifierExists  queryModifier = "exists"
	queryModifierCount   queryModifier = "count"
	queryModifierAny     queryModifier = "any"
	queryModifierValue   queryModifier = "value"
	queryModifierOne     queryModifier = "one"
)

const (
	queryValueFailMessage = "QUERY VALUE expected at least one match"
	queryOneFailMessage   = "QUERY ONE expected exactly one match"
)

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
func (c *ExprCompiler) Compile(ctx fql.IExpressionContext) bytecode.Operand {
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

	return bytecode.NoopOperand
}

func (c *ExprCompiler) CompileIncDec(token antlr.Token, target bytecode.Operand) bytecode.Operand {
	if target.IsConstant() {
		panic("cannot increment/decrement a constant")
	}

	operator := token.GetText()

	switch operator {
	case "++":
		c.ctx.Emitter.EmitA(bytecode.OpIncr, target)
	case "--":
		c.ctx.Emitter.EmitA(bytecode.OpDecr, target)
	default:
		c.ctx.Errors.InvalidToken(token)

		return bytecode.NoopOperand
	}

	return target
}

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
func (c *ExprCompiler) compileUnary(ctx fql.IUnaryOperatorContext, parent fql.IExpressionContext) bytecode.Operand {
	src := c.Compile(parent.GetRight())
	dst := c.ctx.Registers.Allocate()

	var op bytecode.Opcode

	if ctx.Not() != nil {
		op = bytecode.OpNot
	} else if ctx.Minus() != nil {
		op = bytecode.OpFlipNegative
	} else if ctx.Plus() != nil {
		op = bytecode.OpFlipPositive
	} else {
		return bytecode.NoopOperand
	}

	// We do not overwrite the source register
	c.ctx.Emitter.EmitAB(op, dst, src)

	return dst
}

// compileLogicalAnd processes a logical AND operation from the FQL AST.
// It implements short-circuit evaluation: if the left operand is falsy, it returns that value
// without evaluating the right operand. Otherwise, it evaluates and returns the right operand.
// Parameters:
//   - ctx: The expression context from the AST containing both operands
//
// Returns:
//   - An operand representing the result of the logical AND operation
func (c *ExprCompiler) compileLogicalAnd(ctx fql.IExpressionContext) bytecode.Operand {
	left := c.Compile(ctx.GetLeft())

	skip := c.ctx.Emitter.NewLabel("and.false")
	done := c.ctx.Emitter.NewLabel("and.done")
	dst := c.ctx.Registers.Allocate()

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

// compileLogicalOr processes a logical OR operation from the FQL AST.
// It implements short-circuit evaluation: if the left operand is truthy, it returns that value
// without evaluating the right operand. Otherwise, it evaluates and returns the right operand.
// Parameters:
//   - ctx: The expression context from the AST containing both operands
//
// Returns:
//   - An operand representing the result of the logical OR operation
func (c *ExprCompiler) compileLogicalOr(ctx fql.IExpressionContext) bytecode.Operand {
	left := c.Compile(ctx.GetLeft())

	next := c.ctx.Emitter.NewLabel("or.false")
	done := c.ctx.Emitter.NewLabel("or.done")
	dst := c.ctx.Registers.Allocate()

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

// compileTernary processes a ternary conditional operation (condition ? trueExpr : falseExpr) from the FQL AST.
// It evaluates the condition and then either the true expression or the false expression based on the result.
// Parameters:
//   - ctx: The expression context from the AST containing the condition, true expression, and false expression
//
// Returns:
//   - An operand representing the result of either the true or false expression
func (c *ExprCompiler) compileTernary(ctx fql.IExpressionContext) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()

	// Define jump labels
	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	onTrue := ctx.GetOnTrue()
	onFalse := ctx.GetOnFalse()
	cond := ctx.GetCondition()

	// If true branch is omitted, preserve condition value.
	if onTrue == nil && cond != nil {
		condReg := c.Compile(cond)
		c.ctx.Emitter.EmitMove(dst, condReg)
		c.ctx.Emitter.EmitJumpIfFalse(condReg, elseLabel)
	} else if cond != nil {
		// endLabel to 'false' branch if condition is false
		c.emitConditionJump(cond, elseLabel, false)
	}

	// True branch
	if onTrue != nil {
		trueReg := c.Compile(onTrue)
		// Move result of true branch to dst
		c.ctx.Emitter.EmitMove(dst, trueReg)
	}

	// endLabel over false branch
	c.ctx.Emitter.EmitJump(endLabel)
	// Mark label for 'else' branch
	c.ctx.Emitter.MarkLabel(elseLabel)

	// False branch
	if onFalse != nil {
		falseReg := c.Compile(onFalse)
		// Move result of false branch to dst
		c.ctx.Emitter.EmitMove(dst, falseReg)
	}

	// endLabel
	c.ctx.Emitter.MarkLabel(endLabel)

	return dst
}

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
func (c *ExprCompiler) compilePredicate(ctx fql.IPredicateContext) bytecode.Operand {
	if reg, ok := c.compilePredicateAtom(ctx); ok {
		return reg
	}

	opcode, isNegated, ok := c.resolvePredicateOperator(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	left := c.compilePredicate(ctx.Predicate(0))
	right := c.compilePredicate(ctx.Predicate(1))

	return c.emitBinaryPredicate(ctx, opcode, left, right, isNegated)
}

func (c *ExprCompiler) compileLiteralOperand(lit fql.ILiteralContext) bytecode.Operand {
	return compileScalarLiteralOperand(c.ctx, lit)
}

func (c *ExprCompiler) compilePredicateAtom(ctx fql.IPredicateContext) (bytecode.Operand, bool) {
	if ctx == nil {
		return bytecode.NoopOperand, true
	}

	atom := ctx.ExpressionAtom()
	if atom == nil {
		return bytecode.NoopOperand, false
	}

	if atom.ErrorOperator() != nil {
		jumpMode := core.CatchJumpModeNone
		if fe := atom.ForExpression(); fe != nil {
			jumpMode = catchJumpModeForForExpression(fe)
		} else if wfe := atom.WaitForExpression(); wfe != nil {
			jumpMode = catchJumpModeForWaitForExpression(wfe)
		}

		reg := c.ctx.OPCompiler.CompileWithErrorPolicy(core.ErrorPolicySuppress, jumpMode, func() bytecode.Operand {
			return c.compileAtom(atom)
		})

		return reg, true
	}

	if wfe := atom.WaitForExpression(); wfe != nil {
		plan := collectRecoveryPlan(c.ctx, atom, core.RecoveryPlanOptions{
			AllowTimeout: true,
			HasTimeout:   waitForHasExplicitTimeoutClause(wfe),
		})
		if plan.OnError == nil && plan.OnTimeout == nil {
			return c.compileAtom(atom), true
		}

		return c.ctx.WaitCompiler.compileWithOuterRecovery(wfe, plan), true
	}

	if fe := atom.ForExpression(); fe != nil {
		plan := collectRecoveryPlan(c.ctx, atom, core.RecoveryPlanOptions{})
		if plan.OnError == nil {
			return c.compileAtom(atom), true
		}

		return c.ctx.LoopCompiler.compileWithRecovery(fe, plan), true
	}

	plan := collectRecoveryPlan(c.ctx, atom, core.RecoveryPlanOptions{})
	if plan.OnError == nil {
		return c.compileAtom(atom), true
	}

	jumpMode := core.CatchJumpModeNone
	if fe := atom.ForExpression(); fe != nil {
		jumpMode = catchJumpModeForForExpression(fe)
	} else if wfe := atom.WaitForExpression(); wfe != nil {
		jumpMode = catchJumpModeForWaitForExpression(wfe)
	}

	if !hasErrorReturnNoneHandler(plan) {
		jumpMode = core.CatchJumpModeNone
	}

	reg := c.ctx.OPCompiler.CompileWithRecoveryPlan(plan, jumpMode, func() bytecode.Operand {
		return c.compileAtom(atom)
	})

	return reg, true
}

func (c *ExprCompiler) resolvePredicateOperator(ctx fql.IPredicateContext) (bytecode.Opcode, bool, bool) {
	if ctx == nil {
		return bytecode.Opcode(0), false, false
	}

	if op := ctx.EqualityOperator(); op != nil {
		switch op.GetText() {
		case "==":
			return bytecode.OpEq, false, true
		case "!=":
			return bytecode.OpNe, false, true
		case ">":
			return bytecode.OpGt, false, true
		case ">=":
			return bytecode.OpGte, false, true
		case "<":
			return bytecode.OpLt, false, true
		case "<=":
			return bytecode.OpLte, false, true
		default:
			return bytecode.Opcode(0), false, false
		}
	}

	if op := ctx.InOperator(); op != nil {
		return bytecode.OpIn, op.Not() != nil, true
	}

	if op := ctx.LikeOperator(); op != nil {
		return bytecode.OpLike, op.Not() != nil, true
	}

	if op := ctx.ArrayOperator(); op != nil {
		opcode, ok := resolveArrayPredicateOpcode(op)
		return opcode, false, ok
	}

	return bytecode.Opcode(0), false, false
}

func resolveArrayPredicateOpcode(op fql.IArrayOperatorContext) (bytecode.Opcode, bool) {
	if op == nil {
		return bytecode.Opcode(0), false
	}

	var pos int

	switch {
	case op.All() != nil:
		pos = int(bytecode.OpAllEq)
	case op.Any() != nil:
		pos = int(bytecode.OpAnyEq)
	case op.None() != nil:
		pos = int(bytecode.OpNoneEq)
	}

	if eo := op.EqualityOperator(); eo != nil {
		switch eo.GetText() {
		case "!=":
			pos += int(bytecode.OpAllNe) - int(bytecode.OpAllEq)
		case ">":
			pos += int(bytecode.OpAllGt) - int(bytecode.OpAllEq)
		case ">=":
			pos += int(bytecode.OpAllGte) - int(bytecode.OpAllEq)
		case "<":
			pos += int(bytecode.OpAllLt) - int(bytecode.OpAllEq)
		case "<=":
			pos += int(bytecode.OpAllLte) - int(bytecode.OpAllEq)
		default:
			// Keep OpAllEq/OpAnyEq/OpNoneEq fallback for parser-provided defaults.
		}

		return bytecode.Opcode(pos), true
	}

	if op.InOperator() != nil {
		pos += int(bytecode.OpAllIn) - int(bytecode.OpAllEq)

		return bytecode.Opcode(pos), true
	}

	return bytecode.Opcode(0), false
}

func (c *ExprCompiler) emitBinaryPredicate(ctx fql.IPredicateContext, opcode bytecode.Opcode, left, right bytecode.Operand, isNegated bool) bytecode.Operand {
	dest := c.ctx.Registers.Allocate()
	span := source.Span{Start: -1, End: -1}

	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = diagnostics.SpanFromRuleContext(prc)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABC(opcode, dest, left, right)

		if isNegated {
			// If the operator is negated, we need to invert the result
			c.ctx.Emitter.EmitAB(bytecode.OpNot, dest, dest)
		}
	})

	return dest
}

func (c *ExprCompiler) emitPredicateJump(ctx fql.IPredicateContext, label core.Label, jumpOnTrue bool) bool {
	opText, leftCtx, rightCtx, ok := resolvePredicateEqNeJump(ctx)
	if !ok {
		return false
	}

	rightLiteral := c.compilePredicateLiteralOperand(rightCtx)
	if rightLiteral != bytecode.NoopOperand {
		return c.emitPredicateJumpWithLiteralOperand(opText, jumpOnTrue, leftCtx, rightCtx, rightLiteral, true, label)
	}

	leftLiteral := c.compilePredicateLiteralOperand(leftCtx)
	if leftLiteral != bytecode.NoopOperand {
		return c.emitPredicateJumpWithLiteralOperand(opText, jumpOnTrue, leftCtx, rightCtx, leftLiteral, false, label)
	}

	leftOp := c.ensureRegister(c.compilePredicate(leftCtx))
	rightOp := c.ensureRegister(c.compilePredicate(rightCtx))
	c.emitPredicateJumpCompare(opText, jumpOnTrue, leftOp, rightOp, label, false)

	return true
}

func resolvePredicateEqNeJump(ctx fql.IPredicateContext) (string, fql.IPredicateContext, fql.IPredicateContext, bool) {
	if ctx == nil {
		return "", nil, nil, false
	}

	op := ctx.EqualityOperator()
	if op == nil {
		return "", nil, nil, false
	}

	opText := op.GetText()
	if opText != "==" && opText != "!=" {
		return "", nil, nil, false
	}

	leftCtx := ctx.Predicate(0)
	rightCtx := ctx.Predicate(1)
	if leftCtx == nil || rightCtx == nil {
		return "", nil, nil, false
	}

	return opText, leftCtx, rightCtx, true
}

func (c *ExprCompiler) compilePredicateLiteralOperand(ctx fql.IPredicateContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	atom := ctx.ExpressionAtom()
	if atom == nil {
		return bytecode.NoopOperand
	}

	literal := atom.Literal()
	if literal == nil {
		return bytecode.NoopOperand
	}

	return c.compileLiteralOperand(literal)
}

func (c *ExprCompiler) emitPredicateJumpWithLiteralOperand(opText string, jumpOnTrue bool, leftCtx, rightCtx fql.IPredicateContext, literalOp bytecode.Operand, literalOnRight bool, label core.Label) bool {
	if literalOp.IsConstant() {
		exprCtx := leftCtx
		if !literalOnRight {
			exprCtx = rightCtx
		}

		exprOp := c.ensureRegister(c.compilePredicate(exprCtx))
		c.emitPredicateJumpCompare(opText, jumpOnTrue, exprOp, literalOp, label, true)

		return true
	}

	if literalOnRight {
		leftOp := c.ensureRegister(c.compilePredicate(leftCtx))
		rightOp := c.ensureRegister(literalOp)
		c.emitPredicateJumpCompare(opText, jumpOnTrue, leftOp, rightOp, label, false)

		return true
	}

	leftOp := c.ensureRegister(literalOp)
	rightOp := c.ensureRegister(c.compilePredicate(rightCtx))
	c.emitPredicateJumpCompare(opText, jumpOnTrue, leftOp, rightOp, label, false)

	return true
}

func (c *ExprCompiler) emitPredicateJumpCompare(opText string, jumpOnTrue bool, left, right bytecode.Operand, label core.Label, constOperand bool) {
	opcode := resolveEqNeJumpOpcode(opText, jumpOnTrue, constOperand)
	c.ctx.Emitter.EmitJumpCompare(opcode, left, right, label)
}

func resolveEqNeJumpOpcode(opText string, jumpOnTrue, constOperand bool) bytecode.Opcode {
	if constOperand {
		if opText == "==" {
			if jumpOnTrue {
				return bytecode.OpJumpIfEqConst
			}

			return bytecode.OpJumpIfNeConst
		}

		if jumpOnTrue {
			return bytecode.OpJumpIfNeConst
		}

		return bytecode.OpJumpIfEqConst
	}

	if opText == "==" {
		if jumpOnTrue {
			return bytecode.OpJumpIfEq
		}

		return bytecode.OpJumpIfNe
	}

	if jumpOnTrue {
		return bytecode.OpJumpIfNe
	}

	return bytecode.OpJumpIfEq
}

func (c *ExprCompiler) emitConditionJump(expr fql.IExpressionContext, label core.Label, jumpOnTrue bool) {
	if expr == nil {
		return
	}

	if pred := expr.Predicate(); pred != nil {
		if c.emitPredicateJump(pred, label, jumpOnTrue) {
			return
		}
	}

	cond := c.Compile(expr)
	if jumpOnTrue {
		c.ctx.Emitter.EmitJumpIfTrue(cond, label)
	} else {
		c.ctx.Emitter.EmitJumpIfFalse(cond, label)
	}
}

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
func (c *ExprCompiler) compileAtom(ctx fql.IExpressionAtomContext) bytecode.Operand {
	if op, ok := resolveAtomBinaryOperator(ctx); ok {
		return c.compileBinaryAtom(ctx, op)
	}

	return c.compileLeafAtom(ctx)
}

func resolveAtomBinaryOperator(ctx fql.IExpressionAtomContext) (atomBinaryOperator, bool) {
	if op := ctx.MultiplicativeOperator(); op != nil {
		return resolveArithmeticBinaryOperator(op.GetText())
	}

	if op := ctx.AdditiveOperator(); op != nil {
		return resolveArithmeticBinaryOperator(op.GetText())
	}

	if op := ctx.RegexpOperator(); op != nil {
		return atomBinaryOperator{
			opcode:  bytecode.OpRegexp,
			negated: op.GetText() == "!~",
			regexp:  true,
		}, true
	}

	return atomBinaryOperator{}, false
}

func resolveArithmeticBinaryOperator(operator string) (atomBinaryOperator, bool) {
	switch operator {
	case "+", "+=":
		return atomBinaryOperator{opcode: bytecode.OpAdd}, true
	case "-", "-=":
		return atomBinaryOperator{opcode: bytecode.OpSub}, true
	case "*", "*=":
		return atomBinaryOperator{opcode: bytecode.OpMul}, true
	case "/", "/=":
		return atomBinaryOperator{opcode: bytecode.OpDiv}, true
	case "%":
		return atomBinaryOperator{opcode: bytecode.OpMod}, true
	default:
		return atomBinaryOperator{}, false
	}
}

func (c *ExprCompiler) compileBinaryAtom(ctx fql.IExpressionAtomContext, op atomBinaryOperator) bytecode.Operand {
	if op.opcode == bytecode.OpAdd {
		if parts, ok := buildConcatOperandSegmentsFromAtom(c, ctx); ok {
			return emitConcatOperandSegments(c.ctx, parts)
		}
	}

	left := c.compileAtom(ctx.ExpressionAtom(0))
	right := c.compileAtom(ctx.ExpressionAtom(1))
	dst := c.emitBinaryAtomOperation(ctx, op, left, right)

	if op.regexp {
		c.validateRegexpOperand(ctx)
	}

	return dst
}

func (c *ExprCompiler) emitBinaryAtomOperation(ctx fql.IExpressionAtomContext, op atomBinaryOperator, left, right bytecode.Operand) bytecode.Operand {
	prc, _ := ctx.(antlr.ParserRuleContext)

	return emitBinaryOperation(c.ctx, prc, op, left, right)
}

func emitBinaryOperation(ctx *CompilerContext, prc antlr.ParserRuleContext, op atomBinaryOperator, left, right bytecode.Operand) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	dst := ctx.Registers.Allocate()
	span := source.Span{Start: -1, End: -1}

	if prc != nil {
		span = diagnostics.SpanFromRuleContext(prc)
	}

	ctx.Emitter.WithSpan(span, func() {
		ctx.Emitter.EmitABC(op.opcode, dst, left, right)

		if op.negated {
			ctx.Emitter.EmitAB(bytecode.OpNot, dst, dst)
		}
	})

	resultType := inferBinaryResultType(ctx, op, left, right)
	if op.negated {
		resultType = core.TypeBool
	}
	if resultType != core.TypeUnknown {
		ctx.Types.Set(dst, resultType)
	}

	return dst
}

func (c *ExprCompiler) validateRegexpOperand(ctx fql.IExpressionAtomContext) {
	right := ctx.ExpressionAtom(1)
	if right == nil {
		return
	}

	lit := right.Literal()
	if lit == nil {
		return
	}

	if str := lit.StringLiteral(); str != nil {
		exp, ok := parseStringLiteralConst(str)
		if !ok {
			return
		}

		if _, err := regexp.Compile(exp.String()); err != nil {
			c.ctx.Errors.InvalidRegexExpression(ctx, exp.String())
		}

		return
	}

	c.ctx.Errors.InvalidRegexExpression(ctx, lit.GetText())
}

func (c *ExprCompiler) compileLeafAtom(ctx fql.IExpressionAtomContext) bytecode.Operand {
	if fex := ctx.FunctionCallExpression(); fex != nil {
		return c.CompileFunctionCallExpression(fex)
	}

	if mx := ctx.MatchExpression(); mx != nil {
		return c.compileMatchExpression(mx)
	}

	if qx := ctx.QueryExpression(); qx != nil {
		return c.compileQueryExpression(qx)
	}

	if r := ctx.RangeOperator(); r != nil {
		return c.CompileRangeOperator(r)
	}

	if l := ctx.Literal(); l != nil {
		return c.ctx.LiteralCompiler.Compile(l)
	}

	if v := ctx.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if ice := ctx.ImplicitCurrentExpression(); ice != nil {
		return c.CompileImplicitCurrentExpression(ice)
	}

	if ime := ctx.ImplicitMemberExpression(); ime != nil {
		return c.CompileImplicitMemberExpression(ime)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.CompileMemberExpression(me)
	}

	if p := ctx.Param(); p != nil {
		return c.CompileParam(p)
	}

	if de := ctx.DispatchExpression(); de != nil {
		return c.ctx.DispatchCompiler.Compile(de)
	}

	if fe := ctx.ForExpression(); fe != nil {
		return c.ctx.LoopCompiler.Compile(fe)
	}

	if wfe := ctx.WaitForExpression(); wfe != nil {
		return c.ctx.WaitCompiler.Compile(wfe)
	}

	if e := ctx.Expression(); e != nil {
		return c.Compile(e)
	}

	return bytecode.NoopOperand
}

func (c *ExprCompiler) compileMatchExpression(ctx fql.IMatchExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	dst := c.ctx.Registers.Allocate()
	end := c.ctx.Emitter.NewLabel("match.end")

	if arms := ctx.MatchPatternArms(); arms != nil {
		scrutinee := ctx.Expression()
		if scrutinee == nil {
			return bytecode.NoopOperand
		}

		if c.tryCompileMatchConstantFold(scrutinee, arms, dst) {
			c.ctx.Types.Set(dst, core.TypeAny)
			return dst
		}

		scrReg := c.ensureRegister(c.Compile(scrutinee))
		c.compileMatchPatternArms(scrReg, arms, dst, end)
	} else if guards := ctx.MatchGuardArms(); guards != nil {
		c.compileMatchGuardArms(guards, dst, end)
	}

	c.ctx.Emitter.MarkLabel(end)

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeAny)
	}

	return dst
}

func (c *ExprCompiler) compileMatchPatternArms(scrReg bytecode.Operand, ctx fql.IMatchPatternArmsContext, dst bytecode.Operand, end core.Label) {
	if ctx == nil {
		return
	}

	arms := collectMatchPatternArms(ctx)
	mergeLabels, mergeGroups := collectMatchResultMerges(c, arms)
	defaultLabel, hasDefaultLabel := c.matchMergeDefaultLabel(mergeGroups)

	for idx, arm := range arms {
		c.compileMatchPatternArm(scrReg, arm, idx, mergeLabels, dst, end)
	}

	if hasDefaultLabel {
		c.compileMatchMergedResults(mergeGroups, defaultLabel, dst, end)
	}

	c.compileMatchPatternDefaultArm(ctx.MatchDefaultArm(), dst)
}

func collectMatchPatternArms(ctx fql.IMatchPatternArmsContext) []fql.IMatchPatternArmContext {
	if ctx == nil {
		return nil
	}

	list := ctx.MatchPatternArmList()
	if list == nil {
		return nil
	}

	return list.AllMatchPatternArm()
}

func (c *ExprCompiler) matchMergeDefaultLabel(groups []matchResultGroup) (core.Label, bool) {
	if len(groups) == 0 {
		return core.Label{}, false
	}

	return c.ctx.Emitter.NewLabel("match.default"), true
}

func (c *ExprCompiler) compileMatchPatternArm(scrReg bytecode.Operand, arm fql.IMatchPatternArmContext, idx int, mergeLabels map[int]core.Label, dst bytecode.Operand, end core.Label) {
	if arm == nil {
		return
	}

	next := c.ctx.Emitter.NewLabel("match.next")
	c.ctx.Symbols.EnterScope()
	c.compileMatchPatternArmConditions(scrReg, arm, next)
	c.compileMatchPatternArmResult(arm, idx, mergeLabels, dst, end)
	c.ctx.Symbols.ExitScope()
	c.ctx.Emitter.MarkLabel(next)
}

func (c *ExprCompiler) compileMatchPatternArmConditions(scrReg bytecode.Operand, arm fql.IMatchPatternArmContext, next core.Label) {
	if pattern := arm.MatchPattern(); pattern != nil {
		c.compileMatchPatternValue(scrReg, pattern, next)
	}

	guard := arm.MatchPatternGuard()
	if guard == nil {
		return
	}

	expr := guard.Expression()
	if expr == nil {
		return
	}

	c.emitConditionJump(expr, next, false)
}

func (c *ExprCompiler) compileMatchPatternArmResult(arm fql.IMatchPatternArmContext, idx int, mergeLabels map[int]core.Label, dst bytecode.Operand, end core.Label) {
	result := arm.Expression()
	if result == nil {
		c.ctx.Emitter.EmitJump(end)
		return
	}

	if label, ok := mergeLabels[idx]; ok {
		c.ctx.Emitter.EmitJump(label)
		return
	}

	out := c.ensureRegister(c.Compile(result))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Emitter.EmitMove(dst, out)
	}

	c.ctx.Emitter.EmitJump(end)
}

func (c *ExprCompiler) compileMatchMergedResults(groups []matchResultGroup, defaultLabel core.Label, dst bytecode.Operand, end core.Label) {
	c.ctx.Emitter.EmitJump(defaultLabel)

	for _, group := range groups {
		c.ctx.Emitter.MarkLabel(group.label)
		c.ctx.Symbols.EnterScope()
		out := c.ensureRegister(c.Compile(group.result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Emitter.EmitMove(dst, out)
		}
		c.ctx.Symbols.ExitScope()
		c.ctx.Emitter.EmitJump(end)
	}

	c.ctx.Emitter.MarkLabel(defaultLabel)
}

func (c *ExprCompiler) compileMatchPatternDefaultArm(def fql.IMatchDefaultArmContext, dst bytecode.Operand) {
	if def == nil {
		return
	}

	c.ctx.Symbols.EnterScope()
	result := def.Expression()
	if result != nil {
		out := c.ensureRegister(c.Compile(result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Emitter.EmitMove(dst, out)
		}
	}
	c.ctx.Symbols.ExitScope()
}

func (c *ExprCompiler) compileMatchGuardArms(ctx fql.IMatchGuardArmsContext, dst bytecode.Operand, end core.Label) {
	if ctx == nil {
		return
	}

	var arms []fql.IMatchGuardArmContext
	if list := ctx.MatchGuardArmList(); list != nil {
		arms = list.AllMatchGuardArm()
	}

	for _, arm := range arms {
		if arm == nil {
			continue
		}

		next := c.ctx.Emitter.NewLabel("match.next")
		c.ctx.Symbols.EnterScope()

		exprs := arm.AllExpression()
		if len(exprs) > 0 {
			c.emitConditionJump(exprs[0], next, false)
		}

		if len(exprs) > 1 {
			out := c.ensureRegister(c.Compile(exprs[1]))
			if out != bytecode.NoopOperand && out != dst {
				c.ctx.Emitter.EmitMove(dst, out)
			}
		}

		c.ctx.Emitter.EmitJump(end)
		c.ctx.Symbols.ExitScope()
		c.ctx.Emitter.MarkLabel(next)
	}

	if def := ctx.MatchDefaultArm(); def != nil {
		c.ctx.Symbols.EnterScope()
		if result := def.Expression(); result != nil {
			out := c.ensureRegister(c.Compile(result))
			if out != bytecode.NoopOperand && out != dst {
				c.ctx.Emitter.EmitMove(dst, out)
			}
		}
		c.ctx.Symbols.ExitScope()
	}
}

func collectMatchResultMerges(c *ExprCompiler, arms []fql.IMatchPatternArmContext) (map[int]core.Label, []matchResultGroup) {
	if len(arms) == 0 {
		return nil, nil
	}

	groups := make(map[string]*matchResultGroup)
	order := make([]*matchResultGroup, 0)

	for idx, arm := range arms {
		if arm == nil {
			continue
		}
		if arm.MatchPatternGuard() != nil {
			continue
		}
		pattern := arm.MatchPattern()
		if pattern == nil || pattern.MatchLiteralPattern() == nil {
			continue
		}
		result := arm.Expression()
		if result == nil {
			continue
		}
		key, ok := matchPureResultKey(result)
		if !ok {
			continue
		}

		group, exists := groups[key]
		if !exists {
			group = &matchResultGroup{result: result}
			groups[key] = group
			order = append(order, group)
		}
		group.arms = append(group.arms, idx)
	}

	if len(order) == 0 {
		return nil, nil
	}

	labels := make(map[int]core.Label)
	merged := make([]matchResultGroup, 0)

	for _, group := range order {
		if len(group.arms) < 2 {
			continue
		}
		label := c.ctx.Emitter.NewLabel("match.result")
		group.label = label
		for _, idx := range group.arms {
			labels[idx] = label
		}
		merged = append(merged, *group)
	}

	if len(merged) == 0 {
		return nil, nil
	}

	return labels, merged
}

func (c *ExprCompiler) compileMatchPatternValue(valueReg bytecode.Operand, ctx fql.IMatchPatternContext, onFail core.Label) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.MatchLiteralPattern() != nil:
		litOp := c.compileMatchLiteralOperand(ctx.MatchLiteralPattern())
		if litOp == bytecode.NoopOperand {
			return
		}
		if litOp.IsConstant() {
			c.ctx.Emitter.EmitJumpCompare(bytecode.OpJumpIfNeConst, valueReg, litOp, onFail)
		} else {
			c.ctx.Emitter.EmitJumpCompare(bytecode.OpJumpIfNe, valueReg, litOp, onFail)
		}
	case ctx.MatchBindingPattern() != nil:
		binding := ctx.MatchBindingPattern()
		if binding == nil {
			return
		}
		var name string
		if id := binding.Identifier(); id != nil {
			name = id.GetText()
		} else if srw := binding.SafeReservedWord(); srw != nil {
			name = srw.GetText()
		}
		if name != "" {
			c.declareMatchBinding(binding, name, valueReg)
		}
	case ctx.MatchObjectPattern() != nil:
		c.compileMatchObjectPattern(valueReg, ctx.MatchObjectPattern(), onFail)
	}
}

func (c *ExprCompiler) compileMatchLiteralOperand(ctx fql.IMatchLiteralPatternContext) bytecode.Operand {
	return compileScalarLiteralOperand(c.ctx, ctx)
}

func (c *ExprCompiler) tryCompileMatchConstantFold(scrutinee fql.IExpressionContext, armsCtx fql.IMatchPatternArmsContext, dst bytecode.Operand) bool {
	scrutineeVal, arms, ok := matchConstantFoldPreconditions(scrutinee, armsCtx)
	if !ok {
		return false
	}

	selected, ok := selectMatchConstantFoldExpression(scrutineeVal, arms, armsCtx.MatchDefaultArm())
	if !ok {
		return false
	}

	return c.emitMatchConstantFoldExpression(selected, dst)
}

func matchConstantFoldPreconditions(scrutinee fql.IExpressionContext, armsCtx fql.IMatchPatternArmsContext) (runtime.Value, []fql.IMatchPatternArmContext, bool) {
	if scrutinee == nil || armsCtx == nil {
		return nil, nil, false
	}

	scrutineeVal, ok := matchLiteralValueFromExpression(scrutinee)
	if !ok {
		return nil, nil, false
	}

	list := armsCtx.MatchPatternArmList()
	if list == nil {
		return nil, nil, false
	}

	arms := list.AllMatchPatternArm()
	if len(arms) == 0 {
		return nil, nil, false
	}

	return scrutineeVal, arms, true
}

func selectMatchConstantFoldExpression(scrutinee runtime.Value, arms []fql.IMatchPatternArmContext, defaultArm fql.IMatchDefaultArmContext) (fql.IExpressionContext, bool) {
	var selected fql.IExpressionContext

	for _, arm := range arms {
		if arm == nil {
			continue
		}

		patternValue, expression, ok := matchConstantFoldArmExpression(arm)
		if !ok {
			return nil, false
		}

		if runtime.CompareValues(scrutinee, patternValue) == 0 {
			selected = expression
			break
		}
	}

	if selected == nil && defaultArm != nil {
		selected = defaultArm.Expression()
	}

	return selected, selected != nil
}

func matchConstantFoldArmExpression(arm fql.IMatchPatternArmContext) (runtime.Value, fql.IExpressionContext, bool) {
	if arm.MatchPatternGuard() != nil {
		return nil, nil, false
	}

	pattern := arm.MatchPattern()
	if pattern == nil {
		return nil, nil, false
	}

	literalPattern := pattern.MatchLiteralPattern()
	if literalPattern == nil {
		return nil, nil, false
	}

	patternValue, ok := literalValueFromMatchLiteral(literalPattern)
	if !ok {
		return nil, nil, false
	}

	return patternValue, arm.Expression(), true
}

func (c *ExprCompiler) emitMatchConstantFoldExpression(expr fql.IExpressionContext, dst bytecode.Operand) bool {
	if expr == nil {
		return false
	}

	c.ctx.Symbols.EnterScope()
	out := c.ensureRegister(c.Compile(expr))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Emitter.EmitMove(dst, out)
	}
	c.ctx.Symbols.ExitScope()

	return true
}

func matchLiteralValueFromExpression(expr fql.IExpressionContext) (runtime.Value, bool) {
	atom, ok := matchPureResultAtom(expr)
	if !ok {
		return nil, false
	}

	literal := atom.Literal()
	if literal == nil {
		return nil, false
	}

	return literalValueFromLiteral(literal)
}

func literalValueFromLiteral(lit fql.ILiteralContext) (runtime.Value, bool) {
	return scalarLiteralValue(lit)
}

func literalValueFromMatchLiteral(lit fql.IMatchLiteralPatternContext) (runtime.Value, bool) {
	return scalarLiteralValue(lit)
}

func matchPureResultKey(expr fql.IExpressionContext) (string, bool) {
	atom, ok := matchPureResultAtom(expr)
	if !ok {
		return "", false
	}

	return matchPureResultKeyFromAtom(atom)
}

func matchPureResultAtom(expr fql.IExpressionContext) (fql.IExpressionAtomContext, bool) {
	if expr == nil {
		return nil, false
	}

	if expr.UnaryOperator() != nil || expr.LogicalAndOperator() != nil || expr.LogicalOrOperator() != nil || expr.GetTernaryOperator() != nil {
		return nil, false
	}

	pred := expr.Predicate()
	if pred == nil {
		return nil, false
	}

	if pred.EqualityOperator() != nil || pred.ArrayOperator() != nil || pred.InOperator() != nil || pred.LikeOperator() != nil {
		return nil, false
	}

	atom := pred.ExpressionAtom()
	if atom == nil {
		return nil, false
	}

	if atom.ExpressionAtom(0) != nil || atom.ExpressionAtom(1) != nil {
		return nil, false
	}

	if atom.AdditiveOperator() != nil || atom.MultiplicativeOperator() != nil || atom.RegexpOperator() != nil || atom.ErrorOperator() != nil || atom.RecoveryTails() != nil || atom.RangeOperator() != nil {
		return nil, false
	}

	return atom, true
}

func matchPureResultKeyFromAtom(atom fql.IExpressionAtomContext) (string, bool) {
	if atom == nil {
		return "", false
	}

	if inner := atom.Expression(); inner != nil {
		return matchPureResultKey(inner)
	}

	if lit := atom.Literal(); lit != nil {
		val, ok := literalValueFromLiteral(lit)
		if !ok {
			return "", false
		}
		return matchPureLiteralKey(val)
	}

	if v := atom.Variable(); v != nil {
		name := matchVariableName(v)
		if name == "" {
			return "", false
		}
		return "var:" + name, true
	}

	if p := atom.Param(); p != nil {
		name := matchParamName(p)
		if name == "" {
			return "", false
		}
		return "param:" + name, true
	}

	if m := atom.MemberExpression(); m != nil {
		return matchMemberExpressionKey(m)
	}

	if im := atom.ImplicitMemberExpression(); im != nil {
		return matchImplicitMemberExpressionKey(im)
	}

	return "", false
}

func matchPureLiteralKey(val runtime.Value) (string, bool) {
	if val == nil {
		return "", false
	}

	if val == runtime.None {
		return "lit:none", true
	}

	switch v := val.(type) {
	case runtime.Boolean:
		if v {
			return "lit:true", true
		}
		return "lit:false", true
	case runtime.Int:
		return "lit:int:" + strconv.FormatInt(int64(v), 10), true
	case runtime.Float:
		return "lit:float:" + strconv.FormatFloat(float64(v), 'g', -1, 64), true
	case runtime.String:
		return "lit:str:" + strconv.Quote(string(v)), true
	default:
		return "", false
	}
}

func matchVariableName(ctx fql.IVariableContext) string {
	if ctx == nil {
		return ""
	}

	if id := ctx.Identifier(); id != nil {
		return id.GetText()
	}

	if srw := ctx.SafeReservedWord(); srw != nil {
		return srw.GetText()
	}

	return ""
}

func matchParamName(ctx fql.IParamContext) string {
	if ctx == nil {
		return ""
	}

	var name string
	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	}

	if name == "" {
		return ""
	}

	return "@" + name
}

func matchMemberExpressionKey(ctx fql.IMemberExpressionContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	source := ctx.MemberExpressionSource()
	if source == nil {
		return "", false
	}

	var base string
	if v := source.Variable(); v != nil {
		name := matchVariableName(v)
		if name == "" {
			return "", false
		}
		base = "var:" + name
	} else if p := source.Param(); p != nil {
		name := matchParamName(p)
		if name == "" {
			return "", false
		}
		base = "param:" + name
	} else {
		return "", false
	}

	paths := ctx.AllMemberExpressionPath()
	if len(paths) != 1 {
		return "", false
	}

	prop, ok := matchMemberExpressionPathKey(paths[0])
	if !ok {
		return "", false
	}

	return "member:" + base + "." + prop, true
}

func matchImplicitMemberExpressionKey(ctx fql.IImplicitMemberExpressionContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if len(ctx.AllMemberExpressionPath()) > 0 {
		return "", false
	}

	start := ctx.ImplicitMemberExpressionStart()
	if start == nil {
		return "", false
	}

	if start.ErrorOperator() != nil || start.ComputedPropertyName() != nil || start.ArrayExpansion() != nil || start.ArrayContraction() != nil || start.ArrayQuestionMark() != nil || start.ArrayApply() != nil {
		return "", false
	}

	prop, ok := matchPropertyNameKey(start.PropertyName())
	if !ok {
		return "", false
	}

	return "member:implicit:." + prop, true
}

func matchMemberExpressionPathKey(ctx fql.IMemberExpressionPathContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if ctx.ErrorOperator() != nil || ctx.ComputedPropertyName() != nil || ctx.ArrayContraction() != nil || ctx.ArrayExpansion() != nil || ctx.ArrayQuestionMark() != nil || ctx.ArrayApply() != nil {
		return "", false
	}

	return matchPropertyNameKey(ctx.PropertyName())
}

func matchPropertyNameKey(ctx fql.IPropertyNameContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else if urw := ctx.UnsafeReservedWord(); urw != nil {
		name = urw.GetText()
	} else if sl := ctx.StringLiteral(); sl != nil {
		if val, ok := parseStringLiteralConst(sl); ok {
			name = string(val)
		}
	}

	if name == "" {
		return "", false
	}

	return strconv.Quote(name), true
}

func (c *ExprCompiler) compileMatchObjectPattern(valueReg bytecode.Operand, ctx fql.IMatchObjectPatternContext, onFail core.Label) {
	if ctx == nil {
		return
	}

	props := ctx.AllMatchObjectPatternProperty()
	if len(props) == 0 {
		keys := c.emitObjectsKeys(valueReg)
		c.ctx.Emitter.EmitJumpIfNone(keys, onFail)
		return
	}

	for _, prop := range props {
		if prop == nil {
			continue
		}

		keyOp := c.compileMatchObjectPatternKey(prop.MatchObjectPatternKey())
		if keyOp == bytecode.NoopOperand {
			continue
		}

		val := c.ctx.Registers.Allocate()
		if keyOp.IsConstant() {
			c.ctx.Emitter.EmitMatchLoadPropertyConst(val, valueReg, keyOp, onFail)
		} else {
			c.ctx.Emitter.EmitJumpCompare(bytecode.OpJumpIfMissingProperty, valueReg, keyOp, onFail)
			c.ctx.Emitter.EmitABC(bytecode.OpLoadProperty, val, valueReg, keyOp)
		}
		c.compileMatchPatternValue(val, prop.MatchPattern(), onFail)
	}
}

func (c *ExprCompiler) compileMatchObjectPatternKey(ctx fql.IMatchObjectPatternKeyContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if sl := ctx.StringLiteral(); sl != nil {
		if val, ok := parseStringLiteralConst(sl); ok {
			return c.ctx.Symbols.AddConstant(val)
		}

		return c.ctx.LiteralCompiler.CompileStringLiteral(sl)
	}

	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else if urw := ctx.UnsafeReservedWord(); urw != nil {
		name = urw.GetText()
	}

	if name == "" {
		return bytecode.NoopOperand
	}

	return c.ctx.Symbols.AddConstant(runtime.NewString(name))
}

func (c *ExprCompiler) emitObjectsKeys(scrReg bytecode.Operand) bytecode.Operand {
	scrReg = c.ensureRegister(scrReg)
	seq := c.ctx.Registers.AllocateSequence(1)
	c.ctx.Emitter.EmitMove(seq[0], scrReg)

	return c.CompileFunctionCallByNameWith(nil, runtime.NewString("KEYS"), true, seq)
}

func (c *ExprCompiler) declareMatchBinding(ctx antlr.ParserRuleContext, name string, valueReg bytecode.Operand) bytecode.Operand {
	valueReg = c.ensureRegister(valueReg)
	reg, ok := c.ctx.Symbols.DeclareLocal(name, core.TypeAny)
	if ok {
		c.ctx.Emitter.EmitMove(reg, valueReg)
		c.ctx.Types.Set(reg, operandType(c.ctx, valueReg))
		return reg
	}

	if ctx != nil {
		c.ctx.Errors.DuplicateMatchBinding(ctx, name)
	}

	if existing, _, found := c.ctx.Symbols.Resolve(name); found {
		return existing
	}

	return valueReg
}

func (c *ExprCompiler) ensureRegister(op bytecode.Operand) bytecode.Operand {
	if op == bytecode.NoopOperand {
		return op
	}

	if op.IsRegister() {
		return op
	}

	reg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(reg, op)
	c.ctx.Types.Set(reg, operandType(c.ctx, op))

	return reg
}

// CompileImplicitMemberExpression processes an implicit member expression (e.g., .name, .[0], ?.name).
// The implicit current resolution is centralized in resolveImplicitCurrent.
func (c *ExprCompiler) CompileImplicitMemberExpression(ctx fql.IImplicitMemberExpressionContext) bytecode.Operand {
	start := ctx.ImplicitMemberExpressionStart()
	if start == nil {
		return bytecode.NoopOperand
	}

	src, ok := c.resolveImplicitMemberSource(start)
	if !ok {
		return bytecode.NoopOperand
	}

	segments := ctx.AllMemberExpressionPath()
	if arrayResult, handled := c.compileImplicitMemberArrayOperator(src, start, segments); handled {
		return arrayResult
	}

	return c.compileImplicitMemberPath(src, start, segments)
}

func (c *ExprCompiler) resolveImplicitMemberSource(start fql.IImplicitMemberExpressionStartContext) (bytecode.Operand, bool) {
	if start == nil {
		return bytecode.NoopOperand, false
	}

	if c.implicitCurrentDepth == 0 {
		c.resolveImplicitCurrent(getImplicitToken(start))
		return bytecode.NoopOperand, false
	}

	// resolveImplicitCurrent guarantees a register operand on success.
	return c.resolveImplicitCurrent(getImplicitToken(start))
}

func (c *ExprCompiler) compileImplicitMemberArrayOperator(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) (bytecode.Operand, bool) {
	if expansion := start.ArrayExpansion(); expansion != nil {
		return c.compileArrayExpansionChain(src, expansion, segments), true
	}

	if contraction := start.ArrayContraction(); contraction != nil {
		inlineTail, restTail := splitArrayOperatorTail(segments)
		result := c.compileArrayContraction(src, contraction, inlineTail)
		return c.continueImplicitMemberArrayResult(result, restTail), true
	}

	if question := start.ArrayQuestionMark(); question != nil {
		inlineTail, restTail := splitArrayOperatorTail(segments)
		result := c.compileArrayQuestionMark(src, question, inlineTail)
		return c.continueImplicitMemberArrayResult(result, restTail), true
	}

	if apply := start.ArrayApply(); apply != nil {
		return c.compileArrayApply(src, apply, segments), true
	}

	return bytecode.NoopOperand, false
}

func (c *ExprCompiler) continueImplicitMemberArrayResult(result bytecode.Operand, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	if len(tail) == 0 {
		return result
	}

	return c.compileMemberExpressionSegments(result, tail)
}

func (c *ExprCompiler) compileImplicitMemberPath(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	if isSimpleMemberPathChain(segments) {
		return c.compileImplicitSimpleMemberExpressionSegments(src, start, segments)
	}

	return c.compileImplicitGenericMemberExpression(src, start, segments)
}

func (c *ExprCompiler) compileImplicitGenericMemberExpression(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	dst, ok := c.emitImplicitMemberStartLoad(src, start)
	if !ok {
		return bytecode.NoopOperand
	}

	if len(segments) == 0 {
		return dst
	}

	return c.compileMemberExpressionSegments(dst, segments)
}

func (c *ExprCompiler) emitImplicitMemberStartLoad(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext) (bytecode.Operand, bool) {
	operand, constOperand := c.compileImplicitMemberStartOperand(start)
	if operand == bytecode.NoopOperand {
		return bytecode.NoopOperand, false
	}

	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(start.(antlr.ParserRuleContext))
	optional := start.ErrorOperator() != nil

	c.ctx.Emitter.WithSpan(span, func() {
		op := memberLoadOpcode(operandType(c.ctx, src), constOperand, optional)
		c.ctx.Emitter.EmitABC(op, dst, src, operand)
	})

	return dst, true
}

func (c *ExprCompiler) compileWithImplicitCurrent(expr fql.IExpressionContext) bytecode.Operand {
	if expr == nil {
		return bytecode.NoopOperand
	}

	c.implicitCurrentDepth++
	defer func() {
		c.implicitCurrentDepth--
	}()

	return c.Compile(expr)
}

func (c *ExprCompiler) withImplicitCurrent(fn func()) {
	c.implicitCurrentDepth++
	defer func() {
		c.implicitCurrentDepth--
	}()

	fn()
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
func (c *ExprCompiler) CompileMemberExpression(ctx fql.IMemberExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	plan := collectRecoveryPlan(c.ctx, ctx, core.RecoveryPlanOptions{})
	return c.ctx.OPCompiler.CompileWithRecoveryPlan(plan, core.CatchJumpModeNone, func() bytecode.Operand {
		mes := ctx.MemberExpressionSource()
		segments := ctx.AllMemberExpressionPath()

		if mes == nil || len(segments) == 0 {
			return bytecode.NoopOperand
		}

		src := c.compileMemberExpressionSource(mes, segments)

		if src == bytecode.NoopOperand {
			return src
		}

		return c.compileMemberExpressionSegments(src, segments)
	})
}

func (c *ExprCompiler) compileMemberExpressionSource(mes fql.IMemberExpressionSourceContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	if v := mes.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if p := mes.Param(); p != nil {
		return c.CompileParam(p)
	}

	if ol := mes.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	if al := mes.ArrayLiteral(); al != nil {
		return c.ctx.LiteralCompiler.CompileArrayLiteral(al)
	}

	if fc := mes.FunctionCall(); fc != nil {
		// FOO()?.bar
		segment := segments[0]
		return c.CompileFunctionCall(fc, segment.ErrorOperator() != nil)
	}

	if fe := mes.ForExpression(); fe != nil {
		return c.ctx.LoopCompiler.Compile(fe)
	}

	if wfe := mes.WaitForExpression(); wfe != nil {
		return c.ctx.WaitCompiler.Compile(wfe)
	}

	if e := mes.Expression(); e != nil {
		return c.Compile(e)
	}

	return bytecode.NoopOperand
}

// CompileImplicitCurrentExpression processes a bare implicit current shorthand (e.g., .).
func (c *ExprCompiler) CompileImplicitCurrentExpression(ctx fql.IImplicitCurrentExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	src, ok := c.resolveImplicitCurrent(getImplicitToken(ctx))
	if !ok {
		return bytecode.NoopOperand
	}

	return src
}

// resolveImplicitCurrent centralizes implicit current resolution and error reporting.
// It returns a register operand and true on success.
func (c *ExprCompiler) resolveImplicitCurrent(token antlr.Token) (bytecode.Operand, bool) {
	if c.implicitCurrentDepth == 0 {
		c.ctx.Errors.VariableNotFound(token, core.PseudoVariable)
		return bytecode.NoopOperand, false
	}

	binding, found := c.ctx.Symbols.ResolveBinding(core.PseudoVariable)
	if !found {
		c.ctx.Errors.VariableNotFound(token, core.PseudoVariable)
		return bytecode.NoopOperand, false
	}

	src := loadBindingValue(c.ctx, binding)
	src = c.ensureRegister(src)

	return src, true
}

// getImplicitToken picks the span anchor for implicit-current diagnostics.
func getImplicitToken(ctx antlr.ParserRuleContext) antlr.Token {
	switch v := ctx.(type) {
	case fql.IImplicitMemberExpressionStartContext:
		if dot := v.Dot(); dot != nil {
			return dot.GetSymbol()
		}
	case fql.IImplicitCurrentExpressionContext:
		if dot := v.Dot(); dot != nil {
			return dot.GetSymbol()
		}
	}

	return ctx.GetStart()
}

func (c *ExprCompiler) compileMemberExpressionSegments(src bytecode.Operand, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	if len(segments) == 0 {
		return src
	}

	if isSimpleMemberPathChain(segments) {
		return c.compileSimpleMemberExpressionSegments(src, segments)
	}

	for idx, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)

		if contraction := p.ArrayContraction(); contraction != nil {
			inlineTail, restTail := splitArrayOperatorTail(segments[idx+1:])
			result := c.compileArrayContraction(src, contraction, inlineTail)

			if len(restTail) == 0 {
				return result
			}

			return c.compileMemberExpressionSegments(result, restTail)
		}

		if expansion := p.ArrayExpansion(); expansion != nil {
			return c.compileArrayExpansionChain(src, expansion, segments[idx+1:])
		}

		if question := p.ArrayQuestionMark(); question != nil {
			inlineTail, restTail := splitArrayOperatorTail(segments[idx+1:])
			result := c.compileArrayQuestionMark(src, question, inlineTail)

			if len(restTail) == 0 {
				return result
			}

			return c.compileMemberExpressionSegments(result, restTail)
		}

		if apply := p.ArrayApply(); apply != nil {
			return c.compileArrayApply(src, apply, segments[idx+1:])
		}

		src2, constOperand := c.compileMemberPathOperand(p)
		dst := c.ctx.Registers.Allocate()
		span := diagnostics.SpanFromRuleContext(p)

		c.ctx.Emitter.WithSpan(span, func() {
			optional := p.ErrorOperator() != nil
			op := memberLoadOpcode(operandType(c.ctx, src), constOperand, optional)

			c.ctx.Emitter.EmitABC(op, dst, src, src2)
		})

		src = dst
	}

	return src
}

func (c *ExprCompiler) compileImplicitSimpleMemberExpressionSegments(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	startOp, startConst := c.compileImplicitMemberStartOperand(start)
	if startOp == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	state := &optionalMemberChainState{}
	startSpan := diagnostics.SpanFromRuleContext(start.(antlr.ParserRuleContext))
	result := c.emitOptionalMemberLoadSegment(startSpan, src, startOp, startConst, start.ErrorOperator() != nil, state)

	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)
		segmentOp, constOperand := c.compileMemberPathOperand(p)
		span := diagnostics.SpanFromRuleContext(p)
		result = c.emitOptionalMemberLoadSegment(span, result, segmentOp, constOperand, p.ErrorOperator() != nil, state)
	}

	c.finalizeOptionalMemberChain(state)
	return result
}

func (c *ExprCompiler) emitOptionalMemberLoadSegment(span source.Span, src, segmentOp bytecode.Operand, constOperand, optional bool, state *optionalMemberChainState) bytecode.Operand {
	dst := c.allocateOptionalMemberDestination(src, state)

	c.ctx.Emitter.WithSpan(span, func() {
		op := memberLoadOpcode(operandType(c.ctx, src), constOperand, optional)
		c.ctx.Emitter.EmitABC(op, dst, src, segmentOp)

		if optional {
			c.ctx.Emitter.EmitJumpIfNone(dst, c.optionalMemberEndLabel(state))
		}
	})

	if optional {
		state.stickyDst = true
	}

	return dst
}

func (c *ExprCompiler) allocateOptionalMemberDestination(src bytecode.Operand, state *optionalMemberChainState) bytecode.Operand {
	if state != nil && state.stickyDst && src.IsRegister() {
		return src
	}

	return c.ctx.Registers.Allocate()
}

func (c *ExprCompiler) optionalMemberEndLabel(state *optionalMemberChainState) core.Label {
	if state == nil {
		return core.Label{}
	}

	if state.hasJump {
		return state.endLabel
	}

	state.endLabel = c.ctx.Emitter.NewLabel("member", "optional", "end")
	state.hasJump = true

	return state.endLabel
}

func (c *ExprCompiler) finalizeOptionalMemberChain(state *optionalMemberChainState) {
	if state == nil || !state.hasJump {
		return
	}

	c.ctx.Emitter.MarkLabel(state.endLabel)
}

func isSimpleMemberPathChain(segments []fql.IMemberExpressionPathContext) bool {
	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)

		if p.PropertyName() == nil && p.ComputedPropertyName() == nil {
			return false
		}
	}

	return true
}

func (c *ExprCompiler) compileSimpleMemberExpressionSegments(src bytecode.Operand, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	result := src
	stickyDst := false
	hasJump := false
	var endLabel core.Label

	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)
		src2, constOperand := c.compileMemberPathOperand(p)
		optional := p.ErrorOperator() != nil

		dst := result
		if !stickyDst || !dst.IsRegister() {
			dst = c.ctx.Registers.Allocate()
		}

		span := diagnostics.SpanFromRuleContext(p)

		c.ctx.Emitter.WithSpan(span, func() {
			op := memberLoadOpcode(operandType(c.ctx, result), constOperand, optional)
			c.ctx.Emitter.EmitABC(op, dst, result, src2)

			if optional {
				if !hasJump {
					endLabel = c.ctx.Emitter.NewLabel("member", "optional", "end")
					hasJump = true
				}
				c.ctx.Emitter.EmitJumpIfNone(dst, endLabel)
			}
		})

		if optional {
			stickyDst = true
		}

		result = dst
	}

	if hasJump {
		c.ctx.Emitter.MarkLabel(endLabel)
	}

	return result
}

func (c *ExprCompiler) compileMemberPathOperand(p *fql.MemberExpressionPathContext) (bytecode.Operand, bool) {
	if pn := p.PropertyName(); pn != nil {
		if constOp, ok := c.ctx.LiteralCompiler.CompilePropertyNameConst(pn); ok {
			return constOp, true
		}

		return c.ctx.LiteralCompiler.CompilePropertyName(pn), false
	}

	if cpn := p.ComputedPropertyName(); cpn != nil {
		if val, ok := literalValueFromExpression(cpn.Expression()); ok {
			switch val.(type) {
			case *runtime.Array, *runtime.Object:
				// Keep array/object literals dynamic to preserve their stringified key value.
				return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
			default:
				return c.ctx.Symbols.AddConstant(val), true
			}
		}

		return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
	}

	return bytecode.NoopOperand, false
}

func (c *ExprCompiler) compileImplicitMemberStartOperand(start fql.IImplicitMemberExpressionStartContext) (bytecode.Operand, bool) {
	if pn := start.PropertyName(); pn != nil {
		if constOp, ok := c.ctx.LiteralCompiler.CompilePropertyNameConst(pn); ok {
			return constOp, true
		}

		return c.ctx.LiteralCompiler.CompilePropertyName(pn), false
	}

	if cpn := start.ComputedPropertyName(); cpn != nil {
		if val, ok := literalValueFromExpression(cpn.Expression()); ok {
			switch val.(type) {
			case *runtime.Array, *runtime.Object:
				// Keep array/object literals dynamic to preserve their stringified key value.
				return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
			default:
				return c.ctx.Symbols.AddConstant(val), true
			}
		}

		return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
	}

	return bytecode.NoopOperand, false
}

func memberLoadOpcode(srcType core.ValueType, constOperand, optional bool) bytecode.Opcode {
	switch srcType {
	case core.TypeArray:
		if constOperand {
			if optional {
				return bytecode.OpLoadIndexOptionalConst
			}

			return bytecode.OpLoadIndexConst
		}

		if optional {
			return bytecode.OpLoadIndexOptional
		}

		return bytecode.OpLoadIndex
	case core.TypeObject:
		if constOperand {
			if optional {
				return bytecode.OpLoadKeyOptionalConst
			}

			return bytecode.OpLoadKeyConst
		}

		if optional {
			return bytecode.OpLoadKeyOptional
		}

		return bytecode.OpLoadKey
	default:
		if constOperand {
			if optional {
				return bytecode.OpLoadPropertyOptionalConst
			}

			return bytecode.OpLoadPropertyConst
		}

		if optional {
			return bytecode.OpLoadPropertyOptional
		}

		return bytecode.OpLoadProperty
	}
}

func splitArrayOperatorTail(segments []fql.IMemberExpressionPathContext) ([]fql.IMemberExpressionPathContext, []fql.IMemberExpressionPathContext) {
	if len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)

		if p.ArrayContraction() != nil || p.ArrayExpansion() != nil || p.ArrayQuestionMark() != nil {
			return nil, segments
		}
	}

	return segments, nil
}

// splitTerminalArrayContractionTail hoists only the final contraction segment.
// Earlier array operators stay in the per-element tail so existing projection semantics remain unchanged.
func splitTerminalArrayContractionTail(segments []fql.IMemberExpressionPathContext) ([]fql.IMemberExpressionPathContext, fql.IArrayContractionContext) {
	if len(segments) == 0 {
		return nil, nil
	}

	last := segments[len(segments)-1].(*fql.MemberExpressionPathContext)
	contraction := last.ArrayContraction()
	if contraction == nil {
		return segments, nil
	}

	return segments[:len(segments)-1], contraction
}

func (c *ExprCompiler) compileArrayExpansionChain(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	inline := expansion.InlineExpression()

	if inline == nil {
		if next, rest := nextArrayExpansion(tail); next != nil {
			return c.compileArrayExpansionChain(src, next, rest)
		}

		return c.compileArrayExpansionChainWithFilters(src, expansion, tail, nil)
	}

	if !isFilterOnlyInline(inline) {
		tail = dropIdentityExpansions(tail)

		return c.compileArrayExpansionChainWithFilters(src, expansion, tail, nil)
	}

	extraFilters, rest := collectFilterOnlyTail(tail)

	return c.compileArrayExpansionChainWithFilters(src, expansion, rest, extraFilters)
}

func (c *ExprCompiler) compileArrayExpansionWithFilters(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext, extraFilters []fql.IExpressionContext) bytecode.Operand {
	span := diagnostics.SpanFromRuleContext(expansion)
	inline := expansion.InlineExpression()

	return c.compileArrayIteration(src, span, tail, inline, extraFilters)
}

func (c *ExprCompiler) compileArrayExpansionChainWithFilters(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext, extraFilters []fql.IExpressionContext) bytecode.Operand {
	inlineTail, restTail := splitArrayOperatorTail(tail)
	result := c.compileArrayExpansionWithFilters(src, expansion, inlineTail, extraFilters)

	if len(restTail) == 0 {
		return result
	}

	return c.compileMemberExpressionSegments(result, restTail)
}

func isFilterOnlyInline(inline fql.IInlineExpressionContext) bool {
	if inline == nil {
		return false
	}

	return inline.InlineFilter() != nil && inline.InlineLimit() == nil && inline.InlineReturn() == nil
}

func nextArrayExpansion(segments []fql.IMemberExpressionPathContext) (fql.IArrayExpansionContext, []fql.IMemberExpressionPathContext) {
	if len(segments) == 0 {
		return nil, segments
	}

	p := segments[0].(*fql.MemberExpressionPathContext)

	if expansion := p.ArrayExpansion(); expansion != nil {
		return expansion, segments[1:]
	}

	return nil, segments
}

func dropIdentityExpansions(segments []fql.IMemberExpressionPathContext) []fql.IMemberExpressionPathContext {
	for len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()

		if expansion == nil {
			break
		}

		if expansion.InlineExpression() != nil {
			break
		}

		segments = segments[1:]
	}

	return segments
}

func collectFilterOnlyTail(segments []fql.IMemberExpressionPathContext) ([]fql.IExpressionContext, []fql.IMemberExpressionPathContext) {
	extraFilters := make([]fql.IExpressionContext, 0)
	rest := segments

	for len(rest) > 0 {
		p := rest[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()
		if expansion == nil {
			break
		}

		inline := expansion.InlineExpression()
		if inline == nil {
			rest = rest[1:]
			continue
		}

		if !isFilterOnlyInline(inline) {
			break
		}

		filter := inline.InlineFilter()
		if filter != nil {
			extraFilters = append(extraFilters, filter.Expression())
		}

		rest = rest[1:]
	}

	return extraFilters, rest
}

func (c *ExprCompiler) compileArrayQuestionMark(src bytecode.Operand, question fql.IArrayQuestionMarkContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	span := diagnostics.SpanFromRuleContext(question)

	loop := &core.Loop{
		Kind:     core.ForInLoop,
		Type:     core.NormalLoop,
		Distinct: false,
		Allocate: false,
		Dst:      bytecode.NoopOperand,
		Src:      src,
	}

	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	loop.DeclareValueVar(core.PseudoVariable, c.ctx.Symbols, core.TypeAny)

	if loop.Value.IsRegister() {
		c.ctx.Types.Set(loop.Value, core.TypeAny)
	}

	count := c.ctx.Registers.Allocate()
	total := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitA(bytecode.OpLoadZero, count)
		c.ctx.Emitter.EmitA(bytecode.OpLoadZero, total)
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
	})

	// Track total elements
	c.ctx.Emitter.EmitA(bytecode.OpIncr, total)

	// Apply optional filter
	if filter := question.Expression(); filter != nil {
		cond := c.compileWithImplicitCurrent(filter)
		label := c.ctx.Loops.Current().ContinueLabel()
		c.ctx.Emitter.EmitJumpIfFalse(cond, label)
	}

	// Count matches
	c.ctx.Emitter.EmitA(bytecode.OpIncr, count)

	loop.EmitFinalization(c.ctx.Emitter)

	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	result := c.compileArrayQuestionQuantifier(question, count, total)

	if len(tail) > 0 {
		result = c.compileMemberExpressionSegments(result, tail)
	}

	if result.IsRegister() {
		c.ctx.Types.Set(result, core.TypeBool)
	}

	return result
}

func (c *ExprCompiler) compileArrayQuestionQuantifier(question fql.IArrayQuestionMarkContext, count, total bytecode.Operand) bytecode.Operand {
	quant := question.ArrayQuestionQuantifier()
	zero := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitA(bytecode.OpLoadZero, zero)

	if quant == nil || quant.Any() != nil {
		return c.emitComparison(bytecode.OpGt, count, zero)
	}

	if quant.None() != nil {
		return c.emitComparison(bytecode.OpEq, count, zero)
	}

	if quant.All() != nil {
		return c.emitComparison(bytecode.OpEq, count, total)
	}

	values := quant.AllArrayQuestionQuantifierValue()

	if quant.At() != nil {
		if len(values) == 0 {
			return c.emitComparison(bytecode.OpGt, count, zero)
		}

		value := c.compileArrayQuestionQuantifierValue(values[0])

		return c.emitComparison(bytecode.OpGte, count, value)
	}

	if quant.Range() != nil && len(values) >= 2 {
		min := c.compileArrayQuestionQuantifierValue(values[0])
		max := c.compileArrayQuestionQuantifierValue(values[1])

		left := c.emitComparison(bytecode.OpGte, count, min)
		right := c.emitComparison(bytecode.OpLte, count, max)

		return c.emitBooleanAnd(left, right)
	}

	if len(values) > 0 {
		value := c.compileArrayQuestionQuantifierValue(values[0])

		return c.emitComparison(bytecode.OpEq, count, value)
	}

	return c.emitComparison(bytecode.OpGt, count, zero)
}

func (c *ExprCompiler) compileArrayQuestionQuantifierValue(ctx fql.IArrayQuestionQuantifierValueContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	if pm := ctx.Param(); pm != nil {
		return c.CompileParam(pm)
	}

	return bytecode.NoopOperand
}

func (c *ExprCompiler) compileArrayApply(src bytecode.Operand, apply fql.IArrayApplyContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	if apply == nil {
		return src
	}

	query := c.compileQueryLiteral(apply.QueryLiteral())
	if query == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(apply)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABC(bytecode.OpQuery, dst, src, query)
	})

	if len(tail) > 0 {
		return c.compileMemberExpressionSegments(dst, tail)
	}

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeList)
	}

	return dst
}

func (c *ExprCompiler) compileQueryExpression(ctx fql.IQueryExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	plan := collectRecoveryPlan(c.ctx, ctx, core.RecoveryPlanOptions{})
	return c.ctx.OPCompiler.CompileWithRecoveryPlan(plan, core.CatchJumpModeNone, func() bytecode.Operand {
		if ctx == nil {
			return bytecode.NoopOperand
		}

		src, ok := c.compileQueryExpressionSource(ctx)
		if !ok {
			return bytecode.NoopOperand
		}

		span := diagnostics.SpanFromRuleContext(ctx)
		modifier := queryModifierName(ctx.QueryModifier())
		queryReg := c.emitQueryEnvelope(ctx, span)
		queryResult := c.emitApplyQuery(span, src, queryReg)
		dst := c.lowerQueryModifier(span, modifier, queryResult)

		if dst.IsRegister() {
			c.ctx.Types.Set(dst, queryResultTypeForModifier(modifier))
		}

		return dst
	})
}

func (c *ExprCompiler) compileQueryExpressionSource(ctx fql.IQueryExpressionContext) (bytecode.Operand, bool) {
	if ctx == nil {
		return bytecode.NoopOperand, false
	}

	sourceExpr := ctx.Expression()
	if sourceExpr == nil {
		return bytecode.NoopOperand, false
	}

	source := c.Compile(sourceExpr)
	return source, source != bytecode.NoopOperand
}

func (c *ExprCompiler) emitQueryEnvelope(ctx fql.IQueryExpressionContext, span source.Span) bytecode.Operand {
	queryReg := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArray(queryReg, 3)
	})

	kind := c.compileQueryKindOperand(ctx)
	c.emitQueryEnvelopeOperand(span, queryReg, kind)

	payload := c.compileQueryPayloadOperand(ctx.QueryPayload())
	c.emitQueryEnvelopeOperand(span, queryReg, payload)

	options := c.compileQueryOptionsOperand(ctx.QueryWithOpt())
	c.emitQueryEnvelopeOperand(span, queryReg, options)

	return queryReg
}

func (c *ExprCompiler) compileQueryKindOperand(ctx fql.IQueryExpressionContext) bytecode.Operand {
	kind := ""
	if ident := ctx.GetDialect(); ident != nil {
		kind = strings.ToLower(ident.GetText())
	}

	return loadConstant(c.ctx, runtime.NewString(kind))
}

func (c *ExprCompiler) compileQueryPayloadOperand(ctx fql.IQueryPayloadContext) bytecode.Operand {
	if ctx == nil {
		return loadConstant(c.ctx, runtime.EmptyString)
	}

	if literal := ctx.StringLiteral(); literal != nil {
		if value, ok := parseStringLiteralConst(literal); ok {
			return loadConstant(c.ctx, value)
		}

		return c.ctx.LiteralCompiler.CompileStringLiteral(literal)
	}

	if param := ctx.Param(); param != nil {
		return c.CompileParam(param)
	}

	if variable := ctx.Variable(); variable != nil {
		return c.CompileVariable(variable)
	}

	return loadConstant(c.ctx, runtime.EmptyString)
}

func (c *ExprCompiler) compileQueryOptionsOperand(ctx fql.IQueryWithOptContext) bytecode.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return loadConstant(c.ctx, runtime.None)
	}

	return c.Compile(ctx.Expression())
}

func (c *ExprCompiler) emitQueryEnvelopeOperand(span source.Span, queryReg, value bytecode.Operand) {
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(queryReg, value)
	})
}

func (c *ExprCompiler) emitApplyQuery(span source.Span, src, queryReg bytecode.Operand) bytecode.Operand {
	result := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABC(bytecode.OpQuery, result, src, queryReg)
	})

	return result
}

func (c *ExprCompiler) lowerQueryModifier(span source.Span, modifier queryModifier, queryResult bytecode.Operand) bytecode.Operand {
	switch modifier {
	case queryModifierExists:
		dst := c.ctx.Registers.Allocate()
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(bytecode.OpExists, dst, queryResult)
		})
		return dst
	case queryModifierCount:
		dst := c.ctx.Registers.Allocate()
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(bytecode.OpLength, dst, queryResult)
		})
		return dst
	case queryModifierAny:
		dst := c.ctx.Registers.Allocate()
		zero := c.ctx.Symbols.AddConstant(runtime.NewInt(0))
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitABC(bytecode.OpLoadIndexOptionalConst, dst, queryResult, zero)
		})
		return dst
	case queryModifierValue:
		return c.ctx.lowerQueryModifierValue(span, queryResult)
	case queryModifierOne:
		return c.ctx.lowerQueryModifierOne(span, queryResult)
	default:
		return queryResult
	}
}

func queryModifierName(ctx fql.IQueryModifierContext) queryModifier {
	if ctx == nil {
		return queryModifierUnknown
	}

	return parseQueryModifier(ctx.GetText())
}

func queryResultTypeForModifier(modifier queryModifier) core.ValueType {
	switch modifier {
	case queryModifierExists:
		return core.TypeBool
	case queryModifierCount:
		return core.TypeInt
	case queryModifierAny, queryModifierValue, queryModifierOne:
		return core.TypeAny
	default:
		return core.TypeList
	}
}

func parseQueryModifier(text string) queryModifier {
	switch strings.ToLower(text) {
	case string(queryModifierExists):
		return queryModifierExists
	case string(queryModifierCount):
		return queryModifierCount
	case string(queryModifierAny):
		return queryModifierAny
	case string(queryModifierValue):
		return queryModifierValue
	case string(queryModifierOne):
		return queryModifierOne
	default:
		return queryModifierUnknown
	}
}

func (c *CompilerContext) lowerQueryModifierValue(span source.Span, queryResult bytecode.Operand) bytecode.Operand {
	dst := c.Registers.Allocate()
	cond := c.Registers.Allocate()
	zero := c.Symbols.AddConstant(runtime.NewInt(0))
	message := c.Symbols.AddConstant(runtime.NewString(queryValueFailMessage))
	success := c.Emitter.NewLabel("query", string(queryModifierValue), "ok")
	end := c.Emitter.NewLabel("query", string(queryModifierValue), "end")

	c.Emitter.WithSpan(span, func() {
		c.Emitter.EmitAB(bytecode.OpExists, cond, queryResult)
		c.Emitter.EmitJumpIfTrue(cond, success)
		c.Emitter.EmitLoadNone(dst)
		c.Emitter.EmitA(bytecode.OpFail, message)
		c.Emitter.EmitJump(end)
		c.Emitter.MarkLabel(success)
		c.Emitter.EmitABC(bytecode.OpLoadIndexConst, dst, queryResult, zero)
		c.Emitter.MarkLabel(end)
	})

	return dst
}

func (c *CompilerContext) lowerQueryModifierOne(span source.Span, queryResult bytecode.Operand) bytecode.Operand {
	dst := c.Registers.Allocate()
	length := c.Registers.Allocate()
	one := c.Symbols.AddConstant(runtime.NewInt(1))
	zero := c.Symbols.AddConstant(runtime.NewInt(0))
	message := c.Symbols.AddConstant(runtime.NewString(queryOneFailMessage))
	success := c.Emitter.NewLabel("query", string(queryModifierOne), "ok")
	end := c.Emitter.NewLabel("query", string(queryModifierOne), "end")

	c.Emitter.WithSpan(span, func() {
		c.Emitter.EmitAB(bytecode.OpLength, length, queryResult)
		c.Emitter.EmitJumpCompare(bytecode.OpJumpIfEqConst, length, one, success)
		c.Emitter.EmitLoadNone(dst)
		c.Emitter.EmitA(bytecode.OpFail, message)
		c.Emitter.EmitJump(end)
		c.Emitter.MarkLabel(success)
		c.Emitter.EmitABC(bytecode.OpLoadIndexConst, dst, queryResult, zero)
		c.Emitter.MarkLabel(end)
	})

	return dst
}

func (c *ExprCompiler) compileQueryLiteral(ctx fql.IQueryLiteralContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	kind := ""
	if ident := ctx.Identifier(); ident != nil {
		kind = strings.ToLower(ident.GetText())
	}

	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(ctx)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArray(dst, 3)
	})

	kindReg := loadConstant(c.ctx, runtime.NewString(kind))

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, kindReg)
	})

	payloadReg := loadConstant(c.ctx, runtime.EmptyString)
	if str := ctx.StringLiteral(); str != nil {
		if val, ok := parseStringLiteralConst(str); ok {
			payloadReg = loadConstant(c.ctx, val)
		} else {
			payloadReg = c.ctx.LiteralCompiler.CompileStringLiteral(str)
		}
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, payloadReg)
	})

	params := ctx.Expression()
	var paramsReg bytecode.Operand

	if params == nil {
		paramsReg = loadConstant(c.ctx, runtime.None)
	} else {
		paramsReg = c.Compile(params)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, paramsReg)
	})

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeAny)
	}

	return dst
}

func (c *ExprCompiler) emitComparison(op bytecode.Opcode, left, right bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(op, dst, left, right)

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeBool)
	}

	return dst
}

func (c *ExprCompiler) emitBooleanAnd(left, right bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	skip := c.ctx.Emitter.NewLabel("and.false")
	done := c.ctx.Emitter.NewLabel("and.done")

	c.ctx.Emitter.EmitJumpIfFalse(left, skip)
	c.ctx.Emitter.EmitJumpIfFalse(right, skip)
	c.ctx.Emitter.EmitAb(bytecode.OpLoadBool, dst, true)
	c.ctx.Emitter.EmitJump(done)

	c.ctx.Emitter.MarkLabel(skip)
	c.ctx.Emitter.EmitAb(bytecode.OpLoadBool, dst, false)
	c.ctx.Emitter.MarkLabel(done)

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeBool)
	}

	return dst
}

//lint:ignore U1000 Ignore unused method
func (c *ExprCompiler) compileArrayExpansion(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	span := diagnostics.SpanFromRuleContext(expansion)
	inline := expansion.InlineExpression()

	return c.compileArrayIteration(src, span, tail, inline, nil)
}

func (c *ExprCompiler) compileArrayContraction(src bytecode.Operand, contraction fql.IArrayContractionContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	depth := arrayContractionDepth(contraction)

	if depth < 1 {
		depth = 1
	}

	span := diagnostics.SpanFromRuleContext(contraction)
	dst := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABx(bytecode.OpFlatten, dst, src, depth)
	})

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeList)
	}

	inline := contraction.InlineExpression()

	if len(tail) == 0 && inline == nil {
		return dst
	}

	return c.compileArrayIteration(dst, span, tail, inline, nil)
}

func arrayContractionDepth(ctx fql.IArrayContractionContext) int {
	if ctx == nil {
		return 1
	}

	count := len(ctx.GetStars())

	if count > 1 {
		return count - 1
	}

	return 1
}

func (c *ExprCompiler) compileArrayIteration(src bytecode.Operand, span source.Span, tail []fql.IMemberExpressionPathContext, inline fql.IInlineExpressionContext, extraFilters []fql.IExpressionContext) bytecode.Operand {
	tail, postLoopContraction := splitTerminalArrayContractionTail(tail)

	loop := &core.Loop{
		Kind:     core.ForInLoop,
		Type:     core.NormalLoop,
		Distinct: false,
		Allocate: true,
		Dst:      c.ctx.Registers.Allocate(),
		Src:      src,
	}

	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	loop.DeclareValueVar(core.PseudoVariable, c.ctx.Symbols, core.TypeAny)
	if loop.Value.IsRegister() {
		c.ctx.Types.Set(loop.Value, core.TypeAny)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
	})

	if inline != nil {
		c.compileInlineFilter(inline)
	}

	for _, expr := range extraFilters {
		c.compileInlineFilterExpr(expr)
	}

	if inline != nil {
		c.compileInlineLimit(inline)
	}

	projection := loop.Value

	if inline != nil {
		if ret := inline.InlineReturn(); ret != nil {
			projection = c.compileWithImplicitCurrent(ret.Expression())
		}
	}

	if len(tail) > 0 {
		projection = c.compileMemberExpressionSegments(projection, tail)
	}

	c.ctx.Emitter.EmitAB(bytecode.OpPush, loop.Dst, projection)
	loop.EmitFinalization(c.ctx.Emitter)

	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	if loop.Dst.IsRegister() {
		c.ctx.Types.Set(loop.Dst, core.TypeList)
	}

	if postLoopContraction != nil {
		return c.compileArrayContraction(loop.Dst, postLoopContraction, nil)
	}

	return loop.Dst
}

func (c *ExprCompiler) compileInlineFilter(inline fql.IInlineExpressionContext) {
	if inline == nil {
		return
	}

	filter := inline.InlineFilter()

	if filter == nil {
		return
	}

	src := c.compileWithImplicitCurrent(filter.Expression())
	label := c.ctx.Loops.Current().ContinueLabel()
	c.ctx.Emitter.EmitJumpIfFalse(src, label)
}

func (c *ExprCompiler) compileInlineFilterExpr(expr fql.IExpressionContext) {
	if expr == nil {
		return
	}

	src := c.compileWithImplicitCurrent(expr)
	label := c.ctx.Loops.Current().ContinueLabel()
	c.ctx.Emitter.EmitJumpIfFalse(src, label)
}

func (c *ExprCompiler) compileInlineLimit(inline fql.IInlineExpressionContext) {
	if inline == nil {
		return
	}

	limit := inline.InlineLimit()
	if limit == nil {
		return
	}

	clauses := limit.AllLimitClauseValue()
	if len(clauses) == 0 {
		return
	}

	c.withImplicitCurrent(func() {
		if len(clauses) == 1 {
			c.ctx.LoopCompiler.compileLimit(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[0]))
			return
		}

		c.ctx.LoopCompiler.compileOffset(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[0]))
		c.ctx.LoopCompiler.compileLimit(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[1]))
	})
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
func (c *ExprCompiler) CompileVariable(ctx fql.IVariableContext) bytecode.Operand {
	// Check if the context is valid (in case of parser errors)
	var name string
	token := ctx.GetStart()

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
		token = id.GetSymbol()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else {
		return bytecode.NoopOperand
	}

	// Just return the register / constant index
	binding, found := c.ctx.Symbols.ResolveBinding(name)
	if !found {
		c.ctx.Errors.VariableNotFound(token, name)

		return bytecode.NoopOperand
	}

	op := loadBindingValue(c.ctx, binding)

	if op.IsRegister() {
		return op
	}

	return c.ensureRegister(op)
}

// CompileParam processes a parameter reference (e.g., @paramName) from the FQL AST.
// It binds the parameter name in the symbol table and emits instructions to load
// the parameter value at runtime.
// Parameters:
//   - ctx: The parameter context from the AST
//
// Returns:
//   - An operand representing the parameter's value
func (c *ExprCompiler) CompileParam(ctx fql.IParamContext) bytecode.Operand {
	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else {
		return bytecode.NoopOperand
	}

	reg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadParam(reg, c.ctx.Symbols.BindParam(name))
	c.ctx.Types.Set(reg, core.TypeAny)

	return reg
}

// CompileFunctionCallExpression processes a function call expression from the FQL AST.
// It handles both regular function calls and protected function calls (with TRY).
// Parameters:
//   - ctx: The function call expression context from the AST
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	call := ctx.FunctionCall()
	if ctx.ErrorOperator() != nil {
		return c.ctx.OPCompiler.CompileWithErrorPolicy(core.ErrorPolicySuppress, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.CompileFunctionCall(call, true)
		})
	}

	plan := collectRecoveryPlan(c.ctx, ctx, core.RecoveryPlanOptions{})
	if hasErrorReturnNoneHandler(plan) {
		out := c.ctx.OPCompiler.CompileWithErrorPolicy(core.ErrorPolicySuppress, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.CompileFunctionCall(call, true)
		})

		return widenRecoveryResultType(c.ctx, out, plan)
	}

	return c.ctx.OPCompiler.CompileWithRecoveryPlan(plan, core.CatchJumpModeNone, func() bytecode.Operand {
		return c.CompileFunctionCall(call, false)
	})
}

// CompileFunctionCall processes a function call from the FQL AST.
// It compiles the function arguments and delegates to CompileFunctionCallWith.
// Parameters:
//   - ctx: The function call context from the AST
//   - protected: Whether this is a protected call (with TRY)
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCall(ctx fql.IFunctionCallContext, protected bool) bytecode.Operand {
	return c.CompileFunctionCallWith(ctx, protected, c.CompileArgumentList(ctx.ArgumentList()))
}

// CompileFunctionCallWith processes a function call with pre-compiled arguments from the FQL AST.
// It extracts the function name and delegates to CompileFunctionCallByNameWith.
// Parameters:
//   - ctx: The function call context from the AST
//   - protected: Whether this is a protected call
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCallWith(ctx fql.IFunctionCallContext, protected bool, seq core.RegisterSequence) bytecode.Operand {
	name := getFunctionName(ctx, c.ctx.UseAliases)
	span := source.Span{Start: -1, End: -1}

	if ctx != nil {
		if fn := ctx.FunctionName(); fn != nil {
			span = diagnostics.SpanFromRuleContext(fn)

			if ns := ctx.Namespace(); ns != nil && ns.GetStart() != nil {
				span.Start = ns.GetStart().GetStart()
			}
		} else if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			span = diagnostics.SpanFromRuleContext(prc)
		}
	}

	var out bytecode.Operand

	c.ctx.Emitter.WithSpan(span, func() {
		out = c.CompileFunctionCallByNameWith(ctx, name, protected, seq)
	})

	return out
}

// CompileFunctionCallByNameWith processes a function call by name with pre-compiled arguments.
// It resolves UDFs first (if in scope), then built-ins, and finally host functions.
// Parameters:
//   - ctx: The function call context (used for namespace detection and diagnostics)
//   - name: The function name
//   - protected: Whether this is a protected call (with TRY)
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
//
// Panics if a built-in function is called with an incorrect number of arguments.
func (c *ExprCompiler) CompileFunctionCallByNameWith(ctx fql.IFunctionCallContext, name runtime.String, protected bool, seq core.RegisterSequence) bytecode.Operand {
	nameStr := name.String()
	builtinName := strings.ToUpper(nameStr)

	namespaced := strings.Contains(nameStr, runtime.NamespaceSeparator)
	if ctx != nil {
		if ns := ctx.Namespace(); ns != nil && ns.GetText() != "" {
			namespaced = true
		}
	}

	var callCtx antlr.ParserRuleContext
	if ctx != nil {
		if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			callCtx = prc
		}
	}

	if !namespaced && c.ctx.UDFs != nil && c.ctx.UDFScope != nil {
		if udfName, ok := getUDFName(ctx, c.ctx.UseAliases); ok {
			if fn, ok := c.ctx.UDFs.Resolve(udfName, c.ctx.UDFScope); ok {
				return c.compileUdfCallWith(fn, protected, seq, callCtx)
			}
		}
	}

	if !namespaced {
		switch builtinName {
		case runtimeLength:
			dst := c.ctx.Registers.Allocate()

			if seq == nil || len(seq) != 1 {
				panic(runtime.Error(runtime.ErrInvalidArgument, runtimeLength+": expected 1 argument"))
			}

			c.ctx.Emitter.EmitAB(bytecode.OpLength, dst, seq[0])

			return dst
		case runtimeTypename:
			dst := c.ctx.Registers.Allocate()

			if seq == nil || len(seq) != 1 {
				panic(runtime.Error(runtime.ErrInvalidArgument, runtimeTypename+": expected 1 argument"))
			}

			c.ctx.Emitter.EmitAB(bytecode.OpType, dst, seq[0])

			return dst
		case runtimeWait:
			if len(seq) != 1 {
				panic(runtime.Error(runtime.ErrInvalidArgument, runtimeWait+": expected 1 argument"))
			}

			c.ctx.Emitter.EmitA(bytecode.OpSleep, seq[0])

			return seq[0]
		}
	}

	return c.compileHostFunctionCallWith(name, protected, seq)
}

// compileHostFunctionCallWith processes a host function call with pre-compiled arguments.
// It loads the function name as a constant, binds the function in the symbol table,
// and emits the appropriate host call instruction based on the number of arguments and whether
// the call is protected.
// Parameters:
//   - name: The function name
//   - protected: Whether this is a protected call
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) compileHostFunctionCallWith(name runtime.String, protected bool, seq core.RegisterSequence) bytecode.Operand {
	dest := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(name))
	c.ctx.Symbols.BindFunction(name.String(), len(seq))

	opcode := bytecode.OpHCall
	if protected {
		opcode = bytecode.OpProtectedHCall
	}

	c.ctx.Emitter.EmitAs(opcode, dest, seq)

	c.ctx.Types.Set(dest, core.TypeAny)

	return dest
}

// compileUdfCallWith processes a UDF call with pre-compiled arguments.
func (c *ExprCompiler) compileUdfCallWith(fn *core.UDFInfo, protected bool, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) bytecode.Operand {
	args := c.prepareUdfCallArgs(fn, seq, callCtx)

	dest := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(runtime.NewInt(fn.ID)))

	opcode := bytecode.OpCall
	if protected {
		opcode = bytecode.OpProtectedCall
	}

	c.ctx.Emitter.EmitAs(opcode, dest, args)

	c.ctx.Types.Set(dest, core.TypeAny)

	return dest
}

// EmitUdfTailCall emits a tail call to a UDF with pre-compiled arguments.
func (c *ExprCompiler) EmitUdfTailCall(fn *core.UDFInfo, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) {
	args := c.prepareUdfCallArgs(fn, seq, callCtx)

	dest := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(runtime.NewInt(fn.ID)))

	c.ctx.Emitter.EmitAs(bytecode.OpTailCall, dest, args)
}

func (c *ExprCompiler) prepareUdfCallArgs(fn *core.UDFInfo, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) core.RegisterSequence {
	if fn == nil {
		return seq
	}

	if len(seq) != len(fn.Params) && c.ctx.Errors != nil {
		ctx := callCtx
		if ctx == nil && fn.Decl != nil {
			if prc, ok := fn.Decl.(antlr.ParserRuleContext); ok {
				ctx = prc
			}
		}

		if ctx != nil {
			name := fn.DisplayName
			if name == "" {
				name = fn.Name
			}

			c.ctx.Errors.Add(c.ctx.Errors.Create(diagnostics.NameError, ctx, fmt.Sprintf("Function '%s' expects %d arguments, got %d", name, len(fn.Params), len(seq))))
		}
	}

	if len(fn.Captures) == 0 {
		return seq
	}

	total := len(seq) + len(fn.Captures)
	args := c.ctx.Registers.AllocateSequence(total)

	for i, src := range seq {
		c.ctx.Emitter.EmitMove(args[i], src)
		c.ctx.Types.Set(args[i], operandType(c.ctx, src))
	}

	for i, capture := range fn.Captures {
		binding, ok := c.ctx.Symbols.ResolveBinding(capture.Name)
		if !ok {
			if callCtx != nil {
				c.ctx.Errors.VariableNotFound(callCtx.GetStart(), capture.Name)
			}
			continue
		}

		dst := args[len(seq)+i]

		if capture.Storage == core.BindingStorageCell {
			c.ctx.Emitter.EmitPlainMove(dst, binding.Register)
			c.ctx.Types.Set(dst, core.TypeAny)
			continue
		}

		src := loadBindingValue(c.ctx, binding)
		c.ctx.EmitMoveAuto(dst, src)
		c.ctx.Types.Set(dst, operandType(c.ctx, src))
	}

	return args
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
			// Fast path: load scalar/none literals directly into the argument window.
			// This avoids generating LOADC temp + MOVE arg,temp pairs.
			if val, ok := literalValueFromExpression(exp); ok && (bool(runtime.IsScalar(val)) || val == runtime.None) {
				c.ctx.Emitter.EmitLoadConst(seq[i], c.ctx.Symbols.AddConstant(val))
				c.ctx.Types.Set(seq[i], valueTypeFromRuntime(val))
				continue
			}

			// Compile expression and move/load into seq register
			srcReg := c.Compile(exp)

			// The reason we move is that the argument list must be a contiguous sequence of registers
			// Otherwise, we cannot compileInitialization neither a list nor an object literal with arguments
			if srcReg.IsConstant() {
				c.ctx.Emitter.EmitLoadConst(seq[i], srcReg)
			} else {
				c.ctx.Emitter.EmitMove(seq[i], srcReg)
			}
			c.ctx.Types.Set(seq[i], operandType(c.ctx, srcReg))
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
func (c *ExprCompiler) CompileRangeOperator(ctx fql.IRangeOperatorContext) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	start := c.compileRangeOperand(ctx.GetLeft())
	end := c.compileRangeOperand(ctx.GetRight())

	span := source.Span{Start: -1, End: -1}

	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = diagnostics.SpanFromRuleContext(prc)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitRange(dst, start, end)
	})

	c.ctx.Types.Set(dst, core.TypeList)

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
func (c *ExprCompiler) compileRangeOperand(ctx fql.IRangeOperandContext) bytecode.Operand {
	if v := ctx.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.CompileParam(p)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.CompileMemberExpression(me)
	}

	if ice := ctx.ImplicitCurrentExpression(); ice != nil {
		return c.CompileImplicitCurrentExpression(ice)
	}

	if ime := ctx.ImplicitMemberExpression(); ime != nil {
		return c.CompileImplicitMemberExpression(ime)
	}

	if fc := ctx.FunctionCallExpression(); fc != nil {
		return c.CompileFunctionCallExpression(fc)
	}

	return bytecode.NoopOperand
}
