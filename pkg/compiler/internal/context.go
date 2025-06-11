package internal

import "github.com/MontFerret/ferret/pkg/compiler/internal/core"

// FuncContext encapsulates the context and state required for compiling and managing functions during code processing.
type FuncContext struct {
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

// NewFuncContext initializes and returns a new instance of FuncContext, setting up all required components for compilation.
func NewFuncContext() *FuncContext {
	ctx := &FuncContext{
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
