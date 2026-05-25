package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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
	case queryModifierOne:
		return core.TypeAny
	default:
		return core.TypeList
	}
}

func queryOpcodeForModifier(modifier queryModifier) bytecode.Opcode {
	switch modifier {
	case queryModifierExists:
		return bytecode.OpQueryExists
	case queryModifierCount:
		return bytecode.OpQueryCount
	case queryModifierOne:
		return bytecode.OpQueryOne
	default:
		return bytecode.OpQuery
	}
}

func parseQueryModifier(text string) queryModifier {
	switch strings.ToLower(text) {
	case string(queryModifierExists):
		return queryModifierExists
	case string(queryModifierCount):
		return queryModifierCount
	case string(queryModifierOne):
		return queryModifierOne
	default:
		return queryModifierUnknown
	}
}
