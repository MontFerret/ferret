package debugger

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	// Done observes admission to the writer's select without timing assumptions.
	breakpointWaitContext struct {
		context.Context
		waiting chan struct{}
		once    sync.Once
	}

	// The host call gates each loop visit before the two candidate breakpoints.
	liveExecutionGate struct {
		entered chan struct{}
		release chan struct{}
	}

	// Hold a real VM result before native conversion, after the stop is decided.
	deferredStopExecution struct {
		vm.DebugExecution
		ready   chan *vm.DebugExecutionEvent
		release chan struct{}
		reason  vm.DebugStopReason
	}
)

func (c *breakpointWaitContext) Done() <-chan struct{} {
	c.once.Do(func() { close(c.waiting) })

	return c.Context.Done()
}

func (d *deferredStopExecution) Resume(ctx context.Context, mode vm.DebugResumeMode, check vm.DebugBreakpointPredicate) (*vm.DebugExecutionEvent, error) {
	event, err := d.DebugExecution.Resume(ctx, mode, check)
	if err == nil && event.Reason == d.reason {
		select {
		case d.ready <- event:
		case <-ctx.Done():
			return event, nil
		}

		select {
		case <-d.release:
		case <-ctx.Done():
		}
	}

	return event, err
}

func newLiveSession(t *testing.T, hold *deferredStopExecution) (*Session, *liveExecutionGate) {
	t.Helper()

	c, err := compiler.New(compiler.WithDebugInfo())
	if err != nil {
		t.Fatal(err)
	}

	src := source.New("live.fql", "RETURN FOR i IN 1..3\n  LET tick = GATE()\n  LET x = i + 1\n  RETURN x")
	program, err := c.Compile(t.Context(), src)
	if err != nil {
		t.Fatal(err)
	}

	instance, err := vm.New(program)
	if err != nil {
		t.Fatal(err)
	}

	gate := &liveExecutionGate{entered: make(chan struct{}), release: make(chan struct{})}
	env, err := vm.NewEnvironment([]vm.EnvironmentOption{vm.WithFunction("GATE", gate.call)})
	if err != nil {
		t.Fatal(err)
	}

	execution, err := vm.NewDebugExecution(instance, env)
	if err != nil {
		t.Fatal(err)
	}

	if hold != nil {
		hold.DebugExecution = execution
		execution = hold
	}

	session, err := NewSession(Config{
		Execution: execution, Values: vm.NewDebugValueAccess(), Services: &fakeSessionServices{},
		Source: src, DebugPoints: program.Metadata.DebugPoints,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = session.Close() })
	if event, err := session.Start(t.Context()); err != nil || event.Reason != ReasonEntry {
		t.Fatalf("start: %+v, %v", event, err)
	}

	return session, gate
}

func (g *liveExecutionGate) call(ctx context.Context, _ ...runtime.Value) (runtime.Value, error) {
	select {
	case g.entered <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case <-g.release:
		return runtime.None, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (g *liveExecutionGate) visit(t *testing.T) {
	t.Helper()
	waitForSignal(t, g.entered, "debuggee host call")
}

func (g *liveExecutionGate) proceed(t *testing.T) {
	t.Helper()
	select {
	case g.release <- struct{}{}:
	case <-time.After(3 * time.Second):
		t.Fatal("debuggee did not leave host call")
	}
}
