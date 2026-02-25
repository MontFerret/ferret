package sdk

import (
	"context"
	"reflect"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const maxInt64 = ^uint64(0) >> 1

var byteSliceType = reflect.TypeOf([]byte(nil))

// Encode converts a Go value into a runtime Value using ferret tags for structs.
// If "ferret" tag is not present, it falls back to "json" tag, otherwise the field will be ignored.
// It also supports unwrapping values that implement the Unwrappable interface.
func Encode(input any) runtime.Value {
	if input == nil {
		return runtime.None
	}

	if value, ok := input.(runtime.Value); ok {
		return value
	}

	return encodeValue(reflect.ValueOf(input))
}

func encodeValue(v reflect.Value) runtime.Value {
	if !v.IsValid() {
		return runtime.None
	}

	v, ok := derefValue(v)
	if !ok {
		return runtime.None
	}

	if v.CanInterface() {
		if value, ok := v.Interface().(runtime.Value); ok {
			return value
		}
	}

	if value, ok := encodeSpecial(v); ok {
		return value
	}

	if value, ok := encodeScalar(v); ok {
		return value
	}

	return encodeComposite(v)
}

func derefValue(v reflect.Value) (reflect.Value, bool) {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return reflect.Value{}, false
		}
		return derefValue(v.Elem())
	default:
		return v, true
	}
}

func encodeSpecial(v reflect.Value) (runtime.Value, bool) {
	if v.Type() == timeType {
		if v.CanInterface() {
			return runtime.NewDateTime(v.Interface().(time.Time)), true
		}
		return runtime.None, true
	}

	if v.Type() == byteSliceType && v.Kind() == reflect.Slice {
		return runtime.NewBinary(v.Bytes()), true
	}

	return runtime.None, false
}

func encodeScalar(v reflect.Value) (runtime.Value, bool) {
	switch v.Kind() {
	case reflect.Bool:
		return runtime.NewBoolean(v.Bool()), true
	case reflect.String:
		return runtime.NewString(v.String()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return runtime.NewInt64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u := v.Uint()
		if u > maxInt64 {
			return runtime.None, true
		}
		return runtime.NewInt64(int64(u)), true
	case reflect.Float32, reflect.Float64:
		return runtime.NewFloat(v.Float()), true
	default:
		return runtime.None, false
	}
}

func encodeComposite(v reflect.Value) runtime.Value {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return encodeArrayLike(v)
	case reflect.Map:
		return encodeMap(v)
	case reflect.Struct:
		return encodeStruct(v)
	default:
		return runtime.None
	}
}

func encodeArrayLike(v reflect.Value) runtime.Value {
	size := v.Len()
	arr := runtime.NewArray(size)
	ctx := context.Background()

	for i := 0; i < size; i++ {
		_ = arr.Append(ctx, encodeValue(v.Index(i)))
	}

	return arr
}

func encodeMap(v reflect.Value) runtime.Value {
	obj := runtime.NewObject()
	ctx := context.Background()

	for _, key := range v.MapKeys() {
		keyVal := encodeValue(key)
		_ = obj.Set(ctx, runtime.NewString(keyVal.String()), encodeValue(v.MapIndex(key)))
	}

	return obj
}

func encodeStruct(v reflect.Value) runtime.Value {
	obj := runtime.NewObject()
	ctx := context.Background()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name, ok := Tag(field)
		if !ok {
			continue
		}

		_ = obj.Set(ctx, runtime.NewString(name), encodeValue(v.Field(i)))
	}

	return obj
}
