package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
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
		ctx                  *CompilationSession
		front                *CompilationFrontend
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
func NewExprCompiler(ctx *CompilationSession) *ExprCompiler {
	return &ExprCompiler{ctx: ctx}
}

// Compile processes an expression from the FQL AST and delegates to the appropriate
// compilation method based on the expression type (unary, logical, ternary, or predicate).
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

	c.ctx.Emitter.EmitAB(op, dst, src)

	return dst
}

func (c *ExprCompiler) compileLogicalAnd(ctx fql.IExpressionContext) bytecode.Operand {
	left := c.Compile(ctx.GetLeft())

	skip := c.ctx.Emitter.NewLabel("and.false")
	done := c.ctx.Emitter.NewLabel("and.done")
	dst := c.ctx.Registers.Allocate()

	c.ctx.Emitter.EmitJumpIfFalse(left, skip)

	right := c.Compile(ctx.GetRight())
	c.ctx.Emitter.EmitMove(dst, right)
	c.ctx.Emitter.EmitJump(done)

	c.ctx.Emitter.MarkLabel(skip)
	c.ctx.Emitter.EmitMove(dst, left)

	c.ctx.Emitter.MarkLabel(done)

	return dst
}

func (c *ExprCompiler) compileLogicalOr(ctx fql.IExpressionContext) bytecode.Operand {
	left := c.Compile(ctx.GetLeft())

	next := c.ctx.Emitter.NewLabel("or.false")
	done := c.ctx.Emitter.NewLabel("or.done")
	dst := c.ctx.Registers.Allocate()

	c.ctx.Emitter.EmitJumpIfTrue(left, next)

	right := c.Compile(ctx.GetRight())
	c.ctx.Emitter.EmitMove(dst, right)
	c.ctx.Emitter.EmitJump(done)

	c.ctx.Emitter.MarkLabel(next)
	c.ctx.Emitter.EmitMove(dst, left)

	c.ctx.Emitter.MarkLabel(done)

	return dst
}

func (c *ExprCompiler) compileTernary(ctx fql.IExpressionContext) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()

	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	onTrue := ctx.GetOnTrue()
	onFalse := ctx.GetOnFalse()
	cond := ctx.GetCondition()

	if onTrue == nil && cond != nil {
		condReg := c.Compile(cond)
		c.ctx.Emitter.EmitMove(dst, condReg)
		c.ctx.Emitter.EmitJumpIfFalse(condReg, elseLabel)
	} else if cond != nil {
		c.emitConditionJump(cond, elseLabel, false)
	}

	if onTrue != nil {
		trueReg := c.Compile(onTrue)
		c.ctx.Emitter.EmitMove(dst, trueReg)
	}

	c.ctx.Emitter.EmitJump(endLabel)
	c.ctx.Emitter.MarkLabel(elseLabel)

	if onFalse != nil {
		falseReg := c.Compile(onFalse)
		c.ctx.Emitter.EmitMove(dst, falseReg)
	}

	c.ctx.Emitter.MarkLabel(endLabel)

	return dst
}
