package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type (
	// DebugValueKind identifies a built-in collection shape available for safe
	// debugger inspection.
	DebugValueKind uint8

	// DebugValueDescription holds shallow presentation metadata. Length is -1
	// when unavailable. Children identifies an explicitly inspectable shape.
	DebugValueDescription struct {
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

func (a debugValueAccess) Describe(value runtime.Value) DebugValueDescription {
	description := a.classify(value)
	if inspectable, ok := value.(runtime.DebugInspectable); ok {
		description.Hint = a.safeDebugInfo(inspectable)
	}

	return description
}

func (a debugValueAccess) Lookup(value, key runtime.Value) (runtime.Value, error) {
	ctx := context.Background()

	switch value := value.(type) {
	case *runtime.Array:
		index, ok := key.(runtime.Int)
		if !ok {
			return nil, a.lookupTypeError(key, runtime.TypeInt.Name())
		}

		return value.At(ctx, index)
	case *runtime.Object:
		property, ok := key.(runtime.String)
		if !ok {
			return nil, a.lookupTypeError(key, runtime.TypeString.Name())
		}

		return value.Get(ctx, property)
	case *data.FastObject:
		property, ok := key.(runtime.String)
		if !ok {
			return nil, a.lookupTypeError(key, runtime.TypeString.Name())
		}

		return value.Get(ctx, property)
	default:
		return nil, runtime.Errorf(runtime.ErrInvalidOperation, "debugger cannot inspect %s", a.classify(value).TypeName)
	}
}

func (debugValueAccess) Inspect(value runtime.Value, maxItems int) (DebugValueInspection, bool) {
	switch value := value.(type) {
	case *runtime.Array:
		length, _ := value.Length(context.Background())
		out := DebugValueInspection{Kind: DebugValueArray, Length: int(length)}

		if maxItems <= 0 || out.Length > maxItems {
			return out, true
		}

		out.Items = make([]DebugValueItem, 0, out.Length)

		for i := 0; i < out.Length; i++ {
			item, _ := value.At(context.Background(), runtime.Int(i))
			out.Items = append(out.Items, DebugValueItem{Value: item})
		}

		out.Complete = true

		return out, true
	case *runtime.Object:
		return inspectDebugMap(value, maxItems), true
	case *data.FastObject:
		return inspectDebugMap(value, maxItems), true
	default:
		return DebugValueInspection{}, false
	}
}

func (a debugValueAccess) lookupTypeError(value runtime.Value, expected string) error {
	return runtime.Errorf(runtime.ErrInvalidArgument, "debugger lookup requires %s, got %s", expected, a.classify(value).TypeName)
}
