package debugger

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestDebugValueUsesRuntimeDebugInfo(t *testing.T) {
	inspection := &inspector{values: vm.NewDebugValueAccess(), format: DefaultFormatOptions()}
	fallbackType := "Resource"
	tests := []struct {
		name  string
		value runtime.Value
		want  Value
	}{
		{
			name:  "custom type",
			value: hostileDebugValue{info: runtime.DebugInfo{TypeName: "HTML::Node"}},
			want:  Value{Type: "HTML::Node", Display: "<HTML::Node>"},
		},
		{
			name:  "custom display",
			value: hostileDebugValue{info: runtime.DebugInfo{Display: "<div.product-card>"}},
			want:  Value{Type: fallbackType, Display: "<div.product-card>"},
		},
		{
			name: "custom type and display",
			value: hostileDebugValue{info: runtime.DebugInfo{
				TypeName: "SQL::Connection",
				Display:  "open connection",
			}},
			want: Value{Type: "SQL::Connection", Display: "open connection"},
		},
		{
			name:  "empty metadata",
			value: hostileDebugValue{},
			want:  Value{Type: fallbackType, Display: "<Resource>"},
		},
		{
			name:  "panicking metadata",
			value: hostileDebugValue{panicInfo: true},
			want:  Value{Type: fallbackType, Display: "<Resource>"},
		},
		{
			name:  "ordinary built-in",
			value: runtime.NewInt(1),
			want:  Value{Type: runtime.TypeInt.Name(), Display: "1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.value
			calls := 0
			if host, ok := value.(hostileDebugValue); ok {
				host.infoCalls = &calls
				value = host
			}

			if got := inspection.debugValue(value); got != test.want {
				t.Fatalf("unexpected debugger value: got %#v, want %#v", got, test.want)
			}

			if _, ok := value.(hostileDebugValue); ok && calls != 1 {
				t.Fatalf("DebugInfo called %d times, want once", calls)
			}
		})
	}
}

func TestDebugValueBoundsRuntimeDebugInfo(t *testing.T) {
	inspection := &inspector{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxItems: 1, MaxBytes: 4},
	}
	got := inspection.debugValue(hostileDebugValue{info: runtime.DebugInfo{
		TypeName: "HTML::Node",
		Display:  "product card",
	}})
	if got != (Value{Type: "H...", Display: "p..."}) {
		t.Fatalf("unexpected bounded debugger value: %#v", got)
	}
}

func TestDebugValueBoundsRuntimeDebugInfoAtUTF8Boundaries(t *testing.T) {
	inspection := &inspector{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxItems: 1, MaxBytes: 4},
	}
	got := inspection.debugValue(hostileDebugValue{info: runtime.DebugInfo{
		TypeName: "类类",
		Display:  "🙂🙂",
	}})

	if got != (Value{Type: "...", Display: "..."}) {
		t.Fatalf("unexpected bounded debugger value: %#v", got)
	}
	if !utf8.ValidString(got.Type) || !utf8.ValidString(got.Display) {
		t.Fatalf("debugger value is not valid UTF-8: %#v", got)
	}
}

func TestDebugValueNormalizesMalformedRuntimeDebugInfo(t *testing.T) {
	malformed := string([]byte{'a', 0xff, 'b'})
	inspection := &inspector{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxItems: 1, MaxBytes: 8},
	}
	got := inspection.debugValue(hostileDebugValue{info: runtime.DebugInfo{
		TypeName: malformed,
		Display:  malformed,
	}})

	if got != (Value{Type: "Resource", Display: "a\uFFFDb"}) {
		t.Fatalf("unexpected normalized debugger value: %#v", got)
	}
	if !utf8.ValidString(got.Type) || !utf8.ValidString(got.Display) {
		t.Fatalf("debugger value is not valid UTF-8: %#v", got)
	}
}

// Nil embedded capabilities panic if presentation attempts any access.
type (
	debugHostList   struct{ runtime.List }
	namedDebugList  struct{ runtime.List }
	debugHostMap    struct{ runtime.Map }
	debugHostObject struct{ runtime.Map }
	debugHostQuery  struct {
		runtime.Value
		runtime.Queryable
	}
	debugHostIterator struct {
		runtime.Value
		runtime.Iterator
		runtime.Iterable
	}
	debugHostIterable struct {
		runtime.Value
		runtime.Iterable
	}
)

func (namedDebugList) DebugInfo() runtime.DebugInfo {
	return runtime.DebugInfo{TypeName: "DB::Rows"}
}

func (debugHostObject) ObjectLike() { panic("ObjectLike marker invoked") }

func TestDebugValueCapabilitySummaries(t *testing.T) {
	access := &fakeValueAccess{inner: vm.NewDebugValueAccess()}
	s := inspector{values: access, format: DefaultFormatOptions()}
	for _, tc := range []struct {
		value runtime.Value
		name  string
	}{
		{struct{ runtime.Value }{}, "Resource"}, {namedDebugList{}, "DB::Rows"},
		{debugHostList{}, "List"}, {debugHostMap{}, "Map"},
		{debugHostObject{}, "Object"},
		{debugHostQuery{}, "Queryable"}, {debugHostIterator{}, "Iterator"},
		{debugHostIterable{}, "Iterable"},
	} {
		if got := s.debugValue(tc.value); got != (Value{Type: tc.name, Display: "<" + tc.name + ">"}) {
			t.Fatalf("%s description = %+v", tc.name, got)
		}
	}

	if access.inspectCalls != 0 || access.lookupCalls != 0 {
		t.Fatalf("read children while describing host capabilities: %+v", access)
	}
}

func TestDebugValueReadsChildHookOnlyOnExpansion(t *testing.T) {
	calls := 0
	value := runtime.NewArrayWith(hostileDebugValue{infoCalls: &calls, info: runtime.DebugInfo{Display: "child"}})
	s := inspector{
		values:       vm.NewDebugValueAccess(),
		format:       FormatOptions{MaxItems: 1, MaxBytes: 64},
		nextValueRef: 1,
		valueRefs:    make(map[ValueReference]runtime.Value),
	}
	parent := s.debugValue(value)
	if calls != 0 || parent.Display != "Array(1)" || !parent.Reference.Valid() {
		t.Fatalf("parent = %+v, child called %d times", parent, calls)
	}

	children, err := s.variables(parent.Reference)
	if err != nil || len(children) != 1 || calls != 1 || children[0].Value.Display != "child" {
		t.Fatalf("children = %+v, calls = %d, err = %v", children, calls, err)
	}
}

// Wrapped targets may implement their own hooks. The SDK wrapper must not
// forward presentation to those targets.
func TestDebugValueHostWrappers(t *testing.T) {
	calls := 0
	target := hostileDebugValue{infoCalls: &calls, info: runtime.DebugInfo{TypeName: "Private"}}
	s := inspector{values: vm.NewDebugValueAccess(), format: DefaultFormatOptions()}
	for _, tc := range []struct {
		value runtime.Value
		want  Value
	}{
		{sdk.NewHostValue(target), Value{Type: "Resource", Display: "<Resource>"}},
		{sdk.NewHostValueWithType(runtime.NewType("html", "Node", nil), target), Value{Type: "html.Node", Display: "<html.Node>"}},
	} {
		if got := s.debugValue(tc.value); got != tc.want || calls != 0 {
			t.Fatalf("wrapper = %+v, target calls = %d", got, calls)
		}
	}
}

func TestDebugValueLongCustomHints(t *testing.T) {
	text := strings.Repeat("x", 4<<20) + "\xff"
	s := inspector{values: vm.NewDebugValueAccess(), format: FormatOptions{MaxBytes: 8}}
	for _, input := range []string{text, "xxxxx\xffabcdef"} {
		value := hostileDebugValue{info: runtime.DebugInfo{TypeName: input, Display: input}}
		if got := s.debugValue(value); got != (Value{Type: "xxxxx...", Display: "xxxxx..."}) {
			t.Fatalf("truncated hints = %+v", got)
		}
	}
}

func TestEvaluateUsesSessionBudgetForOperandType(t *testing.T) {
	calls := 0
	host := hostileDebugValue{infoCalls: &calls, info: runtime.DebugInfo{TypeName: "HTML::Node"}}
	s := inspector{
		execution: &fakeExecution{locals: []vm.DebugLocal{{Name: "node", Value: host}}},
		values:    vm.NewDebugValueAccess(),
		format:    FormatOptions{MaxItems: 3, MaxBytes: 8},
	}
	shown, err := s.evaluateFrame(t.Context(), 0, "node")
	if err != nil || shown.Type != "HTML:..." || calls != 1 {
		t.Fatalf("node = %+v, err = %v, hook calls = %d", shown, err, calls)
	}

	if _, err := s.evaluateFrame(t.Context(), 0, "node AND true"); !errors.Is(err, runtime.ErrInvalidArgument) ||
		!strings.HasSuffix(err.Error(), "got "+shown.Type) || calls != 2 {
		t.Fatalf("logical expression error = %v, hook calls = %d", err, calls)
	}

	if _, err := s.evaluateFrame(t.Context(), 0, "node.tag"); !errors.Is(err, runtime.ErrInvalidOperation) ||
		strings.Contains(err.Error(), "Resource") || calls != 2 {
		t.Fatalf("member lookup error = %v, hook calls = %d", err, calls)
	}
}

func TestFrameLocalsSnapshotsDeclaredParamsOnce(t *testing.T) {
	for _, tc := range []struct {
		declared []string
		want     []string
		calls    int
	}{
		{nil, []string{"local"}, 0},
		{[]string{"z", "missing", "a"}, []string{"local", "@a", "@z"}, 1},
	} {
		execution := &fakeExecution{
			locals: []vm.DebugLocal{{Name: "local", Value: runtime.True}},
			params: runtime.Params{"a": runtime.True, "z": runtime.False, "undeclared": runtime.None},
		}
		s := newInspector(Config{Execution: execution, Params: tc.declared, Values: vm.NewDebugValueAccess()}, nil, nil, DefaultFormatOptions())
		locals, err := s.frameLocals(0)
		if err != nil || len(locals) != len(tc.want) || execution.paramsCalls != tc.calls {
			t.Fatalf("declared %v: locals = %+v, err = %v, Params calls = %d", tc.declared, locals, err, execution.paramsCalls)
		}

		for i, want := range tc.want {
			if locals[i].Name != want {
				t.Errorf("local %d = %q, want %q", i, locals[i].Name, want)
			}
		}
	}
}
