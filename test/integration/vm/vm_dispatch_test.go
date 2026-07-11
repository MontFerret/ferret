package vm_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"github.com/MontFerret/ferret/v2/test/spec/mock"
)

type testDispatcher struct {
	err    error
	events []runtime.DispatchEvent
}

func (d *testDispatcher) Dispatch(_ context.Context, event runtime.DispatchEvent) error {
	d.events = append(d.events, event)

	if d.err != nil {
		return d.err
	}

	return nil
}

func (d *testDispatcher) MarshalJSON() ([]byte, error) {
	return json.Marshal("dispatcher")
}

func (d *testDispatcher) String() string {
	return "dispatcher"
}

func (d *testDispatcher) Unwrap() interface{} {
	return "dispatcher"
}

func (d *testDispatcher) Hash() uint64 {
	return 0
}

func (d *testDispatcher) Copy() runtime.Value {
	return d
}

func TestDispatch(t *testing.T) {
	dispatcher := &testDispatcher{}

	RunSpecs(t, []spec.Spec{
		S(`
			DISPATCH "click" IN @d
			RETURN 1
		`, 1, "Should dispatch as a statement"),
		S(`RETURN DISPATCH "click" IN @d`, nil, "Should dispatch with default payload and options"),
		S(`LET event = "hover" RETURN DISPATCH event IN @d`, nil, "Should dispatch with variable event name"),
		S(`RETURN DISPATCH @event_name IN @d`, nil, "Should dispatch with param event name"),
		S(`RETURN DISPATCH "input" IN @d WITH "hello"`, nil, "Should dispatch with payload"),
		S(`RETURN DISPATCH "select" IN @d WITH ["1", "2"] OPTIONS { selector: "#a", delay: 50 }`, nil, "Should dispatch with options"),
		S(`RETURN @d <- "focus"`, nil, "Should dispatch shorthand as an expression"),
		S(`
			LET result = @d <- "commit"
			RETURN result
		`, nil, "Should assign NONE from shorthand dispatch"),
		S(`
			LET tag = MATCH @kind (
				"click" => @d <- "press",
				_ => @d <- "hover",
			)
			RETURN tag
		`, nil, "Should allow shorthand dispatch in MATCH arms"),
		S(`
			FUNC fire() => @d <- "blur"
			RETURN fire()
		`, nil, "Should allow shorthand dispatch in UDF arrow bodies"),
	}, vm.WithParams(map[string]runtime.Value{
		"d":          dispatcher,
		"event_name": runtime.NewString("submit"),
		"kind":       runtime.NewString("click"),
	}))

	var hasDefault bool
	var hasVariableEvent bool
	var hasParamEvent bool
	var hasPayload bool
	var hasOptions bool
	var hasShorthandFocus bool
	var hasShorthandCommit bool
	var hasShorthandPress bool
	var hasShorthandBlur bool

	for _, evt := range dispatcher.events {
		switch evt.Name {
		case runtime.NewString("click"):
			if evt.Payload == runtime.None && evt.Options == runtime.None {
				hasDefault = true
			}
		case runtime.NewString("hover"):
			if evt.Payload == runtime.None && evt.Options == runtime.None {
				hasVariableEvent = true
			}
		case runtime.NewString("submit"):
			if evt.Payload == runtime.None && evt.Options == runtime.None {
				hasParamEvent = true
			}
		case runtime.NewString("input"):
			if evt.Payload == runtime.NewString("hello") && evt.Options == runtime.None {
				hasPayload = true
			}
		case runtime.NewString("select"):
			opts, ok := evt.Options.(runtime.Map)
			if !ok {
				continue
			}

			selector, err := opts.Get(context.Background(), runtime.NewString("selector"))
			if err == nil && selector == runtime.NewString("#a") {
				hasOptions = true
			}
		case runtime.NewString("focus"):
			if evt.Payload == runtime.None && evt.Options == runtime.None {
				hasShorthandFocus = true
			}
		case runtime.NewString("commit"):
			if evt.Payload == runtime.None && evt.Options == runtime.None {
				hasShorthandCommit = true
			}
		case runtime.NewString("press"):
			if evt.Payload == runtime.None && evt.Options == runtime.None {
				hasShorthandPress = true
			}
		case runtime.NewString("blur"):
			if evt.Payload == runtime.None && evt.Options == runtime.None {
				hasShorthandBlur = true
			}
		}
	}

	if !hasDefault {
		t.Fatalf("expected default dispatch event with none payload/options")
	}

	if !hasPayload {
		t.Fatalf("expected payload dispatch event")
	}

	if !hasVariableEvent {
		t.Fatalf("expected variable-name dispatch event")
	}

	if !hasParamEvent {
		t.Fatalf("expected param-name dispatch event")
	}

	if !hasOptions {
		t.Fatalf("expected options dispatch event")
	}

	if !hasShorthandFocus {
		t.Fatalf("expected shorthand focus dispatch event")
	}

	if !hasShorthandCommit {
		t.Fatalf("expected shorthand commit dispatch event")
	}

	if !hasShorthandPress {
		t.Fatalf("expected shorthand press dispatch event")
	}

	if !hasShorthandBlur {
		t.Fatalf("expected shorthand blur dispatch event")
	}
}

func TestDispatchGroupedQueryTarget(t *testing.T) {
	dispatcher := &testDispatcher{}
	queryable := mock.NewQueryable(runtime.NewArrayWith(dispatcher))

	RunSpecs(t, []spec.Spec{
		S(`
			RETURN DISPATCH "input" IN (QUERY ONE "#query" IN @page USING css)
				WITH { value: "ferret" }
				OPTIONS { bubbles: true }
		`, nil, "Should dispatch to a value returned by a grouped query"),
	}, vm.WithParam("page", queryable))

	if got := len(queryable.MockQueries()); got != 2 {
		t.Fatalf("expected one query evaluation per optimization level, got %d", got)
	}

	if got := len(dispatcher.events); got != 2 {
		t.Fatalf("expected one dispatched event per optimization level, got %d", got)
	}

	for _, event := range dispatcher.events {
		if event.Name != runtime.NewString("input") {
			t.Fatalf("expected input event, got %q", event.Name)
		}
		if got := queryMapValue(t, event.Payload, "value"); got != runtime.NewString("ferret") {
			t.Fatalf("expected payload value ferret, got %v", got)
		}
		if got := queryMapValue(t, event.Options, "bubbles"); got != runtime.True {
			t.Fatalf("expected bubbles option true, got %v", got)
		}
	}
}

func TestDispatchInForBodies(t *testing.T) {
	dispatcher := &testDispatcher{}

	RunSpecs(t, []spec.Spec{
		Array(`
			VAR i = 0
			FOR WHILE i < 1
				i += 1
				DISPATCH "loop-click" IN @d
				RETURN i
		`, []any{1}, "Should dispatch long form in FOR WHILE body"),
		Array(`
			FOR item IN [1]
				@d <- "loop-shorthand"
				RETURN item
		`, []any{1}, "Should dispatch shorthand in FOR IN body"),
	}, vm.WithParams(map[string]runtime.Value{
		"d": dispatcher,
	}))

	var hasLoopClick bool
	var hasLoopShorthand bool

	for _, evt := range dispatcher.events {
		switch evt.Name {
		case runtime.NewString("loop-click"):
			hasLoopClick = true
		case runtime.NewString("loop-shorthand"):
			hasLoopShorthand = true
		}
	}

	if !hasLoopClick {
		t.Fatalf("expected long-form dispatch event from FOR body")
	}

	if !hasLoopShorthand {
		t.Fatalf("expected shorthand dispatch event from FOR body")
	}
}

func TestDispatchRuntimeErrors(t *testing.T) {
	dispatcher := &testDispatcher{}

	RunSpecs(t, []spec.Spec{
		S(`RETURN DISPATCH "click" IN @value`, "Should fail when target is not a dispatcher").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "invalid type"},
		),
		S("DISPATCH \"click\" IN @value ON ERROR RETURN NONE\nRETURN 1", 1, "Statement suppression should continue after dispatch failure"),
		S(`RETURN DISPATCH @event IN @d`, "Should fail when event name is not a string").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "invalid type"},
		),
		Nil(`RETURN DISPATCH @event IN @d ON ERROR RETURN NONE`, "Expression suppression should return none on dispatch failure"),
		S(`RETURN @value <- "click"`, "Shorthand should fail when target is not a dispatcher").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "invalid type"},
		),
	}, vm.WithParams(map[string]runtime.Value{
		"d":     dispatcher,
		"event": runtime.NewInt(1),
		"value": runtime.NewInt(1),
	}))
}

func TestDispatchShorthandContexts(t *testing.T) {
	dispatcher := &testDispatcher{}

	RunSpecs(t, []spec.Spec{
		// Dispatch produces NONE; the container captures that NONE element/value.
		S(`RETURN [@d <- "click"]`, []interface{}{nil}, "Should allow shorthand dispatch in array literal"),
		S(`RETURN { x: @d <- "click" }`, map[string]interface{}{"x": nil}, "Should allow shorthand dispatch in object literal"),
	}, vm.WithParams(map[string]runtime.Value{
		"d": dispatcher,
	}))
}
