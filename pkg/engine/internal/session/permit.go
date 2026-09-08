package session

import (
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// PermitRelease returns a session permit and, when configured, its borrowed VM.
type PermitRelease func(*vm.VM)

// NewPermitRelease transfers permit ownership to an idempotent callback.
// Normal sessions also use it to return their borrowed VM to the plan pool.
func NewPermitRelease(limiter *Limiter, pool *vm.Pool) PermitRelease {
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
