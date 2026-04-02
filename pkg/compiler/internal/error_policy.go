package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type errorPolicy int

const (
	errorPolicyDefault errorPolicy = iota
	errorPolicySuppress
	errorPolicyFail
)

type catchJumpMode int

const (
	catchJumpNone catchJumpMode = iota
	catchJumpEnd
)

func compileWithErrorPolicy(ctx *CompilerContext, policy errorPolicy, jumpMode catchJumpMode, compile func() bytecode.Operand) bytecode.Operand {
	if ctx == nil || compile == nil || policy != errorPolicySuppress {
		return compile()
	}

	startCatch := ctx.Emitter.Size()
	out := compile()
	endCatchExclusive := ctx.Emitter.Size()
	if endCatchExclusive <= startCatch {
		return out
	}

	endCatch := endCatchExclusive - 1
	jump := -1

	if jumpMode == catchJumpEnd {
		jump = endCatch
	} else {
		endLabel := ctx.Emitter.NewLabel("error", "suppress", "end")
		ctx.Emitter.EmitJump(endLabel)
		jump = ctx.Emitter.Size()
		ctx.Emitter.EmitJump(endLabel)
		ctx.Emitter.MarkLabel(endLabel)
	}

	ctx.CatchTable.Push(startCatch, endCatch, jump)

	return out
}

func catchJumpModeForForExpression(ctx fql.IForExpressionContext) catchJumpMode {
	if ctx != nil && ctx.In() != nil {
		return catchJumpEnd
	}

	return catchJumpNone
}

func catchJumpModeForWaitForExpression(ctx fql.IWaitForExpressionContext) catchJumpMode {
	if ctx != nil && ctx.WaitForEventExpression() != nil {
		return catchJumpEnd
	}

	return catchJumpNone
}
