package internal

import (
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func literalValueOf(ctx fql.ILiteralContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	switch {
	case ctx.NoneLiteral() != nil:
		return runtime.None, true
	case ctx.StringLiteral() != nil:
		return parseStringLiteralConst(ctx.StringLiteral())
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
