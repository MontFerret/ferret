package runtime

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash/fnv"
	"io"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wI2L/jettison"
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
// It supports basic types like bool, string, int, float, time.Time, as well as slices and maps.
// For unsupported types, it returns None.
// It does not use "ferret" tags for struct fields and instead relies on field names directly.
// For more safe and controlled parsing, consider using the "ferret" tags and the Encode function.
func Parse(input interface{}) Value {
	switch value := input.(type) {
	case bool:
		return NewBoolean(value)
	case string:
		return NewString(value)
	case int64:
		return NewInt(int(value))
	case int32:
		return NewInt(int(value))
	case int16:
		return NewInt(int(value))
	case int8:
		return NewInt(int(value))
	case int:
		return NewInt(value)
	case float64:
		return NewFloat(value)
	case float32:
		return NewFloat(float64(value))
	case time.Time:
		return NewDateTime(value)
	case []any:
		ctx := context.Background()
		arr := NewArray(len(value))

		for _, el := range value {
			_ = arr.Append(ctx, Parse(el))
		}

		return arr
	case map[string]any:
		ctx := context.Background()
		obj := NewObject()

		for key, el := range value {
			_ = obj.Set(ctx, NewString(key), Parse(el))
		}

		return obj
	case []byte:
		return NewBinary(value)
	case nil:
		return None
	case Value:
		return value
	default:
		v := reflect.ValueOf(value)
		t := reflect.TypeOf(value)
		kind := t.Kind()

		if kind == reflect.Ptr {
			el := v.Elem()

			if el.Kind() == 0 {
				return None
			}

			return Parse(el.Interface())
		}

		if kind == reflect.Slice || kind == reflect.Array {
			size := v.Len()
			arr := NewArray(size)
			ctx := context.Background()

			for i := 0; i < size; i++ {
				curVal := v.Index(i)
				_ = arr.Append(ctx, Parse(curVal.Interface()))
			}

			return arr
		}

		if kind == reflect.Map {
			keys := v.MapKeys()
			obj := NewObject()
			ctx := context.Background()

			for _, k := range keys {
				key := Parse(k.Interface())
				curVal := v.MapIndex(k)

				_ = obj.Set(ctx, NewString(key.String()), Parse(curVal.Interface()))
			}

			return obj
		}

		if kind == reflect.Struct {
			obj := NewObject()
			size := t.NumField()
			ctx := context.Background()

			for i := 0; i < size; i++ {
				field := t.Field(i)
				if field.PkgPath != "" {
					continue
				}
				fieldValue := v.Field(i)
				if !fieldValue.CanInterface() {
					continue
				}

				_ = obj.Set(ctx, NewString(field.Name), Parse(fieldValue.Interface()))
			}

			return obj
		}

		return None
	}
}

func Unmarshal(value json.RawMessage) (Value, error) {
	var o any

	err := json.Unmarshal(value, &o)

	if err != nil {
		return None, err
	}

	return Parse(o), nil
}

func MustMarshal(value Value) json.RawMessage {
	out, err := value.MarshalJSON()

	if err != nil {
		panic(err)
	}

	return out
}

func MustMarshalAny(input any) json.RawMessage {
	out, err := jettison.MarshalOpts(input, jettison.NoHTMLEscaping())

	if err != nil {
		panic(err)
	}

	return out
}

// IsScalar checks if the input Value is of a scalar type (Int, Float, String, or Boolean).
func IsScalar(input Value) Boolean {
	switch input.(type) {
	case Int, Float, String, Boolean:
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

func ToBoolean(input Value) Boolean {
	if input == None {
		return False
	}

	switch val := input.(type) {
	case Boolean:
		return val
	case String:
		return val != ""
	case Int:
		return val != 0
	case Float:
		return val != 0
	case DateTime:
		return Boolean(!val.IsZero())
	default:
		return True
	}
}

func ToFloat(ctx context.Context, input Value) (Float, error) {
	switch val := input.(type) {
	case Float:
		return val, nil
	case Int:
		return Float(val), nil
	case String:
		i, err := strconv.ParseFloat(string(val), 64)

		if err != nil {
			return ZeroFloat, err
		}

		return Float(i), nil
	case Boolean:
		if val {
			return Float(1), nil
		}

		return Float(0), nil
	case DateTime:
		dt := input.(DateTime)

		if dt.IsZero() {
			return ZeroFloat, nil
		}

		return NewFloat(float64(dt.Unix())), nil
	case List:
		iterator, err := val.Iterate(ctx)

		if err != nil {
			return ZeroFloat, err
		}

		res := ZeroFloat

		for {
			val, _, err := iterator.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}

			if errors.Is(err, ErrTimeout) {
				break
			}

			if err != nil {
				continue
			}

			f, err := ToFloat(ctx, val)

			if err != nil {
				continue
			}

			res += f
		}

		return res, nil
	default:
		return ZeroFloat, TypeErrorOf(input, TypeFloat)
	}
}

func ToString(input Value) String {
	switch val := input.(type) {
	case String:
		return val
	default:
		return NewString(val.String())
	}
}

func ToInt(ctx context.Context, input Value) (Int, error) {
	switch val := input.(type) {
	case Int:
		return val, nil
	case Float:
		return Int(val), nil
	case String:
		i, err := strconv.ParseInt(string(val), 10, 64)

		if err != nil {
			return ZeroInt, err
		}

		return Int(i), nil
	case Boolean:
		if val {
			return Int(1), nil
		}

		return Int(0), nil
	case DateTime:
		dt := input.(DateTime)

		if dt.IsZero() {
			return ZeroInt, nil
		}

		return NewInt(int(dt.Unix())), nil
	case List:
		iterator, err := val.Iterate(ctx)

		if err != nil {
			return ZeroInt, err
		}

		res := ZeroInt

		for {
			item, _, err := iterator.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}

			if errors.Is(err, ErrTimeout) {
				break
			}

			if err != nil {
				continue
			}

			i, err := ToInt(ctx, item)

			if err != nil {
				continue
			}

			res += i
		}

		return res, nil
	default:
		return ZeroInt, TypeErrorOf(input, TypeInt)
	}
}

func ToIntSafe(ctx context.Context, input Value) Int {
	result, err := ToInt(ctx, input)

	if err != nil {
		return ZeroInt
	}

	if result > 0 {
		return result
	}

	return ZeroInt
}

func ToNumberOnly(ctx context.Context, input Value) Value {
	switch value := input.(type) {
	case Int, Float:
		return value
	case String:
		if strings.Contains(value.String(), ".") {
			if val, err := ToFloat(ctx, value); err == nil {
				return val
			}

			return ZeroFloat
		}

		if val, err := ToInt(ctx, value); err == nil {
			return val
		}

		return ZeroFloat
	case Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return ZeroInt
		}

		i := ZeroInt
		f := ZeroFloat

		for {
			val, _, err := iterator.Next(ctx)

			if errors.Is(err, io.EOF) {
				break
			}

			if errors.Is(err, ErrTimeout) {
				break
			}

			if err != nil {
				continue
			}

			out := ToNumberOnly(ctx, val)

			switch item := out.(type) {
			case Int:
				i += item
			case Float:
				f += item
			}
		}

		if f == 0 {
			return i
		}

		return Float(i) + f
	default:
		if val, err := ToFloat(ctx, value); err == nil {
			return val
		}

		return ZeroInt
	}
}

// ToIntDefault attempts to convert an arbitrary Value into an Int.
// If the conversion fails or if the resulting Int is not greater than zero, it returns the provided defaultValue.
// This function is useful for safely converting values to Int while providing a fallback option in case of errors or non-positive results.
func ToIntDefault(ctx context.Context, input Value, defaultValue Int) (Int, error) {
	result, err := ToInt(ctx, input)

	if err != nil {
		return defaultValue, err
	}

	if result > 0 {
		return result, nil
	}

	return defaultValue, nil
}

// ToNumberOrString attempts to convert an arbitrary Value into either a Number (Int or Float) or a String.
// If the input is already an Int, Float, or String, it returns it directly.
// For other types, it converts the input to its string representation and returns it as a String.
// This allows for flexible conversion of various Value types into a format that can be easily used as either a number or a string.
func ToNumberOrString(input Value) Value {
	switch value := input.(type) {
	case Int, Float, String:
		return value
	default:
		return ToString(value)
	}
}

// ToBinary attempts to convert an arbitrary Value into a Binary type.
// If the input is already a Binary, it returns it directly.
// For other types, it converts the input to its string representation and then to a byte slice to create a new Binary.
// This allows for flexible conversion of various Value types into a Binary format, using their string representation as the basis for the binary data.
func ToBinary(input Value) Binary {
	bin, ok := input.(Binary)

	if ok {
		return bin
	}

	return NewBinary([]byte(input.String()))
}

// Hash computes a hash value for the given typename and content using the FNV-1a hashing algorithm.
// It concatenates the typename, a colon, and the content bytes to generate a unique hash value.
// This function can be used to create consistent hash values for different types of content based on their type and content.
func Hash(typename string, content []byte) uint64 {
	h := fnv.New64a()

	h.Write([]byte(typename))
	h.Write([]byte(":"))
	h.Write(content)

	return h.Sum64()
}

// MapHash computes a hash value for a map of string keys to Value values.
// It uses the FNV-1a hashing algorithm to generate a consistent hash based on the keys and their corresponding value hashes.
// The keys are sorted to ensure that the order of key-value pairs does not affect the resulting hash.
func MapHash(input map[string]Value) uint64 {
	h := fnv.New64a()

	keys := make([]string, 0, len(input))

	for key := range input {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	h.Write([]byte("{"))

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := input[key]

		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, el.Hash())

		h.Write(bytes)

		if idx != endIndex {
			h.Write([]byte(","))
		}
	}

	h.Write([]byte("}"))

	return h.Sum64()
}
