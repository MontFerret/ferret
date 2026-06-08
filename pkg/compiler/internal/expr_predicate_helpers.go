package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func resolveArrayPredicateOpcode(op fql.IArrayOperatorContext) (bytecode.Opcode, bool) {
	if op == nil {
		return bytecode.Opcode(0), false
	}

	var pos int

	switch {
	case op.All() != nil:
		pos = int(bytecode.OpAllEq)
	case op.Any() != nil:
		pos = int(bytecode.OpAnyEq)
	case op.None() != nil:
		pos = int(bytecode.OpNoneEq)
	}

	if eo := op.EqualityOperator(); eo != nil {
		switch eo.GetText() {
		case "!=":
			pos += int(bytecode.OpAllNe) - int(bytecode.OpAllEq)
		case ">":
			pos += int(bytecode.OpAllGt) - int(bytecode.OpAllEq)
		case ">=":
			pos += int(bytecode.OpAllGte) - int(bytecode.OpAllEq)
		case "<":
			pos += int(bytecode.OpAllLt) - int(bytecode.OpAllEq)
		case "<=":
			pos += int(bytecode.OpAllLte) - int(bytecode.OpAllEq)
		default:
		}

		return bytecode.Opcode(pos), true
	}

	if op.InOperator() != nil {
		pos += int(bytecode.OpAllIn) - int(bytecode.OpAllEq)

		return bytecode.Opcode(pos), true
	}

	return bytecode.Opcode(0), false
}

func resolvePredicateEqNeJump(ctx fql.IPredicateContext) (string, fql.IPredicateContext, fql.IPredicateContext, bool) {
	if ctx == nil {
		return "", nil, nil, false
	}

	op := ctx.EqualityOperator()
	if op == nil {
		return "", nil, nil, false
	}

	opText := op.GetText()
	if opText != "==" && opText != "!=" {
		return "", nil, nil, false
	}

	leftCtx := ctx.Predicate(0)
	rightCtx := ctx.Predicate(1)
	if leftCtx == nil || rightCtx == nil {
		return "", nil, nil, false
	}

	return opText, leftCtx, rightCtx, true
}

func resolveEqNeJumpOpcode(opText string, jumpOnTrue, constOperand bool) bytecode.Opcode {
	if constOperand {
		if opText == "==" {
			if jumpOnTrue {
				return bytecode.OpJumpIfEqConst
			}

			return bytecode.OpJumpIfNeConst
		}

		if jumpOnTrue {
			return bytecode.OpJumpIfNeConst
		}

		return bytecode.OpJumpIfEqConst
	}

	if opText == "==" {
		if jumpOnTrue {
			return bytecode.OpJumpIfEq
		}

		return bytecode.OpJumpIfNe
	}

	if jumpOnTrue {
		return bytecode.OpJumpIfNe
	}

	return bytecode.OpJumpIfEq
}

func resolveAtomBinaryOperator(ctx fql.IExpressionAtomContext) (atomBinaryOperator, bool) {
	if op := ctx.MultiplicativeOperator(); op != nil {
		return resolveArithmeticBinaryOperator(op.GetText())
	}

	if op := ctx.AdditiveOperator(); op != nil {
		return resolveArithmeticBinaryOperator(op.GetText())
	}

	if op := ctx.RegexpOperator(); op != nil {
		return atomBinaryOperator{
			opcode:  bytecode.OpRegexp,
			text:    op.GetText(),
			negated: op.GetText() == "!~",
			regexp:  true,
		}, true
	}

	return atomBinaryOperator{}, false
}

func resolveArithmeticBinaryOperator(operator string) (atomBinaryOperator, bool) {
	switch operator {
	case "+", "+=":
		return atomBinaryOperator{opcode: bytecode.OpAdd, text: operator}, true
	case "-", "-=":
		return atomBinaryOperator{opcode: bytecode.OpSub, text: operator}, true
	case "*", "*=":
		return atomBinaryOperator{opcode: bytecode.OpMul, text: operator}, true
	case "/", "/=":
		return atomBinaryOperator{opcode: bytecode.OpDiv, text: operator}, true
	case "%":
		return atomBinaryOperator{opcode: bytecode.OpMod, text: operator}, true
	default:
		return atomBinaryOperator{}, false
	}
}

type numericOperandDiagnostic struct {
	operand bytecode.Operand
	span    source.Span
	label   string
}

func emitBinaryOperation(ctx *CompilationSession, facts *TypeFacts, spans binaryOperationSpans, op atomBinaryOperator, left, right bytecode.Operand) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if isStrictNumericBinaryOpcode(op.opcode) && !validateKnownNumericOperands(
		ctx,
		facts,
		spans.operator,
		op.text,
		numericOperandDiagnostic{operand: left, span: spans.leftOperand, label: "left operand"},
		numericOperandDiagnostic{operand: right, span: spans.rightOperand, label: "right operand"},
	) {
		return bytecode.NoopOperand
	}

	dst := ctx.Function.Registers.Allocate()
	ctx.Program.Emitter.WithSpan(spans.expression, func() {
		ctx.Program.Emitter.EmitABC(op.opcode, dst, left, right)

		if op.negated {
			ctx.Program.Emitter.EmitAB(bytecode.OpNot, dst, dst)
		}
	})

	resultType := facts.InferBinaryResultType(op, left, right)
	if op.negated {
		resultType = core.TypeBool
	}
	if resultType != core.TypeUnknown {
		ctx.Function.Types.Set(dst, resultType)
	}

	return dst
}

func isStrictNumericBinaryOpcode(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod:
		return true
	default:
		return false
	}
}

func validateKnownNumericOperands(
	ctx *CompilationSession,
	facts *TypeFacts,
	operatorSpan source.Span,
	operator string,
	operands ...numericOperandDiagnostic,
) bool {
	spans := []pkgdiagnostics.ErrorSpan{
		pkgdiagnostics.NewMainErrorSpan(operatorSpan, ""),
	}
	invalid := false

	for _, operand := range operands {
		typ := facts.OperandType(operand.operand)
		if knownNumericOperandTypeAllowed(typ) {
			continue
		}

		invalid = true
		if operand.label != "" && validDiagnosticSpan(operand.span) {
			spans = append(spans, pkgdiagnostics.NewSecondaryErrorSpan(
				operand.span,
				fmt.Sprintf("%s is %s", operand.label, numericDiagnosticTypeName(typ)),
			))
		}
	}

	if !invalid {
		return true
	}

	ctx.Program.Errors.Add(&pkgdiagnostics.Diagnostic{
		Kind:    parserd.SemanticError,
		Message: fmt.Sprintf("Operator '%s' requires numeric operands", operator),
		Hint:    "Use Int or Float values with this operator.",
		Spans:   spans,
	})

	return false
}

func knownNumericOperandTypeAllowed(typ core.ValueType) bool {
	switch typ {
	case core.TypeUnknown, core.TypeAny, core.TypeInt, core.TypeFloat:
		return true
	default:
		return false
	}
}

func binaryAtomOperationSpans(ctx fql.IExpressionAtomContext) binaryOperationSpans {
	spans := binaryOperationSpans{
		expression:   spanFromParserRuleContext(ctx),
		leftOperand:  spanFromParserRuleContext(ctx.GetLeft()),
		rightOperand: spanFromParserRuleContext(ctx.GetRight()),
	}

	if op := ctx.MultiplicativeOperator(); op != nil {
		spans.operator = parserd.SpanFromRuleContext(op)
	} else if op := ctx.AdditiveOperator(); op != nil {
		spans.operator = parserd.SpanFromRuleContext(op)
	}

	return spans
}

func assignmentOperationSpans(ctx *fql.AssignmentStatementContext) binaryOperationSpans {
	if ctx == nil {
		return binaryOperationSpans{}
	}

	return binaryOperationSpans{
		expression:   parserd.SpanFromRuleContext(ctx),
		operator:     spanFromParserRuleContext(ctx.AssignmentOperator()),
		leftOperand:  spanFromParserRuleContext(ctx.AssignmentTarget()),
		rightOperand: spanFromParserRuleContext(ctx.Expression()),
	}
}

func spanFromParserRuleContext(ctx antlr.ParserRuleContext) source.Span {
	if ctx == nil {
		return source.Span{Start: -1, End: -1}
	}

	return parserd.SpanFromRuleContext(ctx)
}

func validDiagnosticSpan(span source.Span) bool {
	return span.Start >= 0 && span.End > span.Start
}

func numericDiagnosticTypeName(typ core.ValueType) string {
	switch typ {
	case core.TypeNone:
		return "None"
	case core.TypeInt:
		return "Int"
	case core.TypeFloat:
		return "Float"
	case core.TypeString:
		return "String"
	case core.TypeBool:
		return "Bool"
	case core.TypeArray:
		return "Array"
	case core.TypeObject:
		return "Object"
	case core.TypeList:
		return "List"
	case core.TypeMap:
		return "Map"
	case core.TypeAny:
		return "Any"
	default:
		return "Unknown"
	}
}
