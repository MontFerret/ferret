package internal

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"go/types"
	"hash/fnv"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wI2L/jettison"
)

func Parse(input interface{}) core.Value {
	switch value := input.(type) {
	case bool:
		return core.NewBoolean(value)
	case string:
		return core.NewString(value)
	case int64:
		return core.NewInt(int(value))
	case int32:
		return core.NewInt(int(value))
	case int16:
		return core.NewInt(int(value))
	case int8:
		return core.NewInt(int(value))
	case int:
		return core.NewInt(value)
	case float64:
		return core.NewFloat(value)
	case float32:
		return core.NewFloat(float64(value))
	case time.Time:
		return core.NewDateTime(value)
	case []interface{}:
		ctx := context.Background()
		arr := NewArray(len(value))

		for _, el := range value {
			_ = arr.Add(ctx, Parse(el))
		}

		return arr
	case map[string]interface{}:
		ctx := context.Background()
		obj := NewObject()

		for key, el := range value {
			_ = obj.Set(ctx, core.NewString(key), Parse(el))
		}

		return obj
	case []byte:
		return core.NewBinary(value)
	case nil:
		return core.None
	default:
		v := reflect.ValueOf(value)
		t := reflect.TypeOf(value)
		kind := t.Kind()

		if kind == reflect.Ptr {
			el := v.Elem()

			if el.Kind() == 0 {
				return core.None
			}

			return Parse(el.Interface())
		}

		if kind == reflect.Slice || kind == reflect.Array {
			size := v.Len()
			arr := NewArray(size)
			ctx := context.Background()

			for i := 0; i < size; i++ {
				curVal := v.Index(i)
				_ = arr.Add(ctx, Parse(curVal.Interface()))
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

				_ = obj.Set(ctx, core.NewString(key.String()), Parse(curVal.Interface()))
			}

			return obj
		}

		if kind == reflect.Struct {
			obj := NewObject()
			size := t.NumField()
			ctx := context.Background()

			for i := 0; i < size; i++ {
				field := t.Field(i)
				fieldValue := v.Field(i)

				_ = obj.Set(ctx, core.NewString(field.Name), Parse(fieldValue.Interface()))
			}

			return obj
		}

		return core.None
	}
}

func Unmarshal(value json.RawMessage) (core.Value, error) {
	var o interface{}

	err := json.Unmarshal(value, &o)

	if err != nil {
		return core.None, err
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

func IsScalar(input core.Value) core.Boolean {
	switch input.(type) {
	case core.Int, core.Float, core.String, core.Boolean:
		return true
	default:
		return false
	}
}

func IsNumber(input core.Value) core.Boolean {
	switch input.(type) {
	case core.Int, core.Float:
		return true
	default:
		return false
	}
}

func ToBoolean(input core.Value) core.Boolean {
	if input == core.None {
		return core.False
	}

	switch val := input.(type) {
	case core.Boolean:
		return val
	case core.String:
		return val != ""
	case core.Int:
		return val != 0
	case core.Float:
		return val != 0
	case core.DateTime:
		return val.IsZero() == false
	default:
		return core.True
	}
}

func ToFloat(input core.Value) core.Float {
	switch val := input.(type) {
	case core.Float:
		return val
	case core.Int:
		return core.Float(val)
	case core.String:
		i, err := strconv.ParseFloat(string(val), 64)

		if err != nil {
			return core.ZeroFloat
		}

		return core.Float(i)
	case core.Boolean:
		if val {
			return core.Float(1)
		}

		return core.Float(0)
	case core.DateTime:
		dt := input.(core.DateTime)

		if dt.IsZero() {
			return core.ZeroFloat
		}

		return core.NewFloat(float64(dt.Unix()))
	case *Array:
		length := val.Length()

		if length == 0 {
			return core.ZeroFloat
		}

		res := core.ZeroFloat

		for i := 0; i < length; i++ {
			res += ToFloat(val.Get(i))
		}

		return res
	default:
		return core.ZeroFloat
	}
}

func ToString(input core.Value) core.String {
	switch val := input.(type) {
	case core.String:
		return val
	default:
		return core.NewString(val.String())
	}
}

func ToInt(input core.Value) core.Int {
	switch val := input.(type) {
	case core.Int:
		return val
	case core.Float:
		return core.Int(val)
	case core.String:
		i, err := strconv.ParseInt(string(val), 10, 64)

		if err != nil {
			return core.ZeroInt
		}

		return core.Int(i)
	case core.Boolean:
		if val {
			return core.Int(1)
		}

		return core.Int(0)
	case core.DateTime:
		dt := input.(core.DateTime)

		if dt.IsZero() {
			return core.ZeroInt
		}

		return core.NewInt(int(dt.Unix()))
	case *Array:
		length := val.Length()

		if length == 0 {
			return core.ZeroInt
		}

		res := core.ZeroInt

		for i := 0; i < length; i++ {
			res += ToInt(val.Get(i))
		}

		return res
	default:
		return core.ZeroInt
	}
}

func ToIntDefault(input core.Value, defaultValue core.Int) core.Int {
	if result := ToInt(input); result > 0 {
		return result
	}

	return defaultValue
}

func ToArray(ctx context.Context, input core.Value) *Array {
	switch value := input.(type) {
	case core.Boolean,
		core.Int,
		core.Float,
		core.String,
		core.DateTime:

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

		for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
			val, _, err := iterator.Next(ctx)

			if err != nil {
				return arr
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
			obj.Set(ToString(core.Int(idx)), value)

			return true
		})

		return obj
	case core.Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return NewObject()
		}

		obj := NewObject()

		for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
			val, key, err := iterator.Next(ctx)

			if err != nil {
				return obj
			}

			obj.Set(ToString(key), val)
		}

		return obj
	default:
		return NewObject()
	}
}

func ToStrings(input *Array) []core.String {
	res := make([]core.String, input.Length())

	input.ForEach(func(v core.Value, i int) bool {
		res[i] = core.NewString(v.String())

		return true
	})

	return res
}

func ToBinary(input core.Value) core.Binary {
	bin, ok := input.(core.Binary)

	if ok {
		return bin
	}

	return core.NewBinary([]byte(input.String()))
}

func ToRegexp(input core.Value) (*Regexp, error) {
	switch r := input.(type) {
	case *Regexp:
		return r, nil
	case core.String:
		return NewRegexp(r)
	default:
		return nil, core.TypeError(input, types.String, types.Regexp)
	}
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

func UnwrapStrings(values []core.String) []string {
	out := make([]string, len(values))

	for i, v := range values {
		out[i] = v.String()
	}

	return out
}

func Negate(input core.Value) core.Value {
	switch val := input.(type) {
	case core.Int:
		return -val
	case core.Float:
		return -val
	case core.Boolean:
		return !val
	default:
		return core.None
	}
}

func Negative(input core.Value) core.Value {
	switch value := input.(type) {
	case core.Int:
		return -value
	case core.Float:
		return -value
	default:
		// TODO: Maybe we should return None?
		return input
	}
}

func Positive(input core.Value) core.Value {
	switch value := input.(type) {
	case core.Int:
		return +value
	case core.Float:
		return +value
	default:
		// TODO: Maybe we should return None?
		return input
	}
}

func Contains(input core.Value, value core.Value) core.Boolean {
	switch val := input.(type) {
	case *Array:
		return val.Contains(value)
	case core.String:
		return core.Boolean(strings.Contains(val.String(), value.String()))
	default:
		return false
	}
}

func ToNumberOrString(input core.Value) core.Value {
	switch value := input.(type) {
	case core.Int, core.Float, core.String:
		return value
	default:
		return ToString(value)
	}
}

func ToNumberOnly(input core.Value) core.Value {
	switch value := input.(type) {
	case core.Int, core.Float:
		return value
	case core.String:
		if strings.Contains(value.String(), ".") {
			return ToFloat(input)
		}

		return ToInt(input)
	case *Array:
		length := value.Length()

		if length == 0 {
			return core.ZeroInt
		}

		i := core.ZeroInt
		f := core.ZeroFloat

		for y := 0; y < length; y++ {
			out := ToNumberOnly(value.Get(y))

			switch item := out.(type) {
			case core.Int:
				i += item
			case core.Float:
				f += item
			}
		}

		if f == 0 {
			return i
		}

		return core.Float(i) + f
	default:
		return ToInt(input)
	}
}

func CompareStrings(a, b core.String) core.Int {
	return core.Int(strings.Compare(a.String(), b.String()))
}

func Length(value core.Value) (core.Int, error) {
	c, ok := value.(core.Measurable)

	if !ok {
		return 0, core.TypeError(value,
			types.String,
			types.Array,
			types.Object,
			types.Binary,
			types.Measurable,
		)
	}

	return core.Int(c.Length()), nil
}
