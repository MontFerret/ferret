package formatter

import "strings"

// CaseMode controls how the formatter emits FQL keywords.
type CaseMode uint64

const (
	// CaseModeIgnore emits keyword text without applying case conversion.
	CaseModeIgnore CaseMode = iota
	// CaseModeUpper emits uppercase keywords.
	CaseModeUpper
	// CaseModeLower emits lowercase keywords.
	CaseModeLower
)

func applyCase(mode CaseMode, value string) string {
	switch mode {
	case CaseModeUpper:
		return strings.ToUpper(value)
	case CaseModeLower:
		return strings.ToLower(value)
	default:
		return value
	}
}
