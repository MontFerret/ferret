package spec

import (
	"strings"
)

type (
	Skip struct {
		Active bool
		Reason string
	}

	Spec struct {
		Expression  string
		Description string
		SkipInfo    Skip
		RawOutput   bool
		DebugOutput bool
		Compile     Expectation
		Run         Expectation
	}
)

func New(expression string, desc ...string) Spec {
	return Spec{
		Expression:  expression,
		Description: strings.TrimSpace(strings.Join(desc, " ")),
	}
}

func (s Spec) Expect() ExpectationBuilder[Spec] {
	return NewExpectationBuilder(&s, func(s *Spec) *Spec {
		return s
	})
}

func (s Spec) Suffix(suffix string) Spec {
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

func (s Spec) Skip(reason ...string) Spec {
	s.SkipInfo.Active = true
	s.SkipInfo.Reason = strings.TrimSpace(strings.Join(reason, " "))

	return s
}

func (s Spec) Debug() Spec {
	s.DebugOutput = true
	return s
}

func (s Spec) String() string {
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
