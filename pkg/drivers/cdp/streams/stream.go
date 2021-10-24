package streams

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/rpcc"
)

type (
	Decoder func(stream rpcc.Stream) (core.Value, error)

	Factory func(ctx context.Context) (rpcc.Stream, error)
)
