package ferret

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Plan struct {
	prog         *bytecode.Program
	host         *host
	hooks        planHooks
	sessionHooks sessionHooks
	limiter      *sessionLimiter
	pool         *vm.Pool
	mu           sync.RWMutex
	closed       bool
}

func (p *Plan) NewSession(ctx context.Context, setters ...SessionOption) (*Session, error) {
	if p == nil {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	}

	h := p.host
	hooks := p.sessionHooks
	limiter := p.limiter
	pool := p.pool
	p.mu.RUnlock()

	sessionOpts, err := newSessionOptions(setters)
	if err != nil {
		return nil, err
	}

	if err := limiter.Acquire(ctx); err != nil {
		return nil, err
	}

	releaseLimiter := true
	defer func() {
		// Session construction can still fail after the limiter is acquired.
		// Roll back the permit unless ownership is handed to the Session.
		if releaseLimiter {
			limiter.Release()
		}
	}()

	env, err := vm.ExtendEnvironment(&vm.Environment{
		Functions: h.functions,
		Params:    h.params,
		Logging:   h.logging,
	}, sessionOpts.envOptions)

	if err != nil {
		return nil, err
	}

	instance, err := pool.Acquire()
	if err != nil {
		if errors.Is(err, vm.ErrPoolClosed) {
			return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
		}

		return nil, err
	}

	releaseLimiter = false

	return &Session{
		vm:                instance,
		env:               env,
		encoding:          h.encoding,
		outputContentType: sessionOpts.outputContentType,
		hooks:             hooks,
		release:           newSessionRelease(limiter, pool),
	}, nil
}

func (p *Plan) Close() error {
	if p == nil {
		return nil
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}

	p.closed = true
	// Snapshot close dependencies before unlocking so later Session.Close calls can
	// finish independently of the Plan's mutable state.
	hooks := p.hooks
	pool := p.pool
	p.mu.Unlock()

	var err error

	// Plan close hooks follow the hook registry close semantics (LIFO with error aggregation).
	if hookErr := hooks.runCloseHooks(); hookErr != nil {
		err = fmt.Errorf("close hooks: %w", hookErr)
	}

	if poolErr := pool.Close(); poolErr != nil {
		err = errors.Join(err, fmt.Errorf("close pool: %w", poolErr))
	}

	return err
}

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
