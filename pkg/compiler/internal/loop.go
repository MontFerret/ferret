package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopCompiler struct {
	ctx *CompilerContext
}

func NewLoopCompiler(ctx *CompilerContext) *LoopCompiler {
	return &LoopCompiler{ctx: ctx}
}

func (c *LoopCompiler) Compile(ctx fql.IForExpressionContext) vm.Operand {
	returnRuleCtx := c.compileInitialization(ctx)

	// body
	if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
		for _, b := range body {
			if ec := b.ForExpressionStatement(); ec != nil {
				c.compileForExpressionStatement(ec)
			} else if ec := b.ForExpressionClause(); ec != nil {
				c.compileForExpressionClause(ec)
			}
		}
	}

	return c.compileFinalization(returnRuleCtx)
}

func (c *LoopCompiler) compileInitialization(ctx fql.IForExpressionContext) antlr.RuleContext {
	var distinct bool
	var returnRuleCtx antlr.RuleContext
	var loopType core.LoopType
	returnCtx := ctx.ForExpressionReturn()

	if re := returnCtx.ReturnExpression(); re != nil {
		returnRuleCtx = re
		distinct = re.Distinct() != nil
		loopType = core.NormalLoop
	} else if fe := returnCtx.ForExpression(); fe != nil {
		returnRuleCtx = fe
		loopType = core.PassThroughLoop
	}

	src := c.compileForExpressionSource(ctx.ForExpressionSource())
	loop := c.ctx.Loops.CreateFor(loopType, src, distinct)
	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	if val := ctx.GetValueVariable(); val != nil {
		loop.DeclareValueVar(val.GetText(), c.ctx.Symbols)
	}

	if ctr := ctx.GetCounterVariable(); ctr != nil {
		loop.DeclareKeyVar(ctr.GetText(), c.ctx.Symbols)
	}

	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)

	if !loop.Allocate {
		// If the current loop must push distinct items, we must patch the dest dataset
		if loop.Distinct {
			parent := c.ctx.Loops.FindParent(c.ctx.Loops.Depth())

			if parent == nil {
				panic("parent loop not found in loop table")
			}

			c.ctx.Emitter.Patchx(parent.Pos, 1)
		}
	}

	return returnRuleCtx
}

func (c *LoopCompiler) compileFinalization(ctx antlr.RuleContext) vm.Operand {
	loop := c.ctx.Loops.Current()

	// RETURN
	if loop.Type != core.PassThroughLoop {
		re := ctx.(*fql.ReturnExpressionContext)
		expReg := c.ctx.ExprCompiler.Compile(re.Expression())

		c.ctx.Emitter.EmitAB(vm.OpPush, loop.Dst, expReg)
	} else if ctx != nil {
		if fe, ok := ctx.(*fql.ForExpressionContext); ok {
			c.Compile(fe)
		}
	}

	loop.EmitFinalization(c.ctx.Emitter)
	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	// TODO: Free operands

	return loop.Dst
}

func (c *LoopCompiler) compileForExpressionSource(ctx fql.IForExpressionSourceContext) vm.Operand {
	if fce := ctx.FunctionCallExpression(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.ctx.ExprCompiler.CompileParam(p)
	}

	if ro := ctx.RangeOperator(); ro != nil {
		return c.ctx.ExprCompiler.CompileRangeOperator(ro)
	}

	if al := ctx.ArrayLiteral(); al != nil {
		return c.ctx.LiteralCompiler.CompileArrayLiteral(al)
	}

	if ol := ctx.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (c *LoopCompiler) compileForExpressionStatement(ctx fql.IForExpressionStatementContext) {
	if vd := ctx.VariableDeclaration(); vd != nil {
		_ = c.ctx.StmtCompiler.CompileVariableDeclaration(vd)
	} else if fce := ctx.FunctionCallExpression(); fce != nil {
		_ = c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
	}

	// TODO: Free register if needed
}

func (c *LoopCompiler) compileForExpressionClause(ctx fql.IForExpressionClauseContext) {
	if lc := ctx.LimitClause(); lc != nil {
		c.compileLimitClause(lc)
	} else if fc := ctx.FilterClause(); fc != nil {
		c.compileFilterClause(fc)
	} else if sc := ctx.SortClause(); sc != nil {
		c.compileSortClause(sc)
	} else if cc := ctx.CollectClause(); cc != nil {
		c.compileCollectClause(cc)
	}
}

func (c *LoopCompiler) compileLimitClause(ctx fql.ILimitClauseContext) {
	clauses := ctx.AllLimitClauseValue()

	if len(clauses) == 1 {
		c.compileLimit(c.compileLimitClauseValue(clauses[0]))
	} else {
		c.compileOffset(c.compileLimitClauseValue(clauses[0]))
		c.compileLimit(c.compileLimitClauseValue(clauses[1]))
	}
}

func (c *LoopCompiler) compileLimitClauseValue(ctx fql.ILimitClauseValueContext) vm.Operand {
	if pm := ctx.Param(); pm != nil {
		return c.ctx.ExprCompiler.CompileParam(pm)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	if vb := ctx.Variable(); vb != nil {
		return c.ctx.ExprCompiler.CompileVariable(vb)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	if fce := ctx.FunctionCallExpression(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))

}

func (c *LoopCompiler) compileLimit(src vm.Operand) {
	state := c.ctx.Registers.Allocate(core.State)
	c.ctx.Emitter.EmitABx(vm.OpIterLimit, state, src, c.ctx.Loops.Current().Jump)
}

func (c *LoopCompiler) compileOffset(src vm.Operand) {
	state := c.ctx.Registers.Allocate(core.State)
	c.ctx.Emitter.EmitABx(vm.OpIterSkip, state, src, c.ctx.Loops.Current().Jump)
}

func (c *LoopCompiler) compileFilterClause(ctx fql.IFilterClauseContext) {
	src := c.ctx.ExprCompiler.Compile(ctx.Expression())
	c.ctx.Emitter.EmitJumpIfFalse(src, c.ctx.Loops.Current().Jump)
}

func (c *LoopCompiler) compileSortClause(ctx fql.ISortClauseContext) {
	c.ctx.LoopSortCompiler.Compile(ctx)
}

func (c *LoopCompiler) compileCollectClause(ctx fql.ICollectClauseContext) {
	c.ctx.LoopCollectCompiler.Compile(ctx)
}
