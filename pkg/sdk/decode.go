package sdk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var runtimeValueReflectType = reflect.TypeFor[runtime.Value]()

type decodeState struct {
	ctx    context.Context
	config *decodeConfig
	path   conversionPath
}

// Decode binds a Ferret runtime value into a non-nil pointer target.
func Decode(ctx context.Context, src runtime.Value, target any, options ...DecodeOption) error {
	if ctx == nil {
		return newDecodeErrorWithoutPath(
			"$",
			DecodeErrorKindType,
			runtime.Error(runtime.ErrInvalidArgument, "context cannot be nil"),
			true,
		)
	}

	targetValue, err := validateDecodeTarget(target)
	if err != nil {
		return newDecodeErrorWithoutPath("$", DecodeErrorKindType, err, true)
	}

	state := &decodeState{
		ctx:  ctx,
		path: newConversionPath(),
	}

	if len(options) > 0 {
		state.config = &decodeConfig{}

		for _, option := range options {
			if option != nil {
				option(state.config)
			}
		}
	}

	if err := validateDecodeConfig(state.config, targetValue.Elem().Type()); err != nil {
		return newDecodeError("$", DecodeErrorKindType, err, true)
	}

	src = normalizeRuntimeValue(src)

	if src == runtime.None &&
		state.config != nil &&
		state.config.disallowNoneValues &&
		targetValue.Elem().Type() != runtimeValueReflectType {
		return newDecodeError(
			"$",
			DecodeErrorKindNone,
			runtime.Error(runtime.ErrInvalidArgument, "none is not allowed"),
			true,
		)
	}

	if state.config != nil && len(state.config.requiredTypes) > 0 {
		if err := runtime.ValidateType(src, state.config.requiredTypes...); err != nil {
			return newDecodeError("$", DecodeErrorKindType, err, true)
		}
	}

	return bindRuntimeValue(state, src, targetValue.Elem())
}

func bindRuntimeValue(state *decodeState, src runtime.Value, dst reflect.Value) error {
	if err := state.ctx.Err(); err != nil {
		return newDecodeError(state.path.String(), DecodeErrorKindSource, err, false)
	}

	if !dst.CanSet() {
		return newDecodeError(
			state.path.String(),
			DecodeErrorKindType,
			runtime.Error(runtime.ErrInvalidArgumentType, "target is not settable"),
			true,
		)
	}

	if src == runtime.None {
		if dst.Type() == runtimeValueReflectType {
			dst.Set(reflect.ValueOf(src))

			return nil
		}

		if state.config != nil && state.config.disallowNoneValues {
			return newDecodeError(
				state.path.String(),
				DecodeErrorKindNone,
				runtime.Error(runtime.ErrInvalidArgument, "none is not allowed"),
				true,
			)
		}

		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	srcValue := reflect.ValueOf(src)
	if srcValue.IsValid() && srcValue.Type().AssignableTo(dst.Type()) {
		dst.Set(srcValue)
		return nil
	}

	if bindRuntimeUnwrappedExact(src, dst) {
		return nil
	}

	if dst.Kind() == reflect.Pointer {
		return bindRuntimePointer(state, src, dst)
	}

	if handled, err := bindRuntimeScalar(state, src, dst); handled {
		return err
	}

	switch dst.Kind() {
	case reflect.Slice:
		return bindRuntimeSlice(state, src, dst)
	case reflect.Array:
		return bindRuntimeArray(state, src, dst)
	case reflect.Map:
		return bindRuntimeMap(state, src, dst)
	case reflect.Struct:
		if dateTime, ok := src.(runtime.DateTime); ok && dst.Type() == timeType {
			dst.Set(reflect.ValueOf(dateTime.Time))
			return nil
		}

		return bindRuntimeStruct(state, src, dst)
	case reflect.Interface:
		return bindRuntimeInterface(state, src, dst)
	default:
		return bindRuntimeUnwrapped(state, src, dst, true)
	}
}

func bindRuntimePointer(state *decodeState, src runtime.Value, dst reflect.Value) error {
	if src == runtime.None {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	if dst.IsNil() {
		element := reflect.New(dst.Type().Elem())
		if err := bindRuntimeValue(state, src, element.Elem()); err != nil {
			return err
		}

		dst.Set(element)
		return nil
	}

	return bindRuntimeValue(state, src, dst.Elem())
}

func bindRuntimeScalar(state *decodeState, src runtime.Value, dst reflect.Value) (bool, error) {
	switch dst.Kind() {
	case reflect.Bool:
		value, ok := src.(runtime.Boolean)
		if !ok {
			return true, decodeTypeError(&state.path, src, dst.Type())
		}

		dst.SetBool(bool(value))

		return true, nil
	case reflect.String:
		value, ok := src.(runtime.String)
		if !ok {
			return true, decodeTypeError(&state.path, src, dst.Type())
		}

		dst.SetString(value.String())

		return true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, ok := src.(runtime.Int)
		if !ok {
			return true, decodeTypeError(&state.path, src, dst.Type())
		}

		if dst.OverflowInt(int64(value)) {
			return true, newDecodeError(
				state.path.String(),
				DecodeErrorKindRange,
				runtime.Error(runtime.ErrInvalidArgumentType, "integer overflow"),
				true,
			)
		}

		dst.SetInt(int64(value))

		return true, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value, ok := src.(runtime.Int)
		if !ok {
			return true, decodeTypeError(&state.path, src, dst.Type())
		}

		if value < 0 || dst.OverflowUint(uint64(value)) {
			return true, newDecodeError(
				state.path.String(),
				DecodeErrorKindRange,
				runtime.Error(runtime.ErrInvalidArgumentType, "unsigned integer overflow"),
				true,
			)
		}

		dst.SetUint(uint64(value))

		return true, nil
	case reflect.Float32, reflect.Float64:
		var number float64

		switch value := src.(type) {
		case runtime.Float:
			number = float64(value)
		case runtime.Int:
			number = float64(value)
		default:
			return true, decodeTypeError(&state.path, src, dst.Type())
		}

		if dst.OverflowFloat(number) {
			return true, newDecodeError(
				state.path.String(),
				DecodeErrorKindRange,
				runtime.Error(runtime.ErrInvalidArgumentType, "floating-point overflow"),
				true,
			)
		}

		dst.SetFloat(number)

		return true, nil
	default:
		return false, nil
	}
}

func bindRuntimeSlice(state *decodeState, src runtime.Value, dst reflect.Value) (retErr error) {
	iterable, ok := src.(runtime.Iterable)
	if !ok {
		return decodeTypeError(&state.path, src, dst.Type())
	}

	iterator, err := iterable.Iterate(state.ctx)
	if err != nil {
		return newDecodeError(state.path.String(), DecodeErrorKindSource, err, false)
	}

	defer func() {
		retErr = joinDecodeErrors(retErr, closeIteratorAt(iterator, &state.path), &state.path)
	}()

	capacity := 0
	if measurable, ok := src.(runtime.Measurable); ok {
		length, lengthErr := measurable.Length(state.ctx)
		if lengthErr != nil {
			return newDecodeError(state.path.String(), DecodeErrorKindSource, lengthErr, false)
		}

		if length > runtime.Int(int(^uint(0)>>1)) {
			return newDecodeError(
				state.path.String(),
				DecodeErrorKindRange,
				runtime.Error(runtime.ErrInvalidArgumentType, "collection length overflows int"),
				true,
			)
		}

		if length > 0 {
			capacity = int(length)
		}
	}

	output := reflect.MakeSlice(dst.Type(), 0, capacity)

	for index := 0; ; index++ {
		value, _, nextErr := iterator.Next(state.ctx)
		if nextErr == io.EOF {
			break
		}

		if nextErr != nil {
			mark := state.path.PushIndex(index)
			err := newDecodeError(state.path.String(), DecodeErrorKindSource, nextErr, false)
			state.path.Restore(mark)

			return err
		}

		element := reflect.New(dst.Type().Elem()).Elem()
		mark := state.path.PushIndex(index)
		err := bindRuntimeValue(state, value, element)
		state.path.Restore(mark)

		if err != nil {
			return err
		}

		output = reflect.Append(output, element)
	}

	dst.Set(output)

	return nil
}

func bindRuntimeArray(state *decodeState, src runtime.Value, dst reflect.Value) (retErr error) {
	iterable, ok := src.(runtime.Iterable)
	if !ok {
		return decodeTypeError(&state.path, src, dst.Type())
	}

	iterator, err := iterable.Iterate(state.ctx)
	if err != nil {
		return newDecodeError(state.path.String(), DecodeErrorKindSource, err, false)
	}

	defer func() {
		retErr = joinDecodeErrors(retErr, closeIteratorAt(iterator, &state.path), &state.path)
	}()

	for index := 0; ; index++ {
		value, _, nextErr := iterator.Next(state.ctx)
		if nextErr == io.EOF {
			break
		}

		if nextErr != nil {
			mark := state.path.PushIndex(index)
			err := newDecodeError(state.path.String(), DecodeErrorKindSource, nextErr, false)
			state.path.Restore(mark)

			return err
		}

		if index >= dst.Len() {
			return newDecodeError(
				state.path.String(),
				DecodeErrorKindRange,
				runtime.Error(runtime.ErrInvalidArgumentType, "source has more elements than target array"),
				true,
			)
		}

		mark := state.path.PushIndex(index)
		err := bindRuntimeValue(state, value, dst.Index(index))
		state.path.Restore(mark)

		if err != nil {
			return err
		}
	}

	return nil
}

func bindRuntimeMap(state *decodeState, src runtime.Value, dst reflect.Value) error {
	if dst.Type().Key().Kind() != reflect.String {
		return newDecodeError(
			state.path.String(),
			DecodeErrorKindType,
			runtime.Error(runtime.ErrInvalidArgumentType, "map key type must be string"),
			true,
		)
	}

	entries, err := collectRuntimeEntries(state, src)
	if err != nil {
		return err
	}

	output := reflect.MakeMapWithSize(dst.Type(), len(entries))

	for key, value := range entries {
		element := reflect.New(dst.Type().Elem()).Elem()
		mark := state.path.PushKey(key)
		err := bindRuntimeValue(state, value, element)
		state.path.Restore(mark)

		if err != nil {
			return err
		}

		mapKey := reflect.ValueOf(key).Convert(dst.Type().Key())
		output.SetMapIndex(mapKey, element)
	}

	dst.Set(output)

	return nil
}

func bindRuntimeStruct(state *decodeState, src runtime.Value, dst reflect.Value) error {
	entries, err := collectRuntimeEntries(state, src)
	if err != nil {
		return err
	}

	if state.path.IsRoot() && state.config != nil && state.config.onlyFields != nil {
		if err := rejectDisallowedFields(state, entries); err != nil {
			return err
		}
	}

	lowerKeys := buildLowerRuntimeKeyMap(entries)
	used := make(map[string]struct{}, len(entries))
	visiting := make(map[reflect.Type]bool)

	if _, err := bindRuntimeStructEntries(state, dst, entries, lowerKeys, used, visiting); err != nil {
		return err
	}

	if state.config != nil && state.config.disallowUnknownFields {
		unknown := make([]string, 0)

		for key := range entries {
			if _, exists := used[key]; !exists {
				unknown = append(unknown, key)
			}
		}

		if len(unknown) > 0 {
			sort.Strings(unknown)

			return newDecodeError(
				state.path.String(),
				DecodeErrorKindUnknownField,
				runtime.Errorf(runtime.ErrInvalidArgument, "unknown field %q", unknown[0]),
				true,
			)
		}
	}

	return nil
}

func bindRuntimeStructEntries(
	state *decodeState,
	dst reflect.Value,
	entries map[string]runtime.Value,
	lowerKeys map[string]string,
	used map[string]struct{},
	visiting map[reflect.Type]bool,
) (bool, error) {
	typ := dst.Type()
	if visiting[typ] {
		return false, nil
	}

	visiting[typ] = true
	defer delete(visiting, typ)

	matched := false

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name, tagged := Tag(field)
		if !tagged {
			continue
		}

		value, actualKey, found := lookupRuntimeEntry(name, entries, lowerKeys, used)
		if !found {
			continue
		}

		mark := state.path.PushField(name)
		err := bindRuntimeValue(state, value, dst.Field(i))

		state.path.Restore(mark)

		if err != nil {
			return false, err
		}

		used[actualKey] = struct{}{}
		matched = true
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" || !field.Anonymous {
			continue
		}

		if _, tagged := Tag(field); tagged {
			continue
		}

		fieldValue := dst.Field(i)

		switch fieldValue.Kind() {
		case reflect.Struct:
			subMatched, err := bindRuntimeStructEntries(state, fieldValue, entries, lowerKeys, used, visiting)
			if err != nil {
				return false, err
			}

			matched = matched || subMatched
		case reflect.Pointer:
			if fieldValue.Type().Elem().Kind() != reflect.Struct || visiting[fieldValue.Type().Elem()] {
				continue
			}

			if fieldValue.IsNil() {
				element := reflect.New(fieldValue.Type().Elem())
				subMatched, err := bindRuntimeStructEntries(state, element.Elem(), entries, lowerKeys, used, visiting)
				if err != nil {
					return false, err
				}

				if subMatched {
					fieldValue.Set(element)
					matched = true
				}

				continue
			}

			subMatched, err := bindRuntimeStructEntries(state, fieldValue.Elem(), entries, lowerKeys, used, visiting)
			if err != nil {
				return false, err
			}

			matched = matched || subMatched
		}
	}

	return matched, nil
}

func bindRuntimeInterface(state *decodeState, src runtime.Value, dst reflect.Value) error {
	srcValue := reflect.ValueOf(src)

	if srcValue.IsValid() && srcValue.Type().AssignableTo(dst.Type()) {
		dst.Set(srcValue)

		return nil
	}

	return bindRuntimeUnwrapped(state, src, dst, false)
}

func bindRuntimeUnwrappedExact(src runtime.Value, dst reflect.Value) bool {
	unwrappable, ok := src.(runtime.Unwrappable)
	if !ok {
		return false
	}

	unwrapped := unwrappable.Unwrap()
	if unwrapped == nil {
		dst.Set(reflect.Zero(dst.Type()))

		return true
	}

	value := reflect.ValueOf(unwrapped)
	if !value.Type().AssignableTo(dst.Type()) {
		return false
	}

	dst.Set(value)

	return true
}

func bindRuntimeUnwrapped(state *decodeState, src runtime.Value, dst reflect.Value, allowConvert bool) error {
	unwrappable, ok := src.(runtime.Unwrappable)
	if !ok {
		return decodeTypeError(&state.path, src, dst.Type())
	}

	unwrapped := unwrappable.Unwrap()
	if unwrapped == nil {
		dst.Set(reflect.Zero(dst.Type()))

		return nil
	}

	value := reflect.ValueOf(unwrapped)
	if value.Type().AssignableTo(dst.Type()) {
		dst.Set(value)

		return nil
	}

	if allowConvert && value.Type().ConvertibleTo(dst.Type()) {
		dst.Set(value.Convert(dst.Type()))

		return nil
	}

	return decodeTypeError(&state.path, src, dst.Type())
}

func collectRuntimeEntries(state *decodeState, src runtime.Value) (out map[string]runtime.Value, retErr error) {
	input, readable := src.(runtime.KeyReadable)
	iterable, iterableOK := src.(runtime.Iterable)
	if !readable || !iterableOK {
		return nil, decodeTypeError(&state.path, src, reflect.TypeOf(map[string]any{}))
	}

	iterator, err := iterable.Iterate(state.ctx)
	if err != nil {
		return nil, newDecodeError(state.path.String(), DecodeErrorKindSource, err, false)
	}

	defer func() {
		retErr = joinDecodeErrors(retErr, closeIteratorAt(iterator, &state.path), &state.path)
	}()

	out = make(map[string]runtime.Value)
	for {
		_, keyValue, nextErr := iterator.Next(state.ctx)
		if nextErr == io.EOF {
			break
		}

		if nextErr != nil {
			return nil, newDecodeError(state.path.String(), DecodeErrorKindSource, nextErr, false)
		}

		key, ok := keyValue.(runtime.String)
		if !ok {
			return nil, newDecodeError(
				state.path.String(),
				DecodeErrorKindType,
				runtime.Error(runtime.ErrInvalidArgumentType, "map key type must be string"),
				true,
			)
		}

		value, getErr := input.Get(state.ctx, keyValue)
		if getErr != nil {
			mark := state.path.PushKey(key.String())
			err := newDecodeError(state.path.String(), DecodeErrorKindSource, getErr, false)

			state.path.Restore(mark)

			return nil, err
		}

		out[key.String()] = value
	}

	return out, nil
}

func buildLowerRuntimeKeyMap(entries map[string]runtime.Value) map[string]string {
	keys := make([]string, 0, len(entries))

	for key := range entries {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	lowerKeys := make(map[string]string, len(entries))
	for _, key := range keys {
		lower := strings.ToLower(key)

		if _, exists := lowerKeys[lower]; !exists {
			lowerKeys[lower] = key
		}
	}

	return lowerKeys
}

func lookupRuntimeEntry(
	name string,
	entries map[string]runtime.Value,
	lowerKeys map[string]string,
	used map[string]struct{},
) (runtime.Value, string, bool) {
	if _, taken := used[name]; !taken {
		if value, exists := entries[name]; exists {
			return value, name, true
		}
	}

	original, exists := lowerKeys[strings.ToLower(name)]
	if !exists {
		return nil, "", false
	}

	if _, taken := used[original]; taken {
		return nil, "", false
	}

	return entries[original], original, true
}

func closeIteratorAt(iterator runtime.Iterator, path *conversionPath) error {
	closer, ok := iterator.(io.Closer)
	if !ok {
		return nil
	}

	if err := closer.Close(); err != nil {
		return newDecodeError(
			path.String(),
			DecodeErrorKindSource,
			fmt.Errorf("close iterator: %w", err),
			false,
		)
	}

	return nil
}

func decodeTypeError(path *conversionPath, src runtime.Value, target reflect.Type) error {
	return newDecodeError(
		path.String(),
		DecodeErrorKindType,
		runtime.Errorf(runtime.ErrInvalidArgumentType, "cannot bind %s to %s", runtime.TypeOf(src), target),
		true,
	)
}

func validateDecodeConfig(config *decodeConfig, target reflect.Type) error {
	if config == nil {
		return nil
	}

	if config.err != nil {
		return config.err
	}

	if config.onlyFields == nil {
		return nil
	}

	for target.Kind() == reflect.Pointer {
		target = target.Elem()
	}

	if target.Kind() != reflect.Struct {
		return runtime.Error(runtime.ErrInvalidArgument, "OnlyFields requires a struct target")
	}

	return nil
}

func rejectDisallowedFields(state *decodeState, entries map[string]runtime.Value) error {
	unknown := make([]string, 0)

	for key := range entries {
		if _, allowed := state.config.onlyFields[strings.ToLower(key)]; !allowed {
			unknown = append(unknown, key)
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	sort.Strings(unknown)

	return newDecodeError(
		state.path.String(),
		DecodeErrorKindUnknownField,
		runtime.Errorf(runtime.ErrInvalidArgument, "unknown field %q", unknown[0]),
		true,
	)
}

func joinDecodeErrors(primary, secondary error, path *conversionPath) error {
	if secondary == nil {
		return primary
	}

	if primary == nil {
		return secondary
	}

	return newDecodeErrorWithoutPath(
		path.String(),
		DecodeErrorKindSource,
		errors.Join(primary, secondary),
		false,
	)
}

func validateDecodeTarget(target any) (reflect.Value, error) {
	value := reflect.ValueOf(target)
	if !value.IsValid() {
		return reflect.Value{}, runtime.Error(runtime.ErrInvalidArgumentType, "target is invalid")
	}

	if value.Kind() != reflect.Pointer {
		return reflect.Value{}, runtime.Error(runtime.ErrInvalidArgumentType, "target must be a pointer")
	}

	if value.IsNil() {
		return reflect.Value{}, runtime.Error(runtime.ErrInvalidArgumentType, "target must be a non-nil pointer")
	}

	return value, nil
}
