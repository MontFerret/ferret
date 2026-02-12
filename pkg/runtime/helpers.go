package runtime

import (
	"encoding/binary"
	"encoding/json"
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

func Parse(ctx Context, input any) Value {
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
		arr := ctx.Alloc().Array(len(value))

		for _, el := range value {
			_ = arr.Append(ctx, Parse(ctx, el))
		}

		return arr
	case map[string]interface{}:
		obj := ctx.Alloc().Object(len(value))

		for key, el := range value {
			_ = obj.Set(ctx, NewString(key), Parse(ctx, el))
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

			return Parse(ctx, el.Interface())
		}

		if kind == reflect.Slice || kind == reflect.Array {
			size := v.Len()
			arr := ctx.Alloc().Array(size)

			for i := 0; i < size; i++ {
				curVal := v.Index(i)
				_ = arr.Append(ctx, Parse(ctx, curVal.Interface()))
			}

			return arr
		}

		if kind == reflect.Map {
			keys := v.MapKeys()
			obj := ctx.Alloc().Object(len(keys))

			for _, k := range keys {
				key := Parse(ctx, k.Interface())
				curVal := v.MapIndex(k)

				_ = obj.Set(ctx, NewString(key.String()), Parse(ctx, curVal.Interface()))
			}

			return obj
		}

		if kind == reflect.Struct {
			size := t.NumField()
			obj := ctx.Alloc().Object(size)

			for i := 0; i < size; i++ {
				field := t.Field(i)
				fieldValue := v.Field(i)

				_ = obj.Set(ctx, NewString(field.Name), Parse(ctx, fieldValue.Interface()))
			}

			return obj
		}

		return None
	}
}

func Unmarshal(ctx Context, value json.RawMessage) (Value, error) {
	var o interface{}

	err := json.Unmarshal(value, &o)

	if err != nil {
		return None, err
	}

	return Parse(ctx, o), nil
}

func MustMarshal(value Value) json.RawMessage {
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

func IsScalar(input Value) Boolean {
	switch input.(type) {
	case Int, Float, String, Boolean:
		return true
	default:
		return false
	}
}

func IsNumber(input Value) Boolean {
	switch input.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

func ToList(ctx Context, input Value) List {
	switch value := input.(type) {
	case Boolean,
		Int,
		Float,
		String,
		DateTime:

		arr := ctx.Alloc().Array(1)
		_ = arr.Append(ctx, value)

		return arr
	case List:
		cp, err := value.Copy(ctx)

		if err != nil {
			return ctx.Alloc().Array(0)
		}

		list, ok := cp.(List)

		if !ok {
			return ctx.Alloc().Array(0)
		}

		return list
	case Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return ctx.Alloc().Array(0)
		}

		arr := ctx.Alloc().Array(10)

		for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
			val, _, err := iterator.Next(ctx)

			if err != nil {
				return arr
			}

			_ = arr.Append(ctx, val)
		}

		return arr
	default:
		return ctx.Alloc().Array(0)
	}
}

func ToMap(ctx Context, input Value) Map {
	switch value := input.(type) {
	case Map:
		return value
	case *Array:
		size, _ := value.Length(ctx)
		obj := ctx.Alloc().Object(int(size))

		for i, v := range value.data {
			_ = obj.Set(ctx, ToString(Int(i)), v)
		}

		return obj
	case Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return ctx.Alloc().Object(0)
		}

		obj := ctx.Alloc().Object(10)

		for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
			val, key, err := iterator.Next(ctx)

			if err != nil {
				return obj
			}

			_ = obj.Set(ctx, ToString(key), val)
		}

		return obj
	default:
		return ctx.Alloc().Object(0)
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

func ToFloat(ctx Context, input Value) (Float, error) {
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
			hasNext, err := iterator.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			val, _, err := iterator.Next(ctx)

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

func ToInt(ctx Context, input Value) (Int, error) {
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
			hasNext, err := iterator.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			item, _, err := iterator.Next(ctx)

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

func ToIntSafe(ctx Context, input Value) Int {
	result, err := ToInt(ctx, input)

	if err != nil {
		return ZeroInt
	}

	if result > 0 {
		return result
	}

	return ZeroInt
}

func ToNumberOnly(ctx Context, input Value) Value {
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

		for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
			val, _, err := iterator.Next(ctx)

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

func ToIntDefault(ctx Context, input Value, defaultValue Int) (Int, error) {
	result, err := ToInt(ctx, input)

	if err != nil {
		return defaultValue, err
	}

	if result > 0 {
		return result, nil
	}

	return defaultValue, nil
}

func ToNumberOrString(input Value) Value {
	switch value := input.(type) {
	case Int, Float, String:
		return value
	default:
		return ToString(value)
	}
}

func ToStrings(ctx Context, input Collection) []String {
	size, err := input.Length(ctx)

	if err != nil {
		size = 0
	}

	res := make([]String, size)

	iterator, err := input.Iterate(ctx)

	if err != nil {
		return res
	}

	var i int

	for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
		val, _, err := iterator.Next(ctx)

		if err != nil {
			return res
		}

		res[i] = NewString(val.String())
		i++
	}

	if closable, ok := iterator.(io.Closer); ok {
		_ = closable.Close()
	}

	return res
}

func ToBinary(input Value) Binary {
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

func UnwrapStrings(values []String) []string {
	out := make([]string, len(values))

	for i, v := range values {
		out[i] = string(v)
	}

	return out
}

func CompareStrings(a, b String) Int {
	return Int(strings.Compare(a.String(), b.String()))
}
