package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

// Describe reads only cheap native metadata or explicitly opted-in host hooks.
// General host collection capabilities provide names without granting access.
func (debugValueAccess) Describe(value runtime.Value) DebugValueDescription {
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
	case *runtime.Boolean:
		d = describeScalarPointer(v, runtime.TypeBoolean)
	case *runtime.Int:
		d = describeScalarPointer(v, runtime.TypeInt)
	case *runtime.Float:
		d = describeScalarPointer(v, runtime.TypeFloat)
	case *runtime.Duration:
		d = describeScalarPointer(v, runtime.TypeDuration)
	case *runtime.String:
		d = describeScalarPointer(v, runtime.TypeString)
	case *runtime.DateTime:
		d = describeScalarPointer(v, runtime.TypeDateTime)
	case *runtime.Regexp:
		d.TypeName = runtime.TypeRegexp.Name()
	case runtime.Binary:
		d.TypeName = runtime.TypeBinary.Name()
		d.Length = len(v)
	case *runtime.Binary:
		d.TypeName = runtime.TypeBinary.Name()
		if v != nil {
			d.Length = len(*v)
		}
	case *runtime.Array:
		d.TypeName = runtime.TypeArray.Name()
		if v != nil {
			length, _ := v.Length(context.Background())
			d.Length = int(length)
			d.Children = indexedChildren(d.Length)
		}
	case *data.DataSet:
		d.TypeName = runtime.TypeArray.Name()
		d.Length = v.BufferedLength()
		d.Children = indexedChildren(d.Length)
	case *runtime.Object, *data.FastObject:
		d.TypeName = runtime.TypeObject.Name()
		if object := debugObject(value); object != nil {
			length, _ := object.Length(context.Background())
			d.Length = int(length)
			d.Children = namedChildren(d.Length, object.DebugFieldsPageable())
		}
	default:
		describeHost(value, &d)
	}

	return d
}

func describeHost(value runtime.Value, d *DebugValueDescription) {
	switch value.(type) {
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

	if inspectable, ok := value.(runtime.DebugInspectable); ok {
		d.Hint = safeDebugInfo(inspectable)
	}

	if provider, ok := value.(runtime.DebugChildrenProvider); ok {
		d.Children = safeChildrenInfo(provider)
	}
}

func indexedChildren(length int) runtime.DebugChildrenInfo {
	return runtime.DebugChildrenInfo{
		Named:   runtime.DebugChildView{Known: true},
		Indexed: runtime.DebugChildView{Known: length >= 0, Count: int64(max(0, length)), Readable: length > 0},
	}
}

func namedChildren(length int, pageable bool) runtime.DebugChildrenInfo {
	return runtime.DebugChildrenInfo{
		Indexed: runtime.DebugChildView{Known: true},
		Named:   runtime.DebugChildView{Known: true, Count: int64(length), Readable: pageable && length > 0},
	}
}

func safeChildrenInfo(provider runtime.DebugChildrenProvider) (info runtime.DebugChildrenInfo) {
	defer func() {
		if recover() != nil {
			info = runtime.DebugChildrenInfo{}
		}
	}()
	info = provider.DebugChildrenInfo()
	if info.Named.Known && info.Named.Count < 0 || info.Indexed.Known && info.Indexed.Count < 0 {
		return runtime.DebugChildrenInfo{}
	}

	return info
}

func describeScalarPointer[T runtime.Value](value *T, typ runtime.Type) DebugValueDescription {
	d := DebugValueDescription{TypeName: typ.Name(), Length: -1}
	if value != nil {
		d.Scalar = *value
	}

	return d
}

// Recover only around the optional embedder hook.
func safeDebugInfo(value runtime.DebugInspectable) (info runtime.DebugInfo) {
	defer func() {
		if recover() != nil {
			info = runtime.DebugInfo{}
		}
	}()

	return value.DebugInfo()
}
