package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

// Only concrete, buffered implementations may supply a length or children.
// Host collection capabilities provide names without granting access.
func (debugValueAccess) classify(value runtime.Value) DebugValueDescription {
	d := DebugValueDescription{Length: -1}
	if value == nil || value == runtime.None {
		d.TypeName = runtime.TypeNone.Name()

		return d
	}

	switch v := value.(type) {
	case runtime.Boolean:
		d.TypeName = runtime.TypeBoolean.Name()
	case runtime.Int:
		d.TypeName = runtime.TypeInt.Name()
	case runtime.Float:
		d.TypeName = runtime.TypeFloat.Name()
	case runtime.Duration:
		d.TypeName = runtime.TypeDuration.Name()
	case runtime.String:
		d.TypeName = runtime.TypeString.Name()
	case runtime.DateTime:
		d.TypeName = runtime.TypeDateTime.Name()
	case *runtime.Regexp:
		d.TypeName = runtime.TypeRegexp.Name()
	case runtime.Binary:
		d.TypeName = runtime.TypeBinary.Name()
		d.Length = len(v)
	case *runtime.Array:
		d.TypeName = runtime.TypeArray.Name()
		if v != nil {
			length, _ := v.Length(context.Background())
			d.Length, d.Children = int(length), DebugValueArray
		}
	case *data.DataSet:
		d.TypeName = runtime.TypeArray.Name()
		d.Length = v.BufferedLength()
	case *runtime.Object:
		d.TypeName = runtime.TypeObject.Name()
		if v != nil {
			length, _ := v.Length(context.Background())
			d.Length, d.Children = int(length), DebugValueObject
		}
	case *data.FastObject:
		d.TypeName = runtime.TypeObject.Name()
		if v != nil {
			length, _ := v.Length(context.Background())
			d.Length, d.Children = int(length), DebugValueObject
		}
	case runtime.ObjectLike:
		d.TypeName = runtime.TypeObject.Name()
	case runtime.List:
		d.TypeName = runtime.TypeList.Name()
	case runtime.Map:
		d.TypeName = runtime.TypeMap.Name()
	case runtime.Queryable:
		d.TypeName = runtime.TypeQueryable.Name()
	case runtime.Iterator:
		d.TypeName = runtime.TypeIterator.Name()
	case runtime.Iterable:
		d.TypeName = runtime.TypeIterable.Name()
	default:
		d.TypeName = "Resource"
	}

	return d
}

// Recover only around the optional embedder hook.
func (debugValueAccess) safeDebugInfo(value runtime.DebugInspectable) (info runtime.DebugInfo) {
	defer func() {
		if recover() != nil {
			info = runtime.DebugInfo{}
		}
	}()

	return value.DebugInfo()
}
