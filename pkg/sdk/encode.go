package sdk

import (
	"context"
	"errors"
	"io"
	"reflect"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const maxInt64 = ^uint64(0) >> 1

var byteSliceType = reflect.TypeOf([]byte(nil))

// Encode converts a Go value into a runtime Value.
// It handles basic types, slices, maps, and structs (using tags for field names).
// If the input is already a runtime Value, it returns it directly.
// For unsupported types, it returns runtime.None.
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
	return encodeStructWithVisit(v, make(map[reflect.Type]int))
}

func encodeStructWithVisit(v reflect.Value, visiting map[reflect.Type]int) runtime.Value {
	obj := runtime.NewObject()
	ctx := context.Background()
	t := v.Type()
	if visiting[t] > 0 {
		return obj
	}

	visiting[t]++
	defer func() {
		visiting[t]--
		if visiting[t] == 0 {
			delete(visiting, t)
		}
	}()

	used := make(map[string]struct{}, t.NumField())

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
		used[name] = struct{}{}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" || !field.Anonymous {
			continue
		}

		if _, ok := Tag(field); ok {
			continue
		}

		fieldVal := v.Field(i)
		switch fieldVal.Kind() {
		case reflect.Struct:
			mergeObject(obj, encodeStructWithVisit(fieldVal, visiting), used)
		case reflect.Pointer:
			if fieldVal.IsNil() || fieldVal.Elem().Kind() != reflect.Struct {
				continue
			}

			if visiting[fieldVal.Elem().Type()] > 0 {
				continue
			}

			mergeObject(obj, encodeStructWithVisit(fieldVal.Elem(), visiting), used)
		default:
			continue
		}
	}

	return obj
}

func mergeObject(dst *runtime.Object, src runtime.Value, used map[string]struct{}) {
	m, ok := src.(runtime.Map)
	if !ok {
		return
	}

	ctx := context.Background()
	keys, err := m.Keys(ctx)
	if err != nil {
		return
	}

	iter, err := keys.Iterate(ctx)
	if err != nil {
		return
	}

	for {
		keyVal, _, err := iter.Next(ctx)
		if errors.Is(err, io.EOF) || errors.Is(err, runtime.ErrTimeout) {
			break
		}

		if err != nil {
			return
		}

		key, ok := keyVal.(runtime.String)
		if !ok {
			continue
		}

		keyStr := key.String()
		if _, exists := used[keyStr]; exists {
			continue
		}

		val, err := m.Get(ctx, keyVal)
		if err != nil {
			return
		}

		_ = dst.Set(ctx, runtime.NewString(keyStr), val)
		used[keyStr] = struct{}{}
	}
}
