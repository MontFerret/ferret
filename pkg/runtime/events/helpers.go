package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Consumer func(ctx context.Context, data core.Value) error

func Consume(ctx context.Context, ch <-chan Event, consumer Consumer) error {
	for evt := range ch {
		if evt.Err != nil {
			return nil
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if err := consumer(ctx, evt.Data); err != nil {
			if core.IsDone(err) {
				return nil
			}

			return err
		}
	}

	return nil
}
