package diagnostics

import "github.com/MontFerret/ferret/v2/pkg/source"

type ErrorSpan struct {
	Label string
	Span  source.Span
	Main  bool
}

func NewErrorSpan(span source.Span, label string, main bool) ErrorSpan {
	return ErrorSpan{
		Span:  span,
		Label: label,
		Main:  main,
	}
}

func NewMainErrorSpan(span source.Span, label string) ErrorSpan {
	return NewErrorSpan(span, label, true)
}

func NewSecondaryErrorSpan(span source.Span, label string) ErrorSpan {
	return NewErrorSpan(span, label, false)
}
