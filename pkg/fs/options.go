package fs

import "github.com/ziflex/go-options"

type (
	Option = options.Option[config]

	config struct {
		Root     string
		ReadOnly bool
	}
)

func WithRoot(root string) Option {
	return func(opts *config, _ options.Report) {
		opts.Root = root
	}
}

func WithReadOnly(readOnly bool) Option {
	return func(opts *config, _ options.Report) {
		opts.ReadOnly = readOnly
	}
}
