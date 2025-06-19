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

func (lc *LoopCompiler) Compile(ctx fql.IForExpressionContext) vm.Operand {
	returnRuleCtx := lc.compileInitialization(ctx)

	// body
	if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
		for _, b := range body {
			if c := b.ForExpressionStatement(); c != nil {
				lc.compileForExpressionStatement(c)
			} else if c := b.ForExpressionClause(); c != nil {
				lc.compileForExpressionClause(c)
			}
		}
	}

	return lc.compileFinalization(returnRuleCtx)
}

func (lc *LoopCompiler) compileInitialization(ctx fql.IForExpressionContext) antlr.RuleContext {
	var distinct bool
	var returnRuleCtx antlr.RuleContext
	var loopType core.LoopType
	returnCtx := ctx.ForExpressionReturn()

	if c := returnCtx.ReturnExpression(); c != nil {
		returnRuleCtx = c
		distinct = c.Distinct() != nil
		loopType = core.NormalLoop
	} else if c := returnCtx.ForExpression(); c != nil {
		returnRuleCtx = c
		loopType = core.PassThroughLoop
	}

	loop := lc.ctx.Loops.Create(loopType, core.ForLoop, distinct)
	lc.ctx.Symbols.EnterScope()

	if loop.Kind == core.ForLoop {
		loop.Src = lc.compileForExpressionSource(ctx.ForExpressionSource())

		if val := ctx.GetValueVariable(); val != nil {
			loop.DeclareValueVar(val.GetText(), lc.ctx.Symbols)
		}

		if ctr := ctx.GetCounterVariable(); ctr != nil {
			loop.DeclareKeyVar(ctr.GetText(), lc.ctx.Symbols)
		}
	} else {
	}

	loop.EmitInitialization(lc.ctx.Registers, lc.ctx.Emitter)

	return returnRuleCtx
}

func (lc *LoopCompiler) compileFinalization(ctx antlr.RuleContext) vm.Operand {
	loop := lc.ctx.Loops.Current()

	// RETURN
	if loop.Type != core.PassThroughLoop {
		c := ctx.(*fql.ReturnExpressionContext)
		expReg := lc.ctx.ExprCompiler.Compile(c.Expression())

		lc.ctx.Emitter.EmitAB(vm.OpPush, loop.Result, expReg)
	} else if ctx != nil {
		if c, ok := ctx.(*fql.ForExpressionContext); ok {
			lc.Compile(c)
		}
	}

	loop.EmitFinalization(lc.ctx.Emitter)
	lc.ctx.Symbols.ExitScope()
	lc.ctx.Loops.Pop()

	// TODO: Free operands

	return loop.Result
}

func (lc *LoopCompiler) compileForExpressionSource(ctx fql.IForExpressionSourceContext) vm.Operand {
	if c := ctx.FunctionCallExpression(); c != nil {
		return lc.ctx.ExprCompiler.CompileFunctionCallExpression(c)
	}

	if c := ctx.MemberExpression(); c != nil {
		return lc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	if c := ctx.Variable(); c != nil {
		return lc.ctx.ExprCompiler.CompileVariable(c)
	}

	if c := ctx.Param(); c != nil {
		return lc.ctx.ExprCompiler.CompileParam(c)
	}

	if c := ctx.RangeOperator(); c != nil {
		return lc.ctx.ExprCompiler.CompileRangeOperator(c)
	}

	if c := ctx.ArrayLiteral(); c != nil {
		return lc.ctx.LiteralCompiler.CompileArrayLiteral(c)
	}

	if c := ctx.ObjectLiteral(); c != nil {
		return lc.ctx.LiteralCompiler.CompileObjectLiteral(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (lc *LoopCompiler) compileForExpressionStatement(ctx fql.IForExpressionStatementContext) {
	if c := ctx.VariableDeclaration(); c != nil {
		_ = lc.ctx.StmtCompiler.CompileVariableDeclaration(c)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		_ = lc.ctx.ExprCompiler.CompileFunctionCallExpression(c)
	}

	// TODO: Free register if needed
}

func (lc *LoopCompiler) compileForExpressionClause(ctx fql.IForExpressionClauseContext) {
	if c := ctx.LimitClause(); c != nil {
		lc.compileLimitClause(c)
	} else if c := ctx.FilterClause(); c != nil {
		lc.compileFilterClause(c)
	} else if c := ctx.SortClause(); c != nil {
		lc.compileSortClause(c)
	} else if c := ctx.CollectClause(); c != nil {
		lc.compileCollectClause(c)
	}
}

func (lc *LoopCompiler) compileLimitClause(ctx fql.ILimitClauseContext) {
	clauses := ctx.AllLimitClauseValue()

	if len(clauses) == 1 {
		lc.compileLimit(lc.compileLimitClauseValue(clauses[0]))
	} else {
		lc.compileOffset(lc.compileLimitClauseValue(clauses[0]))
		lc.compileLimit(lc.compileLimitClauseValue(clauses[1]))
	}
}

func (lc *LoopCompiler) compileLimitClauseValue(ctx fql.ILimitClauseValueContext) vm.Operand {
	if c := ctx.Param(); c != nil {
		return lc.ctx.ExprCompiler.CompileParam(c)
	}

	if c := ctx.IntegerLiteral(); c != nil {
		return lc.ctx.LiteralCompiler.CompileIntegerLiteral(c)
	}

	if c := ctx.Variable(); c != nil {
		return lc.ctx.ExprCompiler.CompileVariable(c)
	}

	if c := ctx.MemberExpression(); c != nil {
		return lc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return lc.ctx.ExprCompiler.CompileFunctionCallExpression(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))

}

func (lc *LoopCompiler) compileLimit(src vm.Operand) {
	state := lc.ctx.Registers.Allocate(core.State)
	lc.ctx.Emitter.EmitABx(vm.OpIterLimit, state, src, lc.ctx.Loops.Current().Jump)
}

func (lc *LoopCompiler) compileOffset(src vm.Operand) {
	state := lc.ctx.Registers.Allocate(core.State)
	lc.ctx.Emitter.EmitABx(vm.OpIterSkip, state, src, lc.ctx.Loops.Current().Jump)
}

func (lc *LoopCompiler) compileFilterClause(ctx fql.IFilterClauseContext) {
	src := lc.ctx.ExprCompiler.Compile(ctx.Expression())
	lc.ctx.Emitter.EmitJumpIfFalse(src, lc.ctx.Loops.Current().Jump)
}

func (lc *LoopCompiler) compileSortClause(ctx fql.ISortClauseContext) {
	lc.ctx.LoopSortCompiler.Compile(ctx)
}

func (lc *LoopCompiler) compileCollectClause(ctx fql.ICollectClauseContext) {
	lc.ctx.LoopCollectCompiler.Compile(ctx)
}
