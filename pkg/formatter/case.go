package formatter

import "github.com/MontFerret/ferret/v2/pkg/formatter/internal"

// CaseMode controls how the formatter emits FQL keywords.
type CaseMode = internal.CaseMode

const (
	// CaseModeIgnore emits keyword text without applying case conversion.
	CaseModeIgnore = internal.CaseModeIgnore
	// CaseModeUpper emits uppercase keywords.
	CaseModeUpper = internal.CaseModeUpper
	// CaseModeLower emits lowercase keywords.
	CaseModeLower = internal.CaseModeLower
)
