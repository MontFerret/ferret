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
func (c *LoopSortCompiler) Compile(ctx fql.ISortClauseContext) {
	loop := c.ctx.Loops.Current()
	clauses := ctx.AllSortClauseExpression()

	// Compile sort keys and get sort directions
	kvKeyReg, directions := c.compileSortKeys(clauses)

	// Handle the value part of KeyValuePair
	kvValReg := c.resolveValueRegister(loop)

	// Apply the appropriate sorter based on number of sort conditions
	sorterReg := c.compileSorter(loop, clauses, directions)

	// Emit the KeyValuePair and finalize the sorting process
	c.finalizeSorting(loop, core.NewKV(kvKeyReg, kvValReg), sorterReg)
}

// compileSortKeys processes all sort expressions and returns the key register and directions.
// For multiple expressions, it creates an array of keys; for single expression, uses the key directly.
func (c *LoopSortCompiler) compileSortKeys(clauses []fql.ISortClauseExpressionContext) (vm.Operand, []runtime.SortDirection) {
	kvKeyReg := c.ctx.Registers.Allocate(core.Temp)
	directions := make([]runtime.SortDirection, len(clauses))
	isSortMany := len(clauses) > 1

	if isSortMany {
		return c.compileMultipleSortKeys(clauses, kvKeyReg, directions)
	}

	return c.compileSingleSortKey(clauses[0], kvKeyReg, directions)
}

// compileMultipleSortKeys handles compilation when there are multiple sort expressions.
// It creates an array of compiled expressions for multi-key sorting.
func (c *LoopSortCompiler) compileMultipleSortKeys(clauses []fql.ISortClauseExpressionContext, kvKeyReg vm.Operand, directions []runtime.SortDirection) (vm.Operand, []runtime.SortDirection) {
	clausesRegs := make([]vm.Operand, len(clauses))
	keyRegs := c.ctx.Registers.AllocateSequence(len(clauses))

	// Compile each sort expression and store direction
	for i, clause := range clauses {
		clauseReg := c.ctx.ExprCompiler.Compile(clause.Expression())
		c.ctx.Emitter.EmitMove(keyRegs[i], clauseReg)
		clausesRegs[i] = keyRegs[i]
		directions[i] = sortDirection(clause.SortDirection())
		// TODO: Free registers after use
	}

	// CreateFor array of sort keys
	arrReg := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitAs(vm.OpList, arrReg, keyRegs)
	c.ctx.Emitter.EmitAB(vm.OpMove, kvKeyReg, arrReg)
	// TODO: Free registers after use

	return kvKeyReg, directions
}

// compileSingleSortKey handles compilation when there is only one sort expression.
func (c *LoopSortCompiler) compileSingleSortKey(clause fql.ISortClauseExpressionContext, kvKeyReg vm.Operand, directions []runtime.SortDirection) (vm.Operand, []runtime.SortDirection) {
	clauseReg := c.ctx.ExprCompiler.Compile(clause.Expression())
	c.ctx.Emitter.EmitAB(vm.OpMove, kvKeyReg, clauseReg)
	directions[0] = sortDirection(clause.SortDirection())

	return kvKeyReg, directions
}

// resolveValueRegister determines the appropriate register for the value part of KeyValuePair.
// If the loop already has a value name, reuse it; otherwise, allocate a new register
// and load the value from the iterator.
func (c *LoopSortCompiler) resolveValueRegister(loop *core.Loop) vm.Operand {
	// If value is already used in the loop body, reuse the existing register
	if loop.ValueName != "" {
		return loop.Value
	}

	// Otherwise, allocate a new register and load the value from iterator
	kvValReg := c.ctx.Registers.Allocate(core.Temp)
	loop.EmitValue(kvValReg, c.ctx.Emitter)
	return kvValReg
}

// compileSorter configures a sorter for a loop based on provided sort clauses and directions.
// It handles both single-key and multi-key sorting by emitting the appropriate VM operations.
func (c *LoopSortCompiler) compileSorter(loop *core.Loop, clauses []fql.ISortClauseExpressionContext, directions []runtime.SortDirection) vm.Operand {
	isSortMany := len(clauses) > 1

	if isSortMany {
		// Multi-key sorting requires encoded directions and count
		encoded := runtime.EncodeSortDirections(directions)
		count := len(clauses)

		if loop.Allocate {
			c.ctx.Emitter.PatchSwapAxy(loop.Pos, vm.OpDataSetMultiSorter, loop.Dst, encoded, count)

			return loop.Dst
		}

		dst := c.ctx.Registers.Allocate(core.Temp)
		c.ctx.Emitter.PatchInsertAxy(loop.Pos, vm.OpDataSetMultiSorter, loop.Dst, encoded, count)
		loop.Jump++

		return dst
	}

	// Single-key sorting only needs the direction
	dir := sortDirection(clauses[0].SortDirection())

	if loop.Allocate {
		c.ctx.Emitter.PatchSwapAx(loop.Pos, vm.OpDataSetSorter, loop.Dst, int(dir))

		return loop.Dst
	}

	dst := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.PatchInsertAx(loop.Pos, vm.OpDataSetSorter, dst, int(dir))
	loop.Jump++

	return dst
}

// finalizeSorting completes the sorting process by:
// 1. Adding KeyValuePairs to the result dataset
// 2. Finalizing the current loop
// 3. Replacing the loop source with sorted results
// 4. Reinitializing the loop for iteration over sorted data
func (c *LoopSortCompiler) finalizeSorting(loop *core.Loop, kv *core.KV, sorter vm.Operand) {
	// Add the KeyValuePair to the dataset
	c.ctx.Emitter.EmitABC(vm.OpPushKV, sorter, kv.Key, kv.Value)

	// Finalize the current loop iteration
	loop.EmitFinalization(c.ctx.Emitter)

	// Replace the loop source with sorted results
	c.ctx.Emitter.EmitAB(vm.OpMove, loop.Src, sorter)

	if !loop.Allocate {
		c.ctx.Registers.Free(sorter)
	}

	// Reinitialize the loop to iterate over sorted data
	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
}
