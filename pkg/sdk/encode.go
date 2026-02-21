package sdk

import (
	"context"
	"fmt"
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

// EncodeField reads a field from a struct or a value from a map by the provided key and encodes it into a runtime Value.
// It supports unwrapping values that implement the Unwrappable interface.
func EncodeField(ctx context.Context, input any, key runtime.Value) (runtime.Value, error) {
	if input == nil {
		return runtime.None, nil
	}

	wrapper, ok := input.(runtime.Unwrappable)

	if ok {
		input = wrapper.Unwrap()
	}

	t := reflect.TypeOf(input)
	keyStr := key.String()

	switch t.Kind() {
	case reflect.Struct:
		v := reflect.ValueOf(input)

		var name string

		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)

			tagName, ok := runtime.Tag(f)
			if !ok {
				continue
			}

			if tagName == keyStr {
				if !f.IsExported() {
					return runtime.None, nil
				}

				name = f.Name

				break
			}
		}

		if name == "" {
			return runtime.None, nil
		}

		field := v.FieldByName(name)

		if field.IsValid() {
			return Encode(field.Interface()), nil
		}

		return runtime.None, nil
	case reflect.Map:
		v := reflect.ValueOf(input)
		mapKeyType := v.Type().Key()

		var mapKeyVal reflect.Value
		switch mapKeyType.Kind() {
		case reflect.String:
			// Use the string representation of the runtime.Value as the map key.
			mapKeyVal = reflect.ValueOf(keyStr)
		default:
			keyType := reflect.TypeOf(key)
			if keyType == nil || !keyType.AssignableTo(mapKeyType) {
				// Cannot use the provided key for this map type.
				return runtime.None, nil
			}
			mapKeyVal = reflect.ValueOf(key)
		}

		field := v.MapIndex(mapKeyVal)
		if field.IsValid() {
			return Encode(field.Interface()), nil
		}

		return runtime.None, nil
	case reflect.Ptr:
		v := reflect.ValueOf(input)

		if v.IsNil() {
			return runtime.None, nil
		}

		return EncodeField(ctx, v.Elem().Interface(), key)
	default:
		break
	}

	return runtime.None, fmt.Errorf("cannot read field by key from type: %s", t.String())
}

func encodeValue(v reflect.Value) runtime.Value {
	if !v.IsValid() {
		return runtime.None
	}

	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return runtime.None
		}
		return encodeValue(v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			return runtime.None
		}
		return encodeValue(v.Elem())
	}

	if v.CanInterface() {
		if value, ok := v.Interface().(runtime.Value); ok {
			return value
		}
	}

	if v.Type() == timeType {
		if v.CanInterface() {
			return runtime.NewDateTime(v.Interface().(time.Time))
		}
		return runtime.None
	}

	if v.Type() == byteSliceType && v.Kind() == reflect.Slice {
		return runtime.NewBinary(v.Bytes())
	}

	switch v.Kind() {
	case reflect.Bool:
		return runtime.NewBoolean(v.Bool())
	case reflect.String:
		return runtime.NewString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return runtime.NewInt64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u := v.Uint()
		if u > maxInt64 {
			return runtime.None
		}
		return runtime.NewInt64(int64(u))
	case reflect.Float32, reflect.Float64:
		return runtime.NewFloat(v.Float())
	case reflect.Slice, reflect.Array:
		size := v.Len()
		arr := runtime.NewArray(size)
		ctx := context.Background()

		for i := 0; i < size; i++ {
			_ = arr.Append(ctx, encodeValue(v.Index(i)))
		}

		return arr
	case reflect.Map:
		obj := runtime.NewObject()
		ctx := context.Background()

		for _, key := range v.MapKeys() {
			keyVal := encodeValue(key)
			_ = obj.Set(ctx, runtime.NewString(keyVal.String()), encodeValue(v.MapIndex(key)))
		}

		return obj
	case reflect.Struct:
		obj := runtime.NewObject()
		ctx := context.Background()
		t := v.Type()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}

			name, ok := runtime.Tag(field)
			if !ok {
				continue
			}

			_ = obj.Set(ctx, runtime.NewString(name), encodeValue(v.Field(i)))
		}

		return obj
	default:
		return runtime.None
	}
}
