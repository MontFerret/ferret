package runtime

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

const maxInt64 = ^uint64(0) >> 1

var byteSliceType = reflect.TypeOf([]byte(nil))

// Encode converts a Go value into a runtime Value using ferret tags for structs.
// If "ferret" tag is not present, it falls back to "json" tag, otherwise the field will be ignored.
// It also supports unwrapping values that implement the Unwrappable interface.
func Encode(input any) Value {
	if input == nil {
		return None
	}

	if value, ok := input.(Value); ok {
		return value
	}

	return encodeValue(reflect.ValueOf(input))
}

// EncodeField reads a field from a struct or a value from a map by the provided key and encodes it into a runtime Value.
// It supports unwrapping values that implement the Unwrappable interface.
func EncodeField(ctx context.Context, input any, key Value) (Value, error) {
	if input == nil {
		return None, nil
	}

	wrapper, ok := input.(Unwrappable)

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

			tagName, ok := Tag(f)
			if !ok {
				continue
			}

			if tagName == keyStr {
				if !f.IsExported() {
					return None, nil
				}

				name = f.Name

				break
			}
		}

		if name == "" {
			return None, nil
		}

		field := v.FieldByName(name)

		if field.IsValid() {
			return Encode(field.Interface()), nil
		}

		return None, nil
	case reflect.Map:
		v := reflect.ValueOf(input)
		field := v.MapIndex(reflect.ValueOf(key))

		if field.IsValid() {
			return Encode(field.Interface()), nil
		}

		return None, nil
	case reflect.Ptr:
		v := reflect.ValueOf(input)

		if v.IsNil() {
			return None, nil
		}

		return EncodeField(ctx, v.Elem().Interface(), key)
	default:
		break
	}

	return None, fmt.Errorf("cannot read field by key from type: %s", t.String())
}

func encodeValue(v reflect.Value) Value {
	if !v.IsValid() {
		return None
	}

	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return None
		}
		return encodeValue(v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			return None
		}
		return encodeValue(v.Elem())
	}

	if v.CanInterface() {
		if value, ok := v.Interface().(Value); ok {
			return value
		}
	}

	if v.Type() == timeType {
		if v.CanInterface() {
			return NewDateTime(v.Interface().(time.Time))
		}
		return None
	}

	if v.Type() == byteSliceType && v.Kind() == reflect.Slice {
		return NewBinary(v.Bytes())
	}

	switch v.Kind() {
	case reflect.Bool:
		return NewBoolean(v.Bool())
	case reflect.String:
		return NewString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInt64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u := v.Uint()
		if u > maxInt64 {
			return None
		}
		return NewInt64(int64(u))
	case reflect.Float32, reflect.Float64:
		return NewFloat(v.Float())
	case reflect.Slice, reflect.Array:
		size := v.Len()
		arr := NewArray(size)
		ctx := context.Background()

		for i := 0; i < size; i++ {
			_ = arr.Append(ctx, encodeValue(v.Index(i)))
		}

		return arr
	case reflect.Map:
		obj := NewObject()
		ctx := context.Background()

		for _, key := range v.MapKeys() {
			keyVal := encodeValue(key)
			_ = obj.Set(ctx, NewString(keyVal.String()), encodeValue(v.MapIndex(key)))
		}

		return obj
	case reflect.Struct:
		obj := NewObject()
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

			_ = obj.Set(ctx, NewString(name), encodeValue(v.Field(i)))
		}

		return obj
	default:
		return None
	}
}
