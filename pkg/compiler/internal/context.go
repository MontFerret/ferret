package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

// CompilerContext holds the context for the compilation process, including various compilers and allocators.
type CompilerContext struct {
	Source     *file.Source
	Emitter    *core.Emitter
	Registers  *core.RegisterAllocator
	Symbols    *core.SymbolTable
	Types      *core.TypeTracker
	Loops      *core.LoopTable
	CatchTable *core.CatchStack
	Errors     *diagnostics.ErrorHandler
	UseAliases map[string]string

	ExprCompiler        *ExprCompiler
	LiteralCompiler     *LiteralCompiler
	StmtCompiler        *StmtCompiler
	LoopCompiler        *LoopCompiler
	LoopSortCompiler    *LoopSortCompiler
	LoopCollectCompiler *LoopCollectCompiler
	WaitCompiler        *WaitCompiler
	DispatchCompiler    *DispatchCompiler
}

// NewCompilerContext initializes a new CompilerContext with default values.
func NewCompilerContext(src *file.Source, errors *diagnostics.ErrorHandler) *CompilerContext {
	ctx := &CompilerContext{
		Source:     src,
		Errors:     errors,
		Emitter:    core.NewEmitter(),
		Registers:  core.NewRegisterAllocator(),
		Symbols:    nil, // set later
		Loops:      nil, // set later
		CatchTable: core.NewCatchStack(),
		UseAliases: make(map[string]string),
	}
	ctx.Symbols = core.NewSymbolTable(ctx.Registers)
	ctx.Types = core.NewTypeTracker()
	ctx.Loops = core.NewLoopTable(ctx.Registers)

	ctx.ExprCompiler = NewExprCompiler(ctx)
	ctx.LiteralCompiler = NewLiteralCompiler(ctx)
	ctx.StmtCompiler = NewStmtCompiler(ctx)
	ctx.LoopCompiler = NewLoopCompiler(ctx)
	ctx.LoopSortCompiler = NewLoopSortCompiler(ctx)
	ctx.LoopCollectCompiler = NewCollectCompiler(ctx)
	ctx.WaitCompiler = NewWaitCompiler(ctx)
	ctx.DispatchCompiler = NewDispatchCompiler(ctx)

	return ctx
}
