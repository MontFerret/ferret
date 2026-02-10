package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func valueTypeFromRuntime(value runtime.Value) core.ValueType {
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

func operandType(ctx *CompilerContext, op vm.Operand) core.ValueType {
	if ctx == nil {
		return core.TypeUnknown
	}

	if op.IsConstant() {
		return valueTypeFromRuntime(ctx.Symbols.Constant(op))
	}

	return ctx.Types.Resolve(op)
}
