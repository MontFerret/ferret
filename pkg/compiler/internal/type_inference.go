package internal

import (
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func valueTypeFromRuntime(value runtime.Value) core.ValueType {
	if value == runtime.None {
		return core.TypeNone
	}

	switch value.(type) {
	case runtime.Int:
		return core.TypeInt
	case runtime.Float:
		return core.TypeFloat
	case runtime.String:
		return core.TypeString
	case runtime.Boolean:
		return core.TypeBool
	case *runtime.Array:
		return core.TypeArray
	case *runtime.Object:
		return core.TypeObject
	default:
		return core.TypeUnknown
	}
}

func literalType(ctx fql.ILiteralContext) core.ValueType {
	if ctx == nil {
		return core.TypeUnknown
	}

	switch {
	case ctx.NoneLiteral() != nil:
		return core.TypeNone
	case ctx.StringLiteral() != nil:
		return core.TypeString
	case ctx.IntegerLiteral() != nil:
		return core.TypeInt
	case ctx.FloatLiteral() != nil:
		return core.TypeFloat
	case ctx.BooleanLiteral() != nil:
		return core.TypeBool
	case ctx.ArrayLiteral() != nil:
		return core.TypeArray
	case ctx.ObjectLiteral() != nil:
		return core.TypeObject
	default:
		return core.TypeUnknown
	}
}

func literalValue(ctx fql.ILiteralContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	switch {
	case ctx.NoneLiteral() != nil:
		return runtime.None, true
	case ctx.StringLiteral() != nil:
		return parseStringLiteral(ctx.StringLiteral()), true
	case ctx.IntegerLiteral() != nil:
		val, err := strconv.Atoi(ctx.IntegerLiteral().GetText())
		if err != nil {
			return nil, false
		}
		return runtime.NewInt(val), true
	case ctx.FloatLiteral() != nil:
		val, err := strconv.ParseFloat(ctx.FloatLiteral().GetText(), 64)
		if err != nil {
			return nil, false
		}
		return runtime.NewFloat(val), true
	case ctx.BooleanLiteral() != nil:
		switch strings.ToLower(ctx.BooleanLiteral().GetText()) {
		case "true":
			return runtime.True, true
		case "false":
			return runtime.False, true
		}
	case ctx.ArrayLiteral() != nil:
		return runtime.NewArray(0), true
	case ctx.ObjectLiteral() != nil:
		return runtime.NewObject(), true
	}

	return nil, false
}

func literalValueFromExpression(ctx fql.IExpressionContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	if p := ctx.Predicate(); p != nil {
		return literalValueFromPredicate(p)
	}

	return nil, false
}

func literalValueFromPredicate(ctx fql.IPredicateContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	if atom := ctx.ExpressionAtom(); atom != nil {
		return literalValueFromAtom(atom)
	}

	return nil, false
}

func literalValueFromAtom(ctx fql.IExpressionAtomContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	if lit := ctx.Literal(); lit != nil {
		return literalValue(lit)
	}

	return nil, false
}

func operandType(ctx *CompilerContext, op vm.Operand) core.ValueType {
	if ctx == nil {
		return core.TypeUnknown
	}

	if op.IsConstant() {
		return valueTypeFromRuntime(ctx.Symbols.Constant(op))
	}

	return ctx.Types.Resolve(op)
}
