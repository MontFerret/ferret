package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopCompiler struct {
	ctx *FuncContext
}

func NewLoopCompiler(ctx *FuncContext) *LoopCompiler {
	return &LoopCompiler{ctx: ctx}
}

func (lc *LoopCompiler) Compile(ctx fql.IForExpressionContext) vm.Operand {
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

	loop := lc.ctx.Loops.NewLoop(loopType, core.ForLoop, distinct)
	lc.ctx.Symbols.EnterScope()
	lc.ctx.Loops.Push(loop)

	if loop.Kind == core.ForLoop {
		loop.Src = lc.CompileForExpressionSource(ctx.ForExpressionSource())

		if val := ctx.GetValueVariable(); val != nil {
			loop.DeclareValueVar(val.GetText(), lc.ctx.Symbols)
		}

		if ctr := ctx.GetCounterVariable(); ctr != nil {
			loop.DeclareKeyVar(ctr.GetText(), lc.ctx.Symbols)
		}
	} else {
	}

	lc.emitLoopBegin(loop)

	// body
	if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
		for _, b := range body {
			if c := b.ForExpressionStatement(); c != nil {
				lc.CompileForExpressionStatement(c)
			} else if c := b.ForExpressionClause(); c != nil {
				lc.CompileForExpressionClause(c)
			}
		}
	}

	loop = lc.ctx.Loops.Current()

	// RETURN
	if loop.Type != core.PassThroughLoop {
		c := returnRuleCtx.(*fql.ReturnExpressionContext)
		expReg := lc.ctx.ExprCompiler.Compile(c.Expression())

		lc.ctx.Emitter.EmitAB(vm.OpPush, loop.Result, expReg)
	} else if returnRuleCtx != nil {
		if c, ok := returnRuleCtx.(*fql.ForExpressionContext); ok {
			lc.Compile(c)
		}
	}

	res := lc.emitLoopEnd(loop)

	lc.ctx.Symbols.ExitScope()
	lc.ctx.Loops.Pop()

	return res
}

func (lc *LoopCompiler) CompileForExpressionSource(ctx fql.IForExpressionSourceContext) vm.Operand {
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

func (lc *LoopCompiler) CompileForExpressionStatement(ctx fql.IForExpressionStatementContext) {
	if c := ctx.VariableDeclaration(); c != nil {
		_ = lc.ctx.StmtCompiler.CompileVariableDeclaration(c)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		_ = lc.ctx.ExprCompiler.CompileFunctionCallExpression(c)

		// TODO: Free register if needed
	}
}

func (lc *LoopCompiler) CompileForExpressionClause(ctx fql.IForExpressionClauseContext) {
	if c := ctx.LimitClause(); c != nil {
		lc.CompileLimitClause(c)
	} else if c := ctx.FilterClause(); c != nil {
		lc.CompileFilterClause(c)
	} else if c := ctx.SortClause(); c != nil {
		lc.CompileSortClause(c)
	} else if c := ctx.CollectClause(); c != nil {
		lc.CompileCollectClause(c)
	}
}

func (lc *LoopCompiler) CompileLimitClause(ctx fql.ILimitClauseContext) {
	clauses := ctx.AllLimitClauseValue()

	if len(clauses) == 1 {
		lc.CompileLimit(lc.CompileLimitClauseValue(clauses[0]))
	} else {
		lc.CompileOffset(lc.CompileLimitClauseValue(clauses[0]))
		lc.CompileLimit(lc.CompileLimitClauseValue(clauses[1]))
	}
}

func (lc *LoopCompiler) CompileLimitClauseValue(ctx fql.ILimitClauseValueContext) vm.Operand {
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

func (lc *LoopCompiler) CompileLimit(src vm.Operand) {
	state := lc.ctx.Registers.Allocate(core.State)
	lc.ctx.Emitter.EmitABx(vm.OpIterLimit, state, src, lc.ctx.Loops.Current().Jump)
}

func (lc *LoopCompiler) CompileOffset(src vm.Operand) {
	state := lc.ctx.Registers.Allocate(core.State)
	lc.ctx.Emitter.EmitABx(vm.OpIterSkip, state, src, lc.ctx.Loops.Current().Jump)
}

func (lc *LoopCompiler) CompileFilterClause(ctx fql.IFilterClauseContext) {
	src := lc.ctx.ExprCompiler.Compile(ctx.Expression())
	lc.ctx.Emitter.EmitJumpIfFalse(src, lc.ctx.Loops.Current().Jump)
}

func (lc *LoopCompiler) CompileSortClause(ctx fql.ISortClauseContext) {
	loop := lc.ctx.Loops.Current()

	// We collect the sorting conditions (keys
	// And wrap each loop element by a KeyValuePair
	// Where a key is either a single value or a list of values
	// These KeyValuePairs are then added to the dataset
	kvKeyReg := lc.ctx.Registers.Allocate(core.Temp)
	clauses := ctx.AllSortClauseExpression()
	var directions []runtime.SortDirection
	isSortMany := len(clauses) > 1

	if isSortMany {
		clausesRegs := make([]vm.Operand, len(clauses))
		directions = make([]runtime.SortDirection, len(clauses))
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		keyRegs := lc.ctx.Registers.AllocateSequence(len(clauses))

		for i, clause := range clauses {
			clauseReg := lc.ctx.ExprCompiler.Compile(clause.Expression())
			lc.ctx.Emitter.EmitMove(keyRegs[i], clauseReg)
			clausesRegs[i] = keyRegs[i]
			directions[i] = sortDirection(clause.SortDirection())
			// TODO: Free Registers
		}

		arrReg := lc.ctx.Registers.Allocate(core.Temp)
		lc.ctx.Emitter.EmitAs(vm.OpList, arrReg, keyRegs)
		lc.ctx.Emitter.EmitAB(vm.OpMove, kvKeyReg, arrReg) // TODO: Free Registers
	} else {
		clausesReg := lc.ctx.ExprCompiler.Compile(clauses[0].Expression())
		lc.ctx.Emitter.EmitAB(vm.OpMove, kvKeyReg, clausesReg)
	}

	var kvValReg vm.Operand

	// In case the value is not used in the loop body, and only key is used
	if loop.ValueName != "" {
		kvValReg = loop.Value
	} else {
		// If so, we need to load it from the iterator
		kvValReg = lc.ctx.Registers.Allocate(core.Temp)
		loop.EmitValue(kvKeyReg, lc.ctx.Emitter)
	}

	if isSortMany {
		encoded := runtime.EncodeSortDirections(directions)
		count := len(clauses)

		lc.ctx.Emitter.PatchSwapAxy(loop.ResultPos, vm.OpDataSetMultiSorter, loop.Result, encoded, count)
	} else {
		dir := sortDirection(clauses[0].SortDirection())
		lc.ctx.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetSorter, loop.Result, int(dir))
	}

	lc.ctx.Emitter.EmitABC(vm.OpPushKV, loop.Result, kvKeyReg, kvValReg)
	loop.EmitFinalization(lc.ctx.Emitter)

	// Replace source with the Sorter
	lc.ctx.Emitter.EmitAB(vm.OpMove, loop.Src, loop.Result)

	// Create a new loop
	lc.emitLoopBegin(loop)
}

func (lc *LoopCompiler) CompileCollectClause(ctx fql.ICollectClauseContext) {
	lc.ctx.CollectCompiler.Compile(ctx)
}

// emitIterValue emits an instruction to get the value from the iterator
func (lc *LoopCompiler) emitLoopBegin(loop *core.Loop) {
	if loop.Allocate {
		lc.ctx.Emitter.EmitAb(vm.OpDataSet, loop.Result, loop.Distinct)
		loop.ResultPos = lc.ctx.Emitter.Size() - 1
	}

	loop.Iterator = lc.ctx.Registers.Allocate(core.State)

	if loop.Kind == core.ForLoop {
		lc.ctx.Emitter.EmitAB(vm.OpIter, loop.Iterator, loop.Src)
		// core.JumpPlaceholder is a placeholder for the exit jump position
		loop.Jump = lc.ctx.Emitter.EmitJumpc(vm.OpIterNext, core.JumpPlaceholder, loop.Iterator)

		if loop.Value != vm.NoopOperand {
			lc.ctx.Emitter.EmitAB(vm.OpIterValue, loop.Value, loop.Iterator)
		}

		if loop.Key != vm.NoopOperand {
			lc.ctx.Emitter.EmitAB(vm.OpIterKey, loop.Key, loop.Iterator)
		}
	} else {
		//counterReg := lc.ctx.Registers.Allocate(Storage)
		// TODO: Set JumpOffset here
	}
}

// emitPatchLoop replaces the source of the loop with a modified dataset
func (lc *LoopCompiler) emitPatchLoop(loop *core.Loop) {
	// Replace source with sorted array
	lc.ctx.Emitter.EmitAB(vm.OpMove, loop.Src, loop.Result)

	lc.ctx.Symbols.ExitScope()
	lc.ctx.Symbols.EnterScope()

	// Create new for loop
	lc.emitLoopBegin(loop)
}

func (lc *LoopCompiler) emitLoopEnd(loop *core.Loop) vm.Operand {
	lc.ctx.Emitter.EmitJump(loop.Jump - loop.JumpOffset)

	// TODO: Do not allocate for pass-through Loops
	dst := lc.ctx.Registers.Allocate(core.Temp)

	if loop.Allocate {
		// TODO: Reuse the dsReg register
		lc.ctx.Emitter.EmitA(vm.OpClose, loop.Iterator)
		lc.ctx.Emitter.EmitAB(vm.OpMove, dst, loop.Result)

		if loop.Kind == core.ForLoop {
			lc.ctx.Emitter.PatchJump(loop.Jump)
		} else {
			lc.ctx.Emitter.PatchJumpAB(loop.Jump)
		}
	} else {
		if loop.Kind == core.ForLoop {
			lc.ctx.Emitter.PatchJumpNext(loop.Jump)
		} else {
			lc.ctx.Emitter.PatchJumpNextAB(loop.Jump)
		}
	}

	return dst
}

func (lc *LoopCompiler) loopKind(ctx *fql.ForExpressionContext) core.LoopKind {
	if ctx.While() == nil {
		return core.ForLoop
	}

	if ctx.Do() == nil {
		return core.WhileLoop
	}

	return core.DoWhileLoop
}
