package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type LazyValue interface {
	core.Value

	Ready(ctx context.Context) error
}

func Ready(ctx context.Context, value core.Value) error {
	if lazy, ok := value.(LazyValue); ok {
		return lazy.Ready(ctx)
	}

	return nil
}
