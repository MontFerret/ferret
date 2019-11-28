package events

import (
	"context"
	"hash/fnv"
)

func New(name string) ID {
	h := fnv.New32a()

	h.Write([]byte(name))

	return ID(h.Sum32())
}

func isCtxDone(ctx context.Context) bool {
	return ctx.Err() == context.Canceled
}
