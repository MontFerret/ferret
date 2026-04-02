package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type TypeFacts struct {
	session *CompilationSession
}

func NewTypeFacts(session *CompilationSession) *TypeFacts {
	return &TypeFacts{session: session}
}

func (f *TypeFacts) ValueTypeFromRuntime(value runtime.Value) core.ValueType {
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

func (f *TypeFacts) LiteralType(ctx fql.ILiteralContext) core.ValueType {
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

func (f *TypeFacts) LiteralValue(ctx fql.ILiteralContext) (runtime.Value, bool) {
	return literalValueOf(ctx)
}

func (f *TypeFacts) LiteralValueFromExpression(ctx fql.IExpressionContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	if p := ctx.Predicate(); p != nil {
		return f.LiteralValueFromPredicate(p)
	}

	return nil, false
}

func (f *TypeFacts) LiteralValueFromPredicate(ctx fql.IPredicateContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	if atom := ctx.ExpressionAtom(); atom != nil {
		return f.LiteralValueFromAtom(atom)
	}

	return nil, false
}

func (f *TypeFacts) LiteralValueFromAtom(ctx fql.IExpressionAtomContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	if lit := ctx.Literal(); lit != nil {
		return f.LiteralValue(lit)
	}

	return nil, false
}

func (f *TypeFacts) OperandType(op bytecode.Operand) core.ValueType {
	if f == nil || f.session == nil {
		return core.TypeUnknown
	}

	if op.IsConstant() {
		return f.ValueTypeFromRuntime(f.session.Symbols.Constant(op))
	}

	return f.session.Types.Resolve(op)
}

func (f *TypeFacts) InferBinaryResultType(op atomBinaryOperator, left, right bytecode.Operand) core.ValueType {
	leftType := f.OperandType(left)
	rightType := f.OperandType(right)

	switch op.opcode {
	case bytecode.OpAdd:
		switch {
		case leftType == core.TypeString || rightType == core.TypeString:
			return core.TypeString
		case leftType == core.TypeFloat || rightType == core.TypeFloat:
			if leftType == core.TypeInt || leftType == core.TypeFloat {
				if rightType == core.TypeInt || rightType == core.TypeFloat {
					return core.TypeFloat
				}
			}
		case leftType == core.TypeInt && rightType == core.TypeInt:
			return core.TypeInt
		}
	case bytecode.OpRegexp:
		return core.TypeBool
	}

	return core.TypeUnknown
}

func (f *TypeFacts) LoadConstant(value runtime.Value) bytecode.Operand {
	reg := f.session.Registers.Allocate()
	f.session.Emitter.EmitLoadConst(reg, f.session.Symbols.AddConstant(value))
	f.session.Types.Set(reg, f.ValueTypeFromRuntime(value))

	return reg
}

func (f *TypeFacts) EmitMoveAuto(dst, src bytecode.Operand) {
	srcType := f.OperandType(src)

	if srcType.IsUntracked() {
		f.session.Emitter.EmitPlainMove(dst, src)
	} else {
		f.session.Emitter.EmitMoveTracked(dst, src)
	}

	f.session.Types.Set(dst, srcType)
}
