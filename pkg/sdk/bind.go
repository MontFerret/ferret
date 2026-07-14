package sdk

import (
	"context"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Bind0 adapts a typed zero-argument function to a Ferret host function.
func Bind0[R runtime.Value](fn func(context.Context) (R, error)) runtime.Function0 {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context) (runtime.Value, error) {
		result, err := fn(ctx)
		if err != nil {
			return runtime.None, err
		}

		return normalizeBoundResult(result), nil
	}
}

// Bind1 adapts a typed one-argument function to a Ferret host function.
func Bind1[A, R runtime.Value](fn func(context.Context, A) (R, error)) runtime.Function1 {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
		value, err := runtime.CastArg[A](arg, 0)
		if err != nil {
			return runtime.None, err
		}

		result, err := fn(ctx, value)
		if err != nil {
			return runtime.None, err
		}

		return normalizeBoundResult(result), nil
	}
}

// Bind2 adapts a typed two-argument function to a Ferret host function.
func Bind2[A, B, R runtime.Value](fn func(context.Context, A, B) (R, error)) runtime.Function2 {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context, first, second runtime.Value) (runtime.Value, error) {
		firstValue, err := runtime.CastArg[A](first, 0)
		if err != nil {
			return runtime.None, err
		}

		secondValue, err := runtime.CastArg[B](second, 1)
		if err != nil {
			return runtime.None, err
		}

		result, err := fn(ctx, firstValue, secondValue)
		if err != nil {
			return runtime.None, err
		}

		return normalizeBoundResult(result), nil
	}
}

// Bind3 adapts a typed three-argument function to a Ferret host function.
func Bind3[A, B, C, R runtime.Value](fn func(context.Context, A, B, C) (R, error)) runtime.Function3 {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context, first, second, third runtime.Value) (runtime.Value, error) {
		firstValue, err := runtime.CastArg[A](first, 0)
		if err != nil {
			return runtime.None, err
		}

		secondValue, err := runtime.CastArg[B](second, 1)
		if err != nil {
			return runtime.None, err
		}

		thirdValue, err := runtime.CastArg[C](third, 2)
		if err != nil {
			return runtime.None, err
		}

		result, err := fn(ctx, firstValue, secondValue, thirdValue)
		if err != nil {
			return runtime.None, err
		}

		return normalizeBoundResult(result), nil
	}
}

// Bind4 adapts a typed four-argument function to a Ferret host function.
func Bind4[A, B, C, D, R runtime.Value](fn func(context.Context, A, B, C, D) (R, error)) runtime.Function4 {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context, first, second, third, fourth runtime.Value) (runtime.Value, error) {
		firstValue, err := runtime.CastArg[A](first, 0)
		if err != nil {
			return runtime.None, err
		}

		secondValue, err := runtime.CastArg[B](second, 1)
		if err != nil {
			return runtime.None, err
		}

		thirdValue, err := runtime.CastArg[C](third, 2)
		if err != nil {
			return runtime.None, err
		}

		fourthValue, err := runtime.CastArg[D](fourth, 3)
		if err != nil {
			return runtime.None, err
		}

		result, err := fn(ctx, firstValue, secondValue, thirdValue, fourthValue)
		if err != nil {
			return runtime.None, err
		}

		return normalizeBoundResult(result), nil
	}
}

// Bind adapts a typed variadic function to a Ferret host function.
func Bind[T, R runtime.Value](fn func(context.Context, ...T) (R, error)) runtime.Function {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		typedArgs := make([]T, len(args))

		for i, arg := range args {
			value, err := runtime.CastArg[T](arg, i)
			if err != nil {
				return runtime.None, err
			}

			typedArgs[i] = value
		}

		result, err := fn(ctx, typedArgs...)
		if err != nil {
			return runtime.None, err
		}

		return normalizeBoundResult(result), nil
	}
}

func normalizeBoundResult[R runtime.Value](result R) runtime.Value {
	value := runtime.Value(result)
	reflected := reflect.ValueOf(value)

	if !reflected.IsValid() {
		return runtime.None
	}

	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if reflected.IsNil() {
			return runtime.None
		}
	}

	return value
}
