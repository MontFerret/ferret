package runtime

import (
	"context"
	"errors"
	"io"
	"math"
	"math/rand"
	"reflect"
	"time"
)

func IsNil(input any) bool {
	val := reflect.ValueOf(input)
	kind := val.Kind()

	switch kind {
	case reflect.Ptr,
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.Func,
		reflect.Interface,
		reflect.Chan:
		return val.IsNil()
	case reflect.Struct,
		reflect.UnsafePointer:
		return false
	case reflect.Invalid:
		return true
	default:
		return false
	}
}

func NumberBoundaries(input float64) (max float64, min float64) {
	min = input / 2
	max = input * 2

	return
}

func NumberUpperBoundary(input float64) float64 {
	return input * 2
}

func NumberLowerBoundary(input float64) float64 {
	return input / 2
}

func RandomDefault() float64 {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rnd.Float64()
}

func Random(max float64, min float64) float64 {
	r := RandomDefault()
	i := r * (max - min + 1)
	out := math.Floor(i) + min

	return out
}

func Random2(mid float64) float64 {
	randMax, randMin := NumberBoundaries(mid)

	return Random(randMax, randMin)
}

// Parse attempts to convert an arbitrary input into a Value type.
// Deprecated: Use ValueOf for explicit host-to-Ferret value conversion with error handling.
func Parse(input any) Value {
	parsed, err := ValueOf(input)

	if err != nil {
		return None
	}

	return parsed
}

// ValueOf converts a native Go value into a Ferret runtime Value.
// It returns an error if the value cannot be converted.
//
// This is the preferred API for host-to-Ferret value conversion.
//
// For legacy permissive behavior that silently converts unsupported
// values to None, see Parse.
func ValueOf(input any) (Value, error) {
	switch value := input.(type) {
	case nil:
		return None, nil
	case Value:
		return value, nil
	case []Value:
		ctx := context.Background()
		arr := NewArray(len(value))

		for _, el := range value {
			_ = arr.Append(ctx, el)
		}

		return arr, nil
	case bool:
		return NewBoolean(value), nil
	case string:
		return NewString(value), nil
	case int64:
		return NewInt64(value), nil
	case int32:
		return Int(value), nil
	case int16:
		return Int(value), nil
	case int8:
		return Int(value), nil
	case int:
		return NewInt(value), nil
	case uint:
		if uint64(value) > math.MaxInt64 {
			return None, Errorf(ErrRange, "invalid integer %d exceeds runtime.Int range", value)
		}

		return Int(value), nil
	case uint8:
		return Int(value), nil
	case uint16:
		return Int(value), nil
	case uint32:
		return Int(value), nil
	case uint64:
		if value > math.MaxInt64 {
			return None, Errorf(ErrRange, "invalid integer %d exceeds runtime.Int range", value)
		}

		return Int(value), nil
	case float64:
		return NewFloat(value), nil
	case float32:
		return NewFloat(float64(value)), nil
	case time.Time:
		return NewDateTime(value), nil
	case time.Duration:
		return NewDuration(value), nil
	case []any:
		ctx := context.Background()
		arr := NewArray(len(value))

		for idx, el := range value {
			parsed, err := ValueOf(el)

			if err != nil {
				return None, Errorf(err, "at index %d", idx)
			}

			_ = arr.Append(ctx, parsed)
		}

		return arr, nil
	case map[string]any:
		ctx := context.Background()
		obj := NewObject()

		for key, el := range value {
			parsed, err := ValueOf(el)

			if err != nil {
				return None, Errorf(err, "at key %q", key)
			}

			_ = obj.Set(ctx, NewString(key), parsed)
		}

		return obj, nil
	case map[any]any:
		ctx := context.Background()
		obj := NewObject()

		for key, el := range value {
			parsedVal, err := ValueOf(el)

			if err != nil {
				return None, Errorf(err, "at key %v", key)
			}

			_ = obj.Set(ctx, NewStringOf(key), parsedVal)
		}

		return obj, nil
	case []byte:
		return NewBinary(value), nil
	default:
		v := reflect.ValueOf(value)
		t := reflect.TypeOf(value)
		kind := t.Kind()
		ctx := context.Background()

		if kind == reflect.Ptr {
			el := v.Elem()

			if el.Kind() == 0 {
				return None, nil
			}

			return ValueOf(el.Interface())
		}

		if kind == reflect.Slice || kind == reflect.Array {
			size := v.Len()
			arr := NewArray(size)

			for i := 0; i < size; i++ {
				val, err := ValueOf(v.Index(i).Interface())

				if err != nil {
					return None, Errorf(err, "at index %d", i)
				}

				_ = arr.Append(ctx, val)
			}

			return arr, nil
		}

		if kind == reflect.Map {
			keys := v.MapKeys()
			obj := NewObject()

			for _, k := range keys {
				key, err := ValueOf(k.Interface())

				if err != nil {
					return None, Errorf(err, "at key %v", k.Interface())
				}

				val, err := ValueOf(v.MapIndex(k).Interface())

				if err != nil {
					return None, Errorf(err, "at key %v", k.Interface())
				}

				_ = obj.Set(ctx, NewString(key.String()), val)
			}

			return obj, nil
		}

		if kind == reflect.Struct {
			obj := NewObject()
			size := t.NumField()

			for i := 0; i < size; i++ {
				field := t.Field(i)
				if field.PkgPath != "" {
					continue
				}

				fieldValue := v.Field(i)
				if !fieldValue.CanInterface() {
					continue
				}

				parsed, err := ValueOf(fieldValue.Interface())

				if err != nil {
					return None, Errorf(err, "at field %q", field.Name)
				}

				_ = obj.Set(ctx, NewString(field.Name), parsed)
			}

			return obj, nil
		}

		return None, Errorf(ErrInvalidType, "cannot parse type %T", input)
	}
}

// IsScalar checks if the input Value is a scalar value.
func IsScalar(input Value) Boolean {
	switch input.(type) {
	case Int, Float, Duration, String, Boolean:
		return true
	default:
		return false
	}
}

// IsNumber checks if the input Value is of type Int or Float, indicating that it is a numeric type.
func IsNumber(input Value) Boolean {
	switch input.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

// ToList attempts to convert an arbitrary Value into a List type.
// It supports basic types like Boolean, Int, Float, String, DateTime by wrapping them into a single-element List.
// For List types, it returns a copy of the List.
// For Iterable types, it iterates through the elements and appends them to a new List.
// For unsupported types, it returns an empty List.
func ToList(ctx context.Context, input Value) (List, error) {
	switch value := input.(type) {
	case Boolean,
		Int,
		Float,
		Duration,
		String,
		DateTime:

		return NewArrayWith(value), nil
	case List:
		return value.Copy().(List), nil
	case Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return nil, err
		}

		arr := NewArray(10)

		for {
			val, _, err := iterator.Next(ctx)

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return nil, err
			}

			_ = arr.Append(ctx, val)
		}

		return arr, nil
	default:
		return EmptyArray(), nil
	}
}

// ToMap attempts to convert an arbitrary Value into a Map type.
// It supports Map, Array, and Iterable types, converting them into a Map format.
// For unsupported types, it returns an empty Map.
// The function uses the string representation of keys for Arrays and Iterables, with array indices as keys for Arrays.
func ToMap(ctx context.Context, input Value) (Map, error) {
	switch value := input.(type) {
	case Map:
		return value, nil
	case *Array:
		obj := NewObject()

		for i, v := range value.data {
			_ = obj.Set(ctx, ToString(Int(i)), v)
		}

		return obj, nil
	case Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return nil, err
		}

		obj := NewObject()

		for {
			val, key, err := iterator.Next(ctx)

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return nil, err
			}

			_ = obj.Set(ctx, ToString(key), val)
		}

		return obj, nil
	default:
		return NewObject(), nil
	}
}
