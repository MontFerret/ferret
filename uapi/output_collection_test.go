package uapi

import (
	"context"
	"errors"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	collectionSession interface {
		Run(context.Context) (encoding.Output, error)
	}

	collectionEngine interface {
		Run(context.Context, source.Source, ...engine.SessionOption) (encoding.Output, error)
	}

	collectionRuntime interface {
		Run(context.Context, api.Source, ...api.SessionOption) (api.Output, error)
	}
)

func collectSession(session collectionSession, ctx context.Context) (*encoding.Content, error) {
	output, err := session.Run(ctx)

	return collectHandle(ctx, output, err)
}

func collectEngine(engine collectionEngine, ctx context.Context, src source.Source, opts ...engine.SessionOption) (*encoding.Content, error) {
	output, err := engine.Run(ctx, src, opts...)

	return collectHandle(ctx, output, err)
}

func collectRuntime(runtime collectionRuntime, ctx context.Context, src api.Source, opts ...api.SessionOption) (*api.Content, error) {
	output, err := runtime.Run(ctx, src, opts...)

	return collectHandle(ctx, output, err)
}

// Existing representation tests collect explicitly; lifecycle tests use live handles.
func collectHandle(ctx context.Context, output encoding.Output, err error) (content *encoding.Content, resultErr error) {
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := output.Close(); closeErr != nil && !errors.Is(resultErr, closeErr) {
			resultErr = errors.Join(resultErr, closeErr)
		}
	}()

	return output.Collect(ctx)
}
