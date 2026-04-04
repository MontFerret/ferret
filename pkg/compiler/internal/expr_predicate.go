package internal

import (
	"regexp"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

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
	return compileScalarLiteralOperand(c.ctx, c.literals, lit)
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

		reg := c.recovery.CompileWithErrorPolicy(core.ErrorPolicySuppress, jumpMode, func() bytecode.Operand {
			return c.compileAtom(atom)
		})

		return reg, true
	}

	if wfe := atom.WaitForExpression(); wfe != nil {
		outerPlan := c.recovery.CollectPlan(atom, core.RecoveryPlanOptions{
			AllowTimeout: true,
			HasTimeout:   waitForHasExplicitTimeoutClause(wfe),
		})

		return c.wait.CompileWithOuterRecoveryPlan(wfe, outerPlan), true
	}

	if fe := atom.ForExpression(); fe != nil {
		outerPlan := c.recovery.CollectPlan(atom, core.RecoveryPlanOptions{})
		return c.loops.CompileWithOuterRecoveryPlan(fe, outerPlan), true
	}

	reg := c.recovery.CompileOperation(OperationRecoverySpec{
		Owner:    atom,
		JumpMode: core.CatchJumpModeNone,
		CompilePlain: func() bytecode.Operand {
			return c.compileAtom(atom)
		},
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

func (c *ExprCompiler) emitBinaryPredicate(ctx fql.IPredicateContext, opcode bytecode.Opcode, left, right bytecode.Operand, isNegated bool) bytecode.Operand {
	dest := c.ctx.Function.Registers.Allocate()
	span := source.Span{Start: -1, End: -1}

	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = diagnostics.SpanFromRuleContext(prc)
	}

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitABC(opcode, dest, left, right)

		if isNegated {
			c.ctx.Program.Emitter.EmitAB(bytecode.OpNot, dest, dest)
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
	c.ctx.Program.Emitter.EmitJumpCompare(opcode, left, right, label)
}

// EmitConditionJump is the supported cross-compiler entrypoint for branching on
// expression truthiness while preserving expression-level predicate lowering.
func (c *ExprCompiler) EmitConditionJump(expr fql.IExpressionContext, label core.Label, jumpOnTrue bool) {
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
		c.ctx.Program.Emitter.EmitJumpIfTrue(cond, label)
	} else {
		c.ctx.Program.Emitter.EmitJumpIfFalse(cond, label)
	}
}

func (c *ExprCompiler) compileAtom(ctx fql.IExpressionAtomContext) bytecode.Operand {
	if op, ok := resolveAtomBinaryOperator(ctx); ok {
		return c.compileBinaryAtom(ctx, op)
	}

	return c.compileLeafAtom(ctx)
}

func (c *ExprCompiler) compileBinaryAtom(ctx fql.IExpressionAtomContext, op atomBinaryOperator) bytecode.Operand {
	if op.opcode == bytecode.OpAdd {
		if parts, ok := buildConcatOperandSegmentsFromAtom(c, ctx); ok {
			return emitConcatOperandSegments(c.ctx, c.facts, parts)
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

	return emitBinaryOperation(c.ctx, c.facts, prc, op, left, right)
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
			c.ctx.Program.Errors.InvalidRegexExpression(ctx, exp.String())
		}

		return
	}

	c.ctx.Program.Errors.InvalidRegexExpression(ctx, lit.GetText())
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
		return c.literals.Compile(l)
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
		return c.dispatch.Compile(de)
	}

	if fe := ctx.ForExpression(); fe != nil {
		return c.loops.Compile(fe)
	}

	if wfe := ctx.WaitForExpression(); wfe != nil {
		return c.wait.Compile(wfe)
	}

	if e := ctx.Expression(); e != nil {
		return c.Compile(e)
	}

	return bytecode.NoopOperand
}
