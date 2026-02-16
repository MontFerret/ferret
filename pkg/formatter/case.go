package formatter

type CaseMode uint64

const (
	CaseModeIgnore CaseMode = iota
	CaseModeUpper
	CaseModeLower
)
