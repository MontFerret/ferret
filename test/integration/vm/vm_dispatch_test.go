package vm_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type testDispatcher struct {
	result runtime.Value
	err    error
	events []runtime.DispatchEvent
}

func (d *testDispatcher) Dispatch(_ context.Context, event runtime.DispatchEvent) (runtime.Value, error) {
	d.events = append(d.events, event)

	if d.err != nil {
		return runtime.None, d.err
	}

	if d.result != nil {
		return d.result, nil
	}

	return runtime.None, nil
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
	dispatcher := &testDispatcher{
		result: runtime.NewString("ok"),
	}

	RunUseCases(t, []UseCase{
		Case(`
			DISPATCH "click" IN @d
			RETURN 1
		`, 1, "Should dispatch as a statement"),
		Case(`RETURN DISPATCH "click" IN @d`, "ok", "Should dispatch with default payload and options"),
		Case(`LET event = "hover" RETURN DISPATCH event IN @d`, "ok", "Should dispatch with variable event name"),
		Case(`RETURN DISPATCH @event_name IN @d`, "ok", "Should dispatch with param event name"),
		Case(`RETURN DISPATCH "input" IN @d WITH "hello"`, "ok", "Should dispatch with payload"),
		Case(`RETURN DISPATCH "select" IN @d WITH ["1", "2"] OPTIONS { selector: "#a", delay: 50 }`, "ok", "Should dispatch with options"),
	}, vm.WithParams(map[string]runtime.Value{
		"d":          dispatcher,
		"event_name": runtime.NewString("submit"),
	}))

	var hasDefault bool
	var hasVariableEvent bool
	var hasParamEvent bool
	var hasPayload bool
	var hasOptions bool

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
}

func TestDispatchRuntimeErrors(t *testing.T) {
	dispatcher := &testDispatcher{result: runtime.NewString("ok")}

	RunUseCases(t, []UseCase{
		RuntimeErrorCase(`RETURN DISPATCH "click" IN @value`, ExpectedRuntimeError{Message: "Invalid type"}, "Should fail when target is not a dispatcher"),
		RuntimeErrorCase(`RETURN DISPATCH @event IN @d`, ExpectedRuntimeError{Message: "Invalid type"}, "Should fail when event name is not a string"),
	}, vm.WithParams(map[string]runtime.Value{
		"d":     dispatcher,
		"event": runtime.NewInt(1),
		"value": runtime.NewInt(1),
	}))
}
