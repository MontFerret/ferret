package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	// DebugServices implements debugger.SessionServices for native executions.
	// The debugger serializes its use and owns once-only closure; retained VM state
	// is owned by the debugger, independently of these host resources and permit.
	DebugServices struct {
		logger            logging.Logger
		hooks             *host.SessionHooks
		filesystem        fs.FileSystem
		network           ferretnet.Network
		encoding          *encoding.Registry
		limiter           *Limiter
		resources         *resource.Manager
		outputContentType string
	}

	// DebugServicesConfig supplies the borrowed services and output settings for
	// one debug session. Resource ownership is tracked separately by its manager.
	DebugServicesConfig struct {
		Logger            logging.Logger
		Hooks             *host.SessionHooks
		FileSystem        fs.FileSystem
		Network           ferretnet.Network
		Encoding          *encoding.Registry
		OutputContentType string
	}
)

// NewDebugServices binds session services to an already-acquired limiter permit
// and host-resource manager. Construction does not acquire or release resources;
// the caller retains rollback until the debugger session is successfully built.
func NewDebugServices(config DebugServicesConfig, limiter *Limiter, resources *resource.Manager) *DebugServices {
	return &DebugServices{
		logger:            config.Logger,
		hooks:             config.Hooks,
		filesystem:        config.FileSystem,
		network:           config.Network,
		encoding:          config.Encoding,
		outputContentType: config.OutputContentType,
		limiter:           limiter,
		resources:         resources,
	}
}

// BeforeRun runs the configured before-run hook chain.
func (s *DebugServices) BeforeRun(ctx context.Context) (context.Context, error) {
	return s.hooks.RunBeforeRun(ctx)
}

// AfterRun runs the configured after-run hooks with the primary execution error.
func (s *DebugServices) AfterRun(ctx context.Context, runErr error) error {
	return s.hooks.RunAfterRun(ctx, runErr)
}

// ExtendContext injects the session host services in their established order.
func (s *DebugServices) ExtendContext(ctx context.Context) context.Context {
	ctx = s.logger.WithContext(ctx)
	ctx = encoding.WithRegistry(ctx, s.encoding)
	ctx = fs.WithFileSystem(ctx, s.filesystem)

	return ferretnet.WithNetwork(ctx, s.network)
}

// Materialize encodes output while leaving result cleanup to the debugger.
func (s *DebugServices) Materialize(result *vm.Result) (*encoding.Output, error) {
	return Materialize(s.encoding, s.outputContentType, result)
}

// Close runs close hooks, releases host resources, and finally returns the permit.
func (s *DebugServices) Close() error {
	var err error

	if s.hooks != nil {
		if hookErr := s.hooks.RunClose(); hookErr != nil {
			err = fmt.Errorf("close hooks: %w", hookErr)
		}
	}

	if closeErr := s.resources.Close(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}

	s.filesystem = nil

	if s.limiter != nil {
		s.limiter.Release()
		s.limiter = nil
	}

	return err
}
