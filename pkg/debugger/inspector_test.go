package debugger

import (
	"testing"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestDebugValueUsesRuntimeDebugInfo(t *testing.T) {
	inspection := &inspector{values: vm.NewDebugValueAccess(), format: DefaultFormatOptions()}
	fallbackType := "debugger.hostileDebugValue"
	tests := []struct {
		name  string
		value runtime.Value
		want  Value
	}{
		{
			name:  "custom type",
			value: hostileDebugValue{info: runtime.DebugInfo{TypeName: "HTML::Node"}},
			want:  Value{Type: "HTML::Node", Display: "HostValue(HTML::Node)"},
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
			want:  Value{Type: fallbackType, Display: "HostValue(" + fallbackType + ")"},
		},
		{
			name:  "panicking metadata",
			value: hostileDebugValue{panicInfo: true},
			want:  Value{Type: fallbackType, Display: "HostValue(" + fallbackType + ")"},
		},
		{
			name:  "ordinary built-in",
			value: runtime.NewInt(1),
			want:  Value{Type: runtime.TypeInt.Name(), Display: "1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := inspection.debugValue(test.value); got != test.want {
				t.Fatalf("unexpected debugger value: got %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestDebugValueReadsCustomTypeInfoOnce(t *testing.T) {
	calls := 0
	inspection := &inspector{values: vm.NewDebugValueAccess(), format: DefaultFormatOptions()}
	value := hostileDebugValue{
		info:      runtime.DebugInfo{TypeName: "HTML::Node"},
		infoCalls: &calls,
	}

	if got := inspection.debugValue(value); got != (Value{Type: "HTML::Node", Display: "HostValue(HTML::Node)"}) {
		t.Fatalf("unexpected debugger value: %#v", got)
	}
	if calls != 1 {
		t.Fatalf("unexpected DebugInfo calls: got %d, want 1", calls)
	}
}

func TestDebugValueBoundsRuntimeDebugInfo(t *testing.T) {
	inspection := &inspector{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxDepth: 1, MaxItems: 1, MaxBytes: 4},
	}
	got := inspection.debugValue(hostileDebugValue{info: runtime.DebugInfo{
		TypeName: "HTML::Node",
		Display:  "product card",
	}})
	if got != (Value{Type: "HTML...", Display: "prod..."}) {
		t.Fatalf("unexpected bounded debugger value: %#v", got)
	}
}

func TestDebugValueBoundsRuntimeDebugInfoAtUTF8Boundaries(t *testing.T) {
	inspection := &inspector{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxDepth: 1, MaxItems: 1, MaxBytes: 4},
	}
	got := inspection.debugValue(hostileDebugValue{info: runtime.DebugInfo{
		TypeName: "类类",
		Display:  "🙂🙂",
	}})

	if got != (Value{Type: "类...", Display: "🙂..."}) {
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
		format: FormatOptions{MaxDepth: 1, MaxItems: 1, MaxBytes: 8},
	}
	got := inspection.debugValue(hostileDebugValue{info: runtime.DebugInfo{
		TypeName: malformed,
		Display:  malformed,
	}})

	if got != (Value{Type: "a\uFFFDb", Display: "a\uFFFDb"}) {
		t.Fatalf("unexpected normalized debugger value: %#v", got)
	}
	if !utf8.ValidString(got.Type) || !utf8.ValidString(got.Display) {
		t.Fatalf("debugger value is not valid UTF-8: %#v", got)
	}
}
