package internal

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
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

func resolveErrorPolicy(ctx *CompilerContext, optional bool, tail fql.IErrorPolicyTailContext) errorPolicy {
	if optional {
		return errorPolicySuppress
	}

	return resolveErrorPolicyTail(ctx, tail)
}

func resolveErrorPolicyTail(ctx *CompilerContext, tail fql.IErrorPolicyTailContext) errorPolicy {
	if tail == nil {
		return errorPolicyDefault
	}

	if !strings.EqualFold(tail.OnKeyword().GetText(), "ON") {
		return errorPolicyDefault
	}

	if tail.ErrorKeyword() == nil || !strings.EqualFold(tail.ErrorKeyword().GetText(), "ERROR") {
		reportInvalidErrorPolicyTail(ctx, tail.OnKeyword(), "Expected ERROR after 'ON' in error policy tail", "Complete the tail as ON ERROR SUPPRESS or ON ERROR FAIL.")
		return errorPolicyDefault
	}

	policy := tail.GetPolicy()
	if policy == nil {
		reportInvalidErrorPolicyTail(ctx, tail.ErrorKeyword(), "Expected SUPPRESS or FAIL after 'ON ERROR'", "Use ON ERROR SUPPRESS to swallow failures or ON ERROR FAIL to propagate them.")
		return errorPolicyDefault
	}

	switch text := policy.GetText(); {
	case strings.EqualFold(text, "FAIL"):
		return errorPolicyFail
	case strings.EqualFold(text, "SUPPRESS"):
		return errorPolicySuppress
	}

	reportInvalidErrorPolicyTail(ctx, tail.ErrorKeyword(), "Expected SUPPRESS or FAIL after 'ON ERROR'", "Use ON ERROR SUPPRESS to swallow failures or ON ERROR FAIL to propagate them.")

	return errorPolicyDefault
}

func reportInvalidErrorPolicyTail(ctx *CompilerContext, node antlr.ParserRuleContext, message, hint string) {
	if ctx == nil || ctx.Errors == nil || node == nil {
		return
	}

	err := ctx.Errors.Create(parserd.SyntaxError, node, message)
	err.Hint = hint
	ctx.Errors.Add(err)
}

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

func allowsTailCallPolicy(policy errorPolicy) bool {
	return policy != errorPolicySuppress
}
