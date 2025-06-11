package internal

// FuncContext encapsulates the context and state required for compiling and managing functions during code processing.
type FuncContext struct {
	Emitter    *Emitter
	Registers  *RegisterAllocator
	Symbols    *SymbolTable
	Loops      *LoopTable
	CatchTable *CatchStack

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
		Emitter:    NewEmitter(),
		Registers:  NewRegisterAllocator(),
		Symbols:    nil, // set later
		Loops:      nil, // set later
		CatchTable: NewCatchStack(),
	}
	ctx.Symbols = NewSymbolTable(ctx.Registers)
	ctx.Loops = NewLoopTable(ctx.Registers)

	ctx.ExprCompiler = NewExprCompiler(ctx)
	ctx.LiteralCompiler = NewLiteralCompiler(ctx)
	ctx.StmtCompiler = NewStmtCompiler(ctx)
	ctx.LoopCompiler = NewLoopCompiler(ctx)
	ctx.CollectCompiler = NewCollectCompiler(ctx)
	ctx.WaitCompiler = NewWaitCompiler(ctx)

	return ctx
}
