package debugger

import (
	"context"
	"errors"
	"math"
	"testing"

	apidebugger "github.com/MontFerret/api/debugger"
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type inspectionProvider struct {
	read func(context.Context, VariableFilter, int64) (string, runtime.Value, bool, error)
	hostileDebugValue
	info runtime.DebugChildrenInfo
}

func (p *inspectionProvider) DebugChildrenInfo() runtime.DebugChildrenInfo { return p.info }
func (p *inspectionProvider) DebugChildAt(ctx context.Context, kind VariableFilter, pos int64) (string, runtime.Value, bool, error) {
	return p.read(ctx, kind, pos)
}

func pausedInspection(t *testing.T, value runtime.Value) (*Session, ValueReference) {
	t.Helper()
	exec := &fakeExecution{startEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry}, resumeEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopStep}, locals: []vm.DebugLocal{{Name: "root", Value: value}}, params: runtime.Params{"input": value}}
	s, err := NewSession(Config{Values: vm.NewDebugValueAccess(), Execution: exec, Services: &fakeSessionServices{}, Source: source.NewAnonymous("RETURN @input"), Params: []string{"input"}})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Error(err)
		}
	})
	if _, err := s.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	locals, err := s.Locals(t.Context())
	if err != nil || len(locals) != 2 || !locals[0].Value.Reference.Valid() {
		t.Fatalf("locals: %+v, %v", locals, err)
	}

	return s, locals[0].Value.Reference
}

func TestInspectionOccurrenceLifetime(t *testing.T) {
	nested := runtime.NewArrayWith(runtime.Int(42))
	children := make([]runtime.Value, 30)

	for i := range children {
		children[i] = nested
	}

	s, root := pausedInspection(t, runtime.NewArrayOf(children))
	roots, err := s.FrameLocals(t.Context(), 0)
	if err != nil || roots[0].Value.Reference != root || roots[1].Value.Reference == root {
		t.Fatalf("root bindings: %+v, %v", roots, err)
	}

	for range 3 {
		locals, err := s.Locals(t.Context())
		if err != nil || locals[1].Value.Reference != roots[1].Value.Reference || len(s.inspector.valueRefs) != 2 {
			t.Fatalf("refresh: %v, %v", locals, err)
		}
	}

	first, err := s.Variables(t.Context(), root, VariablePage{Count: 20})
	if err != nil || len(first) != 20 {
		t.Fatalf("%v, %v", first, err)
	}

	for range 3 {
		next, err := s.Variables(t.Context(), root, VariablePage{Filter: VariableFilterIndexed, Start: 10, Count: 20})
		if err != nil || len(next) != 20 {
			t.Fatalf("%v, %v", next, err)
		}

		for i := range 10 {
			if !next[i].Value.Reference.Valid() || next[i].Value.Reference != first[i+10].Value.Reference {
				t.Fatalf("overlap %d", i)
			}
		}

		if len(s.inspector.valueRefs) != 32 {
			t.Fatalf("retained %d handles", len(s.inspector.valueRefs))
		}
	}

	leaf, err := s.Variables(t.Context(), first[10].Value.Reference, VariablePage{})
	if err != nil || len(leaf) != 1 || leaf[0].Value.Display != "42" {
		t.Fatalf("nested: %v %v", leaf, err)
	}

	if _, err := s.StepIn(t.Context()); err != nil {
		t.Fatal(err)
	}

	for _, ref := range []ValueReference{0, -1, root, first[10].Value.Reference, 1000} {
		if got, err := s.Variables(t.Context(), ref, VariablePage{}); got != nil || !errors.Is(err, apidebugger.ErrInvalidValueReference) {
			t.Fatalf("expired %d: %v %v", ref, got, err)
		}
	}

	locals, err := s.Locals(t.Context())
	if err != nil || locals[0].Value.Reference <= 32 {
		t.Fatalf("reused ID: %v %v", locals, err)
	}

	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	if got, err := s.Variables(t.Context(), root, VariablePage{}); got != nil || !errors.Is(err, &StateError{}) {
		t.Fatalf("closed: %v %v", got, err)
	}

	if first[10].Name != "10" || first[10].Value.Display != "Array(1)" || *first[10].Value.IndexedCount != 1 || *first[10].Value.NamedCount != 0 {
		t.Fatalf("retained snapshot: %+v", first[10])
	}
}

func TestInspectionReorderedLocals(t *testing.T) {
	s, root := pausedInspection(t, runtime.NewArrayWith(runtime.Int(11)))
	exec := s.lifecycle.execution.(*fakeExecution)
	exec.locals = append(exec.locals, vm.DebugLocal{Name: "input", Value: runtime.NewArrayWith(runtime.Int(22))})
	first, err := s.Locals(t.Context())
	if err != nil || len(first) != 3 {
		t.Fatalf("locals: %v, %v", first, err)
	}

	exec.locals[0], exec.locals[1] = exec.locals[1], exec.locals[0]
	again, err := s.FrameLocals(t.Context(), 0)
	if err != nil || len(again) != 3 {
		t.Fatalf("frame locals: %v, %v", again, err)
	}

	for i, want := range []struct {
		name    string
		display string
		ref     ValueReference
	}{
		{name: "input", display: "22", ref: first[1].Value.Reference},
		{name: "root", display: "11", ref: root},
		{name: "@input", display: "11", ref: first[2].Value.Reference},
	} {
		v := again[i]
		if v.Name != want.name || v.Value.Reference != want.ref {
			t.Fatalf("binding %d: %+v, want %+v", i, v, want)
		}

		page, err := s.Variables(t.Context(), v.Value.Reference, VariablePage{})
		if err != nil || len(page) != 1 || page[0].Value.Display != want.display {
			t.Fatalf("binding %s: %v, %v", v.Name, page, err)
		}
	}

	if len(s.inspector.valueRefs) != 3 {
		t.Fatalf("expected three distinct bindings, got %d", len(s.inspector.valueRefs))
	}
}

func TestAmbiguousLocalNamesFailAtomically(t *testing.T) {
	s, root := pausedInspection(t, runtime.NewArrayWith(runtime.True))
	exec := s.lifecycle.execution.(*fakeExecution)
	child := runtime.NewArrayWith(runtime.False)
	exec.locals = []vm.DebugLocal{{Name: "x", Value: child}, {Name: "x", Value: child}}
	last := s.inspector.lastValueRef

	got, err := s.Locals(t.Context())
	if got != nil || !errors.Is(err, runtime.ErrUnexpected) || s.inspector.lastValueRef != last {
		t.Fatalf("ambiguous bindings: %v, %v, last reference %d", got, err, s.inspector.lastValueRef)
	}

	page, err := s.Variables(t.Context(), root, VariablePage{})
	if err != nil || len(page) != 1 || page[0].Value.Display != "true" {
		t.Fatalf("previous binding: %+v, %v", page, err)
	}
}

func TestMixedChildOccurrences(t *testing.T) {
	known := runtime.DebugChildView{Known: true, Count: 1, Readable: true}
	p := &inspectionProvider{info: runtime.DebugChildrenInfo{Named: known, Indexed: known}}
	p.read = func(context.Context, VariableFilter, int64) (string, runtime.Value, bool, error) {
		return "0", runtime.NewArrayWith(runtime.True), true, nil
	}

	s, root := pausedInspection(t, p)
	all, err := s.Variables(t.Context(), root, VariablePage{})
	if err != nil || len(all) != 2 || all[0].Value.Reference == all[1].Value.Reference {
		t.Fatalf("mixed page: %+v, %v", all, err)
	}

	for i, kind := range []VariableFilter{VariableFilterNamed, VariableFilterIndexed} {
		result, err := s.Variables(t.Context(), root, VariablePage{Filter: kind})
		if err != nil || len(result) != 1 || result[0].Name != "0" || result[0].Value.Reference != all[i].Value.Reference {
			t.Fatalf("%s: %+v, %v", kind, result, err)
		}
	}
}

func TestInspectionAtomicFailure(t *testing.T) {
	for _, failure := range []string{"read", "canceled", "limit", "overflow"} {
		t.Run(failure, func(t *testing.T) {
			inner := runtime.NewArrayWith(runtime.True)
			backing := []runtime.Value{inner, inner, inner}
			p := &inspectionProvider{info: runtime.DebugChildrenInfo{Named: runtime.DebugChildView{Known: true}, Indexed: runtime.DebugChildView{Known: true, Readable: true, Count: 3}}}
			injected := errors.New("child read failed")
			failRead := false
			p.read = func(_ context.Context, _ VariableFilter, i int64) (string, runtime.Value, bool, error) {
				if failRead && i == 2 {
					return "", nil, false, injected
				}

				return "", backing[i], true, nil
			}

			s, root := pausedInspection(t, p)
			old, err := s.Variables(t.Context(), root, VariablePage{Count: 1})
			if err != nil || len(old) != 1 {
				t.Fatalf("initial page: %v %v", old, err)
			}

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			want := ErrValueReferencesExhausted
			switch failure {
			case "read":
				failRead = true
				want = injected
			case "canceled":
				// The first new occurrence has been staged when the second description runs.
				descs := 0
				innerAccess := s.inspector.values
				s.inspector.values = &fakeValueAccess{inner: innerAccess, describe: func(v runtime.Value) vm.DebugValueDescription {
					descs++
					if descs == 2 {
						cancel()
					}

					return innerAccess.Describe(v)
				}}
				want = context.Canceled
			case "limit":
				s.inspector.handleLimit = len(s.inspector.valueRefs) + 1
			case "overflow":
				s.inspector.lastValueRef = math.MaxInt - 1
			}

			last, live := s.inspector.lastValueRef, len(s.inspector.valueRefs)
			got, err := s.Variables(ctx, root, VariablePage{Start: 1, Count: 2})
			if got != nil || !errors.Is(err, want) || s.inspector.lastValueRef != last || len(s.inspector.valueRefs) != live || len(s.inspector.occurrences) != live {
				t.Fatalf("partial commit: page=%+v err=%v last=%d live=%d", got, err, s.inspector.lastValueRef, len(s.inspector.valueRefs))
			}

			got, err = s.Variables(t.Context(), old[0].Value.Reference, VariablePage{})
			if err != nil || len(got) != 1 || got[0].Value.Display != "true" {
				t.Fatalf("prior handle: %v %v", got, err)
			}

			failRead = false
			got, err = s.Variables(t.Context(), root, VariablePage{Start: 1, Count: 1})
			if err != nil || len(got) != 1 || got[0].Value.Reference != last+1 {
				t.Fatalf("retry: %+v, %v", got, err)
			}

			again, err := s.Variables(t.Context(), root, VariablePage{Start: 1, Count: 1})
			if err != nil || len(again) != 1 || again[0].Value.Reference != got[0].Value.Reference {
				t.Fatalf("reuse at limit: %v %v", again, err)
			}
		})
	}
}

func TestPageOrderedWithExecution(t *testing.T) {
	for _, command := range []string{"close", "resume"} {
		t.Run(command, func(t *testing.T) {
			entered, release := make(chan struct{}), make(chan struct{})
			p := &inspectionProvider{info: runtime.DebugChildrenInfo{Indexed: runtime.DebugChildView{Readable: true, Known: true, Count: 1}, Named: runtime.DebugChildView{Known: true}}}
			p.read = func(ctx context.Context, _ VariableFilter, _ int64) (string, runtime.Value, bool, error) {
				close(entered)
				select {
				case <-ctx.Done():
					return "", nil, false, ctx.Err()
				case <-release:
					return "", runtime.NewArrayWith(runtime.True), true, nil
				}
			}

			s, ref := pausedInspection(t, p)
			type response struct {
				err  error
				page []Variable
			}

			done := make(chan response, 1)
			go func() {
				page, err := s.Variables(t.Context(), ref, VariablePage{})
				done <- response{page: page, err: err}
			}()
			<-entered
			if command == "close" {
				last := s.inspector.lastValueRef
				if err := s.Close(); err != nil {
					t.Fatal(err)
				}

				out := <-done
				if out.page != nil || out.err == nil || s.inspector.lastValueRef != last || len(s.inspector.valueRefs) != 0 {
					t.Fatalf("close: %+v", out)
				}

				return
			}

			resumed := make(chan error, 1)
			checked := make(chan struct{}, 1)
			observed := &contextObserver{Context: t.Context(), check: func() {
				select {
				case checked <- struct{}{}:
				default:
				}
			}}
			go func() { _, err := s.Continue(observed); resumed <- err }()
			<-checked
			select {
			case err := <-resumed:
				t.Fatalf("resume overtook inspection: %v", err)
			default:
			}

			close(release)
			out := <-done
			if out.err != nil || len(out.page) != 1 || !out.page[0].Value.Reference.Valid() {
				t.Fatalf("page: %+v", out)
			}

			if err := <-resumed; err != nil {
				t.Fatal(err)
			}

			if _, err := s.Variables(t.Context(), out.page[0].Value.Reference, VariablePage{}); !errors.Is(err, apidebugger.ErrInvalidValueReference) {
				t.Fatalf("old page: %v", err)
			}
		})
	}
}

func TestCloseWinsStagedPagePublication(t *testing.T) {
	child := runtime.NewArrayWith(runtime.True)
	children := make([]runtime.Value, 1000)
	for i := range children {
		children[i] = child
	}
	session, root := pausedInspection(t, runtime.NewArrayOf(children))
	last := session.inspector.lastValueRef
	entered, release := make(chan struct{}), make(chan struct{})
	var pageCtx context.Context
	inner := session.inspector.values
	descriptions := 0
	session.inspector.values = &fakeValueAccess{inner: inner, read: func(ctx context.Context, value runtime.Value, page VariablePage) ([]vm.DebugChild, error) {
		pageCtx = ctx
		return inner.ReadChildren(ctx, value, page)
	}, describe: func(raw runtime.Value) vm.DebugValueDescription {
		descriptions++
		if descriptions == 2 {
			close(entered)
			<-release
		}

		return inner.Describe(raw)
	}}
	done := make(chan error, 1)
	go func() {
		page, err := session.Variables(t.Context(), root, VariablePage{Count: int64(len(children))})
		if page != nil {
			t.Errorf("partial page: %+v", page)
		}

		done <- err
	}()
	<-entered
	closed := make(chan error, 1)
	go func() { closed <- session.Close() }()
	<-pageCtx.Done()
	close(release)
	if err := <-done; err == nil {
		t.Fatal("staged inspection committed after close")
	}

	if err := <-closed; err != nil {
		t.Fatal(err)
	}

	if descriptions != 2 {
		t.Fatalf("Close allowed %d child descriptions, want 2", descriptions)
	}

	if session.inspector.lastValueRef != last || len(session.inspector.valueRefs) != 0 {
		t.Fatal("staged references survived close")
	}
}

func TestRejectedResumePreservesStopReferences(t *testing.T) {
	session, root := pausedInspection(t, runtime.NewArrayWith(runtime.True))
	exec := session.lifecycle.execution.(*fakeExecution)
	exec.resumeErr = context.Canceled
	exec.errorStatus = vm.DebugExecutionPaused
	if _, err := session.Continue(t.Context()); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}

	page, err := session.Variables(t.Context(), root, VariablePage{})
	if err != nil || len(page) != 1 || page[0].Value.Display != "true" {
		t.Fatalf("unchanged stop: %v, %v", page, err)
	}
}

func TestAdmittedResumeReleasesStoppedValues(t *testing.T) {
	exec := newBlockingExecution(bytecode.DebugPoint{})
	exec.params = runtime.Params{"input": runtime.NewArrayWith(runtime.True)}
	s, err := NewSession(Config{Execution: exec, Values: vm.NewDebugValueAccess(), Services: &fakeSessionServices{}, Source: source.NewAnonymous("RETURN 1")})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if _, err := s.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	// Evaluate supplies a root independently of the execution's own locals.
	value, err := s.Evaluate(t.Context(), "@input")
	if err != nil || !value.Reference.Valid() {
		t.Fatalf("evaluation: %+v, %v", value, err)
	}

	done := make(chan error, 1)
	go func() { _, err := s.Continue(t.Context()); done <- err }()
	<-exec.resumeStarted
	// The backend is parked until release. No other command can publish.
	retained := len(s.inspector.valueRefs)
	exec.release()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if retained != 0 {
		t.Fatalf("running execution retained %d stopped values", retained)
	}
}
