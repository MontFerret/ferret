package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	HostConfigurer interface {
		Params() runtime.Params
		Library() runtime.Library
		Encoding() encoding.CodecRegistrar
	}

	host struct {
		functions *runtime.Functions
		params    runtime.Params
		encoding  *encoding.Registry
		logging   runtime.LogSettings
	}

	hostBuilder struct {
		library  runtime.Library
		params   runtime.Params
		encoding *encoding.Registry
		logging  runtime.LogSettings
	}
)

func newHostBuilder(opts *options) *hostBuilder {
	return &hostBuilder{
		library:  opts.library,
		params:   opts.params,
		logging:  opts.logging,
		encoding: opts.encoding,
	}
}

func (h *hostBuilder) Params() runtime.Params {
	return h.params
}

func (h *hostBuilder) Library() runtime.Library {
	return h.library
}

func (h *hostBuilder) Encoding() encoding.CodecRegistrar {
	return h.encoding
}

func (h *hostBuilder) Build() (*host, error) {
	funcs, err := h.library.Build()

	if err != nil {
		return nil, err
	}

	return &host{
		logging:   h.logging,
		functions: funcs,
		params:    h.params.Clone(),
		encoding:  h.encoding.Clone(),
	}, nil
}
