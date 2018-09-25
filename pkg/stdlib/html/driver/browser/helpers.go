package browser

import (
	"golang.org/x/sync/errgroup"
)

func PointerInt(input int) *int {
	return &input
}

type BatchFunc = func() error

func RunBatch(funcs ...BatchFunc) error {
	eg := errgroup.Group{}

	for _, f := range funcs {
		eg.Go(f)
	}

	return eg.Wait()
}
