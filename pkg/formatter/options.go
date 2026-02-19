package formatter

import "github.com/MontFerret/ferret/v2/pkg/formatter/internal"

type Option = internal.Option

func WithPrintWidth(size uint64) Option {
	return internal.WithPrintWidth(size)
}

func WithTabWidth(size uint64) Option {
	return internal.WithTabWidth(size)
}

func WithSingleQuote(val bool) Option {
	return internal.WithSingleQuote(val)
}

func WithBracketSpacing(val bool) Option {
	return internal.WithBracketSpacing(val)
}

func WithCaseMode(mode CaseMode) Option {
	return internal.WithCaseMode(mode)
}
