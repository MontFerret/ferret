package ferret

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type debugSessionServices struct {
	hooks             sessionHooks
	limiter           *sessionLimiter
	encoding          *encoding.Registry
	logger            logging.Logger
	fs                fs.FileSystem
	outputContentType string
}

func (s *debugSessionServices) BeforeRun(ctx context.Context) (context.Context, error) {
	return s.hooks.runBeforeRunHooks(ctx)
}

func (s *debugSessionServices) AfterRun(ctx context.Context, runErr error) error {
	return s.hooks.runAfterRunHooks(ctx, runErr)
}

func (s *debugSessionServices) ExtendContext(ctx context.Context) context.Context {
	ctx = s.logger.WithContext(ctx)
	ctx = encoding.WithRegistry(ctx, s.encoding)
	return fs.WithFileSystem(ctx, s.fs)
}

func (s *debugSessionServices) Materialize(result *vm.Result) (*encoding.Output, error) {
	return newOutput(s.encoding, s.outputContentType, result)
}

func (s *debugSessionServices) Close() error {
	var err error
	if s.hooks != nil {
		if hookErr := s.hooks.runCloseHooks(); hookErr != nil {
			err = fmt.Errorf("close hooks: %w", hookErr)
		}
	}
	if s.limiter != nil {
		s.limiter.Release()
		s.limiter = nil
	}
	return err
}
