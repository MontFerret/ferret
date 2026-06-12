package debugger

import (
	"testing"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type hostileDebugValue struct {
	infoCalls *int
	info      runtime.DebugInfo
	panicInfo bool
}

func (hostileDebugValue) String() string      { panic("String called") }
func (hostileDebugValue) Hash() uint64        { panic("Hash called") }
func (hostileDebugValue) Copy() runtime.Value { panic("Copy called") }
func (hostileDebugValue) Type() runtime.Type  { panic("Type called") }
func (v hostileDebugValue) DebugInfo() runtime.DebugInfo {
	if v.infoCalls != nil {
		(*v.infoCalls)++
	}
	if v.panicInfo {
		panic("DebugInfo called")
	}

	return v.info
}

func TestFormatValueDoesNotInvokeOpaqueHostValue(t *testing.T) {
	value := hostileDebugValue{}
	access := vm.NewDebugValueAccess()
	if got := formatValue(value, access, DefaultFormatOptions()); got != "HostValue(debugger.hostileDebugValue)" {
		t.Fatalf("unexpected host summary: %q", got)
	}
	if got := access.TypeName(value); got != "debugger.hostileDebugValue" {
		t.Fatalf("unexpected host type name: %q", got)
	}
}

func TestFormatValueBoundsStrings(t *testing.T) {
	const maxBytes = 8
	got := formatValue(runtime.NewString("abcdefghijklmnopqrstuvwxyz"), vm.NewDebugValueAccess(), FormatOptions{
		MaxDepth: 1,
		MaxItems: 1,
		MaxBytes: maxBytes,
	})
	if len(got) > maxBytes+3 {
		t.Fatalf("formatted value exceeds limit: %q", got)
	}
}

func TestBoundedText(t *testing.T) {
	tests := []struct {
		name  string
		value string
		max   int
		want  string
	}{
		{name: "ascii", value: "product", max: 4, want: "prod..."},
		{name: "unicode", value: "éclair", max: 3, want: "éc..."},
		{name: "emoji", value: "🙂x", max: 3, want: "..."},
		{name: "zero limit", value: "value", max: 0, want: "..."},
		{name: "negative limit", value: "value", max: -1, want: "..."},
		{name: "empty with negative limit", value: "", max: -1, want: ""},
		{name: "within limit", value: "é", max: 2, want: "é"},
		{name: "malformed UTF-8", value: string([]byte{'a', 0xff, 'b'}), max: 8, want: "a\uFFFDb"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := boundedText(test.value, test.max)

			if got != test.want {
				t.Fatalf("unexpected bounded text: got %q, want %q", got, test.want)
			}
			if !utf8.ValidString(got) {
				t.Fatalf("bounded text is not valid UTF-8: %q", got)
			}
		})
	}
}

func TestFormatValueFinalTruncationIsValidUTF8(t *testing.T) {
	got := formatValue(
		runtime.NewArrayWith(hostileDebugValue{info: runtime.DebugInfo{Display: "🙂🙂"}}),
		vm.NewDebugValueAccess(),
		FormatOptions{MaxDepth: 1, MaxItems: 1, MaxBytes: 5},
	)

	if !utf8.ValidString(got) {
		t.Fatalf("formatted value is not valid UTF-8: %q", got)
	}
}

func TestDebugValueUsesRuntimeDebugInfo(t *testing.T) {
	session := &Session{values: vm.NewDebugValueAccess(), format: DefaultFormatOptions()}
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
			if got := session.debugValue(test.value); got != test.want {
				t.Fatalf("unexpected debugger value: got %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestFormatValueUsesNestedRuntimeDebugDisplay(t *testing.T) {
	value := runtime.NewArrayWith(hostileDebugValue{
		info: runtime.DebugInfo{Display: "<div.product-card>"},
	})

	if got := formatValue(value, vm.NewDebugValueAccess(), DefaultFormatOptions()); got != "[<div.product-card>]" {
		t.Fatalf("unexpected nested custom display: %q", got)
	}
}

func TestDebugValueReadsCustomTypeInfoOnce(t *testing.T) {
	calls := 0
	session := &Session{values: vm.NewDebugValueAccess(), format: DefaultFormatOptions()}
	value := hostileDebugValue{
		info:      runtime.DebugInfo{TypeName: "HTML::Node"},
		infoCalls: &calls,
	}

	if got := session.debugValue(value); got != (Value{Type: "HTML::Node", Display: "HostValue(HTML::Node)"}) {
		t.Fatalf("unexpected debugger value: %#v", got)
	}
	if calls != 1 {
		t.Fatalf("unexpected DebugInfo calls: got %d, want 1", calls)
	}
}

func TestDebugValueBoundsRuntimeDebugInfo(t *testing.T) {
	session := &Session{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxDepth: 1, MaxItems: 1, MaxBytes: 4},
	}
	got := session.debugValue(hostileDebugValue{info: runtime.DebugInfo{
		TypeName: "HTML::Node",
		Display:  "product card",
	}})
	if got != (Value{Type: "HTML...", Display: "prod..."}) {
		t.Fatalf("unexpected bounded debugger value: %#v", got)
	}
}

func TestDebugValueBoundsRuntimeDebugInfoAtUTF8Boundaries(t *testing.T) {
	session := &Session{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxDepth: 1, MaxItems: 1, MaxBytes: 4},
	}
	got := session.debugValue(hostileDebugValue{info: runtime.DebugInfo{
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
	session := &Session{
		values: vm.NewDebugValueAccess(),
		format: FormatOptions{MaxDepth: 1, MaxItems: 1, MaxBytes: 8},
	}
	got := session.debugValue(hostileDebugValue{info: runtime.DebugInfo{
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

func TestValueAccessRejectsOpaqueKeys(t *testing.T) {
	value := runtime.NewObject()
	if _, err := vm.NewDebugValueAccess().Lookup(value, hostileDebugValue{}); err == nil {
		t.Fatal("expected opaque object key to be rejected")
	}
}
