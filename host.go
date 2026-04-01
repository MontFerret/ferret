package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// HostContext exposes host-scoped registries and services during engine bootstrap.
	HostContext interface {
		// Library returns the runtime library registry being assembled for the engine.
		Library() runtime.Library
		// Params returns the default parameter set inherited by new sessions.
		Params() runtime.Params
		// Encoding returns the codec registrar used for output encoding.
		Encoding() encoding.CodecRegistrar
		// Logger returns the engine logger used for derived sessions.
		Logger() logging.Logger
		// FileSystem returns the file system configured for the engine.
		FileSystem() fs.FileSystem
	}

	host struct {
		functions *runtime.Functions
		params    runtime.Params
		encoding  *encoding.Registry
		logger    logging.Logger
		fs        fs.FileSystem
	}

	hostContext struct {
		library  runtime.Library
		params   runtime.Params
		encoding *encoding.Registry
		logger   logging.Logger
		fs       fs.FileSystem
	}
)

func newHostContext(opts *options) (*hostContext, error) {
	rootFs, err := fs.New(fs.WithRoot(opts.fsRoot), fs.WithReadOnly(opts.fsReadOnly))

	if err != nil {
		return nil, err
	}

	return &hostContext{
		library:  opts.library,
		params:   opts.params,
		encoding: opts.encoding,
		logger:   logging.New(opts.logger...),
		fs:       rootFs,
	}, nil
}

func (h *hostContext) Logger() logging.Logger {
	return h.logger
}

func (h *hostContext) FileSystem() fs.FileSystem {
	return h.fs
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
	}, nil
}
