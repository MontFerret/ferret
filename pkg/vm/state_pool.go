package vm

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

// statePool manages reusable execution states for VM runs.
type statePool struct {
	program   *bytecode.Program
	catchByPC []int
	states    []*execState
}

// Init prepares the pool and optionally prewarms it with initialized states.
func (p *statePool) Init(program *bytecode.Program, catchByPC []int, prewarmCount int) {
	if prewarmCount < 0 {
		prewarmCount = 0
	}

	p.program = program
	p.catchByPC = catchByPC
	p.states = make([]*execState, 0, prewarmCount)

	for i := 0; i < prewarmCount; i++ {
		state := &execState{}
		state.init(program, catchByPC)
		p.states = append(p.states, state)
	}
}

// Get returns a pooled state or allocates and initializes a new one.
func (p *statePool) Get() *execState {
	n := len(p.states)
	if n > 0 {
		state := p.states[n-1]
		p.states = p.states[:n-1]
		return state
	}

	state := &execState{}
	state.init(p.program, p.catchByPC)
	return state
}

// Put cleans and stores a state for reuse.
func (p *statePool) Put(state *execState) {
	if state == nil {
		return
	}

	state.cleanupForPool()
	p.states = append(p.states, state)
}
