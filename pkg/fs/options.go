package fs

import "github.com/ziflex/go-options"

type (
	config struct {
		Root     string
		ReadOnly bool
	}

	Option = options.Option[config]
)

func defaultConfig() config {
	return config{
		Root:     "",
		ReadOnly: false,
	}
}

func WithRoot(root string) Option {
	return func(config *config) error {
		config.Root = root

		return nil
	}
}

func WithReadOnly(readOnly bool) Option {
	return func(config *config) error {
		config.ReadOnly = readOnly

		return nil
	}
}
