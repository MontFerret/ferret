package spec

import "strings"

type (
	Skip struct {
		Reason string
		Active bool
	}

	BaseSpec struct {
		Input       Input
		Description string
		SkipInfo    Skip
		DebugOutput bool
	}
)

func NewBaseSpec(expr string, desc ...string) BaseSpec {
	return NewBaseSpecWith(NewExpressionInput(expr), desc...)
}

func NewBaseSpecWith(input Input, desc ...string) BaseSpec {
	return BaseSpec{
		Input:       input,
		Description: strings.Join(desc, " "),
	}
}

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
	return s.Input.String()
}

func (s BaseSpec) Merge(other BaseSpec) BaseSpec {
	out := s

	out.Input = s.Input.Merge(other.Input)

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
