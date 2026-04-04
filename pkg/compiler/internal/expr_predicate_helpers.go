package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
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
			negated: op.GetText() == "!~",
			regexp:  true,
		}, true
	}

	return atomBinaryOperator{}, false
}

func resolveArithmeticBinaryOperator(operator string) (atomBinaryOperator, bool) {
	switch operator {
	case "+", "+=":
		return atomBinaryOperator{opcode: bytecode.OpAdd}, true
	case "-", "-=":
		return atomBinaryOperator{opcode: bytecode.OpSub}, true
	case "*", "*=":
		return atomBinaryOperator{opcode: bytecode.OpMul}, true
	case "/", "/=":
		return atomBinaryOperator{opcode: bytecode.OpDiv}, true
	case "%":
		return atomBinaryOperator{opcode: bytecode.OpMod}, true
	default:
		return atomBinaryOperator{}, false
	}
}

func emitBinaryOperation(ctx *CompilationSession, facts *TypeFacts, prc antlr.ParserRuleContext, op atomBinaryOperator, left, right bytecode.Operand) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	dst := ctx.Function.Registers.Allocate()
	span := source.Span{Start: -1, End: -1}

	if prc != nil {
		span = diagnostics.SpanFromRuleContext(prc)
	}

	ctx.Program.Emitter.WithSpan(span, func() {
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
