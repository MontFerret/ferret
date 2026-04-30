package ferret

import "github.com/MontFerret/ferret/v2/pkg/module"

type bootstrap struct {
	host  *hostContext
	hooks *hookRegistry
}

func newBootstrap(opts *options) (*bootstrap, error) {
	hostCtx, err := newHostContext(opts)
	if err != nil {
		return nil, err
	}

	return &bootstrap{
		host:  hostCtx,
		hooks: opts.hooks,
	}, nil
}

func (b *bootstrap) Host() module.HostContext {
	return b.host
}

func (b *bootstrap) Hooks() module.HookRegistrar {
	return b.hooks
}
