package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/api/debugger"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type (
	// DebugValueDescription contains shallow, in-memory presentation and child metadata.
	DebugValueDescription struct {
		Scalar   runtime.Value
		Hint     runtime.DebugInfo
		TypeName string
		Length   int // presentation length, or -1 when unavailable
		Children runtime.DebugChildrenInfo
	}

	// DebugChild identifies an immediate child by its position in its own view.
	DebugChild struct {
		Value    runtime.Value
		Name     string
		Kind     debugger.VariableFilter
		Position int64
	}

	// DebugValueAccess separates shallow descriptions, explicit evaluation
	// lookup and bounded reads of immediate children.
	DebugValueAccess interface {
		Describe(runtime.Value) DebugValueDescription
		Lookup(runtime.Value, runtime.Value) (runtime.Value, error)
		ReadChildren(context.Context, runtime.Value, debugger.VariablePage) ([]DebugChild, error)
	}

	debugValueAccess struct{}
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

// ReadChildren validates the range before touching the backing value.
func (debugValueAccess) ReadChildren(ctx context.Context, value runtime.Value, page debugger.VariablePage) ([]DebugChild, error) {
	if err := validateOperationContext(ctx); err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	page, err := page.Normalize()
	if err != nil {
		return nil, err
	}

	var children []DebugChild
	switch v := value.(type) {
	case *runtime.Array:
		if v == nil {
			return nil, errors.ErrUnsupported
		}

		length, _ := v.Length(ctx)
		children, err = readBufferedPage(ctx, page, int64(length), v.LookupAt)
	case *data.DataSet:
		length := v.BufferedLength()
		if length < 0 {
			return nil, errors.ErrUnsupported
		}

		children, err = readBufferedPage(ctx, page, int64(length), v.BufferedAt)
	default:
		if object := debugObject(value); object != nil {
			children, err = readObjectPage(ctx, object, page)
		} else if provider, ok := value.(runtime.DebugChildrenProvider); ok {
			children, err = readProviderPage(ctx, provider, page)
		} else {
			return nil, errors.ErrUnsupported
		}
	}

	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return children, nil
}
