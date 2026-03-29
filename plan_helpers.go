package ferret

import (
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func newSessionRelease(limiter *sessionLimiter, pool *vm.Pool) vmReleaseFunc {
	var once sync.Once

	return func(instance *vm.VM) {
		once.Do(func() {
			// Release the engine-wide session slot even if the plan has already been closed.
			limiter.Release()
			pool.Release(instance)
		})
	}
}
