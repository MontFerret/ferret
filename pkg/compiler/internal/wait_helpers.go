package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func resolveRetryBackoff(ctx *CompilerContext, clause fql.IRecoveryRetryBackoffClauseContext) (core.RetryBackoff, bool) {
	if clause == nil {
		return core.RetryBackoffNone, true
	}

	kind := clause.RecoveryRetryBackoffKind()
	if kind == nil {
		reportInvalidRecoveryTail(ctx, clause, "Expected backoff kind after 'BACKOFF' in retry policy", "Use BACKOFF CONSTANT, BACKOFF LINEAR, or BACKOFF EXPONENTIAL.")
		return core.RetryBackoffNone, false
	}

	raw := ""

	switch {
	case kind.Identifier() != nil:
		raw = kind.Identifier().GetText()
	case kind.StringLiteral() != nil:
		if parsed, ok := parseStringLiteralConst(kind.StringLiteral()); ok {
			raw = parsed.String()
		}
	case kind.None() != nil:
		raw = kind.None().GetText()
	}

	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "CONSTANT":
		return core.RetryBackoffNone, true
	case "LINEAR":
		return core.RetryBackoffLinear, true
	case "EXPONENTIAL":
		return core.RetryBackoffExponential, true
	default:
		reportInvalidRecoveryTail(ctx, kind, "Unknown BACKOFF strategy", "Use one of: CONSTANT, LINEAR, EXPONENTIAL.")
		return core.RetryBackoffNone, false
	}
}
