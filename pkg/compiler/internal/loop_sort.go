package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopSortCompiler struct {
	ctx *CompilerContext
}

func NewLoopSortCompiler(ctx *CompilerContext) *LoopSortCompiler {
	return &LoopSortCompiler{ctx: ctx}
}

func (lc *LoopSortCompiler) Compile(ctx fql.ISortClauseContext) {
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
	loop.EmitInitialization(lc.ctx.Registers, lc.ctx.Emitter)
}
