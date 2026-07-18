package ferret

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	host struct {
		functions *runtime.Functions
		params    runtime.Params
		encoding  *encoding.Registry
		logger    logging.Logger
		fs        fs.FileSystem
		network   ferretnet.Network
	}

	hostContext struct {
		library  runtime.Library
		params   runtime.Params
		encoding *encoding.Registry
		logger   logging.Logger
		fs       fs.FileSystem
		network  ferretnet.Network
	}
)

func newHostContext(opts *options) (*hostContext, error) {
	rootFs, err := fs.New(fs.WithRoot(opts.fsRoot), fs.WithReadOnly(opts.fsReadOnly))

	if err != nil {
		return nil, err
	}

	network := opts.network
	if network == nil {
		network, err = ferretnet.New()
		if err != nil {
			return nil, fmt.Errorf("network: %w", err)
		}
	}

	return &hostContext{
		library:  opts.library,
		params:   opts.params,
		encoding: opts.encoding,
		logger:   logging.New(opts.logger...),
		fs:       rootFs,
		network:  network,
	}, nil
}

func (h *hostContext) Logger() logging.Logger {
	return h.logger
}

func (h *hostContext) FileSystem() fs.FileSystem {
	return h.fs
}

func (h *hostContext) Network() ferretnet.Network {
	return h.network
}

func (h *hostContext) Params() runtime.Params {
	return h.params
}

func (h *hostContext) Library() runtime.Library {
	return h.library
}

func (h *hostContext) Encoding() encoding.CodecRegistrar {
	return h.encoding
}

func (h *hostContext) Build() (*host, error) {
	funcs, err := h.library.Build()

	if err != nil {
		return nil, err
	}

	return &host{
		functions: funcs,
		params:    h.params.Clone(),
		encoding:  h.encoding.Clone(),
		logger:    h.logger,
		fs:        h.fs,
		network:   h.network,
	}, nil
}
