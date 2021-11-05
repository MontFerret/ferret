package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"sync"
)

type mrger struct {
	inputs []Stream
}

func (m *mrger) Close(ctx context.Context) error {
	errs := make([]error, 0, len(m.inputs))

	for _, s := range m.inputs {
		if err := s.Close(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return core.Errors(errs...)
}

func (m *mrger) Read(ctx context.Context) <-chan Message {
	var wg sync.WaitGroup
	wg.Add(len(m.inputs))

	out := make(chan Message)
	consume := func(c context.Context, input Stream) {
		for evt := range input.Read(c) {
			out <- evt
		}

		wg.Done()
	}

	for _, ch := range m.inputs {
		go consume(ctx, ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Merge(inputs ...Stream) Stream {
	return &mrger{inputs}
}
