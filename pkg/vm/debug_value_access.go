package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// DebugValueKind identifies a built-in collection shape available for safe
	// debugger inspection.
	DebugValueKind uint8

	// DebugValueDescription holds shallow presentation metadata. Length is -1
	// when unavailable. Children identifies an explicitly inspectable shape.
	DebugValueDescription struct {
		// Scalar is the value behind a non-nil pointer to a native scalar.
		// Binary provides only its byte length.
		Scalar   runtime.Value
		Hint     runtime.DebugInfo
		TypeName string
		Length   int
		Children DebugValueKind
	}

	// DebugValueItem is one safely inspected collection item. Key is populated
	// for object values.
	DebugValueItem struct {
		Value runtime.Value
		Key   string
	}

	// DebugValueInspection is a bounded snapshot of a built-in or VM-owned
	// collection.
	DebugValueInspection struct {
		Items    []DebugValueItem
		Length   int
		Kind     DebugValueKind
		Complete bool
	}

	// DebugValueAccess exposes only side-effect-free VM value operations needed
	// by source-level debugger policy.
	DebugValueAccess interface {
		Describe(runtime.Value) DebugValueDescription
		Lookup(runtime.Value, runtime.Value) (runtime.Value, error)
		Inspect(runtime.Value, int) (DebugValueInspection, bool)
	}

	debugValueAccess struct{}
)

const (
	DebugValueArray DebugValueKind = iota + 1
	DebugValueObject
)

// NewDebugValueAccess creates a side-effect-free accessor for built-in and
// VM-owned values.
func NewDebugValueAccess() DebugValueAccess {
	return debugValueAccess{}
}

func (debugValueAccess) Lookup(value, key runtime.Value) (runtime.Value, error) {
	ctx := context.Background()

	if array, ok := value.(*runtime.Array); ok && array != nil {
		index, ok := key.(runtime.Int)
		if !ok {
			return nil, runtime.Error(runtime.ErrInvalidArgument, "debugger array lookup requires an Int index")
		}

		return array.At(ctx, index)
	}

	if object := debugObject(value); object != nil {
		property, ok := key.(runtime.String)
		if !ok {
			return nil, runtime.Error(runtime.ErrInvalidArgument, "debugger object lookup requires a String key")
		}

		return object.Get(ctx, property)
	}

	return nil, runtime.Error(runtime.ErrInvalidOperation, "value does not support debugger member lookup")
}

func (debugValueAccess) Inspect(value runtime.Value, maxItems int) (DebugValueInspection, bool) {
	if array, ok := value.(*runtime.Array); ok && array != nil {
		length, _ := array.Length(context.Background())
		out := DebugValueInspection{Kind: DebugValueArray, Length: int(length)}

		if maxItems <= 0 || out.Length > maxItems {
			return out, true
		}

		out.Items = make([]DebugValueItem, 0, out.Length)

		for i := 0; i < out.Length; i++ {
			item, _ := array.At(context.Background(), runtime.Int(i))
			out.Items = append(out.Items, DebugValueItem{Value: item})
		}

		out.Complete = true

		return out, true
	}

	if object := debugObject(value); object != nil {
		return inspectDebugMap(object, maxItems), true
	}

	return DebugValueInspection{}, false
}
