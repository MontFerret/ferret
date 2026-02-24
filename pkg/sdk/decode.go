package sdk

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var timeType = reflect.TypeOf(time.Time{})

// Decode binds a runtime Value into the provided target.
// Target must be a non-nil pointer.
// It uses tags (ferret or json) for struct fields and supports unwrapping values that implement the Unwrappable interface.
func Decode(src runtime.Value, target any) error {
	if src == nil {
		src = runtime.None
	}

	targetValue := reflect.ValueOf(target)
	if !targetValue.IsValid() {
		return runtime.Error(runtime.ErrInvalidArgumentType, "target is invalid")
	}

	if targetValue.Kind() != reflect.Pointer {
		return runtime.Error(runtime.ErrInvalidArgumentType, "target must be a pointer")
	}

	if targetValue.IsNil() {
		return runtime.Error(runtime.ErrInvalidArgumentType, "target must be a non-nil pointer")
	}

	return bindValue(src, targetValue.Elem())
}

func bindValue(src runtime.Value, dst reflect.Value) error {
	if !dst.CanSet() {
		return runtime.Error(runtime.ErrInvalidArgumentType, "target is not settable")
	}

	if src == nil {
		src = runtime.None
	}

	if dst.Kind() == reflect.Pointer {
		if src == runtime.None {
			dst.Set(reflect.Zero(dst.Type()))
			return nil
		}

		if dst.IsNil() {
			elem := reflect.New(dst.Type().Elem())
			if err := bindValue(src, elem.Elem()); err != nil {
				return err
			}
			dst.Set(elem)
			return nil
		}

		return bindValue(src, dst.Elem())
	}

	if src == runtime.None {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.IsValid() && srcVal.Type().AssignableTo(dst.Type()) {
		dst.Set(srcVal)
		return nil
	}

	switch dst.Kind() {
	case reflect.Bool:
		val, ok := src.(runtime.Boolean)
		if !ok {
			return bindTypeError(src, dst.Type())
		}
		dst.SetBool(bool(val))
		return nil
	case reflect.String:
		val, ok := src.(runtime.String)
		if !ok {
			return bindTypeError(src, dst.Type())
		}
		dst.SetString(val.String())
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := src.(type) {
		case runtime.Int:
			return setInt(dst, int64(v))
		case runtime.Float:
			return bindTypeError(src, dst.Type())
		default:
			return bindTypeError(src, dst.Type())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v := src.(type) {
		case runtime.Int:
			return setUint(dst, int64(v))
		case runtime.Float:
			return bindTypeError(src, dst.Type())
		default:
			return bindTypeError(src, dst.Type())
		}
	case reflect.Float32, reflect.Float64:
		switch v := src.(type) {
		case runtime.Float:
			dst.SetFloat(float64(v))
			return nil
		case runtime.Int:
			dst.SetFloat(float64(v))
			return nil
		default:
			return bindTypeError(src, dst.Type())
		}
	case reflect.Slice:
		return bindSlice(src, dst)
	case reflect.Array:
		return bindArray(src, dst)
	case reflect.Map:
		return bindMap(src, dst)
	case reflect.Struct:
		if dt, ok := src.(runtime.DateTime); ok && dst.Type() == timeType {
			dst.Set(reflect.ValueOf(dt.Time))
			return nil
		}

		return bindStruct(src, dst)
	case reflect.Interface:
		return bindInterface(src, dst)
	default:
		return bindFallback(src, dst)
	}
}

func bindInterface(src runtime.Value, dst reflect.Value) error {
	srcVal := reflect.ValueOf(src)
	if srcVal.IsValid() && srcVal.Type().AssignableTo(dst.Type()) {
		dst.Set(srcVal)
		return nil
	}

	unwrappable, ok := src.(runtime.Unwrappable)

	if !ok {
		return nil
	}

	unwrapped := unwrappable.Unwrap()
	if unwrapped == nil {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	unwrapVal := reflect.ValueOf(unwrapped)
	if unwrapVal.Type().AssignableTo(dst.Type()) {
		dst.Set(unwrapVal)
		return nil
	}

	return bindTypeError(src, dst.Type())
}

func bindFallback(src runtime.Value, dst reflect.Value) error {
	unwrappable, ok := src.(runtime.Unwrappable)

	if !ok {
		return nil
	}

	unwrapped := unwrappable.Unwrap()
	if unwrapped == nil {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	unwrapVal := reflect.ValueOf(unwrapped)
	if unwrapVal.Type().AssignableTo(dst.Type()) {
		dst.Set(unwrapVal)
		return nil
	}

	if unwrapVal.Type().ConvertibleTo(dst.Type()) {
		dst.Set(unwrapVal.Convert(dst.Type()))
		return nil
	}

	return bindTypeError(src, dst.Type())
}

func bindSlice(src runtime.Value, dst reflect.Value) error {
	iterable, ok := src.(runtime.Iterable)
	if !ok {
		return bindTypeError(src, dst.Type())
	}

	ctx := context.Background()
	iter, err := iterable.Iterate(ctx)
	if err != nil {
		return err
	}

	elemType := dst.Type().Elem()
	out := reflect.MakeSlice(dst.Type(), 0, 0)
	index := 0

	for {
		val, _, err := iter.Next(ctx)
		if errors.Is(err, io.EOF) || errors.Is(err, runtime.ErrTimeout) {
			break
		}
		if err != nil {
			return err
		}

		elem := reflect.New(elemType).Elem()
		if err := bindValue(val, elem); err != nil {
			return err
		}

		out = reflect.Append(out, elem)
		index++
	}

	dst.Set(out)
	return nil
}

func bindArray(src runtime.Value, dst reflect.Value) error {
	iterable, ok := src.(runtime.Iterable)
	if !ok {
		return bindTypeError(src, dst.Type())
	}

	ctx := context.Background()
	iter, err := iterable.Iterate(ctx)
	if err != nil {
		return err
	}

	elemType := dst.Type().Elem()
	index := 0

	for {
		if index >= dst.Len() {
			return runtime.Error(runtime.ErrInvalidArgumentType, "source has more elements than target array")
		}

		val, _, err := iter.Next(ctx)
		if errors.Is(err, io.EOF) || errors.Is(err, runtime.ErrTimeout) {
			break
		}
		if err != nil {
			return err
		}

		elem := reflect.New(elemType).Elem()
		if err := bindValue(val, elem); err != nil {
			return err
		}

		dst.Index(index).Set(elem)
		index++
	}

	return nil
}

func bindMap(src runtime.Value, dst reflect.Value) error {
	if dst.Type().Key().Kind() != reflect.String {
		return runtime.Error(runtime.ErrInvalidArgumentType, "map key type must be string")
	}

	entries, err := collectEntries(src)
	if err != nil {
		return err
	}

	elemType := dst.Type().Elem()
	out := reflect.MakeMap(dst.Type())

	for key, value := range entries {
		elem := reflect.New(elemType).Elem()
		if err := bindValue(value, elem); err != nil {
			return err
		}

		out.SetMapIndex(reflect.ValueOf(key), elem)
	}

	dst.Set(out)
	return nil
}

func bindStruct(src runtime.Value, dst reflect.Value) error {
	entries, err := collectEntries(src)
	if err != nil {
		return err
	}

	lowerKeys := make(map[string]string, len(entries))
	for key := range entries {
		lower := strings.ToLower(key)
		if _, exists := lowerKeys[lower]; !exists {
			lowerKeys[lower] = key
		}
	}

	dstType := dst.Type()
	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name, ok := Tag(field)
		if !ok {
			continue
		}

		value, ok := entries[name]
		if !ok {
			if original, found := lowerKeys[strings.ToLower(name)]; found {
				value = entries[original]
				ok = true
			}
		}

		if !ok {
			continue
		}

		if err := bindValue(value, dst.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func collectEntries(src runtime.Value) (map[string]runtime.Value, error) {
	m, ok := src.(runtime.Map)
	if !ok {
		return nil, bindTypeError(src, reflect.TypeOf(map[string]any{}))
	}

	ctx := context.Background()
	keys, err := m.Keys(ctx)
	if err != nil {
		return nil, err
	}

	iter, err := keys.Iterate(ctx)
	if err != nil {
		return nil, err
	}

	out := make(map[string]runtime.Value)
	for {
		keyVal, _, err := iter.Next(ctx)
		if errors.Is(err, io.EOF) || errors.Is(err, runtime.ErrTimeout) {
			break
		}
		if err != nil {
			return nil, err
		}

		key, ok := keyVal.(runtime.String)
		if !ok {
			return nil, runtime.Error(runtime.ErrInvalidArgumentType, "map key type must be string")
		}

		val, err := m.Get(ctx, keyVal)
		if err != nil {
			return nil, err
		}

		out[key.String()] = val
	}

	return out, nil
}

func setInt(dst reflect.Value, value int64) error {
	switch dst.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.OverflowInt(value) {
			return runtime.Error(runtime.ErrInvalidArgumentType, "integer overflow")
		}
		dst.SetInt(value)
		return nil
	case reflect.Float32, reflect.Float64:
		dst.SetFloat(float64(value))
		return nil
	default:
		return runtime.Error(runtime.ErrInvalidArgumentType, "invalid integer target type")
	}
}

func setUint(dst reflect.Value, value int64) error {
	if value < 0 {
		return runtime.Error(runtime.ErrInvalidArgumentType, "negative value for unsigned target")
	}

	u := uint64(value)
	switch dst.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if dst.OverflowUint(u) {
			return runtime.Error(runtime.ErrInvalidArgumentType, "unsigned integer overflow")
		}
		dst.SetUint(u)
		return nil
	default:
		return runtime.Error(runtime.ErrInvalidArgumentType, "invalid unsigned target type")
	}
}

func bindTypeError(src runtime.Value, target reflect.Type) error {
	return runtime.Errorf(runtime.ErrInvalidArgumentType, "cannot bind %s to %s", runtime.TypeOf(src), target.String())
}
