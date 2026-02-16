package formatter

type (
	options struct {
		printWidth     uint64
		tabWidth       uint64
		singleQuote    bool
		bracketSpacing bool
		caseMode       CaseMode
	}

	Option func(*options)
)

func defaultOptions() *options {
	return &options{
		printWidth:     80, // Characters
		tabWidth:       4,  // Spaces
		singleQuote:    false,
		bracketSpacing: true,
		caseMode:       CaseModeUpper,
	}
}

func WithPrintWidth(size uint64) Option {
	return func(opts *options) {
		opts.printWidth = size
	}
}

func WithTabWidth(size uint64) Option {
	return func(opts *options) {
		opts.tabWidth = size
	}
}

func WithSingleQuote(val bool) Option {
	return func(opts *options) {
		opts.singleQuote = val
	}
}

func WithBracketSpacing(val bool) Option {
	return func(opts *options) {
		opts.bracketSpacing = val
	}
}

func WithCaseMode(mode CaseMode) Option {
	return func(opts *options) {
		opts.caseMode = mode
	}
}
