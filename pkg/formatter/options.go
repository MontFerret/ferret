package formatter

import (
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/formatter/internal"
)

// WithPrintWidth configures the maximum preferred line width. The width must
// be greater than zero.
func WithPrintWidth(size uint64) Option {
	return options.New(
		func(config *internal.Config, size uint64) {
			config.PrintWidth = size
		},
	).
		Named("print width").
		Validators(
			options.NotZero[uint64](),
		).
		Value(size).
		Build()
}

// WithTabWidth configures the number of spaces per indentation level. The
// width must be greater than zero.
func WithTabWidth(size uint64) Option {
	return options.New(
		func(config *internal.Config, size uint64) {
			config.TabWidth = size
		},
	).
		Named("tab width").
		Validators(
			options.NotZero[uint64](),
		).
		Value(size).
		Build()
}

// WithSingleQuote configures whether string literals use single quotes.
func WithSingleQuote(val bool) Option {
	return func(config *internal.Config) error {
		config.SingleQuote = val

		return nil
	}
}

// WithBracketSpacing configures spacing inside object and destructuring braces.
func WithBracketSpacing(val bool) Option {
	return func(config *internal.Config) error {
		config.BracketSpacing = val

		return nil
	}
}

// WithCaseMode configures keyword casing. The mode must be one of the
// supported CaseMode values.
func WithCaseMode(mode CaseMode) Option {
	return options.New(
		func(config *internal.Config, mode CaseMode) {
			config.CaseMode = mode
		},
	).
		Named("case mode").
		Validators(
			options.OneOf(
				CaseModeIgnore,
				CaseModeUpper,
				CaseModeLower,
			),
		).
		Value(mode).
		Build()
}
