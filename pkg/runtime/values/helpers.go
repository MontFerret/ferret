package values

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func Parse(input interface{}) core.Value {
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
	case []interface{}:
		arr := NewArray(len(value))

		for _, el := range value {
			arr.Push(Parse(el))
		}

		return arr
	case map[string]interface{}:
		obj := NewObject()

		for key, el := range value {
			obj.Set(NewString(key), Parse(el))
		}

		return obj
	case []byte:
		return NewBinary(value)
	case nil:
		return None
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

			for i := 0; i < size; i++ {
				curVal := v.Index(i)
				arr.Push(Parse(curVal.Interface()))
			}

			return arr
		}

		if kind == reflect.Map {
			keys := v.MapKeys()
			obj := NewObject()

			for _, k := range keys {
				key := Parse(k.Interface())
				curVal := v.MapIndex(k)

				obj.Set(NewString(key.String()), Parse(curVal.Interface()))
			}

			return obj
		}

		if kind == reflect.Struct {
			obj := NewObject()
			size := t.NumField()

			for i := 0; i < size; i++ {
				field := t.Field(i)
				fieldValue := v.Field(i)

				obj.Set(NewString(field.Name), Parse(fieldValue.Interface()))
			}

			return obj
		}

		return None
	}
}

func Unmarshal(value json.RawMessage) (core.Value, error) {
	var o interface{}

	err := json.Unmarshal(value, &o)

	if err != nil {
		return None, err
	}

	return Parse(o), nil
}

func MustMarshal(value core.Value) json.RawMessage {
	out, err := value.MarshalJSON()

	if err != nil {
		panic(err)
	}

	return out
}

func MustMarshalAny(input interface{}) json.RawMessage {
	out, err := jettison.MarshalOpts(input, jettison.NoHTMLEscaping())

	if err != nil {
		panic(err)
	}

	return out
}

func IsScalar(input core.Value) Boolean {
	switch input.(type) {
	case Int, Float, String, Boolean:
		return true
	default:
		return false
	}
}

func IsNumber(input core.Value) Boolean {
	switch input.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

func ToBoolean(input core.Value) Boolean {
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
		return val.IsZero() == false
	case *none:
		return False
	default:
		return True
	}
}

func ToFloat(input core.Value) Float {
	switch val := input.(type) {
	case Float:
		return val
	case Int:
		return Float(val)
	case String:
		i, err := strconv.ParseFloat(string(val), 64)

		if err != nil {
			return ZeroFloat
		}

		return Float(i)
	case Boolean:
		if val {
			return Float(1)
		}

		return Float(0)
	case DateTime:
		dt := input.(DateTime)

		if dt.IsZero() {
			return ZeroFloat
		}

		return NewFloat(float64(dt.Unix()))
	case *Array:
		length := val.Length()

		if length == 0 {
			return ZeroFloat
		}

		res := ZeroFloat

		for i := Int(0); i < length; i++ {
			res += ToFloat(val.Get(i))
		}

		return res
	default:
		return ZeroFloat
	}
}

func ToString(input core.Value) String {
	switch val := input.(type) {
	case String:
		return val
	default:
		return NewString(val.String())
	}
}

func ToInt(input core.Value) Int {
	switch val := input.(type) {
	case Int:
		return val
	case Float:
		return Int(val)
	case String:
		i, err := strconv.ParseInt(string(val), 10, 64)

		if err != nil {
			return ZeroInt
		}

		return Int(i)
	case Boolean:
		if val {
			return Int(1)
		}

		return Int(0)
	case DateTime:
		dt := input.(DateTime)

		if dt.IsZero() {
			return ZeroInt
		}

		return NewInt(int(dt.Unix()))
	case *Array:
		length := val.Length()

		if length == 0 {
			return ZeroInt
		}

		res := ZeroInt

		for i := Int(0); i < length; i++ {
			res += ToInt(val.Get(i))
		}

		return res
	default:
		return ZeroInt
	}
}

func ToIntDefault(input core.Value, defaultValue Int) Int {
	if result := ToInt(input); result > 0 {
		return result
	}

	return defaultValue
}

func ToArray(ctx context.Context, input core.Value) *Array {
	switch value := input.(type) {
	case Boolean,
		Int,
		Float,
		String,
		DateTime:

		return NewArrayWith(value)
	case *Array:
		return value.Copy().(*Array)
	case *Object:
		arr := NewArray(int(value.Length()))

		value.ForEach(func(value core.Value, key string) bool {
			arr.Push(value)

			return true
		})

		return arr
	case core.Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return NewArray(0)
		}

		arr := NewArray(10)

		for {
			val, _, err := iterator.Next(ctx)

			if err != nil {
				return NewArray(0)
			}

			if val == None {
				break
			}

			arr.Push(val)
		}

		return arr
	default:
		return EmptyArray()
	}
}

func ToObject(ctx context.Context, input core.Value) *Object {
	switch value := input.(type) {
	case *Object:
		return value
	case *Array:
		obj := NewObject()

		value.ForEach(func(value core.Value, idx int) bool {
			obj.Set(ToString(Int(idx)), value)

			return true
		})

		return obj
	case core.Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return NewObject()
		}

		obj := NewObject()

		for {
			val, key, err := iterator.Next(ctx)

			if err != nil {
				return obj
			}

			if val == None {
				break
			}

			obj.Set(String(key.String()), val)
		}

		return obj
	default:
		return NewObject()
	}
}

func ToStrings(input *Array) []String {
	res := make([]String, input.Length())

	input.ForEach(func(v core.Value, i int) bool {
		res[i] = NewString(v.String())

		return true
	})

	return res
}

func ToBinary(input core.Value) Binary {
	bin, ok := input.(Binary)

	if ok {
		return bin
	}

	return NewBinary([]byte(input.String()))
}

func Hash(typename string, content []byte) uint64 {
	h := fnv.New64a()

	h.Write([]byte(typename))
	h.Write([]byte(":"))
	h.Write(content)

	return h.Sum64()
}

func MapHash(input map[string]core.Value) uint64 {
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

func UnwrapStrings(values []String) []string {
	out := make([]string, len(values))

	for i, v := range values {
		out[i] = v.String()
	}

	return out
}

func Negate(input core.Value) core.Value {
	switch val := input.(type) {
	case Int:
		return -val
	case Float:
		return -val
	case Boolean:
		return !val
	default:
		return None
	}
}

func Negative(input core.Value) core.Value {
	switch value := input.(type) {
	case Int:
		return -value
	case Float:
		return -value
	default:
		// TODO: Maybe we should return None?
		return input
	}
}

func Positive(input core.Value) core.Value {
	switch value := input.(type) {
	case Int:
		return +value
	case Float:
		return +value
	default:
		// TODO: Maybe we should return None?
		return input
	}
}

func Contains(input core.Value, value core.Value) Boolean {
	switch val := input.(type) {
	case *Array:
		return val.Contains(value)
	case String:
		return Boolean(strings.Contains(val.String(), value.String()))
	default:
		return false
	}
}

func ToNumberOrString(input core.Value) core.Value {
	switch value := input.(type) {
	case Int, Float, String:
		return value
	default:
		return ToString(value)
	}
}

func ToNumberOnly(input core.Value) core.Value {
	switch value := input.(type) {
	case Int, Float:
		return value
	case String:
		if strings.Contains(value.String(), ".") {
			return ToFloat(input)
		}

		return ToInt(input)
	case *Array:
		length := value.Length()

		if length == 0 {
			return ZeroInt
		}

		i := ZeroInt
		f := ZeroFloat

		for y := Int(0); y < length; y++ {
			out := ToNumberOnly(value.Get(y))

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
		return ToInt(input)
	}
}

func CompareStrings(a, b String) Int {
	return Int(strings.Compare(a.String(), b.String()))
}
