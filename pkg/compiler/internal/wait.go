package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
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
func (c *WaitCompiler) Compile(ctx fql.IWaitForExpressionContext) vm.Operand {
	// Create a new scope for the WAIT FOR expression
	c.ctx.Symbols.EnterScope()

	// Compile the event source and event name expressions
	srcReg := c.CompileWaitForEventSource(ctx.WaitForEventSource())
	eventReg := c.CompileWaitForEventName(ctx.WaitForEventName())

	var optsReg vm.Operand

	// Compile options clause if present
	if opts := ctx.OptionsClause(); opts != nil {
		optsReg = c.CompileOptionsClause(opts)
	}

	var timeoutReg vm.Operand

	// Compile timeout clause if present
	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = c.CompileTimeoutClauseContext(timeout)
	}

	// Allocate a register for the event stream
	streamReg := c.ctx.Registers.Allocate()

	// Move the source object to the stream register in order to re-use it in OpStream
	c.ctx.Emitter.EmitMove(streamReg, srcReg)
	// Create a stream from the source, listening for the specified event
	c.ctx.Emitter.EmitABC(vm.OpStream, streamReg, eventReg, optsReg)
	// Set up iteration over the stream with optional timeout
	c.ctx.Emitter.EmitAB(vm.OpStreamIter, streamReg, timeoutReg)

	var valReg vm.Operand
	// Create labels for the start and end of the iteration loop
	start := c.ctx.Emitter.NewLabel()
	end := c.ctx.Emitter.NewLabel()

	// Mark the start of the iteration loop
	c.ctx.Emitter.MarkLabel(start)
	// Emit instruction to get the next event from the stream
	// If no more events, jump to the end label
	c.ctx.Emitter.EmitIterNext(streamReg, end)

	// Handle filter clause if present
	if filter := ctx.FilterClause(); filter != nil {
		// Declare a local variable to hold the event value
		valReg, _ = c.ctx.Symbols.DeclareLocal(core.PseudoVariable, core.TypeUnknown)
		// Load the current event value into the variable
		c.ctx.Emitter.EmitAB(vm.OpIterValue, valReg, streamReg)

		// Compile the filter expression
		cond := c.ctx.ExprCompiler.Compile(filter.Expression())

		// If the filter condition is false, jump back to the start to get the next event
		c.ctx.Emitter.EmitJumpIfFalse(cond, start)

		// TODO: Do we need to use timeout here too? We can really get stuck in the loop if no event satisfies the filter
	}

	// Mark the end of the iteration loop
	c.ctx.Emitter.MarkLabel(end)
	// Clean up the stream by closing it
	c.ctx.Emitter.EmitA(vm.OpClose, streamReg)
	// Exit the scope created for the WAIT FOR expression
	c.ctx.Symbols.ExitScope()

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
func (c *WaitCompiler) CompileWaitForEventName(ctx fql.IWaitForEventNameContext) vm.Operand {
	// Handle string literal event names (e.g., WAIT FOR doc ON "click")
	if sl := ctx.StringLiteral(); sl != nil {
		return c.ctx.LiteralCompiler.CompileStringLiteral(sl)
	}

	// Handle variable event names (e.g., WAIT FOR doc ON eventName)
	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	// Handle parameter event names (e.g., WAIT FOR doc ON @eventParam)
	if p := ctx.Param(); p != nil {
		return c.ctx.ExprCompiler.CompileParam(p)
	}

	// Handle member expression event names (e.g., WAIT FOR doc ON events.click)
	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	// Handle function call expression event names (e.g., WAIT FOR doc ON getEventName())
	if fce := ctx.FunctionCall(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCall(fce, false)
	}

	return vm.NoopOperand
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
func (c *WaitCompiler) CompileWaitForEventSource(ctx fql.IWaitForEventSourceContext) vm.Operand {
	// Handle variable event sources (e.g., WAIT FOR document ON "click")
	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	// Handle member expression event sources (e.g., WAIT FOR page.document ON "click")
	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	// Handle function call expression event sources (e.g., WAIT FOR getDocument() ON "click")
	if fce := ctx.FunctionCallExpression(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
	}

	return vm.NoopOperand
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
func (c *WaitCompiler) CompileOptionsClause(ctx fql.IOptionsClauseContext) vm.Operand {
	// Handle object literal options (e.g., WAIT FOR doc ON "click" OPTIONS { timeout: 5000 })
	if ol := ctx.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	return vm.NoopOperand
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
func (c *WaitCompiler) CompileTimeoutClauseContext(ctx fql.ITimeoutClauseContext) vm.Operand {
	// Handle integer literal timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT 5000)
	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	// Handle variable timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT timeout)
	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	// Handle parameter timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT @timeout)
	if p := ctx.Param(); p != nil {
		return c.ctx.ExprCompiler.CompileParam(p)
	}

	// Handle member expression timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT config.timeout)
	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	// Handle function call timeouts (e.g., WAIT FOR doc ON "click" TIMEOUT getTimeout())
	if fc := ctx.FunctionCall(); fc != nil {
		return c.ctx.ExprCompiler.CompileFunctionCall(fc, false)
	}

	return vm.NoopOperand
}
