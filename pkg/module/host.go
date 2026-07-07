package module

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// HostContext exposes the host-scoped registries and services that modules can
// configure during engine bootstrap.
type HostContext interface {
	// Library returns the runtime library registry being assembled for the
	// engine.
	Library() runtime.Library
	// Params returns the default parameter set inherited by sessions created
	// from the engine.
	Params() runtime.Params
	// Encoding returns the codec registrar used to configure output encoding for
	// the engine.
	Encoding() encoding.CodecRegistrar
	// Logger returns the engine logger used by the engine and inherited by
	// sessions created from it.
	Logger() logging.Logger
	// FileSystem returns the file system used by the engine and derived
	// executions.
	FileSystem() fs.FileSystem
	// Network returns the network service used by the engine and derived
	// executions.
	Network() ferretnet.Network
}
