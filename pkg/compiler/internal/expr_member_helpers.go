package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

// getImplicitToken picks the span anchor for implicit-current diagnostics.
func getImplicitToken(ctx antlr.ParserRuleContext) antlr.Token {
	switch v := ctx.(type) {
	case fql.IImplicitMemberExpressionStartContext:
		if dot := v.Dot(); dot != nil {
			return dot.GetSymbol()
		}
	case fql.IImplicitCurrentExpressionContext:
		if dot := v.Dot(); dot != nil {
			return dot.GetSymbol()
		}
	}

	return ctx.GetStart()
}

func isSimpleMemberPathChain(segments []fql.IMemberExpressionPathContext) bool {
	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)

		if p.PropertyName() == nil && p.ComputedPropertyName() == nil {
			return false
		}
	}

	return true
}

func memberLoadOpcode(srcType core.ValueType, constOperand, optional bool) bytecode.Opcode {
	switch srcType {
	case core.TypeArray:
		if constOperand {
			if optional {
				return bytecode.OpLoadIndexOptionalConst
			}

			return bytecode.OpLoadIndexConst
		}

		if optional {
			return bytecode.OpLoadIndexOptional
		}

		return bytecode.OpLoadIndex
	case core.TypeObject:
		if constOperand {
			if optional {
				return bytecode.OpLoadKeyOptionalConst
			}

			return bytecode.OpLoadKeyConst
		}

		if optional {
			return bytecode.OpLoadKeyOptional
		}

		return bytecode.OpLoadKey
	default:
		if constOperand {
			if optional {
				return bytecode.OpLoadPropertyOptionalConst
			}

			return bytecode.OpLoadPropertyConst
		}

		if optional {
			return bytecode.OpLoadPropertyOptional
		}

		return bytecode.OpLoadProperty
	}
}

func splitArrayOperatorTail(segments []fql.IMemberExpressionPathContext) ([]fql.IMemberExpressionPathContext, []fql.IMemberExpressionPathContext) {
	if len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)

		if p.ArrayContraction() != nil || p.ArrayExpansion() != nil || p.ArrayQuestionMark() != nil {
			return nil, segments
		}
	}

	return segments, nil
}

// splitTerminalArrayContractionTail hoists only the final contraction segment.
// Earlier array operators stay in the per-element tail so existing projection semantics remain unchanged.
func splitTerminalArrayContractionTail(segments []fql.IMemberExpressionPathContext) ([]fql.IMemberExpressionPathContext, fql.IArrayContractionContext) {
	if len(segments) == 0 {
		return nil, nil
	}

	last := segments[len(segments)-1].(*fql.MemberExpressionPathContext)
	contraction := last.ArrayContraction()
	if contraction == nil {
		return segments, nil
	}

	return segments[:len(segments)-1], contraction
}

func isFilterOnlyInline(inline fql.IInlineExpressionContext) bool {
	if inline == nil {
		return false
	}

	return inline.InlineFilter() != nil && inline.InlineLimit() == nil && inline.InlineReturn() == nil
}

func nextArrayExpansion(segments []fql.IMemberExpressionPathContext) (fql.IArrayExpansionContext, []fql.IMemberExpressionPathContext) {
	if len(segments) == 0 {
		return nil, segments
	}

	p := segments[0].(*fql.MemberExpressionPathContext)

	if expansion := p.ArrayExpansion(); expansion != nil {
		return expansion, segments[1:]
	}

	return nil, segments
}

func dropIdentityExpansions(segments []fql.IMemberExpressionPathContext) []fql.IMemberExpressionPathContext {
	for len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()

		if expansion == nil {
			break
		}

		if expansion.InlineExpression() != nil {
			break
		}

		segments = segments[1:]
	}

	return segments
}

func collectFilterOnlyTail(segments []fql.IMemberExpressionPathContext) ([]fql.IExpressionContext, []fql.IMemberExpressionPathContext) {
	extraFilters := make([]fql.IExpressionContext, 0)
	rest := segments

	for len(rest) > 0 {
		p := rest[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()
		if expansion == nil {
			break
		}

		inline := expansion.InlineExpression()
		if inline == nil {
			rest = rest[1:]
			continue
		}

		if !isFilterOnlyInline(inline) {
			break
		}

		filter := inline.InlineFilter()
		if filter != nil {
			extraFilters = append(extraFilters, filter.Expression())
		}

		rest = rest[1:]
	}

	return extraFilters, rest
}

func arrayContractionDepth(ctx fql.IArrayContractionContext) int {
	if ctx == nil {
		return 1
	}

	count := len(ctx.GetStars())

	if count > 1 {
		return count - 1
	}

	return 1
}
