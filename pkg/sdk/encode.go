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

type encodeState struct {
	embedVisiting map[reflect.Type]int
	ptrVisiting   map[uintptr]int
}

func newEncodeState() *encodeState {
	return &encodeState{
		embedVisiting: make(map[reflect.Type]int),
		ptrVisiting:   make(map[uintptr]int),
	}
}

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

	return encodeValueWithState(reflect.ValueOf(input), newEncodeState())
}

func encodeValueWithState(v reflect.Value, state *encodeState) runtime.Value {
	if !v.IsValid() {
		return runtime.None
	}

	visitedPtrs := make([]uintptr, 0, 2)

	defer func() {
		for i := len(visitedPtrs) - 1; i >= 0; i-- {
			ptr := visitedPtrs[i]
			state.ptrVisiting[ptr]--

			if state.ptrVisiting[ptr] == 0 {
				delete(state.ptrVisiting, ptr)
			}
		}
	}()

	for {
		switch v.Kind() {
		case reflect.Interface:
			if v.IsNil() {
				return runtime.None
			}

			v = v.Elem()
		case reflect.Pointer:
			if v.IsNil() {
				return runtime.None
			}

			ptr := v.Pointer()
			if state.ptrVisiting[ptr] > 0 {
				return runtime.None
			}

			state.ptrVisiting[ptr]++
			visitedPtrs = append(visitedPtrs, ptr)
			v = v.Elem()
		default:
			goto resolved
		}
	}

resolved:
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

	return encodeComposite(v, state)
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

func encodeComposite(v reflect.Value, state *encodeState) runtime.Value {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return encodeArrayLike(v, state)
	case reflect.Map:
		return encodeMap(v, state)
	case reflect.Struct:
		return encodeStruct(v, state)
	default:
		return runtime.None
	}
}

func encodeArrayLike(v reflect.Value, state *encodeState) runtime.Value {
	size := v.Len()
	arr := runtime.NewArray(size)
	ctx := context.Background()

	for i := 0; i < size; i++ {
		_ = arr.Append(ctx, encodeValueWithState(v.Index(i), state))
	}

	return arr
}

func encodeMap(v reflect.Value, state *encodeState) runtime.Value {
	obj := runtime.NewObject()
	ctx := context.Background()

	for _, key := range v.MapKeys() {
		keyVal := encodeValueWithState(key, state)
		_ = obj.Set(ctx, runtime.NewString(keyVal.String()), encodeValueWithState(v.MapIndex(key), state))
	}

	return obj
}

func encodeStruct(v reflect.Value, state *encodeState) runtime.Value {
	return encodeStructWithVisit(v, state)
}

func encodeStructWithVisit(v reflect.Value, state *encodeState) runtime.Value {
	obj := runtime.NewObject()
	ctx := context.Background()
	t := v.Type()
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

		_ = obj.Set(ctx, runtime.NewString(name), encodeValueWithState(v.Field(i), state))
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
		embedded, ok := resolveEmbeddedStruct(fieldVal)
		if !ok {
			continue
		}

		embedType := embedded.Type()
		if !enterTypeVisit(state.embedVisiting, embedType) {
			continue
		}

		mergeObject(obj, encodeStructWithVisit(embedded, state), used)
		leaveTypeVisit(state.embedVisiting, embedType)
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

	defer func() {
		_ = closeIter(iter)
	}()

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

func resolveEmbeddedStruct(fieldVal reflect.Value) (reflect.Value, bool) {
	switch fieldVal.Kind() {
	case reflect.Struct:
		return fieldVal, true
	case reflect.Pointer:
		if fieldVal.IsNil() || fieldVal.Elem().Kind() != reflect.Struct {
			return reflect.Value{}, false
		}

		return fieldVal.Elem(), true
	default:
		return reflect.Value{}, false
	}
}

func enterTypeVisit(visiting map[reflect.Type]int, t reflect.Type) bool {
	if visiting[t] > 0 {
		return false
	}

	visiting[t]++
	return true
}

func leaveTypeVisit(visiting map[reflect.Type]int, t reflect.Type) {
	visiting[t]--

	if visiting[t] == 0 {
		delete(visiting, t)
	}
}
