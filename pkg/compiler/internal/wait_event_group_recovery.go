package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func (c *WaitCompiler) buildProtectedEventGroupRecovery(
	ctx fql.IWaitForEventGroupExpressionContext,
	recoveryLabel, timeoutLabel, endLabel core.Label,
) ProtectedRecoveryRegion {
	hasTimeout := waitForEventGroupTimeoutClause(ctx) != nil
	streamReg := c.ctx.Function.Registers.Allocate()
	resultReg := c.ctx.Function.Registers.Allocate()
	errorStateReg := c.ctx.Function.Registers.Allocate()
	streamReadyReg := c.ctx.Function.Registers.Allocate()
	timeoutStateReg := bytecode.NoopOperand

	if hasTimeout {
		timeoutStateReg = c.ctx.Function.Registers.Allocate()
		c.ctx.Program.Emitter.EmitBoolean(timeoutStateReg, false)
	}

	c.ctx.Program.Emitter.EmitLoadNone(resultReg)
	c.ctx.Program.Emitter.EmitBoolean(errorStateReg, false)
	c.ctx.Program.Emitter.EmitBoolean(streamReadyReg, false)

	startCatch := c.ctx.Program.Emitter.Size()
	state, ok := c.buildWaitEventGroupState(ctx)
	if !ok {
		return ProtectedRecoveryRegion{Result: bytecode.NoopOperand}
	}

	c.emitWaitEventGroupStreamSetupWithReady(state, streamReg, streamReadyReg)
	c.compileWaitEventGroupTrigger(ctx)

	start := c.ctx.Program.Emitter.NewLabel()
	done := c.ctx.Program.Emitter.NewLabel()
	cleanup := c.ctx.Program.Emitter.NewLabel()
	routeRecovery := c.ctx.Program.Emitter.NewLabel("waitfor", "event", "group", "recover")

	c.ctx.Program.Emitter.MarkLabel(start)
	c.emitWaitEventGroupIteration(state, streamReg, resultReg, timeoutStateReg, start, done)
	c.ctx.Program.Emitter.EmitJump(cleanup)
	c.ctx.Program.Emitter.MarkLabel(done)
	c.ctx.Program.Emitter.EmitJump(cleanup)
	c.ctx.Program.Emitter.MarkLabel(cleanup)
	c.emitWaitEventGroupCleanupIfReady(state, streamReg, streamReadyReg)

	endCatchExclusive := c.ctx.Program.Emitter.Size()
	if hasTimeout {
		c.ctx.Program.Emitter.EmitJumpIfTrue(timeoutStateReg, timeoutLabel)
	}

	c.ctx.Program.Emitter.EmitJumpIfTrue(errorStateReg, routeRecovery)
	c.ctx.Program.Emitter.EmitJump(endLabel)

	errorPreludePC := c.ctx.Program.Emitter.Size()
	c.ctx.Program.Emitter.EmitBoolean(errorStateReg, true)

	if hasTimeout {
		c.ctx.Program.Emitter.EmitBoolean(timeoutStateReg, false)
	}

	c.ctx.Program.Emitter.EmitJump(cleanup)

	c.ctx.Program.Emitter.MarkLabel(routeRecovery)
	c.ctx.Program.Emitter.EmitBoolean(errorStateReg, false)

	if hasTimeout {
		c.ctx.Program.Emitter.EmitBoolean(timeoutStateReg, false)
	}

	c.ctx.Program.Emitter.EmitJump(recoveryLabel)

	return ProtectedRecoveryRegion{
		Result:            resultReg,
		StartCatch:        startCatch,
		EndCatchExclusive: endCatchExclusive,
		CatchHandlerPC:    errorPreludePC,
		HasTimeout:        hasTimeout,
	}
}
