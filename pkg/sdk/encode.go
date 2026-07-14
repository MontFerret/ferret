package sdk

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const maxInt64 = ^uint64(0) >> 1

var (
	byteSliceType = reflect.TypeOf([]byte(nil))
	timeType      = reflect.TypeOf(time.Time{})
)

type (
	encodeState struct {
		ctx      context.Context
		visiting map[encodeVisit]struct{}
		path     conversionPath
	}

	encodeVisit struct {
		typ reflect.Type
		ptr uintptr
	}
)

// Encode converts a Go value into a Ferret runtime value.
// Unsupported values, cycles, overflows, and canceled contexts are reported as errors.
func Encode(ctx context.Context, input any) (runtime.Value, error) {
	if ctx == nil {
		return runtime.None, runtime.Error(runtime.ErrInvalidArgument, "context cannot be nil")
	}

	if input == nil {
		return runtime.None, nil
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, fmt.Errorf("$: %w", err)
	}

	if value, ok := input.(runtime.Value); ok {
		return normalizeRuntimeValue(value), nil
	}

	reflected := reflect.ValueOf(input)

	if special, ok := encodeSpecialValue(reflected); ok {
		return special, nil
	}

	if scalar, ok, err := encodeScalarValue(reflected); ok {
		if err != nil {
			return runtime.None, fmt.Errorf("$: %w", err)
		}

		return scalar, nil
	}

	state := &encodeState{
		ctx:      ctx,
		path:     newConversionPath(),
		visiting: make(map[encodeVisit]struct{}),
	}

	return encodeReflectValue(state, reflected)
}

func encodeReflectValue(state *encodeState, value reflect.Value) (runtime.Value, error) {
	if err := state.ctx.Err(); err != nil {
		return runtime.None, fmt.Errorf("%s: %w", state.path, err)
	}

	if !value.IsValid() {
		return runtime.None, nil
	}

	if value.CanInterface() {
		if runtimeValue, ok := value.Interface().(runtime.Value); ok {
			return normalizeRuntimeValue(runtimeValue), nil
		}
	}

	visits := make([]encodeVisit, 0, 2)
	defer func() {
		for i := len(visits) - 1; i >= 0; i-- {
			delete(state.visiting, visits[i])
		}
	}()

	for {
		switch value.Kind() {
		case reflect.Interface:
			if value.IsNil() {
				return runtime.None, nil
			}

			value = value.Elem()
		case reflect.Pointer:
			if value.IsNil() {
				return runtime.None, nil
			}

			visit := encodeVisit{typ: value.Type(), ptr: value.Pointer()}
			if _, exists := state.visiting[visit]; exists {
				return runtime.None, encodeCycleError(&state.path)
			}

			state.visiting[visit] = struct{}{}
			visits = append(visits, visit)
			value = value.Elem()
		default:
			goto resolved
		}

		if value.CanInterface() {
			if runtimeValue, ok := value.Interface().(runtime.Value); ok {
				return normalizeRuntimeValue(runtimeValue), nil
			}
		}
	}

resolved:
	if special, ok := encodeSpecialValue(value); ok {
		return special, nil
	}

	if scalar, ok, err := encodeScalarValue(value); ok {
		if err != nil {
			return runtime.None, fmt.Errorf("%s: %w", state.path, err)
		}

		return scalar, nil
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		if value.Kind() == reflect.Slice && !value.IsNil() {
			visit := encodeVisit{typ: value.Type(), ptr: value.Pointer()}
			if _, exists := state.visiting[visit]; exists {
				return runtime.None, encodeCycleError(&state.path)
			}

			state.visiting[visit] = struct{}{}
			defer delete(state.visiting, visit)
		}

		return encodeArrayValue(state, value)
	case reflect.Map:
		if value.IsNil() {
			return runtime.None, nil
		}

		visit := encodeVisit{typ: value.Type(), ptr: value.Pointer()}
		if _, exists := state.visiting[visit]; exists {
			return runtime.None, encodeCycleError(&state.path)
		}

		state.visiting[visit] = struct{}{}
		defer delete(state.visiting, visit)

		return encodeMapValue(state, value)
	case reflect.Struct:
		entries, err := encodeStructEntries(state, value)
		if err != nil {
			return runtime.None, err
		}

		return runtime.NewObjectWith(entries), nil
	default:
		return runtime.None, fmt.Errorf(
			"%s: %w",
			state.path,
			runtime.Errorf(runtime.ErrInvalidArgumentType, "cannot encode %s", value.Type()),
		)
	}
}

func encodeSpecialValue(value reflect.Value) (runtime.Value, bool) {
	if value.Type() == timeType && value.CanInterface() {
		return runtime.NewDateTime(value.Interface().(time.Time)), true
	}

	if value.Type() == byteSliceType && value.Kind() == reflect.Slice {
		return runtime.NewBinary(value.Bytes()), true
	}

	return runtime.None, false
}

func encodeScalarValue(value reflect.Value) (runtime.Value, bool, error) {
	switch value.Kind() {
	case reflect.Bool:
		return runtime.NewBoolean(value.Bool()), true, nil
	case reflect.String:
		return runtime.NewString(value.String()), true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return runtime.NewInt64(value.Int()), true, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		unsigned := value.Uint()
		if unsigned > maxInt64 {
			return runtime.None, true, runtime.Error(runtime.ErrInvalidArgumentType, "unsigned integer overflows Ferret Int")
		}

		return runtime.NewInt64(int64(unsigned)), true, nil
	case reflect.Float32, reflect.Float64:
		return runtime.NewFloat(value.Float()), true, nil
	default:
		return runtime.None, false, nil
	}
}

func encodeArrayValue(state *encodeState, value reflect.Value) (runtime.Value, error) {
	items := make([]runtime.Value, value.Len())

	for i := 0; i < value.Len(); i++ {
		mark := state.path.PushIndex(i)
		item, err := encodeReflectValue(state, value.Index(i))

		state.path.Restore(mark)

		if err != nil {
			return runtime.None, err
		}

		items[i] = item
	}

	return runtime.NewArrayWith(items...), nil
}

func encodeMapValue(state *encodeState, value reflect.Value) (runtime.Value, error) {
	if value.Type().Key().Kind() != reflect.String {
		return runtime.None, fmt.Errorf(
			"%s: %w",
			state.path,
			runtime.Errorf(runtime.ErrInvalidArgumentType, "map key type must be string, got %s", value.Type().Key()),
		)
	}

	entries := make(map[string]runtime.Value, value.Len())
	iterator := value.MapRange()

	for iterator.Next() {
		key := iterator.Key().String()
		mark := state.path.PushKey(key)
		encoded, err := encodeReflectValue(state, iterator.Value())

		state.path.Restore(mark)

		if err != nil {
			return runtime.None, err
		}

		entries[key] = encoded
	}

	return runtime.NewObjectWith(entries), nil
}

func encodeStructEntries(state *encodeState, value reflect.Value) (map[string]runtime.Value, error) {
	typ := value.Type()
	entries := make(map[string]runtime.Value, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name, tagged := Tag(field)
		if !tagged {
			continue
		}

		mark := state.path.PushField(name)
		encoded, err := encodeReflectValue(state, value.Field(i))
		state.path.Restore(mark)

		if err != nil {
			return nil, err
		}

		entries[name] = encoded
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.PkgPath != "" || !field.Anonymous {
			continue
		}

		if _, tagged := Tag(field); tagged {
			continue
		}

		mark := state.path.PushField(field.Name)
		embedded, ok, err := encodeEmbeddedStructEntries(state, value.Field(i))
		state.path.Restore(mark)

		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		for name, encoded := range embedded {
			if _, exists := entries[name]; !exists {
				entries[name] = encoded
			}
		}
	}

	return entries, nil
}

func encodeEmbeddedStructEntries(state *encodeState, value reflect.Value) (map[string]runtime.Value, bool, error) {
	visits := make([]encodeVisit, 0, 1)

	defer func() {
		for i := len(visits) - 1; i >= 0; i-- {
			delete(state.visiting, visits[i])
		}
	}()

	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, false, nil
		}

		if value.Kind() == reflect.Pointer {
			visit := encodeVisit{typ: value.Type(), ptr: value.Pointer()}
			if _, exists := state.visiting[visit]; exists {
				return nil, false, encodeCycleError(&state.path)
			}

			state.visiting[visit] = struct{}{}
			visits = append(visits, visit)
		}

		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, false, nil
	}

	entries, err := encodeStructEntries(state, value)

	return entries, true, err
}

func encodeCycleError(path *conversionPath) error {
	return fmt.Errorf(
		"%s: %w",
		path,
		runtime.Error(runtime.ErrInvalidArgumentType, "cycle detected; value is already being encoded"),
	)
}

func normalizeRuntimeValue(value runtime.Value) runtime.Value {
	if value == nil {
		return runtime.None
	}

	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if reflected.IsNil() {
			return runtime.None
		}
	}

	return value
}
