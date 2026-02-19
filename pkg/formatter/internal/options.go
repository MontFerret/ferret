package internal

type (
	Options struct {
		printWidth     uint64
		tabWidth       uint64
		singleQuote    bool
		bracketSpacing bool
		caseMode       CaseMode
	}

	Option func(*Options)
)

func DefaultOptions() *Options {
	return &Options{
		printWidth:     80, // Characters
		tabWidth:       4,  // Spaces
		singleQuote:    false,
		bracketSpacing: true,
		caseMode:       CaseModeUpper,
	}
}

func WithPrintWidth(size uint64) Option {
	return func(opts *Options) {
		opts.printWidth = size
	}
}

func WithTabWidth(size uint64) Option {
	return func(opts *Options) {
		opts.tabWidth = size
	}
}

func WithSingleQuote(val bool) Option {
	return func(opts *Options) {
		opts.singleQuote = val
	}
}

func WithBracketSpacing(val bool) Option {
	return func(opts *Options) {
		opts.bracketSpacing = val
	}
}

func WithCaseMode(mode CaseMode) Option {
	return func(opts *Options) {
		opts.caseMode = mode
	}
}
