package internal

import "github.com/MontFerret/ferret/pkg/compiler/internal/core"

// CompilerContext holds the context for the compilation process, including various compilers and allocators.
type CompilerContext struct {
	Emitter    *core.Emitter
	Registers  *core.RegisterAllocator
	Symbols    *core.SymbolTable
	Loops      *core.LoopTable
	CatchTable *core.CatchStack

	ExprCompiler    *ExprCompiler
	LiteralCompiler *LiteralCompiler
	StmtCompiler    *StmtCompiler
	LoopCompiler    *LoopCompiler
	CollectCompiler *CollectCompiler
	WaitCompiler    *WaitCompiler
}

// NewCompilerContext initializes a new CompilerContext with default values.
func NewCompilerContext() *CompilerContext {
	ctx := &CompilerContext{
		Emitter:    core.NewEmitter(),
		Registers:  core.NewRegisterAllocator(),
		Symbols:    nil, // set later
		Loops:      nil, // set later
		CatchTable: core.NewCatchStack(),
	}
	ctx.Symbols = core.NewSymbolTable(ctx.Registers)
	ctx.Loops = core.NewLoopTable(ctx.Registers)

	ctx.ExprCompiler = NewExprCompiler(ctx)
	ctx.LiteralCompiler = NewLiteralCompiler(ctx)
	ctx.StmtCompiler = NewStmtCompiler(ctx)
	ctx.LoopCompiler = NewLoopCompiler(ctx)
	ctx.CollectCompiler = NewCollectCompiler(ctx)
	ctx.WaitCompiler = NewWaitCompiler(ctx)

	return ctx
}
