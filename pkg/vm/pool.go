package vm

import (
	"errors"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// Pool represents a pool of reusable virtual machines for executing bytecode programs with configurable options.
// It allows for efficient management of resources by reusing VM instances, reducing the overhead of creating new VMs for each execution.
type Pool struct {
	program  *bytecode.Program
	opts     []Option
	idle     []*VM
	inUse    map[*VM]struct{}
	total    int
	max      int
	maxTotal int
	mu       sync.Mutex
	closed   bool
}

func NewPool(program *bytecode.Program, max int, opts ...Option) *Pool {
	return NewPoolWithLimits(program, max, 0, opts...)
}

func NewPoolWithLimits(program *bytecode.Program, maxIdle, maxTotal int, opts ...Option) *Pool {
	if maxIdle < 0 {
		maxIdle = 0
	}

	if maxTotal < 0 {
		maxTotal = 0
	}

	return &Pool{
		program:  program,
		opts:     opts,
		inUse:    make(map[*VM]struct{}),
		max:      maxIdle,
		maxTotal: maxTotal,
	}
}

func (p *Pool) Acquire() (*VM, error) {
	p.mu.Lock()
	if p.inUse == nil {
		p.inUse = make(map[*VM]struct{})
	}

	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}

	if len(p.idle) > 0 {
		vm := p.idle[len(p.idle)-1]
		p.idle = p.idle[:len(p.idle)-1]
		p.inUse[vm] = struct{}{}
		p.mu.Unlock()

		return vm, nil
	}

	if p.maxTotal > 0 && p.total >= p.maxTotal {
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	}

	program := p.program
	opts := p.opts
	p.total++
	p.mu.Unlock()

	instance, err := NewWith(program, opts...)
	if err != nil {
		p.mu.Lock()
		if p.total > 0 {
			p.total--
		}
		p.mu.Unlock()

		return nil, err
	}

	p.mu.Lock()
	if p.closed {
		if p.total > 0 {
			p.total--
		}
		p.mu.Unlock()

		return nil, errors.Join(ErrPoolClosed, instance.Close())
	}

	p.inUse[instance] = struct{}{}
	p.mu.Unlock()

	return instance, nil
}

func (p *Pool) Release(vm *VM) {
	if vm == nil {
		return
	}

	shouldClose := false

	p.mu.Lock()
	if _, ok := p.inUse[vm]; !ok {
		p.mu.Unlock()
		return
	}

	delete(p.inUse, vm)

	switch {
	case p.closed || vm.closed:
		shouldClose = !vm.closed

		if p.total > 0 {
			p.total--
		}
	case p.max == 0 || len(p.idle) >= p.max:
		// Idle retention is capped separately from total live capacity.
		shouldClose = true

		if p.total > 0 {
			p.total--
		}
	default:
		p.idle = append(p.idle, vm)
	}

	p.mu.Unlock()

	if shouldClose {
		_ = vm.Close()
	}
}

func (p *Pool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}

	p.closed = true
	idle := p.idle
	p.idle = nil
	p.program = nil
	p.opts = nil

	if len(idle) > 0 {
		if p.total >= len(idle) {
			p.total -= len(idle)
		} else {
			p.total = 0
		}
	}

	p.mu.Unlock()

	var err error
	for _, vm := range idle {
		err = errors.Join(err, vm.Close())
	}

	return err
}
