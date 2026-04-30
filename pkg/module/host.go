package module

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// HostContext exposes host-scoped registries and services during engine bootstrap.
type HostContext interface {
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
