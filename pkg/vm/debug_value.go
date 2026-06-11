package vm

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

// DebugFormatOptions bounds debugger value traversal and rendered output.
type DebugFormatOptions struct {
	MaxDepth int
	MaxItems int
	MaxBytes int
}

// DefaultDebugFormatOptions returns conservative debugger formatting limits.
func DefaultDebugFormatOptions() DebugFormatOptions {
	return DebugFormatOptions{MaxDepth: 3, MaxItems: 8, MaxBytes: 1024}
}

// DebugFormatValue formats built-in values without invoking arbitrary host
// value methods.
func DebugFormatValue(value runtime.Value, options DebugFormatOptions) string {
	if options.MaxDepth <= 0 {
		options.MaxDepth = 3
	}
	if options.MaxItems <= 0 {
		options.MaxItems = 8
	}
	if options.MaxBytes <= 0 {
		options.MaxBytes = 1024
	}
	var b strings.Builder
	formatDebugValue(&b, value, options, 0)
	out := b.String()
	if len(out) > options.MaxBytes {
		return out[:options.MaxBytes] + "..."
	}
	return out
}

func formatDebugValue(b *strings.Builder, value runtime.Value, options DebugFormatOptions, depth int) {
	if value == nil || reflect.TypeOf(value) == reflect.TypeOf(runtime.None) {
		b.WriteString("NONE")
		return
	}
	switch value := value.(type) {
	case runtime.Boolean:
		if value {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case runtime.Int:
		b.WriteString(strconv.FormatInt(int64(value), 10))
	case runtime.Float:
		b.WriteString(strconv.FormatFloat(float64(value), 'g', -1, 64))
	case runtime.String:
		b.WriteString(strconv.Quote(boundedDebugText(value.String(), options.MaxBytes)))
	case runtime.DateTime:
		b.WriteString(value.String())
	case runtime.Binary:
		fmt.Fprintf(b, "Binary(%d)", len(value))
	case *runtime.Array:
		length, _ := value.Length(context.Background())
		if depth >= options.MaxDepth || int(length) > options.MaxItems {
			fmt.Fprintf(b, "Array(%d)", length)
			return
		}
		b.WriteByte('[')
		for i := 0; i < int(length); i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			item, _ := value.At(context.Background(), runtime.Int(i))
			formatDebugValue(b, item, options, depth+1)
		}
		b.WriteByte(']')
	case *runtime.Object:
		formatDebugMap(b, value, options, depth)
	case *data.FastObject:
		formatDebugMapValue(b, value, options, depth)
	default:
		fmt.Fprintf(b, "HostValue(%s)", DebugValueTypeName(value))
	}
}

// DebugValueTypeName returns a stable debugger type name without invoking
// opaque host-value methods.
func DebugValueTypeName(value runtime.Value) string {
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

func formatDebugMapValue(b *strings.Builder, value runtime.Map, options DebugFormatOptions, depth int) {
	length, _ := value.Length(context.Background())
	if depth >= options.MaxDepth || int(length) > options.MaxItems {
		fmt.Fprintf(b, "Object(%d)", length)
		return
	}
	keys, _ := value.Keys(context.Background())
	names := make([]string, 0, int(length))
	for i := runtime.Int(0); i < length; i++ {
		key, _ := keys.At(context.Background(), i)
		names = append(names, key.String())
	}
	sort.Strings(names)
	b.WriteByte('{')
	for i, name := range names {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Quote(boundedDebugText(name, options.MaxBytes)))
		b.WriteString(": ")
		item, _ := value.Get(context.Background(), runtime.NewString(name))
		formatDebugValue(b, item, options, depth+1)
	}
	b.WriteByte('}')
}

func formatDebugMap(b *strings.Builder, value *runtime.Object, options DebugFormatOptions, depth int) {
	length, _ := value.Length(context.Background())
	if depth >= options.MaxDepth || int(length) > options.MaxItems {
		fmt.Fprintf(b, "Object(%d)", length)
		return
	}
	keys, _ := value.Keys(context.Background())
	names := make([]string, 0, int(length))
	for i := runtime.Int(0); i < length; i++ {
		key, _ := keys.At(context.Background(), i)
		names = append(names, key.String())
	}
	sort.Strings(names)
	b.WriteByte('{')
	for i, name := range names {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Quote(boundedDebugText(name, options.MaxBytes)))
		b.WriteString(": ")
		item, _ := value.Get(context.Background(), runtime.NewString(name))
		formatDebugValue(b, item, options, depth+1)
	}
	b.WriteByte('}')
}

func boundedDebugText(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}

// DebugLookupValue performs side-effect-free reads on VM-owned and built-in
// collections. Opaque host collections are rejected.
func DebugLookupValue(value, key runtime.Value) (runtime.Value, error) {
	ctx := context.Background()
	switch value := value.(type) {
	case *runtime.Array:
		index, ok := key.(runtime.Int)
		if !ok {
			return nil, debugLookupTypeError(key, runtime.TypeInt.Name())
		}
		return value.At(ctx, index)
	case *runtime.Object:
		property, ok := key.(runtime.String)
		if !ok {
			return nil, debugLookupTypeError(key, runtime.TypeString.Name())
		}
		return value.Get(ctx, property)
	case *data.FastObject:
		property, ok := key.(runtime.String)
		if !ok {
			return nil, debugLookupTypeError(key, runtime.TypeString.Name())
		}
		return value.Get(ctx, property)
	default:
		return nil, runtime.Errorf(runtime.ErrInvalidOperation, "debugger cannot inspect %s", DebugValueTypeName(value))
	}
}

func debugLookupTypeError(value runtime.Value, expected string) error {
	return runtime.Errorf(runtime.ErrInvalidArgument, "debugger lookup requires %s, got %s", expected, DebugValueTypeName(value))
}
