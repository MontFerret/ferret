package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func CollectUDFs(ctx *CompilerContext, program *fql.ProgramContext) *core.UDFTable {
	table := core.NewUDFTable()
	table.GlobalScope = core.NewUDFScope(nil)

	if program == nil || program.Body() == nil {
		return table
	}

	body, ok := program.Body().(*fql.BodyContext)
	if !ok {
		return table
	}

	top := collectScopeFunctionsFromBody(ctx, table, table.GlobalScope, body)

	for _, fn := range top {
		collectNestedFunctions(ctx, table, fn)
	}

	if ctx != nil && ctx.OptimizationLevel > optimization.LevelNone {
		pruneUnusedUDFs(ctx, table, body)
	}

	analyzeCaptures(ctx, table, body)

	return table
}
