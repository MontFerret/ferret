package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	concatAtomPart struct {
		atom    fql.IExpressionAtomContext
		literal runtime.String
	}

	concatOperandSegment struct {
		literal runtime.String
		operand bytecode.Operand
	}
)

func buildConcatOperandSegmentsFromExpression(c *ExprCompiler, expr fql.IExpressionContext) []concatOperandSegment {
	if c == nil || expr == nil {
		return nil
	}

	if val, ok := tryConstConcatStringFromExpression(expr); ok {
		return []concatOperandSegment{{literal: val}}
	}

	if predicate := expr.Predicate(); predicate != nil {
		if atom := predicate.ExpressionAtom(); atom != nil {
			if segments, ok := buildConcatOperandSegmentsFromAtom(c, atom); ok {
				return segments
			}
		}
	}

	return []concatOperandSegment{{operand: c.Compile(expr)}}
}

func buildConcatOperandSegmentsFromAtom(c *ExprCompiler, atom fql.IExpressionAtomContext) ([]concatOperandSegment, bool) {
	if c == nil || !isConcatCompatibleAdditive(c.ctx, c.facts, atom) {
		return nil, false
	}

	parts := make([]concatAtomPart, 0, 4)
	collectConcatAtomParts(c.ctx, c.facts, atom, &parts)

	segments := make([]concatOperandSegment, 0, len(parts))

	for _, part := range parts {
		if part.atom != nil {
			segments = append(segments, concatOperandSegment{
				operand: c.compileAtom(part.atom),
			})

			continue
		}

		segments = append(segments, concatOperandSegment{literal: part.literal})
	}

	return segments, true
}

func collectConcatAtomParts(ctx *CompilationSession, facts *TypeFacts, atom fql.IExpressionAtomContext, parts *[]concatAtomPart) {
	if atom == nil || parts == nil {
		return
	}

	if !isConcatCompatibleAdditive(ctx, facts, atom) {
		appendConcatAtomPart(ctx, facts, atom, parts)
		return
	}

	appendConcatAtomPart(ctx, facts, atom.ExpressionAtom(0), parts)
	appendConcatAtomPart(ctx, facts, atom.ExpressionAtom(1), parts)
}

func appendConcatAtomPart(ctx *CompilationSession, facts *TypeFacts, atom fql.IExpressionAtomContext, parts *[]concatAtomPart) {
	if atom == nil || parts == nil {
		return
	}

	if isConcatCompatibleAdditive(ctx, facts, atom) {
		collectConcatAtomParts(ctx, facts, atom, parts)
		return
	}

	if val, ok := tryConstConcatStringFromAtom(atom); ok {
		*parts = append(*parts, concatAtomPart{literal: val})
		return
	}

	*parts = append(*parts, concatAtomPart{atom: atom})
}

func isConcatCompatibleAdditive(ctx *CompilationSession, facts *TypeFacts, atom fql.IExpressionAtomContext) bool {
	if ctx == nil || atom == nil {
		return false
	}

	op := atom.AdditiveOperator()
	if op == nil || op.GetText() != "+" {
		return false
	}

	return inferConcatAtomType(ctx, facts, atom) == core.TypeString
}

func inferConcatExpressionType(ctx *CompilationSession, facts *TypeFacts, expr fql.IExpressionContext) core.ValueType {
	if expr == nil {
		return core.TypeUnknown
	}

	if predicate := expr.Predicate(); predicate != nil {
		return inferConcatPredicateType(ctx, facts, predicate)
	}

	return core.TypeUnknown
}

func inferConcatPredicateType(ctx *CompilationSession, facts *TypeFacts, predicate fql.IPredicateContext) core.ValueType {
	if predicate == nil {
		return core.TypeUnknown
	}

	if atom := predicate.ExpressionAtom(); atom != nil {
		return inferConcatAtomType(ctx, facts, atom)
	}

	return core.TypeUnknown
}

func inferConcatAtomType(ctx *CompilationSession, facts *TypeFacts, atom fql.IExpressionAtomContext) core.ValueType {
	if ctx == nil || atom == nil {
		return core.TypeUnknown
	}

	if lit := atom.Literal(); lit != nil {
		return facts.LiteralType(lit)
	}

	if inner := atom.Expression(); inner != nil {
		return inferConcatExpressionType(ctx, facts, inner)
	}

	if v := atom.Variable(); v != nil {
		if binding, ok := ctx.Function.Symbols.ResolveBinding(v.GetText()); ok {
			return binding.Type
		}

		return core.TypeUnknown
	}

	if atom.Param() != nil || atom.FunctionCallExpression() != nil || atom.MatchExpression() != nil ||
		atom.MemberExpression() != nil || atom.ImplicitMemberExpression() != nil ||
		atom.ImplicitCurrentExpression() != nil || atom.DispatchExpression() != nil ||
		atom.WaitForExpression() != nil {
		return core.TypeAny
	}

	if atom.RangeOperator() != nil || atom.ForExpression() != nil {
		return core.TypeList
	}

	if op := atom.AdditiveOperator(); op != nil {
		left := atom.ExpressionAtom(0)
		right := atom.ExpressionAtom(1)

		leftType := inferConcatAtomType(ctx, facts, left)
		rightType := inferConcatAtomType(ctx, facts, right)

		if op.GetText() == "+" {
			switch {
			case leftType == core.TypeString || rightType == core.TypeString:
				return core.TypeString
			case concatAnchorFromAtom(left) || concatAnchorFromAtom(right):
				return core.TypeString
			case leftType == core.TypeFloat || rightType == core.TypeFloat:
				if isNumericType(leftType) && isNumericType(rightType) {
					return core.TypeFloat
				}
			case leftType == core.TypeInt && rightType == core.TypeInt:
				return core.TypeInt
			}
		}

		if op.GetText() == "-" {
			if leftType == core.TypeFloat || rightType == core.TypeFloat {
				if isNumericType(leftType) && isNumericType(rightType) {
					return core.TypeFloat
				}
			}

			if leftType == core.TypeInt && rightType == core.TypeInt {
				return core.TypeInt
			}
		}

		return core.TypeUnknown
	}

	if op := atom.MultiplicativeOperator(); op != nil {
		left := inferConcatAtomType(ctx, facts, atom.ExpressionAtom(0))
		right := inferConcatAtomType(ctx, facts, atom.ExpressionAtom(1))

		if left == core.TypeFloat || right == core.TypeFloat {
			if isNumericType(left) && isNumericType(right) {
				return core.TypeFloat
			}
		}

		if left == core.TypeInt && right == core.TypeInt {
			return core.TypeInt
		}

		return core.TypeUnknown
	}

	return core.TypeUnknown
}

func concatAnchorFromAtom(atom fql.IExpressionAtomContext) bool {
	val, ok := tryConcatConstValueFromAtom(atom)
	if !ok {
		return false
	}

	return guaranteesConcatStringResult(val)
}

func isNumericType(typ core.ValueType) bool {
	return typ == core.TypeInt || typ == core.TypeFloat
}

func guaranteesConcatStringResult(val runtime.Value) bool {
	if val == runtime.None {
		return true
	}

	switch val.(type) {
	case runtime.String, runtime.Boolean:
		return true
	default:
		return false
	}
}

func tryConcatConstValueFromExpression(expr fql.IExpressionContext) (runtime.Value, bool) {
	if expr == nil {
		return nil, false
	}

	if predicate := expr.Predicate(); predicate != nil {
		return tryConcatConstValueFromPredicate(predicate)
	}

	return nil, false
}

func tryConcatConstValueFromPredicate(predicate fql.IPredicateContext) (runtime.Value, bool) {
	if predicate == nil {
		return nil, false
	}

	if atom := predicate.ExpressionAtom(); atom != nil {
		return tryConcatConstValueFromAtom(atom)
	}

	return nil, false
}

func tryConcatConstValueFromAtom(atom fql.IExpressionAtomContext) (runtime.Value, bool) {
	if atom == nil {
		return nil, false
	}

	if lit := atom.Literal(); lit != nil {
		return literalValueOf(lit)
	}

	if inner := atom.Expression(); inner != nil {
		return tryConcatConstValueFromExpression(inner)
	}

	return nil, false
}

func tryConstConcatStringFromExpression(expr fql.IExpressionContext) (runtime.String, bool) {
	val, ok := tryConcatConstValueFromExpression(expr)
	if !ok {
		return runtime.EmptyString, false
	}

	return tryConstConcatStringFromValue(val)
}

func tryConstConcatStringFromAtom(atom fql.IExpressionAtomContext) (runtime.String, bool) {
	val, ok := tryConcatConstValueFromAtom(atom)
	if !ok {
		return runtime.EmptyString, false
	}

	return tryConstConcatStringFromValue(val)
}

func tryConstConcatStringFromValue(val runtime.Value) (runtime.String, bool) {
	if val == nil {
		return runtime.EmptyString, false
	}

	if val == runtime.None {
		return runtime.EmptyString, true
	}

	switch val.(type) {
	case runtime.String, runtime.Int, runtime.Float, runtime.Boolean:
		return runtime.ToString(val), true
	default:
		return runtime.EmptyString, false
	}
}

func emitConcatOperandSegments(ctx *CompilationSession, facts *TypeFacts, parts []concatOperandSegment) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	merged := mergeConcatOperandSegments(parts)
	if len(merged) == 0 {
		return facts.LoadConstant(runtime.EmptyString)
	}

	if len(merged) == 1 {
		part := merged[0]
		if part.operand == bytecode.NoopOperand {
			return facts.LoadConstant(part.literal)
		}

		start := ensureConcatRegister(ctx, facts, part.operand)
		if !start.IsRegister() {
			return bytecode.NoopOperand
		}

		ctx.Program.Emitter.EmitABC(bytecode.OpConcat, start, start, bytecode.Operand(1))
		ctx.Function.Types.Set(start, core.TypeString)

		return start
	}

	if len(merged) == 2 && merged[0].operand != bytecode.NoopOperand && merged[1].operand == bytecode.NoopOperand {
		left := ensureConcatRegister(ctx, facts, merged[0].operand)
		if !left.IsRegister() {
			return bytecode.NoopOperand
		}

		dst := ctx.Function.Registers.Allocate()
		ctx.Program.Emitter.EmitABC(bytecode.OpAddConst, dst, left, ctx.Function.Symbols.AddConstant(merged[1].literal))
		ctx.Function.Types.Set(dst, core.TypeString)

		return dst
	}

	seq := ctx.Function.Registers.AllocateSequence(len(merged))

	for i, part := range merged {
		target := seq[i]

		if part.operand != bytecode.NoopOperand {
			if !loadConcatOperandIntoRegister(ctx, facts, target, part.operand) {
				return bytecode.NoopOperand
			}

			continue
		}

		ctx.Program.Emitter.EmitLoadConst(target, ctx.Function.Symbols.AddConstant(part.literal))
		ctx.Function.Types.Set(target, core.TypeString)
	}

	dst := seq[0]
	ctx.Program.Emitter.EmitABC(bytecode.OpConcat, dst, seq[0], bytecode.Operand(len(seq)))
	ctx.Function.Types.Set(dst, core.TypeString)

	return dst
}

func mergeConcatOperandSegments(parts []concatOperandSegment) []concatOperandSegment {
	if len(parts) == 0 {
		return nil
	}

	merged := make([]concatOperandSegment, 0, len(parts))
	var literal strings.Builder

	flushLiteral := func() {
		if literal.Len() == 0 {
			return
		}

		merged = append(merged, concatOperandSegment{
			literal: runtime.NewString(literal.String()),
		})
		literal.Reset()
	}

	for _, part := range parts {
		if part.operand == bytecode.NoopOperand {
			if part.literal == runtime.EmptyString {
				continue
			}

			literal.WriteString(part.literal.String())
			continue
		}

		flushLiteral()
		merged = append(merged, concatOperandSegment{operand: part.operand})
	}

	flushLiteral()

	return merged
}

func ensureConcatRegister(ctx *CompilationSession, facts *TypeFacts, op bytecode.Operand) bytecode.Operand {
	if op == bytecode.NoopOperand || ctx == nil {
		return bytecode.NoopOperand
	}

	if op.IsRegister() {
		return op
	}

	reg := ctx.Function.Registers.Allocate()
	ctx.Program.Emitter.EmitLoadConst(reg, op)
	ctx.Function.Types.Set(reg, facts.OperandType(op))

	return reg
}

func loadConcatOperandIntoRegister(ctx *CompilationSession, facts *TypeFacts, target, op bytecode.Operand) bool {
	if ctx == nil || target == bytecode.NoopOperand || op == bytecode.NoopOperand {
		return false
	}

	if op.IsConstant() {
		ctx.Program.Emitter.EmitLoadConst(target, op)
		ctx.Function.Types.Set(target, facts.OperandType(op))
		return true
	}

	if !op.IsRegister() {
		return false
	}

	if !op.Equals(target) {
		ctx.Program.Emitter.EmitMove(target, op)
	}

	ctx.Function.Types.Set(target, facts.OperandType(op))

	return true
}
