package debugger

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type hostileDebugValue struct{}

func (hostileDebugValue) String() string      { panic("String called") }
func (hostileDebugValue) Hash() uint64        { panic("Hash called") }
func (hostileDebugValue) Copy() runtime.Value { panic("Copy called") }
func (hostileDebugValue) Type() runtime.Type  { panic("Type called") }

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

func TestValueAccessRejectsOpaqueKeys(t *testing.T) {
	value := runtime.NewObject()
	if _, err := vm.NewDebugValueAccess().Lookup(value, hostileDebugValue{}); err == nil {
		t.Fatal("expected opaque object key to be rejected")
	}
}
