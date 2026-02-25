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
// It uses tags (ferret or json) for struct fields.
func Decode(src runtime.Value, target any) error {
	src = normalizeSrc(src)

	targetValue, err := validateTarget(target)
	if err != nil {
		return err
	}

	return bindValue(src, targetValue.Elem())
}

func bindValue(src runtime.Value, dst reflect.Value) error {
	if !dst.CanSet() {
		return runtime.Error(runtime.ErrInvalidArgumentType, "target is not settable")
	}

	src = normalizeSrc(src)

	if dst.Kind() == reflect.Pointer {
		return bindPointer(src, dst)
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

	if handled, err := bindScalar(src, dst); handled {
		return err
	}

	return bindComposite(src, dst)
}

func bindPointer(src runtime.Value, dst reflect.Value) error {
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

func bindScalar(src runtime.Value, dst reflect.Value) (bool, error) {
	switch dst.Kind() {
	case reflect.Bool:
		val, ok := src.(runtime.Boolean)

		if !ok {
			return true, bindTypeError(src, dst.Type())
		}

		dst.SetBool(bool(val))

		return true, nil
	case reflect.String:
		val, ok := src.(runtime.String)

		if !ok {
			return true, bindTypeError(src, dst.Type())
		}

		dst.SetString(val.String())

		return true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := src.(type) {
		case runtime.Int:
			return true, setInt(dst, int64(v))
		case runtime.Float:
			return true, bindTypeError(src, dst.Type())
		default:
			return true, bindTypeError(src, dst.Type())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v := src.(type) {
		case runtime.Int:
			return true, setUint(dst, int64(v))
		case runtime.Float:
			return true, bindTypeError(src, dst.Type())
		default:
			return true, bindTypeError(src, dst.Type())
		}
	case reflect.Float32, reflect.Float64:
		switch v := src.(type) {
		case runtime.Float:
			dst.SetFloat(float64(v))
			return true, nil
		case runtime.Int:
			dst.SetFloat(float64(v))
			return true, nil
		default:
			return true, bindTypeError(src, dst.Type())
		}
	default:
		return false, nil
	}
}

func bindComposite(src runtime.Value, dst reflect.Value) error {
	switch dst.Kind() {
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

	return bindUnwrapped(src, dst, false)
}

func bindFallback(src runtime.Value, dst reflect.Value) error {
	return bindUnwrapped(src, dst, true)
}

func bindUnwrapped(src runtime.Value, dst reflect.Value, allowConvert bool) error {
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

	if allowConvert && unwrapVal.Type().ConvertibleTo(dst.Type()) {
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

	lowerKeys := buildLowerKeyMap(entries)
	used := make(map[string]struct{}, len(entries))

	_, err = bindStructEntries(dst, entries, lowerKeys, used)
	return err
}

func bindStructEntries(dst reflect.Value, entries map[string]runtime.Value, lowerKeys map[string]string, used map[string]struct{}) (bool, error) {
	dstType := dst.Type()
	matched := false

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)

		if field.PkgPath != "" {
			continue
		}

		name, ok := Tag(field)
		if !ok {
			continue
		}

		value, key, ok := lookupEntry(name, entries, lowerKeys, used)
		if !ok {
			continue
		}

		if err := bindValue(value, dst.Field(i)); err != nil {
			return false, err
		}

		used[key] = struct{}{}
		matched = true
	}

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		if field.PkgPath != "" || !field.Anonymous {
			continue
		}

		if _, ok := Tag(field); ok {
			continue
		}

		fieldVal := dst.Field(i)
		fieldType := fieldVal.Type()

		switch fieldType.Kind() {
		case reflect.Struct:
			subMatched, err := bindStructEntries(fieldVal, entries, lowerKeys, used)
			if err != nil {
				return false, err
			}
			if subMatched {
				matched = true
			}
		case reflect.Pointer:
			if fieldType.Elem().Kind() != reflect.Struct {
				continue
			}

			elem := reflect.New(fieldType.Elem())
			subMatched, err := bindStructEntries(elem.Elem(), entries, lowerKeys, used)

			if err != nil {
				return false, err
			}

			if subMatched {
				fieldVal.Set(elem)
				matched = true
			}
		default:
			continue
		}
	}

	return matched, nil
}

func buildLowerKeyMap(entries map[string]runtime.Value) map[string]string {
	lowerKeys := make(map[string]string, len(entries))
	for key := range entries {
		lower := strings.ToLower(key)

		if _, exists := lowerKeys[lower]; !exists {
			lowerKeys[lower] = key
		}
	}

	return lowerKeys
}

func lookupEntry(name string, entries map[string]runtime.Value, lowerKeys map[string]string, used map[string]struct{}) (runtime.Value, string, bool) {
	if _, taken := used[name]; !taken {
		if value, ok := entries[name]; ok {
			return value, name, true
		}
	}

	lower := strings.ToLower(name)
	if original, ok := lowerKeys[lower]; ok {
		if _, taken := used[original]; !taken {
			return entries[original], original, true
		}
	}

	return nil, "", false
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

func normalizeSrc(src runtime.Value) runtime.Value {
	if src == nil {
		return runtime.None
	}

	return src
}

func validateTarget(target any) (reflect.Value, error) {
	targetValue := reflect.ValueOf(target)
	if !targetValue.IsValid() {
		return reflect.Value{}, runtime.Error(runtime.ErrInvalidArgumentType, "target is invalid")
	}

	if targetValue.Kind() != reflect.Pointer {
		return reflect.Value{}, runtime.Error(runtime.ErrInvalidArgumentType, "target must be a pointer")
	}

	if targetValue.IsNil() {
		return reflect.Value{}, runtime.Error(runtime.ErrInvalidArgumentType, "target must be a non-nil pointer")
	}

	return targetValue, nil
}
