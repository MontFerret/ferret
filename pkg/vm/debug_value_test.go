package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type hostileDebugValue struct{}

func (hostileDebugValue) String() string      { panic("String called") }
func (hostileDebugValue) Hash() uint64        { panic("Hash called") }
func (hostileDebugValue) Copy() runtime.Value { panic("Copy called") }
func (hostileDebugValue) Type() runtime.Type  { panic("Type called") }

func TestDebugFormatValueDoesNotInvokeOpaqueHostValue(t *testing.T) {
	value := hostileDebugValue{}
	if got := DebugFormatValue(value, DefaultDebugFormatOptions()); got != "HostValue(vm.hostileDebugValue)" {
		t.Fatalf("unexpected host summary: %q", got)
	}
	if got := DebugValueTypeName(value); got != "vm.hostileDebugValue" {
		t.Fatalf("unexpected host type name: %q", got)
	}
}

func TestDebugFormatValueBoundsStrings(t *testing.T) {
	const maxBytes = 8
	got := DebugFormatValue(runtime.NewString("abcdefghijklmnopqrstuvwxyz"), DebugFormatOptions{
		MaxDepth: 1,
		MaxItems: 1,
		MaxBytes: maxBytes,
	})
	if len(got) > maxBytes+3 {
		t.Fatalf("formatted value exceeds limit: %q", got)
	}
}

func TestDebugLookupValueRejectsOpaqueKeys(t *testing.T) {
	value := runtime.NewObject()
	if _, err := DebugLookupValue(value, hostileDebugValue{}); err == nil {
		t.Fatal("expected opaque object key to be rejected")
	}
}
