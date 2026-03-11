package vm

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func testRunStatePoolFixture() (*bytecode.Program, []int) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
	}

	return program, buildCatchByPC(len(program.Bytecode), program.CatchTable)
}

func TestRunStatePoolGet_InitializesStateWhenPoolIsEmpty(t *testing.T) {
	program, catchByPC := testRunStatePoolFixture()

	var pool statePool
	pool.Init(program, catchByPC, 0)

	state := pool.Get()
	if state == nil {
		t.Fatal("expected state from empty pool acquisition")
	}

	if state.program != program {
		t.Fatal("expected state program to be initialized from pool context")
	}

	if got, want := len(state.registers.Values), program.Registers; got != want {
		t.Fatalf("unexpected register file size: got %d, want %d", got, want)
	}

	if got, want := len(state.catchByPC), len(catchByPC); got != want {
		t.Fatalf("unexpected catch index size: got %d, want %d", got, want)
	}
}

func TestRunStatePoolPutGet_ReusesStateAfterCleanup(t *testing.T) {
	program, catchByPC := testRunStatePoolFixture()

	var pool statePool
	pool.Init(program, catchByPC, 0)

	state := pool.Get()
	state.prepareRun(NewDefaultEnvironment())
	state.raiseRuntime(errors.New("boom"), recoverDefault, bytecode.NoopOperand, nil, false)
	if !state.hasFailure() {
		t.Fatal("expected staged failure before pool cleanup")
	}

	pool.Put(state)

	reused := pool.Get()
	if reused != state {
		t.Fatal("expected the same state pointer to be reused from pool")
	}

	if reused.env != nil {
		t.Fatal("expected pool cleanup to clear environment")
	}

	if reused.hasFailure() {
		t.Fatal("expected pool cleanup to clear pending failure")
	}

	if got, want := reused.pc, 0; got != want {
		t.Fatalf("unexpected cleaned program counter: got %d, want %d", got, want)
	}

	if got, want := reused.lastPC, -1; got != want {
		t.Fatalf("unexpected cleaned last program counter: got %d, want %d", got, want)
	}
}

func TestRunStatePoolInitPrewarmAndLIFOReuse(t *testing.T) {
	program, catchByPC := testRunStatePoolFixture()

	var pool statePool
	pool.Init(program, catchByPC, 1)

	first := pool.Get()
	second := pool.Get()
	if first == second {
		t.Fatal("expected second acquisition to allocate a distinct state while first is active")
	}

	pool.Put(first)
	pool.Put(second)

	reused1 := pool.Get()
	reused2 := pool.Get()

	if reused1 != second {
		t.Fatal("expected LIFO order for first reused state")
	}

	if reused2 != first {
		t.Fatal("expected LIFO order for second reused state")
	}
}

func TestRunStatePoolPutNilIsNoOp(t *testing.T) {
	program, catchByPC := testRunStatePoolFixture()

	var pool statePool
	pool.Init(program, catchByPC, 0)

	pool.Put(nil)
	if got, want := len(pool.states), 0; got != want {
		t.Fatalf("unexpected pool size after nil put: got %d, want %d", got, want)
	}

	state := pool.Get()
	if state == nil {
		t.Fatal("expected state after nil put no-op")
	}
}
