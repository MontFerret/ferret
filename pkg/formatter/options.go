package formatter

import "github.com/ziflex/go-options"

// WithPrintWidth configures the maximum preferred line width. The width must
// be greater than zero.
func WithPrintWidth(size uint64) Option {
	return options.New(
		func(config *config, size uint64) {
			config.printWidth = size
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
		func(config *config, size uint64) {
			config.tabWidth = size
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
	return func(config *config) error {
		config.singleQuote = val

		return nil
	}
}

// WithBracketSpacing configures spacing inside object and destructuring braces.
func WithBracketSpacing(val bool) Option {
	return func(config *config) error {
		config.bracketSpacing = val

		return nil
	}
}

// WithCaseMode configures keyword casing. The mode must be one of the
// supported CaseMode values.
func WithCaseMode(mode CaseMode) Option {
	return options.New(
		func(config *config, mode CaseMode) {
			config.caseMode = mode
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
