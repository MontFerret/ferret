package debugger

import (
	"testing"
	"time"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestFormatDurationValue(t *testing.T) {
	t.Parallel()

	value := runtime.NewDuration(1500 * time.Millisecond)
	access := vm.NewDebugValueAccess()
	if got := formatValue(value, access, DefaultFormatOptions()); got != "1.5s" {
		t.Fatalf("formatted duration = %q", got)
	}
	if got := access.TypeName(value); got != runtime.TypeDuration.Name() {
		t.Fatalf("duration type name = %q", got)
	}
}

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
		want  string
		max   int
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

func TestFormatValueUsesNestedRuntimeDebugDisplay(t *testing.T) {
	value := runtime.NewArrayWith(hostileDebugValue{
		info: runtime.DebugInfo{Display: "<div.product-card>"},
	})

	if got := formatValue(value, vm.NewDebugValueAccess(), DefaultFormatOptions()); got != "[<div.product-card>]" {
		t.Fatalf("unexpected nested custom display: %q", got)
	}
}

func TestValueAccessRejectsOpaqueKeys(t *testing.T) {
	value := runtime.NewObject()
	if _, err := vm.NewDebugValueAccess().Lookup(value, hostileDebugValue{}); err == nil {
		t.Fatal("expected opaque object key to be rejected")
	}
}
