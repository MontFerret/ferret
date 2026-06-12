package vm

import (
	"context"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type (
	// DebugValueKind identifies a built-in collection shape available for safe
	// debugger inspection.
	DebugValueKind uint8
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
		TypeName(runtime.Value) string
		DebugInfo(runtime.Value) (runtime.DebugInfo, bool)
		Lookup(runtime.Value, runtime.Value) (runtime.Value, error)
		Inspect(runtime.Value, int) (DebugValueInspection, bool)
	}
)

const (
	DebugValueArray DebugValueKind = iota + 1
	DebugValueObject
)

type debugValueAccess struct{}

// NewDebugValueAccess creates a side-effect-free accessor for built-in and
// VM-owned values.
func NewDebugValueAccess() DebugValueAccess {
	return debugValueAccess{}
}

func (a debugValueAccess) TypeName(value runtime.Value) string {
	if info, ok := a.DebugInfo(value); ok && info.TypeName != "" {
		return info.TypeName
	}

	return a.defaultTypeName(value)
}

// DebugInfo returns optional presentation metadata supplied by a runtime
// value.
func (a debugValueAccess) DebugInfo(value runtime.Value) (runtime.DebugInfo, bool) {
	inspectable, ok := value.(runtime.DebugInspectable)
	if !ok {
		return runtime.DebugInfo{}, false
	}

	return a.safeDebugInfo(inspectable)
}

// safeDebugInfo isolates panic recovery to implementations of the optional
// runtime hook.
func (debugValueAccess) safeDebugInfo(inspectable runtime.DebugInspectable) (info runtime.DebugInfo, ok bool) {
	defer func() {
		if recover() != nil {
			info = runtime.DebugInfo{}
			ok = false
		}
	}()

	return inspectable.DebugInfo(), true
}

func (debugValueAccess) defaultTypeName(value runtime.Value) string {
	if value == nil || reflect.TypeOf(value) == reflect.TypeOf(runtime.None) {
		return runtime.TypeNone.Name()
	}

	switch value.(type) {
	case runtime.Boolean:
		return runtime.TypeBoolean.Name()
	case runtime.Int:
		return runtime.TypeInt.Name()
	case runtime.Float:
		return runtime.TypeFloat.Name()
	case runtime.String:
		return runtime.TypeString.Name()
	case runtime.DateTime:
		return runtime.TypeDateTime.Name()
	case runtime.Binary:
		return runtime.TypeBinary.Name()
	case *runtime.Array:
		return runtime.TypeArray.Name()
	case *runtime.Object, *data.FastObject:
		return runtime.TypeObject.Name()
	default:
		return reflect.TypeOf(value).String()
	}
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
		return nil, runtime.Errorf(runtime.ErrInvalidOperation, "debugger cannot inspect %s", a.TypeName(value))
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
	return runtime.Errorf(runtime.ErrInvalidArgument, "debugger lookup requires %s, got %s", expected, a.TypeName(value))
}
