package spec

import "strings"

type (
	Skip struct {
		Reason string
		Active bool
	}

	BaseSpec struct {
		Expression  string
		Description string
		SkipInfo    Skip
		DebugOutput bool
	}
)

func (s BaseSpec) Suffix(suffix string) BaseSpec {
	suffix = strings.TrimSpace(suffix)

	if suffix == "" {
		return s
	}

	if s.Description == "" {
		s.Description = suffix
		return s
	}

	s.Description = s.Description + " - " + suffix
	return s
}

func (s BaseSpec) Skip(reason ...string) BaseSpec {
	s.SkipInfo.Active = true
	s.SkipInfo.Reason = strings.TrimSpace(strings.Join(reason, " "))

	return s
}

func (s BaseSpec) Debug() BaseSpec {
	s.DebugOutput = true
	return s
}

func (s BaseSpec) SuiteName(suite string) string {
	return suite + "/" + s.String()
}

func (s BaseSpec) String() string {
	if s.Description != "" {
		return strings.TrimSpace(s.Description)
	}

	exp := strings.TrimSpace(s.Expression)
	exp = strings.ReplaceAll(exp, "\n", " ")
	exp = strings.ReplaceAll(exp, "\t", " ")
	// Replace multiple spaces with a single space
	exp = strings.Join(strings.Fields(exp), " ")

	return exp
}

func (s BaseSpec) Merge(other BaseSpec) BaseSpec {
	out := s

	if other.Expression != "" {
		out.Expression = other.Expression
	}

	if other.Description != "" {
		out.Description = other.Description
	}

	if other.SkipInfo.Active {
		out.SkipInfo = other.SkipInfo
	}

	if other.DebugOutput {
		out.DebugOutput = other.DebugOutput
	}

	return out
}
