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
		bindings             *BindingCompiler
		calls                *CallResolver
		callCompiler         *exprCallCompiler
		dispatch             *DispatchCompiler
		literals             *LiteralCompiler
		loops                *LoopCompiler
		matchCompiler        *exprMatchCompiler
		queryCompiler        *exprQueryCompiler
		recovery             *RecoveryCompiler
		facts                *TypeFacts
		wait                 *WaitCompiler
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
	c := &ExprCompiler{ctx: ctx}
	c.callCompiler = newExprCallCompiler(ctx, exprCallCallbacks{
		compileExpr:            c.Compile,
		compileMember:          c.CompileMemberExpression,
		compileImplicitCurrent: c.CompileImplicitCurrentExpression,
		compileImplicitMember:  c.CompileImplicitMemberExpression,
	})
	c.queryCompiler = newExprQueryCompiler(ctx, exprQueryCallbacks{
		compileExpr:     c.Compile,
		compileParam:    c.CompileParam,
		compileVariable: c.CompileVariable,
	})
	c.matchCompiler = newExprMatchCompiler(ctx, exprMatchCallbacks{
		compileExpr:                   c.Compile,
		emitConditionJump:             c.EmitConditionJump,
		compileFunctionCallByNameWith: c.CompileFunctionCallByNameWith,
	})

	return c
}

func (c *ExprCompiler) bind(
	bindings *BindingCompiler,
	calls *CallResolver,
	dispatch *DispatchCompiler,
	literals *LiteralCompiler,
	loops *LoopCompiler,
	recovery *RecoveryCompiler,
	facts *TypeFacts,
	wait *WaitCompiler,
) {
	if c == nil {
		return
	}

	c.bindings = bindings
	c.calls = calls
	c.dispatch = dispatch
	c.literals = literals
	c.loops = loops
	c.recovery = recovery
	c.facts = facts
	c.wait = wait

	c.callCompiler.bind(bindings, calls, literals, recovery, facts)
	c.queryCompiler.bind(literals, recovery, facts)
	c.matchCompiler.bind(literals, facts)
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
		core.PanicInvariant("cannot increment/decrement a constant")
	}

	operator := token.GetText()

	switch operator {
	case "++":
		c.ctx.Program.Emitter.EmitA(bytecode.OpIncr, target)
	case "--":
		c.ctx.Program.Emitter.EmitA(bytecode.OpDecr, target)
	default:
		c.ctx.Program.Errors.InvalidToken(token)

		return bytecode.NoopOperand
	}

	return target
}

func (c *ExprCompiler) compileUnary(ctx fql.IUnaryOperatorContext, parent fql.IExpressionContext) bytecode.Operand {
	src := c.Compile(parent.GetRight())
	dst := c.ctx.Function.Registers.Allocate()

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

	c.ctx.Program.Emitter.EmitAB(op, dst, src)

	return dst
}

func (c *ExprCompiler) compileLogicalAnd(ctx fql.IExpressionContext) bytecode.Operand {
	left := c.Compile(ctx.GetLeft())

	skip := c.ctx.Program.Emitter.NewLabel("and.false")
	done := c.ctx.Program.Emitter.NewLabel("and.done")
	dst := c.ctx.Function.Registers.Allocate()

	c.ctx.Program.Emitter.EmitJumpIfFalse(left, skip)

	right := c.Compile(ctx.GetRight())
	c.ctx.Program.Emitter.EmitMove(dst, right)
	c.ctx.Program.Emitter.EmitJump(done)

	c.ctx.Program.Emitter.MarkLabel(skip)
	c.ctx.Program.Emitter.EmitMove(dst, left)

	c.ctx.Program.Emitter.MarkLabel(done)

	return dst
}

func (c *ExprCompiler) compileLogicalOr(ctx fql.IExpressionContext) bytecode.Operand {
	left := c.Compile(ctx.GetLeft())

	next := c.ctx.Program.Emitter.NewLabel("or.false")
	done := c.ctx.Program.Emitter.NewLabel("or.done")
	dst := c.ctx.Function.Registers.Allocate()

	c.ctx.Program.Emitter.EmitJumpIfTrue(left, next)

	right := c.Compile(ctx.GetRight())
	c.ctx.Program.Emitter.EmitMove(dst, right)
	c.ctx.Program.Emitter.EmitJump(done)

	c.ctx.Program.Emitter.MarkLabel(next)
	c.ctx.Program.Emitter.EmitMove(dst, left)

	c.ctx.Program.Emitter.MarkLabel(done)

	return dst
}

func (c *ExprCompiler) compileTernary(ctx fql.IExpressionContext) bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()

	elseLabel := c.ctx.Program.Emitter.NewLabel()
	endLabel := c.ctx.Program.Emitter.NewLabel()

	onTrue := ctx.GetOnTrue()
	onFalse := ctx.GetOnFalse()
	cond := ctx.GetCondition()

	if onTrue == nil && cond != nil {
		condReg := c.Compile(cond)
		c.ctx.Program.Emitter.EmitMove(dst, condReg)
		c.ctx.Program.Emitter.EmitJumpIfFalse(condReg, elseLabel)
	} else if cond != nil {
		c.EmitConditionJump(cond, elseLabel, false)
	}

	if onTrue != nil {
		trueReg := c.Compile(onTrue)
		c.ctx.Program.Emitter.EmitMove(dst, trueReg)
	}

	c.ctx.Program.Emitter.EmitJump(endLabel)
	c.ctx.Program.Emitter.MarkLabel(elseLabel)

	if onFalse != nil {
		falseReg := c.Compile(onFalse)
		c.ctx.Program.Emitter.EmitMove(dst, falseReg)
	}

	c.ctx.Program.Emitter.MarkLabel(endLabel)

	return dst
}
