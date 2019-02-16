package drivers

type (
	options struct {
		defaultDriver string
	}

	Option func(drv Driver, opts *options)
)

func AsDefault() Option {
	return func(drv Driver, opts *options) {
		opts.defaultDriver = drv.Name()
	}
}
