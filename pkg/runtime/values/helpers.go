package values

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// GetIn checks that from implements core.Getter interface. If it implements,
// GetIn call from.GetIn method, otherwise return error.
func GetIn(ctx context.Context, from core.Value, byPath []core.Value) (core.Value, error) {
	getter, ok := from.(core.Getter)
	if !ok {
		return None, core.TypeError(
			from.Type(),
			core.NewType("Getter"),
		)
	}

	return getter.GetIn(ctx, byPath)
}

func SetIn(ctx context.Context, to core.Value, byPath []core.Value, value core.Value) error {
	if len(byPath) == 0 {
		return nil
	}

	var parent core.Value
	var current = to
	target := len(byPath) - 1

	for idx, segment := range byPath {
		parent = current
		isTarget := target == idx
		segmentType := segment.Type()

		switch parVal := parent.(type) {
		case *Object:
			if segmentType != types.String {
				return core.TypeError(segmentType, types.String)
			}

			if !isTarget {
				current, _ = parVal.Get(segment.(String))
			} else {
				parVal.Set(segment.(String), value)
			}
		case *Array:
			if segmentType != types.Int {
				return core.TypeError(segmentType, types.Int)
			}

			if !isTarget {
				current = parVal.Get(segment.(Int))
			} else {
				if err := parVal.Set(segment.(Int), value); err != nil {
					return err
				}
			}
		case core.Setter:
			return parVal.SetIn(ctx, byPath[idx:], value)
		default:
			// redefine parent
			isArray := segmentType.Equals(types.Int)

			// it's not an index
			if !isArray {
				obj := NewObject()
				parent = obj

				if segmentType != types.String {
					return core.TypeError(segmentType, types.String)
				}

				if isTarget {
					obj.Set(segment.(String), value)
				}
			} else {
				arr := NewArray(10)
				parent = arr

				if isTarget {
					if err := arr.Set(segment.(Int), value); err != nil {
						return err
					}
				}
			}

			// set new parent
			if err := SetIn(ctx, to, byPath[0:idx-1], parent); err != nil {
				return err
			}

			if !isTarget {
				current = None
			}
		}
	}

	return nil
}

func Parse(input interface{}) core.Value {
	switch value := input.(type) {
	case bool:
		return NewBoolean(value)
	case string:
		return NewString(value)
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
			return Parse(v.Elem().Interface())
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

func ToBoolean(input core.Value) Boolean {
	switch input.Type() {
	case types.Boolean:
		return input.(Boolean)
	case types.String:
		return NewBoolean(input.(String) != "")
	case types.Int:
		return NewBoolean(input.(Int) != 0)
	case types.Float:
		return NewBoolean(input.(Float) != 0)
	case types.DateTime:
		return NewBoolean(!input.(DateTime).IsZero())
	case types.None:
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
		return NewArray(0)
	}
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

func IsNumber(input core.Value) Boolean {
	t := input.Type()

	return t == types.Int || t == types.Float
}
