package internal

import (
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	scalarLiteralNode interface {
		NoneLiteral() fql.INoneLiteralContext
		BooleanLiteral() fql.IBooleanLiteralContext
		StringLiteral() fql.IStringLiteralContext
		FloatLiteral() fql.IFloatLiteralContext
		IntegerLiteral() fql.IIntegerLiteralContext
	}

	operandBranch struct {
		compile func() bytecode.Operand
		enabled bool
	}
)

func compileScalarLiteralOperand(ctx *CompilationSession, literals *LiteralCompiler, lit scalarLiteralNode) bytecode.Operand {
	if lit == nil {
		return bytecode.NoopOperand
	}

	if lit.NoneLiteral() != nil {
		return ctx.Symbols.AddConstant(runtime.None)
	}

	if bl := lit.BooleanLiteral(); bl != nil {
		if val, ok := literalBooleanValue(bl.GetText()); ok {
			return ctx.Symbols.AddConstant(val)
		}

		return bytecode.NoopOperand
	}

	if sl := lit.StringLiteral(); sl != nil {
		if val, ok := parseStringLiteralConst(sl); ok {
			return ctx.Symbols.AddConstant(val)
		}

		return literals.CompileStringLiteral(sl)
	}

	if fl := lit.FloatLiteral(); fl != nil {
		if val, ok := literalFloatValue(fl.GetText()); ok {
			return ctx.Symbols.AddConstant(val)
		}

		return literals.CompileFloatLiteral(fl)
	}

	if il := lit.IntegerLiteral(); il != nil {
		if val, ok := literalIntValue(il.GetText()); ok {
			return ctx.Symbols.AddConstant(val)
		}

		return literals.CompileIntegerLiteral(il)
	}

	return bytecode.NoopOperand
}

func scalarLiteralValue(lit scalarLiteralNode) (runtime.Value, bool) {
	if lit == nil {
		return nil, false
	}

	if lit.NoneLiteral() != nil {
		return runtime.None, true
	}

	if bl := lit.BooleanLiteral(); bl != nil {
		return literalBooleanValue(bl.GetText())
	}

	if sl := lit.StringLiteral(); sl != nil {
		if val, ok := parseStringLiteralConst(sl); ok {
			return val, true
		}

		return nil, false
	}

	if fl := lit.FloatLiteral(); fl != nil {
		return literalFloatValue(fl.GetText())
	}

	if il := lit.IntegerLiteral(); il != nil {
		return literalIntValue(il.GetText())
	}

	return nil, false
}

func literalBooleanValue(text string) (runtime.Value, bool) {
	switch strings.ToLower(text) {
	case "true":
		return runtime.True, true
	case "false":
		return runtime.False, true
	default:
		return nil, false
	}
}

func literalFloatValue(text string) (runtime.Value, bool) {
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return nil, false
	}

	return runtime.NewFloat(val), true
}

func literalIntValue(text string) (runtime.Value, bool) {
	val, err := strconv.Atoi(text)
	if err != nil {
		return nil, false
	}

	return runtime.NewInt(val), true
}

func newOperandBranch(enabled bool, compile func() bytecode.Operand) operandBranch {
	return operandBranch{
		enabled: enabled,
		compile: compile,
	}
}

func compileFirstOperand(branches ...operandBranch) bytecode.Operand {
	for _, branch := range branches {
		if !branch.enabled || branch.compile == nil {
			continue
		}

		if op := branch.compile(); op != bytecode.NoopOperand {
			return op
		}
	}

	return bytecode.NoopOperand
}
