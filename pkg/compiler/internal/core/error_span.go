package core

import "github.com/MontFerret/ferret/pkg/file"

type ErrorSpan struct {
	Span  file.Span
	Label string
	Main  bool
}

func NewErrorSpan(span file.Span, label string, main bool) ErrorSpan {
	return ErrorSpan{
		Span:  span,
		Label: label,
		Main:  main,
	}
}

func NewMainErrorSpan(span file.Span, label string) ErrorSpan {
	return NewErrorSpan(span, label, true)
}

func NewSecondaryErrorSpan(span file.Span, label string) ErrorSpan {
	return NewErrorSpan(span, label, false)
}
