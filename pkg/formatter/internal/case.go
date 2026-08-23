package internal

import "strings"

type CaseMode uint64

const (
	CaseModeIgnore CaseMode = iota
	CaseModeUpper
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
