package internal

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

// emitMoveAuto emits OpMove when the source is known to be untracked.
// Tracked values use OpMoveTracked so ownership-sensitive paths preserve value semantics.
func emitMoveAuto(ctx *CompilerContext, dst, src bytecode.Operand) {
	srcType := operandType(ctx, src)

	if srcType.IsUntracked() {
		ctx.Emitter.EmitPlainMove(dst, src)
	} else {
		ctx.Emitter.EmitMoveTracked(dst, src)
	}

	ctx.Types.Set(dst, srcType)
}
