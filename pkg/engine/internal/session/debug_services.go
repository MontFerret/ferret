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

// DebugServices implements debugger.SessionServices for native executions.
// The debugger serializes its use and owns once-only closure; retained VM state
// is owned by the debugger, independently of these host resources and permit.
type DebugServices struct {
	Logger            logging.Logger
	Hooks             *host.SessionHooks
	FileSystem        fs.FileSystem
	Network           ferretnet.Network
	Encoding          *encoding.Registry
	ReleasePermit     PermitRelease
	Resources         *resource.Manager
	OutputContentType string
}

// BeforeRun runs the configured before-run hook chain.
func (s *DebugServices) BeforeRun(ctx context.Context) (context.Context, error) {
	return s.Hooks.RunBeforeRun(ctx)
}

// AfterRun runs the configured after-run hooks with the primary execution error.
func (s *DebugServices) AfterRun(ctx context.Context, runErr error) error {
	return s.Hooks.RunAfterRun(ctx, runErr)
}

// ExtendContext injects the session host services in their established order.
func (s *DebugServices) ExtendContext(ctx context.Context) context.Context {
	ctx = s.Logger.WithContext(ctx)
	ctx = encoding.WithRegistry(ctx, s.Encoding)
	ctx = fs.WithFileSystem(ctx, s.FileSystem)

	return ferretnet.WithNetwork(ctx, s.Network)
}

// Materialize encodes output while leaving result cleanup to the debugger.
func (s *DebugServices) Materialize(result *vm.Result) (*encoding.Output, error) {
	return Materialize(s.Encoding, s.OutputContentType, result)
}

// Close runs close hooks, releases host resources, and finally returns the permit.
func (s *DebugServices) Close() error {
	var err error

	if s.Hooks != nil {
		if hookErr := s.Hooks.RunClose(); hookErr != nil {
			err = fmt.Errorf("close hooks: %w", hookErr)
		}
	}

	if closeErr := s.Resources.Close(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}

	s.FileSystem = nil

	if s.ReleasePermit != nil {
		s.ReleasePermit(nil)
		s.ReleasePermit = nil
	}

	return err
}
