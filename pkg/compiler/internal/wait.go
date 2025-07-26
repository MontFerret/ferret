package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// WaitCompiler handles the compilation of WAIT FOR expressions in FQL queries.
// It transforms wait operations into VM instructions for event streaming and processing.
type WaitCompiler struct {
	ctx *CompilerContext
}

// NewWaitCompiler creates a new instance of WaitCompiler with the given compiler context.
func NewWaitCompiler(ctx *CompilerContext) *WaitCompiler {
	return &WaitCompiler{
		ctx: ctx,
	}
}

// Compile processes a WAIT FOR expression from the FQL AST and generates the appropriate VM instructions.
// It sets up event streaming, handles optional timeouts and filters, and manages the event iteration loop.
// Parameters:
//   - ctx: The WAIT FOR expression context from the AST
//
// Returns:
//   - A no-op operand as the WAIT FOR expression doesn't produce a value directly
func (wc *WaitCompiler) Compile(ctx fql.IWaitForExpressionContext) vm.Operand {
	// Create a new scope for the WAIT FOR expression
	wc.ctx.Symbols.EnterScope()

	// Compile the event source and event name expressions
	srcReg := wc.CompileWaitForEventSource(ctx.WaitForEventSource())
	eventReg := wc.CompileWaitForEventName(ctx.WaitForEventName())

	var optsReg vm.Operand

	// Compile options clause if present
	if opts := ctx.OptionsClause(); opts != nil {
		optsReg = wc.CompileOptionsClause(opts)
	}

	var timeoutReg vm.Operand

	// Compile timeout clause if present
	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = wc.CompileTimeoutClauseContext(timeout)
	}

	// Allocate a register for the event stream
	streamReg := wc.ctx.Registers.Allocate(core.Temp)

	// Move the source object to the stream register in order to re-use it in OpStream
	wc.ctx.Emitter.EmitMove(streamReg, srcReg)
	// Create a stream from the source, listening for the specified event
	wc.ctx.Emitter.EmitABC(vm.OpStream, streamReg, eventReg, optsReg)
	// Set up iteration over the stream with optional timeout
	wc.ctx.Emitter.EmitAB(vm.OpStreamIter, streamReg, timeoutReg)

	var valReg vm.Operand
	// Create labels for the start and end of the iteration loop
	start := wc.ctx.Emitter.NewLabel()
	end := wc.ctx.Emitter.NewLabel()

	// Mark the start of the iteration loop
	wc.ctx.Emitter.MarkLabel(start)
	// Emit instruction to get the next event from the stream
	// If no more events, jump to the end label
	wc.ctx.Emitter.EmitIterNext(streamReg, end)

	// Handle filter clause if present
	if filter := ctx.FilterClause(); filter != nil {
		// Declare a local variable to hold the event value
		valReg = wc.ctx.Symbols.DeclareLocal(core.PseudoVariable, core.TypeUnknown)
		// Load the current event value into the variable
		wc.ctx.Emitter.EmitAB(vm.OpIterValue, valReg, streamReg)

		// Compile the filter expression
		cond := wc.ctx.ExprCompiler.Compile(filter.Expression())

		// If the filter condition is false, jump back to the start to get the next event
		wc.ctx.Emitter.EmitJumpIfFalse(cond, start)

		// TODO: Do we need to use timeout here too? We can really get stuck in the loop if no event satisfies the filter
	}

	// Mark the end of the iteration loop
	wc.ctx.Emitter.MarkLabel(end)
	// Clean up the stream by closing it
	wc.ctx.Emitter.EmitA(vm.OpClose, streamReg)
	// Exit the scope created for the WAIT FOR expression
	wc.ctx.Symbols.ExitScope()

	// WAIT FOR doesn't produce a value directly
	return vm.NoopOperand
}

// CompileWaitForEventName processes the event name expression in a WAIT FOR statement.
// It handles various types of expressions that can be used as event names,
// such as string literals, variables, parameters, member expressions, and function calls.
// Parameters:
//   - ctx: The event name context from the AST
//
// Returns:
//   - An operand representing the compiled event name expression
//
// Panics if the event name expression type is not recognized.
func (wc *WaitCompiler) CompileWaitForEventName(ctx fql.IWaitForEventNameContext) vm.Operand {
	// Handle string literal event names (e.g., WAIT FOR doc ON "click")
	if c := ctx.StringLiteral(); c != nil {
		return wc.ctx.LiteralCompiler.CompileStringLiteral(c)
	}

	// Handle variable event names (e.g., WAIT FOR doc ON eventName)
	if c := ctx.Variable(); c != nil {
		return wc.ctx.ExprCompiler.CompileVariable(c)
	}

	// Handle parameter event names (e.g., WAIT FOR doc ON @eventParam)
	if c := ctx.Param(); c != nil {
		return wc.ctx.ExprCompiler.CompileParam(c)
	}

	// Handle member expression event names (e.g., WAIT FOR doc ON events.click)
	if c := ctx.MemberExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	// Handle function call expression event names (e.g., WAIT FOR doc ON getEventName())
	if c := ctx.FunctionCallExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileFunctionCallExpression(c)
	}

	// If none of the above, the event name expression is invalid
	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

// CompileWaitForEventSource processes the event source expression in a WAIT FOR statement.
// It handles various types of expressions that can be used as event sources,
// such as variables, member expressions, and function calls.
// Parameters:
//   - ctx: The event source context from the AST
//
// Returns:
//   - An operand representing the compiled event source expression
//
// Panics if the event source expression type is not recognized.
func (wc *WaitCompiler) CompileWaitForEventSource(ctx fql.IWaitForEventSourceContext) vm.Operand {
	// Handle variable event sources (e.g., WAIT FOR document ON "click")
	if c := ctx.Variable(); c != nil {
		return wc.ctx.ExprCompiler.CompileVariable(c)
	}

	// Handle member expression event sources (e.g., WAIT FOR page.document ON "click")
	if c := ctx.MemberExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	// Handle function call expression event sources (e.g., WAIT FOR getDocument() ON "click")
	if c := ctx.FunctionCallExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileFunctionCallExpression(c)
	}

	// If none of the above, the event source expression is invalid
	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

// CompileOptionsClause processes the options clause in a WAIT FOR statement.
// It compiles the object literal that contains configuration options for the event stream.
// Parameters:
//   - ctx: The options clause context from the AST
//
// Returns:
//   - An operand representing the compiled options object
//
// Panics if the options expression is not an object literal.
func (wc *WaitCompiler) CompileOptionsClause(ctx fql.IOptionsClauseContext) vm.Operand {
	// Handle object literal options (e.g., WAIT FOR doc ON "click" OPTIONS { timeout: 5000 })
	if c := ctx.ObjectLiteral(); c != nil {
		return wc.ctx.LiteralCompiler.CompileObjectLiteral(c)
	}

	// If not an object literal, the options expression is invalid
	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

// CompileTimeoutClauseContext processes the timeout clause in a WAIT FOR statement.
// It handles various types of expressions that can be used as timeout values,
// such as integer literals, variables, parameters, member expressions, and function calls.
// Parameters:
//   - ctx: The timeout clause context from the AST
//
// Returns:
//   - An operand representing the compiled timeout expression
//
// Panics if the timeout expression type is not recognized.
func (wc *WaitCompiler) CompileTimeoutClauseContext(ctx fql.ITimeoutClauseContext) vm.Operand {
	// Handle integer literal timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT 5000)
	if c := ctx.IntegerLiteral(); c != nil {
		return wc.ctx.LiteralCompiler.CompileIntegerLiteral(c)
	}

	// Handle variable timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT timeout)
	if c := ctx.Variable(); c != nil {
		return wc.ctx.ExprCompiler.CompileVariable(c)
	}

	// Handle parameter timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT @timeout)
	if c := ctx.Param(); c != nil {
		return wc.ctx.ExprCompiler.CompileParam(c)
	}

	// Handle member expression timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT config.timeout)
	if c := ctx.MemberExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	// Handle function call timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT getTimeout())
	if c := ctx.FunctionCall(); c != nil {
		return wc.ctx.ExprCompiler.CompileFunctionCall(c, false)
	}

	// If none of the above, the timeout expression is invalid
	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}
