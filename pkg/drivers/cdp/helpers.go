package cdp

import "golang.org/x/sync/errgroup"

type (
	batchFunc = func() error
)

func runBatch(funcs ...batchFunc) error {
	eg := errgroup.Group{}

	for _, f := range funcs {
		eg.Go(f)
	}

	return eg.Wait()
}
