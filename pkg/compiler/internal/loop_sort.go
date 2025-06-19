package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// LoopSortCompiler handles compilation of SORT clauses within loops.
// It transforms loop elements into KeyValuePairs where keys are sort expressions
// and values are the original loop elements.
type LoopSortCompiler struct {
	ctx *CompilerContext
}

func NewLoopSortCompiler(ctx *CompilerContext) *LoopSortCompiler {
	return &LoopSortCompiler{ctx: ctx}
}

// Compile processes a SORT clause by:
// 1. Extracting sort expressions and directions
// 2. Creating KeyValuePairs for sorting
// 3. Patching the loop with appropriate sorter operations
// 4. Reinitializing the loop with sorted data
func (lc *LoopSortCompiler) Compile(ctx fql.ISortClauseContext) {
	loop := lc.ctx.Loops.Current()
	clauses := ctx.AllSortClauseExpression()

	// Compile sort keys and get sort directions
	kvKeyReg, directions := lc.compileSortKeys(clauses)

	// Handle the value part of KeyValuePair
	kvValReg := lc.resolveValueRegister(loop)

	// Apply the appropriate sorter based on number of sort conditions
	lc.applySorter(loop, clauses, directions)

	// Emit the KeyValuePair and finalize the sorting process
	lc.finalizeSorting(loop, kvKeyReg, kvValReg)
}

// compileSortKeys processes all sort expressions and returns the key register and directions.
// For multiple expressions, it creates an array of keys; for single expression, uses the key directly.
func (lc *LoopSortCompiler) compileSortKeys(clauses []fql.ISortClauseExpressionContext) (vm.Operand, []runtime.SortDirection) {
	kvKeyReg := lc.ctx.Registers.Allocate(core.Temp)
	directions := make([]runtime.SortDirection, len(clauses))
	isSortMany := len(clauses) > 1

	if isSortMany {
		return lc.compileMultipleSortKeys(clauses, kvKeyReg, directions)
	}

	return lc.compileSingleSortKey(clauses[0], kvKeyReg, directions)
}

// compileMultipleSortKeys handles compilation when there are multiple sort expressions.
// It creates an array of compiled expressions for multi-key sorting.
func (lc *LoopSortCompiler) compileMultipleSortKeys(clauses []fql.ISortClauseExpressionContext, kvKeyReg vm.Operand, directions []runtime.SortDirection) (vm.Operand, []runtime.SortDirection) {
	clausesRegs := make([]vm.Operand, len(clauses))
	keyRegs := lc.ctx.Registers.AllocateSequence(len(clauses))

	// Compile each sort expression and store direction
	for i, clause := range clauses {
		clauseReg := lc.ctx.ExprCompiler.Compile(clause.Expression())
		lc.ctx.Emitter.EmitMove(keyRegs[i], clauseReg)
		clausesRegs[i] = keyRegs[i]
		directions[i] = sortDirection(clause.SortDirection())
		// TODO: Free registers after use
	}

	// CreateFor array of sort keys
	arrReg := lc.ctx.Registers.Allocate(core.Temp)
	lc.ctx.Emitter.EmitAs(vm.OpList, arrReg, keyRegs)
	lc.ctx.Emitter.EmitAB(vm.OpMove, kvKeyReg, arrReg)
	// TODO: Free registers after use

	return kvKeyReg, directions
}

// compileSingleSortKey handles compilation when there is only one sort expression.
func (lc *LoopSortCompiler) compileSingleSortKey(clause fql.ISortClauseExpressionContext, kvKeyReg vm.Operand, directions []runtime.SortDirection) (vm.Operand, []runtime.SortDirection) {
	clauseReg := lc.ctx.ExprCompiler.Compile(clause.Expression())
	lc.ctx.Emitter.EmitAB(vm.OpMove, kvKeyReg, clauseReg)
	directions[0] = sortDirection(clause.SortDirection())

	return kvKeyReg, directions
}

// resolveValueRegister determines the appropriate register for the value part of KeyValuePair.
// If the loop already has a value name, reuse it; otherwise, allocate a new register
// and load the value from the iterator.
func (lc *LoopSortCompiler) resolveValueRegister(loop *core.Loop) vm.Operand {
	// If value is already used in the loop body, reuse the existing register
	if loop.ValueName != "" {
		return loop.Value
	}

	// Otherwise, allocate a new register and load the value from iterator
	kvValReg := lc.ctx.Registers.Allocate(core.Temp)
	loop.EmitValue(kvValReg, lc.ctx.Emitter)
	return kvValReg
}

// applySorter patches the loop with the appropriate sorter operation based on
// whether we have single or multiple sort conditions.
func (lc *LoopSortCompiler) applySorter(loop *core.Loop, clauses []fql.ISortClauseExpressionContext, directions []runtime.SortDirection) {
	isSortMany := len(clauses) > 1

	if isSortMany {
		// Multi-key sorting requires encoded directions and count
		encoded := runtime.EncodeSortDirections(directions)
		count := len(clauses)
		lc.ctx.Emitter.PatchSwapAxy(loop.DstPos, vm.OpDataSetMultiSorter, loop.Dst, encoded, count)
	} else {
		// Single-key sorting only needs the direction
		dir := sortDirection(clauses[0].SortDirection())
		lc.ctx.Emitter.PatchSwapAx(loop.DstPos, vm.OpDataSetSorter, loop.Dst, int(dir))
	}
}

// finalizeSorting completes the sorting process by:
// 1. Adding KeyValuePairs to the result dataset
// 2. Finalizing the current loop
// 3. Replacing the loop source with sorted results
// 4. Reinitializing the loop for iteration over sorted data
func (lc *LoopSortCompiler) finalizeSorting(loop *core.Loop, kvKeyReg, kvValReg vm.Operand) {
	// Add the KeyValuePair to the dataset
	lc.ctx.Emitter.EmitABC(vm.OpPushKV, loop.Dst, kvKeyReg, kvValReg)

	// Finalize the current loop iteration
	loop.EmitFinalization(lc.ctx.Emitter)

	// Replace the loop source with sorted results
	lc.ctx.Emitter.EmitAB(vm.OpMove, loop.Src, loop.Dst)

	// Reinitialize the loop to iterate over sorted data
	loop.EmitInitialization(lc.ctx.Registers, lc.ctx.Emitter)
}
