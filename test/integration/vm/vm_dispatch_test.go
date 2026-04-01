package vm_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
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
		S(`RETURN "focus" -> @d`, nil, "Should dispatch shorthand as an expression"),
		S(`
			LET result = "commit" -> @d
			RETURN result
		`, nil, "Should assign NONE from shorthand dispatch"),
		S(`
			LET tag = MATCH @kind (
				"click" => "press" -> @d,
				_ => "hover" -> @d,
			)
			RETURN tag
		`, nil, "Should allow shorthand dispatch in MATCH arms"),
		S(`
			FUNC fire() => "blur" -> @d
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

func TestDispatchRuntimeErrors(t *testing.T) {
	dispatcher := &testDispatcher{}

	RunSpecs(t, []spec.Spec{
		S(`RETURN DISPATCH "click" IN @value`, "Should fail when target is not a dispatcher").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "Invalid type"},
		),
		S(`RETURN DISPATCH @event IN @d`, "Should fail when event name is not a string").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "Invalid type"},
		),
	}, vm.WithParams(map[string]runtime.Value{
		"d":     dispatcher,
		"event": runtime.NewInt(1),
		"value": runtime.NewInt(1),
	}))
}
