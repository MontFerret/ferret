package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type debugMetadataValue struct {
	info      runtime.DebugInfo
	panicInfo bool
}

var _ runtime.DebugInspectable = debugMetadataValue{}

func (v debugMetadataValue) String() string {
	return "debug metadata value"
}

func (v debugMetadataValue) Hash() uint64 {
	return 0
}

func (v debugMetadataValue) Copy() runtime.Value {
	return v
}

func (v debugMetadataValue) DebugInfo() runtime.DebugInfo {
	if v.panicInfo {
		panic("debug info failed")
	}

	return v.info
}

func TestDebugValueAccessUsesRuntimeDebugInfo(t *testing.T) {
	access := NewDebugValueAccess()
	fallbackType := "vm.debugMetadataValue"
	tests := []struct {
		name     string
		value    runtime.Value
		info     runtime.DebugInfo
		typeName string
		ok       bool
	}{
		{
			name:     "custom type",
			value:    debugMetadataValue{info: runtime.DebugInfo{TypeName: "SQL::Connection"}},
			info:     runtime.DebugInfo{TypeName: "SQL::Connection"},
			typeName: "SQL::Connection",
			ok:       true,
		},
		{
			name:     "custom display",
			value:    debugMetadataValue{info: runtime.DebugInfo{Display: "open connection"}},
			info:     runtime.DebugInfo{Display: "open connection"},
			typeName: fallbackType,
			ok:       true,
		},
		{
			name: "custom type and display",
			value: debugMetadataValue{info: runtime.DebugInfo{
				TypeName: "SQL::Connection",
				Display:  "open connection",
			}},
			info: runtime.DebugInfo{
				TypeName: "SQL::Connection",
				Display:  "open connection",
			},
			typeName: "SQL::Connection",
			ok:       true,
		},
		{
			name:     "empty metadata",
			value:    debugMetadataValue{},
			typeName: fallbackType,
			ok:       true,
		},
		{
			name:     "panicking metadata",
			value:    debugMetadataValue{panicInfo: true},
			typeName: fallbackType,
		},
		{
			name:     "ordinary built-in",
			value:    runtime.NewString("value"),
			typeName: runtime.TypeString.Name(),
		},
		{
			name:     "ordinary opaque",
			value:    distinctCollisionValue{label: "opaque"},
			typeName: "vm.distinctCollisionValue",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info, ok := access.DebugInfo(test.value)
			if ok != test.ok {
				t.Fatalf("unexpected metadata availability: got %t, want %t", ok, test.ok)
			}
			if info != test.info {
				t.Fatalf("unexpected metadata: got %#v, want %#v", info, test.info)
			}
			if got := access.TypeName(test.value); got != test.typeName {
				t.Fatalf("unexpected type name: got %q, want %q", got, test.typeName)
			}
		})
	}
}

func TestDebugValueAccessInspectRemainsBuiltInOnly(t *testing.T) {
	access := NewDebugValueAccess()
	array := runtime.NewArrayWith(runtime.NewInt(1))

	inspection, ok := access.Inspect(array, 1)
	if !ok || !inspection.Complete || inspection.Kind != DebugValueArray || inspection.Length != 1 {
		t.Fatalf("unexpected built-in inspection: %#v, available=%t", inspection, ok)
	}
	if inspection, ok := access.Inspect(debugMetadataValue{
		info: runtime.DebugInfo{Display: "custom"},
	}, 1); ok || inspection.Items != nil || inspection.Length != 0 || inspection.Kind != 0 || inspection.Complete {
		t.Fatalf("custom metadata unexpectedly enabled expansion: %#v, available=%t", inspection, ok)
	}
}
