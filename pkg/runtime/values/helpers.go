package values

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"reflect"
	"sort"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func GetIn(ctx context.Context, from core.Value, byPath []core.Value) (core.Value, error) {
	if byPath == nil || len(byPath) == 0 {
		return None, nil
	}

	var result = from

	for i, segment := range byPath {
		if result == None || result == nil {
			break
		}

		segmentType := segment.Type()
		resultType := result.Type()

		if resultType == types.Object {
			obj := result.(*Object)

			if segmentType != types.String {
				return nil, core.TypeError(segmentType, types.String)
			}

			result, _ = obj.Get(segment.(String))
		} else if resultType == types.Array {
			arr := result.(*Array)

			if segmentType != types.Int {
				return nil, core.TypeError(segmentType, types.Int)
			}

			result = arr.Get(segment.(Int))
		} else {
			getter, ok := result.(core.Getter)

			if ok {
				return getter.GetIn(ctx, byPath[i:])
			}

			return None, core.TypeError(
				from.Type(),
				types.Array,
				types.Object,
				core.NewType("Getter"),
			)
		}
	}

	return result, nil
}

func SetIn(ctx context.Context, to core.Value, byPath []core.Value, value core.Value) error {
	if byPath == nil || len(byPath) == 0 {
		return nil
	}

	var parent core.Value
	var current = to
	target := len(byPath) - 1

	for idx, segment := range byPath {
		parent = current
		isTarget := target == idx
		segmentType := segment.Type()
		parentType := parent.Type()

		if parentType == types.Object {
			parent := parent.(*Object)

			if segmentType != types.String {
				return core.TypeError(segmentType, types.String)
			}

			if isTarget == false {
				current, _ = parent.Get(segment.(String))
			} else {
				parent.Set(segment.(String), value)
			}
		} else if parentType == types.Array {
			if segmentType != types.Int {
				return core.TypeError(segmentType, types.Int)
			}

			parent := parent.(*Array)

			if isTarget == false {
				current = parent.Get(segment.(Int))
			} else {
				if err := parent.Set(segment.(Int), value); err != nil {
					return err
				}
			}
		} else {
			setter, ok := parent.(core.Setter)

			if ok {
				return setter.SetIn(ctx, byPath[idx:], value)
			}

			// redefine parent
			isArray := segmentType.Equals(types.Int)

			// it's not an index
			if isArray == false {
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

			if isTarget == false {
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

func ToBoolean(input core.Value) core.Value {
	switch input.Type() {
	case types.Boolean:
		return input
	case types.None:
		return False
	case types.String:
		return NewBoolean(input.String() != "")
	case types.Int:
		return NewBoolean(input.(Int) != 0)
	case types.Float:
		return NewBoolean(input.(Float) != 0)
	default:
		return True
	}
}

func ToArray(ctx context.Context, input core.Value) (core.Value, error) {
	switch input.Type() {
	case types.Boolean,
		types.Int,
		types.Float,
		types.String,
		types.DateTime:

		return NewArrayWith(input), nil
	case types.Array:
		return input.Copy(), nil
	case types.Object:
		obj, ok := input.(*Object)

		if !ok {
			return NewArray(0), nil
		}

		arr := NewArray(int(obj.Length()))

		obj.ForEach(func(value core.Value, key string) bool {
			arr.Push(value)

			return true
		})

		return obj, nil
	default:
		iterable, ok := input.(core.Iterable)

		if !ok {
			return NewArray(0), nil
		}

		iterator, err := iterable.Iterate(ctx)

		if err != nil {
			return None, err
		}

		arr := NewArray(10)

		for {
			val, _, err := iterator.Next(ctx)

			if err != nil {
				return None, err
			}

			if val == None {
				break
			}

			arr.Push(val)
		}

		return arr, nil
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
