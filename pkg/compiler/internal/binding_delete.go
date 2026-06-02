package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func (c *BindingCompiler) CompileDeleteStatement(ctx fql.IDeleteStatementContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	stmt, ok := ctx.(*fql.DeleteStatementContext)
	if !ok || stmt == nil {
		return bytecode.NoopOperand
	}

	target := stmt.AssignmentTarget()
	if target == nil {
		c.reportInvalidDeleteTarget(stmt)
		return bytecode.NoopOperand
	}

	deletion, ok := newAssignmentTarget(target)
	if !ok {
		c.reportInvalidDeleteTarget(stmt)
		return bytecode.NoopOperand
	}

	if deletion.Root == "" || deletion.Root == core.IgnorePseudoVariable || len(deletion.Segments) == 0 {
		c.reportInvalidDeleteTarget(stmt)
		return bytecode.NoopOperand
	}

	binding, found := c.ctx.Function.Symbols.ResolveBinding(deletion.Root)
	if !found {
		if c.reportFunctionDeleteTarget(deletion) {
			return bytecode.NoopOperand
		}

		reportVariableNotFound(c.ctx, deletion.RootTok, deletion.Root)
		return bytecode.NoopOperand
	}

	return c.compilePathDeleteStatement(stmt, deletion, binding)
}

func (c *BindingCompiler) compilePathDeleteStatement(_ *fql.DeleteStatementContext, target assignmentTarget, binding *core.Variable) bytecode.Operand {
	endLabel := c.ctx.Program.Emitter.NewLabel("delete", "end")
	defer c.ctx.Program.Emitter.MarkLabel(endLabel)

	parent, final, ok := c.compileAssignmentParent(target, binding, endLabel)
	if !ok {
		return bytecode.NoopOperand
	}

	if final.Safe {
		c.ctx.Program.Emitter.EmitJumpIfNone(parent, endLabel)
	}

	key, constKey := c.compileAssignmentSegmentOperand(final)
	if key == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	c.emitDelete(parent, final, key, constKey)

	return parent
}

func (c *BindingCompiler) emitDelete(parent bytecode.Operand, segment assignmentTargetSegment, key bytecode.Operand, constKey bool) {
	span := parserd.SpanFromRuleContext(segment.Context)
	op := c.deleteOpcode(c.facts.OperandType(parent), segment.Computed != nil, constKey)

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitAB(op, parent, key)
	})
}

func (c *BindingCompiler) deleteOpcode(parentType core.ValueType, computed bool, constKey bool) bytecode.Opcode {
	if computed || parentType == core.TypeObject {
		if constKey {
			return bytecode.OpDeleteKeyConst
		}

		return bytecode.OpDeleteKey
	}

	if constKey {
		return bytecode.OpDeletePropertyConst
	}

	return bytecode.OpDeleteProperty
}
