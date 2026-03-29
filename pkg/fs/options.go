package fs

type (
	Option func(*options)

	options struct {
		Root     string
		ReadOnly bool
	}
)

func WithRoot(root string) Option {
	return func(opts *options) {
		opts.Root = root
	}
}

func WithReadOnly(readOnly bool) Option {
	return func(opts *options) {
		opts.ReadOnly = readOnly
	}
}
