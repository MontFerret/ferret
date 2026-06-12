package ferret

import (
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// newSessionPermitRelease transfers permit ownership to an idempotent callback.
// Normal sessions also use it to return their borrowed VM to the plan pool.
func newSessionPermitRelease(limiter *sessionLimiter, pool *vm.Pool) sessionPermitRelease {
	var once sync.Once

	return func(instance *vm.VM) {
		once.Do(func() {
			limiter.Release()
			if pool != nil {
				// Return borrowed VMs even if the plan has already been closed.
				pool.Release(instance)
			}
		})
	}
}
