// Package host assembles native engine host services and lifecycle hooks.
package host

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Config contains only the inputs needed to construct host services.
	// Resource ownership remains with the manager supplied by the engine.
	Config struct {
		Library    runtime.Library
		Network    ferretnet.Network
		Params     runtime.Params
		Encoding   *encoding.Registry
		FSRoot     string
		Logger     []logging.Option
		FSReadOnly bool
	}

	// Host is the finalized service snapshot shared by an engine and its plans.
	// Consumers retain the existing borrowing rules; Host does not close services.
	Host struct {
		Functions  *runtime.Functions
		Params     runtime.Params
		Encoding   *encoding.Registry
		Logger     logging.Logger
		FileSystem fs.FileSystem
		Network    ferretnet.Network
		FSReadOnly bool
	}

	hostContext struct {
		library    runtime.Library
		params     runtime.Params
		encoding   *encoding.Registry
		logger     logging.Logger
		fs         fs.FileSystem
		network    ferretnet.Network
		fsReadOnly bool
	}
)

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

func (h *hostContext) Build() (*Host, error) {
	funcs, err := h.library.Build()
	if err != nil {
		return nil, err
	}

	return &Host{
		Functions:  funcs,
		Params:     h.params.Clone(),
		Encoding:   h.encoding.Clone(),
		Logger:     h.logger,
		FileSystem: h.fs,
		Network:    h.network,
		FSReadOnly: h.fsReadOnly,
	}, nil
}

func newHostContext(opts Config, resources *resource.Manager) (*hostContext, error) {
	logger, err := logging.New(opts.Logger...)
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}

	rootFs, err := fs.New(fs.WithRoot(opts.FSRoot), fs.WithReadOnly(opts.FSReadOnly))
	if err != nil {
		return nil, err
	}

	if err := resources.Own(resource.FileSystem, rootFs.Close); err != nil {
		return nil, err
	}

	network := opts.Network
	if network == nil {
		network, err = ferretnet.New()
		if err != nil {
			return nil, fmt.Errorf("network: %w", err)
		}

		if err := resources.Own(resource.Network, func() error {
			ferretnet.CloseIdleNetworkConnections(network)

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &hostContext{
		library:    opts.Library,
		params:     opts.Params,
		encoding:   opts.Encoding,
		logger:     logger,
		fs:         rootFs,
		network:    network,
		fsReadOnly: opts.FSReadOnly,
	}, nil
}
