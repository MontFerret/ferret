package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func queryModifierName(ctx fql.IQueryModifierContext) queryModifier {
	if ctx == nil {
		return queryModifierUnknown
	}

	return parseQueryModifier(ctx.GetText())
}

func queryResultTypeForModifier(modifier queryModifier) core.ValueType {
	switch modifier {
	case queryModifierExists:
		return core.TypeBool
	case queryModifierCount:
		return core.TypeInt
	case queryModifierAny, queryModifierValue, queryModifierOne:
		return core.TypeAny
	default:
		return core.TypeList
	}
}

func parseQueryModifier(text string) queryModifier {
	switch strings.ToLower(text) {
	case string(queryModifierExists):
		return queryModifierExists
	case string(queryModifierCount):
		return queryModifierCount
	case string(queryModifierAny):
		return queryModifierAny
	case string(queryModifierValue):
		return queryModifierValue
	case string(queryModifierOne):
		return queryModifierOne
	default:
		return queryModifierUnknown
	}
}
