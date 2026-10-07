package uapi

import (
	"context"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type session struct {
	native *engine.Session
}

var _ api.Session = (*session)(nil)

func (s *session) Run(ctx context.Context) (api.Output, error) {
	output, err := s.native.Run(ctx)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	return &outputAdapter{native: output}, nil
}

func (s *session) Close() error {
	return wrapDiagnosticError(s.native.Close())
}
