package universal

import (
	"context"
	"sync"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type session struct {
	native    *engine.Session
	closeErr  error
	closeOnce sync.Once
}

var _ api.Session = (*session)(nil)

func (s *session) Run(ctx context.Context) (api.Output, error) {
	output, err := s.native.Run(ctx)

	return convertOutput(output), wrapDiagnosticError(err)
}

func (s *session) Close() error {
	s.closeOnce.Do(func() {
		s.closeErr = wrapDiagnosticError(s.native.Close())
	})

	return s.closeErr
}
