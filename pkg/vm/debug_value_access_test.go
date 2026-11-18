package vm

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type debugMetadataValue struct {
	info runtime.DebugInfo
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
	return v.info
}

func TestDebugValueBufferedCollections(t *testing.T) {
	access := NewDebugValueAccess()
	set := data.NewDataSet(false)
	if err := set.Append(t.Context(), runtime.True); err != nil {
		t.Fatal(err)
	}

	if got := access.Describe(set); got != (DebugValueDescription{TypeName: "Array", Length: 1}) {
		t.Fatalf("dataset description = %+v", got)
	}

	if _, ok := access.Inspect(set, 8); ok {
		t.Fatal("dataset unexpectedly exposed children")
	}

	for _, threshold := range []int{0, 1} {
		object := data.NewFastObject(nil, threshold)
		object.SetStringCached("b", runtime.False)
		object.SetStringCached("a", runtime.True)

		if got := access.Describe(object); got != (DebugValueDescription{TypeName: "Object", Length: 2, Children: DebugValueObject}) {
			t.Fatalf("fast object description = %+v", got)
		}

		children, ok := access.Inspect(object, 2)
		if !ok || !children.Complete || len(children.Items) != 2 {
			t.Fatalf("fast object children = %+v", children)
		}
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

func TestDebugValueAccessRejectsNilCollections(t *testing.T) {
	access := NewDebugValueAccess()
	for _, tc := range []struct{ value, key runtime.Value }{
		{(*runtime.Array)(nil), runtime.Int(0)},
		{(*runtime.Object)(nil), runtime.String("name")},
		{(*data.FastObject)(nil), runtime.String("name")},
	} {
		if _, err := access.Lookup(tc.value, tc.key); !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Errorf("lookup on %T: got %v, want ErrInvalidOperation", tc.value, err)
		}

		if _, ok := access.Inspect(tc.value, 8); ok {
			t.Errorf("inspection available for nil %T", tc.value)
		}
	}
}
