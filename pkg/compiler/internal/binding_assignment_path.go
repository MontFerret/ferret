package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *BindingCompiler) compilePathAssignmentStatement(stmt *fql.AssignmentStatementContext, target assignmentTarget, binding *core.Variable) bytecode.Operand {
	if len(target.Segments) == 0 {
		return bytecode.NoopOperand
	}

	endLabel := c.ctx.Program.Emitter.NewLabel("assignment", "end")
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

	value := bytecode.NoopOperand
	operator := assignmentOperatorText(stmt)

	if operator == "=" {
		if final.Safe {
			guard := c.emitAssignmentLoadSegment(parent, final, key, constKey, endLabel, true)
			if guard == bytecode.NoopOperand {
				return bytecode.NoopOperand
			}
		}

		value = c.exprs.Compile(stmt.Expression())
	} else {
		current := c.emitAssignmentLoadSegment(parent, final, key, constKey, endLabel, final.Safe)
		if current == bytecode.NoopOperand {
			return bytecode.NoopOperand
		}

		if !augmentedAssignmentKnownTypeAllowed(operator, c.facts.OperandType(current)) {
			c.reportInvalidAugmentedAssignment(stmt, operator)
			return bytecode.NoopOperand
		}

		op, ok := resolveArithmeticBinaryOperator(operator)
		if !ok {
			return bytecode.NoopOperand
		}

		right := c.exprs.Compile(stmt.Expression())
		if right == bytecode.NoopOperand {
			return bytecode.NoopOperand
		}

		value = emitBinaryOperation(c.ctx, c.facts, stmt, op, current, right)
	}

	if value == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	c.emitAssignmentStore(parent, final, key, constKey, value)

	return parent
}

func (c *BindingCompiler) compileAssignmentParent(target assignmentTarget, binding *core.Variable, endLabel core.Label) (bytecode.Operand, assignmentTargetSegment, bool) {
	segments := target.Segments
	final := segments[len(segments)-1]
	parent := c.LoadBindingValue(binding)
	parent = ensureOperandRegister(c.ctx, c.facts, parent)

	for _, segment := range segments[:len(segments)-1] {
		if segment.Safe {
			c.ctx.Program.Emitter.EmitJumpIfNone(parent, endLabel)
		}

		key, constKey := c.compileAssignmentSegmentOperand(segment)
		if key == bytecode.NoopOperand {
			return bytecode.NoopOperand, assignmentTargetSegment{}, false
		}

		parent = c.emitAssignmentLoadSegment(parent, segment, key, constKey, endLabel, false)
		if parent == bytecode.NoopOperand {
			return bytecode.NoopOperand, assignmentTargetSegment{}, false
		}
	}

	return parent, final, true
}

func (c *BindingCompiler) compileAssignmentSegmentOperand(segment assignmentTargetSegment) (bytecode.Operand, bool) {
	if segment.Property != nil {
		if constOp, ok := c.exprs.literals.CompilePropertyNameConst(segment.Property); ok {
			return constOp, true
		}

		return c.exprs.literals.CompilePropertyName(segment.Property), false
	}

	if segment.Computed != nil {
		if val, ok := c.facts.LiteralValueFromExpression(segment.Computed.Expression()); ok {
			switch val.(type) {
			case *runtime.Array, *runtime.Object:
				return c.exprs.literals.CompileComputedPropertyName(segment.Computed), false
			default:
				return c.ctx.Function.Symbols.AddConstant(val), true
			}
		}

		return c.exprs.literals.CompileComputedPropertyName(segment.Computed), false
	}

	return bytecode.NoopOperand, false
}

func (c *BindingCompiler) emitAssignmentLoadSegment(parent bytecode.Operand, segment assignmentTargetSegment, key bytecode.Operand, constKey bool, endLabel core.Label, jumpOnNone bool) bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()
	span := parserd.SpanFromRuleContext(segment.Context)

	c.ctx.Program.Emitter.WithSpan(span, func() {
		op := memberLoadOpcode(c.facts.OperandType(parent), constKey, segment.Safe)
		c.ctx.Program.Emitter.EmitABC(op, dst, parent, key)

		if jumpOnNone {
			c.ctx.Program.Emitter.EmitJumpIfNone(dst, endLabel)
		}
	})

	return dst
}

func (c *BindingCompiler) emitAssignmentStore(parent bytecode.Operand, segment assignmentTargetSegment, key bytecode.Operand, constKey bool, value bytecode.Operand) {
	span := parserd.SpanFromRuleContext(segment.Context)
	op := assignmentSetOpcode(c.facts.OperandType(parent), constKey)

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitABC(op, parent, key, value)
	})
}

func assignmentSetOpcode(parentType core.ValueType, constKey bool) bytecode.Opcode {
	switch parentType {
	case core.TypeArray:
		if constKey {
			return bytecode.OpSetIndexConst
		}

		return bytecode.OpSetIndex
	case core.TypeObject:
		if constKey {
			return bytecode.OpSetKeyConst
		}

		return bytecode.OpSetKey
	default:
		if constKey {
			return bytecode.OpSetPropertyConst
		}

		return bytecode.OpSetProperty
	}
}
